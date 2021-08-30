// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

package containermeta

import (
	"context"
	"sync"
	"time"

	"github.com/DataDog/datadog-agent/pkg/errors"
	"github.com/DataDog/datadog-agent/pkg/status/health"
	"github.com/DataDog/datadog-agent/pkg/util/log"
	"github.com/DataDog/datadog-agent/pkg/util/retry"
)

const (
	retryCollectorInterval = 30 * time.Second
	pullCollectorInterval  = 5 * time.Second
)

// Store contains metadata for container workloads.
type Store struct {
	mu         sync.RWMutex
	containers map[string]Container
	kubePods   map[string]KubernetesPod
	ecsTasks   map[string]ECSTask

	candidates map[string]Collector
	collectors []Collector

	eventCh chan []Event

	subscribers map[chan []Event]struct{}
}

// NewStore creates a new container metadata store. Call Run to start it.
func NewStore() *Store {
	candidates := make(map[string]Collector)
	for id, c := range collectorCatalog {
		candidates[id] = c()
	}

	return &Store{
		candidates:  candidates,
		collectors:  make([]Collector, len(candidates)),
		containers:  make(map[string]Container),
		kubePods:    make(map[string]KubernetesPod),
		ecsTasks:    make(map[string]ECSTask),
		subscribers: make(map[chan []Event]struct{}),
	}
}

// Run starts the container metadata store.
func (s *Store) Run(ctx context.Context) {
	retryTicker := time.NewTicker(retryCollectorInterval)
	pullTicker := time.NewTicker(pullCollectorInterval)
	health := health.RegisterLiveness("containermeta-store")

	// Dummy ctx and cancel func until the first pull starts
	pullCtx, pullCancel := context.WithCancel(ctx)

	go func() {
		for {
			select {
			case <-health.C:

			case <-pullTicker.C:
				// pullCtx will always be expired at this point
				// if pullTicker has the same duration as
				// pullCtx, so we cancel just as good practice
				pullCancel()

				pullCtx, pullCancel = context.WithTimeout(ctx, pullCollectorInterval)
				s.pull(pullCtx)

			case evs := <-s.eventCh:
				s.handleEvents(evs)

			case <-retryTicker.C:
				s.startCandidates(ctx)

				if len(s.candidates) == 0 {
					retryTicker.Stop()
				}

			case <-ctx.Done():
				retryTicker.Stop()
				pullTicker.Stop()

				err := health.Deregister()
				if err != nil {
					log.Warnf("error de-registering health check: %s", err)
				}

				return
			}
		}
	}()
}

// Subscribe returns a channel where container metadata events will be streamed
// as they happen.
func (s *Store) Subscribe() chan []Event {
	// this buffer size is an educated guess, as we know the rate of
	// updates, but not how fast these can be streamed out yet. it most
	// likely should be configurable.
	const bufferSize = 100

	// this is a `ch []Event` instead of a `ch Event` to improve
	// throughput, as bursts of events are as likely to occur as isolated
	// events, especially at startup or with collectors that periodically
	// pull changes.
	ch := make(chan []Event, bufferSize)

	s.mu.Lock()
	defer s.mu.Unlock()

	s.subscribers[ch] = struct{}{}

	// TODO(juliogreff): do the initial notifying of existing entities

	return ch
}

// Unsubscribe ends a subscription to entity events and closes its channel.
func (s *Store) Unsubscribe(ch chan []Event) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.subscribers, ch)
	close(ch)
}

// GetContainer returns metadata about a container.
func (s *Store) GetContainer(id string) (Container, error) {
	s.mu.RLock()
	defer s.mu.Unlock()

	c, ok := s.containers[id]
	if !ok {
		return c, errors.NewNotFound(id)
	}

	return c, nil
}

// GetKubernetesPod returns metadata about a Kubernetes pod.
func (s *Store) GetKubernetesPod(id string) (KubernetesPod, error) {
	s.mu.RLock()
	defer s.mu.Unlock()

	p, ok := s.kubePods[id]
	if !ok {
		return p, errors.NewNotFound(id)
	}

	return p, nil
}

// GetECSTask returns metadata about an ECS task.
func (s *Store) GetECSTask(id string) (ECSTask, error) {
	s.mu.RLock()
	defer s.mu.Unlock()

	t, ok := s.ecsTasks[id]
	if !ok {
		return t, errors.NewNotFound(id)
	}

	return t, nil
}

func (s *Store) startCandidates(ctx context.Context) {
	for id, c := range s.candidates {
		err := c.Start(ctx, s.eventCh)

		// Leave candidates that returned a retriable error to be
		// re-started in the next tick
		if err != nil && retry.IsErrWillRetry(err) {
			continue
		}

		// Store successfully started collectors for future reference
		if err == nil {
			s.collectors = append(s.collectors, c)
		}

		// Remove non-retriable and successfully started collectors
		// from the list of candidates so they're not retried in the
		// next tick
		delete(s.candidates, id)
	}
}

func (s *Store) pull(ctx context.Context) {
	// NOTE: s.collectors is not guarded by a mutex as it's only called by
	// the store itself, and the store runs on a single goroutine. If this
	// method is made public in the future, we need to guard it.
	for _, c := range s.collectors {
		// Run each pull in its own separate goroutine to reduce
		// latency and unlock the main goroutine to do other work.
		go func(c Collector) {
			err := c.Pull(ctx)
			if err != nil {
				log.Warnf("error pulling from collector: %s", err.Error())
			}
		}(c)
	}
}

func (s *Store) handleEvents(evs []Event) {
	for _, ev := range evs {
		meta := ev.Entity.GetMeta()

		switch meta.Kind {
		case KindContainer:
			s.handleContainerEvent(ev)
		case KindKubernetesPod:
			s.handleKubePodEvent(ev)
		case KindECSTask:
			s.handleECSTask(ev)
		default:
			log.Errorf("cannot handle event for entity %q with kind %q", meta.ID, meta.Kind)
			continue
		}
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	for ch, _ := range s.subscribers {
		ch <- evs
	}
}

func (s *Store) handleContainerEvent(ev Event) {
	s.mu.Lock()
	defer s.mu.Unlock()

	c := ev.Entity.(Container)

	switch ev.Type {
	case EventTypeSet:
		s.containers[c.ID] = c
	case EventTypeUnset:
		delete(s.containers, c.ID)
	}
}

func (s *Store) handleKubePodEvent(ev Event) {
	s.mu.Lock()
	defer s.mu.Unlock()

	p := ev.Entity.(KubernetesPod)

	switch ev.Type {
	case EventTypeSet:
		s.kubePods[p.ID] = p
	case EventTypeUnset:
		delete(s.kubePods, p.ID)
	}
}

func (s *Store) handleECSTask(ev Event) {
	s.mu.Lock()
	defer s.mu.Unlock()

	t := ev.Entity.(ECSTask)

	switch ev.Type {
	case EventTypeSet:
		s.ecsTasks[t.ID] = t
	case EventTypeUnset:
		delete(s.kubePods, t.ID)
	}
}
