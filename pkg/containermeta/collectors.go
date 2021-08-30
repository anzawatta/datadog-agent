// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

package containermeta

import "context"

type Collector interface {
	Pull(context.Context) error
	Start(context.Context, chan []Event) error
}

type collectorFactory func() Collector

var collectorCatalog map[string]collectorFactory

func RegisterCollector(id string, c collectorFactory) {
	collectorCatalog[id] = c
}
