// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

package kubelet

import (
	"context"
	"errors"
	"time"

	"github.com/DataDog/datadog-agent/pkg/config"
	"github.com/DataDog/datadog-agent/pkg/containermeta"
	"github.com/DataDog/datadog-agent/pkg/util/kubernetes/kubelet"
)

const (
	collectorID = "kubelet"
	expireFreq  = 5 * time.Minute
)

type collector struct {
	watcher    *kubelet.PodWatcher
	ch         chan []containermeta.Event
	lastExpire time.Time
	expireFreq time.Duration
}

func init() {
	containermeta.RegisterCollector(collectorID, func() containermeta.Collector {
		return &collector{}
	})
}

func (c *collector) Start(ctx context.Context, ch chan []containermeta.Event) error {
	if !config.IsKubernetes() {
		return errors.New("the Agent is not running in Kubernetes")
	}

	var err error

	c.ch = ch
	c.lastExpire = time.Now()
	c.expireFreq = expireFreq
	c.watcher, err = kubelet.NewPodWatcher(5*time.Minute, true)
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			}
		}
	}()

	return nil
}

func (c *collector) Pull(ctx context.Context) error {
	events := []containermeta.Event{}

	updatedPods, err := c.watcher.PullChanges(ctx)
	if err != nil {
		return err
	}

	events = append(events, c.parsePods(updatedPods)...)

	if time.Now().Sub(c.lastExpire) >= c.expireFreq {
		expireList, err := c.watcher.Expire()
		if err != nil {
			return err
		}

		events = append(events, c.parseExpires(expireList)...)

		c.lastExpire = time.Now()
	}

	if len(events) > 0 {
		c.ch <- events
	}

	return nil
}

func (c *collector) parsePods(pods []*kubelet.Pod) []containermeta.Event {
	return nil
}

func (c *collector) parseExpires(pods []string) []containermeta.Event {
	return nil
}
