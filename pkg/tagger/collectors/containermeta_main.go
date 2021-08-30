// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

package collectors

import (
	"context"

	"github.com/DataDog/datadog-agent/pkg/containermeta"
	"github.com/DataDog/datadog-agent/pkg/status/health"
	"github.com/DataDog/datadog-agent/pkg/util/log"
)

const (
	containermetaCollectorName = "containermeta"
)

type ContainerMetaCollector struct {
	store *containermeta.Store
	out   chan<- []*TagInfo
	ch    chan []containermeta.Event
	stop  chan struct{}
}

func (c *ContainerMetaCollector) Detect(ctx context.Context, out chan<- []*TagInfo) (CollectionMode, error) {
	c.out = out
	c.stop = make(chan struct{})

	return StreamCollection, nil
}

func (c *ContainerMetaCollector) Stream() error {
	ctx, cancel := context.WithCancel(context.Background())
	health := health.RegisterLiveness("tagger-containermeta")

	c.ch = c.store.Subscribe()

	for {
		select {
		case evs := <-c.ch:
			c.processEvents(ctx, evs)

		case <-health.C:

		case <-c.stop:
			cancel()

			err := health.Deregister()
			if err != nil {
				log.Warnf("error de-registering health check: %s", err)
			}

			return nil
		}
	}
}

func (c *ContainerMetaCollector) Stop() error {
	c.store.Unsubscribe(c.ch)
	c.stop <- struct{}{}
	return nil
}

func (c *ContainerMetaCollector) Fetch(ctx context.Context, entity string) ([]string, []string, []string, error) {
	// ContainerMetaCollector does not implement Fetch to prevent expensive
	// and race-condition prone forcing of pulls from upstream collectors.
	// Since containermeta.Store will eventually own notifying all
	// downstream consumers, this codepath should never trigger anyway.
	return nil, nil, nil, nil
}

func containermetaFactory() Collector {
	// TODO(juliogreff): use a global instance of containermeta.Store
	// instead of instantiating a new one.
	return &ContainerMetaCollector{
		store: containermeta.NewStore(),
	}
}

func init() {
	// NOTE: ContainerMetaCollector is meant to be used as the single
	// collector, so priority doesn't matter and should be removed entirely
	// after migration is done.

	// NOTE: uncomment this and remove all other collector registrations
	// once ContainerMetaCollector is ready for release.
	registerCollector(containermetaCollectorName, containermetaFactory, NodeRuntime)
}
