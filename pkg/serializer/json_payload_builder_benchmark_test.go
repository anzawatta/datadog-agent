// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-present Datadog, Inc.

//+build zlib

package serializer

import (
	"fmt"
	"testing"

	"github.com/DataDog/datadog-agent/pkg/forwarder"
	"github.com/DataDog/datadog-agent/pkg/metrics"
	"github.com/DataDog/datadog-agent/pkg/serializer/marshaler"
	"github.com/DataDog/datadog-agent/pkg/serializer/stream"
	"github.com/stretchr/testify/require"
)

func generateData(points int, items int, tags int) metrics.Series {
	series := metrics.Series{}
	for i := 0; i < items; i++ {
		series = append(series, &metrics.Serie{
			Points: func() []metrics.Point {
				ps := make([]metrics.Point, points)
				for p := 0; p < points; p++ {
					ps[p] = metrics.Point{Ts: float64(p * i), Value: float64(p + i)}
				}
				return ps
			}(),
			MType:    metrics.APIGaugeType,
			Name:     "test.metrics",
			Interval: 15,
			Host:     "localHost",
			Tags: func() []string {
				ts := make([]string, tags)
				for t := 0; t < tags; t++ {
					ts[t] = fmt.Sprintf("tag%d:foobar", t)
				}
				return ts
			}(),
		})
	}
	return series
}

var payloads forwarder.Payloads

func benchmarkSeriesJSerializationUsage(b *testing.B, points int, items int, tags int) {
	series := generateData(points, items, tags)
	payloadBuilder := stream.NewJSONPayloadBuilder(true)

	b.ResetTimer()
	b.ReportAllocs()

	var payloadCount int
	var payloadCompressedSize uint64
	for n := 0; n < b.N; n++ {
		var err error
		payloads, err = payloadBuilder.Build(series)
		payloadCount += len(payloads)
		for _, pl := range payloads {
			payloadCompressedSize += uint64(len(*pl))
		}
		require.NoError(b, err)
	}
	b.ReportMetric(float64(payloadCount)/float64(b.N), "payloads")
	b.ReportMetric(float64(payloadCompressedSize)/float64(b.N), "compressed-payload-bytes")
}

func benchmarkSeriesPSerializationUsage(b *testing.B, points int, items int, tags int) {
	series := generateData(points, items, tags)

	b.ResetTimer()
	b.ReportAllocs()

	var payloadCount int
	var payloadCompressedSize uint64
	for n := 0; n < b.N; n++ {
		var err error
		payloads, err = series.MarshalSplitCompress(marshaler.DefaultBufferContext())
		payloadCount += len(payloads)
		for _, pl := range payloads {
			payloadCompressedSize += uint64(len(*pl))
		}
		require.NoError(b, err)
	}
	b.ReportMetric(float64(payloadCount)/float64(b.N), "payloads")
	b.ReportMetric(float64(payloadCompressedSize)/float64(b.N), "compressed-payload-bytes")
}
func BenchmarkSeriesP_000005Items_000005Points_000010Tags(b *testing.B) {
	benchmarkSeriesPSerializationUsage(b, 5, 5, 10)
}
func BenchmarkSeriesP_000005Items_000005Points_000050Tags(b *testing.B) {
	benchmarkSeriesPSerializationUsage(b, 5, 5, 50)
}
func BenchmarkSeriesP_000005Items_000010Points_000010Tags(b *testing.B) {
	benchmarkSeriesPSerializationUsage(b, 10, 5, 10)
}
func BenchmarkSeriesP_000005Items_000010Points_000050Tags(b *testing.B) {
	benchmarkSeriesPSerializationUsage(b, 10, 5, 50)
}
func BenchmarkSeriesP_000005Items_000100Points_000010Tags(b *testing.B) {
	benchmarkSeriesPSerializationUsage(b, 100, 5, 10)
}
func BenchmarkSeriesP_000005Items_000100Points_000050Tags(b *testing.B) {
	benchmarkSeriesPSerializationUsage(b, 100, 5, 50)
}
func BenchmarkSeriesP_000005Items_001000Points_000010Tags(b *testing.B) {
	benchmarkSeriesPSerializationUsage(b, 1000, 5, 10)
}
func BenchmarkSeriesP_000005Items_001000Points_000050Tags(b *testing.B) {
	benchmarkSeriesPSerializationUsage(b, 1000, 5, 50)
}
func BenchmarkSeriesP_000010Items_000005Points_000010Tags(b *testing.B) {
	benchmarkSeriesPSerializationUsage(b, 5, 10, 10)
}
func BenchmarkSeriesP_000010Items_000005Points_000050Tags(b *testing.B) {
	benchmarkSeriesPSerializationUsage(b, 5, 10, 50)
}
func BenchmarkSeriesP_000010Items_000010Points_000010Tags(b *testing.B) {
	benchmarkSeriesPSerializationUsage(b, 10, 10, 10)
}
func BenchmarkSeriesP_000010Items_000010Points_000050Tags(b *testing.B) {
	benchmarkSeriesPSerializationUsage(b, 10, 10, 50)
}
func BenchmarkSeriesP_000010Items_000100Points_000010Tags(b *testing.B) {
	benchmarkSeriesPSerializationUsage(b, 100, 10, 10)
}
func BenchmarkSeriesP_000010Items_000100Points_000050Tags(b *testing.B) {
	benchmarkSeriesPSerializationUsage(b, 100, 10, 50)
}
func BenchmarkSeriesP_000010Items_001000Points_000010Tags(b *testing.B) {
	benchmarkSeriesPSerializationUsage(b, 1000, 10, 10)
}
func BenchmarkSeriesP_000010Items_001000Points_000050Tags(b *testing.B) {
	benchmarkSeriesPSerializationUsage(b, 1000, 10, 50)
}
func BenchmarkSeriesP_000100Items_000005Points_000010Tags(b *testing.B) {
	benchmarkSeriesPSerializationUsage(b, 5, 100, 10)
}
func BenchmarkSeriesP_000100Items_000005Points_000050Tags(b *testing.B) {
	benchmarkSeriesPSerializationUsage(b, 5, 100, 50)
}
func BenchmarkSeriesP_000100Items_000010Points_000010Tags(b *testing.B) {
	benchmarkSeriesPSerializationUsage(b, 10, 100, 10)
}
func BenchmarkSeriesP_000100Items_000010Points_000050Tags(b *testing.B) {
	benchmarkSeriesPSerializationUsage(b, 10, 100, 50)
}
func BenchmarkSeriesP_000100Items_000100Points_000010Tags(b *testing.B) {
	benchmarkSeriesPSerializationUsage(b, 100, 100, 10)
}
func BenchmarkSeriesP_000100Items_000100Points_000050Tags(b *testing.B) {
	benchmarkSeriesPSerializationUsage(b, 100, 100, 50)
}
func BenchmarkSeriesP_000100Items_001000Points_000010Tags(b *testing.B) {
	benchmarkSeriesPSerializationUsage(b, 1000, 100, 10)
}
func BenchmarkSeriesP_000100Items_001000Points_000050Tags(b *testing.B) {
	benchmarkSeriesPSerializationUsage(b, 1000, 100, 50)
}
func BenchmarkSeriesP_000500Items_000005Points_000010Tags(b *testing.B) {
	benchmarkSeriesPSerializationUsage(b, 5, 100, 10)
}
func BenchmarkSeriesP_000500Items_000005Points_000050Tags(b *testing.B) {
	benchmarkSeriesPSerializationUsage(b, 5, 100, 50)
}
func BenchmarkSeriesP_000500Items_000010Points_000010Tags(b *testing.B) {
	benchmarkSeriesPSerializationUsage(b, 10, 100, 10)
}
func BenchmarkSeriesP_000500Items_000010Points_000050Tags(b *testing.B) {
	benchmarkSeriesPSerializationUsage(b, 10, 100, 50)
}
func BenchmarkSeriesP_000500Items_000100Points_000010Tags(b *testing.B) {
	benchmarkSeriesPSerializationUsage(b, 100, 100, 10)
}
func BenchmarkSeriesP_000500Items_000100Points_000050Tags(b *testing.B) {
	benchmarkSeriesPSerializationUsage(b, 100, 100, 50)
}
func BenchmarkSeriesP_000500Items_001000Points_000010Tags(b *testing.B) {
	benchmarkSeriesPSerializationUsage(b, 1000, 100, 10)
}
func BenchmarkSeriesP_000500Items_001000Points_000050Tags(b *testing.B) {
	benchmarkSeriesPSerializationUsage(b, 1000, 100, 50)
}
func BenchmarkSeriesP_001000Items_000005Points_000010Tags(b *testing.B) {
	benchmarkSeriesPSerializationUsage(b, 5, 1000, 10)
}
func BenchmarkSeriesP_001000Items_000005Points_000050Tags(b *testing.B) {
	benchmarkSeriesPSerializationUsage(b, 5, 1000, 50)
}
func BenchmarkSeriesP_001000Items_000010Points_000010Tags(b *testing.B) {
	benchmarkSeriesPSerializationUsage(b, 10, 1000, 10)
}
func BenchmarkSeriesP_001000Items_000010Points_000050Tags(b *testing.B) {
	benchmarkSeriesPSerializationUsage(b, 10, 1000, 50)
}
func BenchmarkSeriesP_001000Items_000100Points_000010Tags(b *testing.B) {
	benchmarkSeriesPSerializationUsage(b, 100, 1000, 10)
}
func BenchmarkSeriesP_001000Items_000100Points_000050Tags(b *testing.B) {
	benchmarkSeriesPSerializationUsage(b, 100, 1000, 50)
}
func BenchmarkSeriesP_001000Items_001000Points_000010Tags(b *testing.B) {
	benchmarkSeriesPSerializationUsage(b, 1000, 1000, 10)
}
func BenchmarkSeriesP_001000Items_001000Points_000050Tags(b *testing.B) {
	benchmarkSeriesPSerializationUsage(b, 1000, 1000, 50)
}
func BenchmarkSeriesP_010000Items_000005Points_000010Tags(b *testing.B) {
	benchmarkSeriesPSerializationUsage(b, 5, 10000, 10)
}
func BenchmarkSeriesP_010000Items_000005Points_000050Tags(b *testing.B) {
	benchmarkSeriesPSerializationUsage(b, 5, 10000, 50)
}
func BenchmarkSeriesP_010000Items_000010Points_000010Tags(b *testing.B) {
	benchmarkSeriesPSerializationUsage(b, 10, 10000, 10)
}
func BenchmarkSeriesP_010000Items_000010Points_000050Tags(b *testing.B) {
	benchmarkSeriesPSerializationUsage(b, 10, 10000, 50)
}
func BenchmarkSeriesP_010000Items_000100Points_000010Tags(b *testing.B) {
	benchmarkSeriesPSerializationUsage(b, 100, 10000, 10)
}
func BenchmarkSeriesP_010000Items_000100Points_000050Tags(b *testing.B) {
	benchmarkSeriesPSerializationUsage(b, 100, 10000, 50)
}
func BenchmarkSeriesP_010000Items_001000Points_000010Tags(b *testing.B) {
	benchmarkSeriesPSerializationUsage(b, 1000, 10000, 10)
}
func BenchmarkSeriesP_010000Items_001000Points_000050Tags(b *testing.B) {
	benchmarkSeriesPSerializationUsage(b, 1000, 10000, 50)
}
func BenchmarkSeriesP_100000Items_000005Points_000010Tags(b *testing.B) {
	benchmarkSeriesPSerializationUsage(b, 5, 100000, 10)
}
func BenchmarkSeriesP_100000Items_000005Points_000050Tags(b *testing.B) {
	benchmarkSeriesPSerializationUsage(b, 5, 100000, 50)
}
func BenchmarkSeriesP_100000Items_000010Points_000010Tags(b *testing.B) {
	benchmarkSeriesPSerializationUsage(b, 10, 100000, 10)
}
func BenchmarkSeriesP_100000Items_000010Points_000050Tags(b *testing.B) {
	benchmarkSeriesPSerializationUsage(b, 10, 100000, 50)
}
func BenchmarkSeriesP_100000Items_000100Points_000010Tags(b *testing.B) {
	benchmarkSeriesPSerializationUsage(b, 100, 100000, 10)
}
func BenchmarkSeriesP_100000Items_000100Points_000050Tags(b *testing.B) {
	benchmarkSeriesPSerializationUsage(b, 100, 100000, 50)
}
func BenchmarkSeriesP_100000Items_001000Points_000010Tags(b *testing.B) {
	benchmarkSeriesPSerializationUsage(b, 1000, 100000, 10)
}

func BenchmarkSeriesJ_000005Items_000005Points_000010Tags(b *testing.B) {
	benchmarkSeriesJSerializationUsage(b, 5, 5, 10)
}
func BenchmarkSeriesJ_000005Items_000005Points_000050Tags(b *testing.B) {
	benchmarkSeriesJSerializationUsage(b, 5, 5, 50)
}
func BenchmarkSeriesJ_000005Items_000010Points_000010Tags(b *testing.B) {
	benchmarkSeriesJSerializationUsage(b, 10, 5, 10)
}
func BenchmarkSeriesJ_000005Items_000010Points_000050Tags(b *testing.B) {
	benchmarkSeriesJSerializationUsage(b, 10, 5, 50)
}
func BenchmarkSeriesJ_000005Items_000100Points_000010Tags(b *testing.B) {
	benchmarkSeriesJSerializationUsage(b, 100, 5, 10)
}
func BenchmarkSeriesJ_000005Items_000100Points_000050Tags(b *testing.B) {
	benchmarkSeriesJSerializationUsage(b, 100, 5, 50)
}
func BenchmarkSeriesJ_000005Items_001000Points_000010Tags(b *testing.B) {
	benchmarkSeriesJSerializationUsage(b, 1000, 5, 10)
}
func BenchmarkSeriesJ_000005Items_001000Points_000050Tags(b *testing.B) {
	benchmarkSeriesJSerializationUsage(b, 1000, 5, 50)
}
func BenchmarkSeriesJ_000010Items_000005Points_000010Tags(b *testing.B) {
	benchmarkSeriesJSerializationUsage(b, 5, 10, 10)
}
func BenchmarkSeriesJ_000010Items_000005Points_000050Tags(b *testing.B) {
	benchmarkSeriesJSerializationUsage(b, 5, 10, 50)
}
func BenchmarkSeriesJ_000010Items_000010Points_000010Tags(b *testing.B) {
	benchmarkSeriesJSerializationUsage(b, 10, 10, 10)
}
func BenchmarkSeriesJ_000010Items_000010Points_000050Tags(b *testing.B) {
	benchmarkSeriesJSerializationUsage(b, 10, 10, 50)
}
func BenchmarkSeriesJ_000010Items_000100Points_000010Tags(b *testing.B) {
	benchmarkSeriesJSerializationUsage(b, 100, 10, 10)
}
func BenchmarkSeriesJ_000010Items_000100Points_000050Tags(b *testing.B) {
	benchmarkSeriesJSerializationUsage(b, 100, 10, 50)
}
func BenchmarkSeriesJ_000010Items_001000Points_000010Tags(b *testing.B) {
	benchmarkSeriesJSerializationUsage(b, 1000, 10, 10)
}
func BenchmarkSeriesJ_000010Items_001000Points_000050Tags(b *testing.B) {
	benchmarkSeriesJSerializationUsage(b, 1000, 10, 50)
}
func BenchmarkSeriesJ_000100Items_000005Points_000010Tags(b *testing.B) {
	benchmarkSeriesJSerializationUsage(b, 5, 100, 10)
}
func BenchmarkSeriesJ_000100Items_000005Points_000050Tags(b *testing.B) {
	benchmarkSeriesJSerializationUsage(b, 5, 100, 50)
}
func BenchmarkSeriesJ_000100Items_000010Points_000010Tags(b *testing.B) {
	benchmarkSeriesJSerializationUsage(b, 10, 100, 10)
}
func BenchmarkSeriesJ_000100Items_000010Points_000050Tags(b *testing.B) {
	benchmarkSeriesJSerializationUsage(b, 10, 100, 50)
}
func BenchmarkSeriesJ_000100Items_000100Points_000010Tags(b *testing.B) {
	benchmarkSeriesJSerializationUsage(b, 100, 100, 10)
}
func BenchmarkSeriesJ_000100Items_000100Points_000050Tags(b *testing.B) {
	benchmarkSeriesJSerializationUsage(b, 100, 100, 50)
}
func BenchmarkSeriesJ_000100Items_001000Points_000010Tags(b *testing.B) {
	benchmarkSeriesJSerializationUsage(b, 1000, 100, 10)
}
func BenchmarkSeriesJ_000100Items_001000Points_000050Tags(b *testing.B) {
	benchmarkSeriesJSerializationUsage(b, 1000, 100, 50)
}
func BenchmarkSeriesJ_000500Items_000005Points_000010Tags(b *testing.B) {
	benchmarkSeriesJSerializationUsage(b, 5, 100, 10)
}
func BenchmarkSeriesJ_000500Items_000005Points_000050Tags(b *testing.B) {
	benchmarkSeriesJSerializationUsage(b, 5, 100, 50)
}
func BenchmarkSeriesJ_000500Items_000010Points_000010Tags(b *testing.B) {
	benchmarkSeriesJSerializationUsage(b, 10, 100, 10)
}
func BenchmarkSeriesJ_000500Items_000010Points_000050Tags(b *testing.B) {
	benchmarkSeriesJSerializationUsage(b, 10, 100, 50)
}
func BenchmarkSeriesJ_000500Items_000100Points_000010Tags(b *testing.B) {
	benchmarkSeriesJSerializationUsage(b, 100, 100, 10)
}
func BenchmarkSeriesJ_000500Items_000100Points_000050Tags(b *testing.B) {
	benchmarkSeriesJSerializationUsage(b, 100, 100, 50)
}
func BenchmarkSeriesJ_000500Items_001000Points_000010Tags(b *testing.B) {
	benchmarkSeriesJSerializationUsage(b, 1000, 100, 10)
}
func BenchmarkSeriesJ_000500Items_001000Points_000050Tags(b *testing.B) {
	benchmarkSeriesJSerializationUsage(b, 1000, 100, 50)
}
func BenchmarkSeriesJ_001000Items_000005Points_000010Tags(b *testing.B) {
	benchmarkSeriesJSerializationUsage(b, 5, 1000, 10)
}
func BenchmarkSeriesJ_001000Items_000005Points_000050Tags(b *testing.B) {
	benchmarkSeriesJSerializationUsage(b, 5, 1000, 50)
}
func BenchmarkSeriesJ_001000Items_000010Points_000010Tags(b *testing.B) {
	benchmarkSeriesJSerializationUsage(b, 10, 1000, 10)
}
func BenchmarkSeriesJ_001000Items_000010Points_000050Tags(b *testing.B) {
	benchmarkSeriesJSerializationUsage(b, 10, 1000, 50)
}
func BenchmarkSeriesJ_001000Items_000100Points_000010Tags(b *testing.B) {
	benchmarkSeriesJSerializationUsage(b, 100, 1000, 10)
}
func BenchmarkSeriesJ_001000Items_000100Points_000050Tags(b *testing.B) {
	benchmarkSeriesJSerializationUsage(b, 100, 1000, 50)
}
func BenchmarkSeriesJ_001000Items_001000Points_000010Tags(b *testing.B) {
	benchmarkSeriesJSerializationUsage(b, 1000, 1000, 10)
}
func BenchmarkSeriesJ_001000Items_001000Points_000050Tags(b *testing.B) {
	benchmarkSeriesJSerializationUsage(b, 1000, 1000, 50)
}
func BenchmarkSeriesJ_010000Items_000005Points_000010Tags(b *testing.B) {
	benchmarkSeriesJSerializationUsage(b, 5, 10000, 10)
}
func BenchmarkSeriesJ_010000Items_000005Points_000050Tags(b *testing.B) {
	benchmarkSeriesJSerializationUsage(b, 5, 10000, 50)
}
func BenchmarkSeriesJ_010000Items_000010Points_000010Tags(b *testing.B) {
	benchmarkSeriesJSerializationUsage(b, 10, 10000, 10)
}
func BenchmarkSeriesJ_010000Items_000010Points_000050Tags(b *testing.B) {
	benchmarkSeriesJSerializationUsage(b, 10, 10000, 50)
}
func BenchmarkSeriesJ_010000Items_000100Points_000010Tags(b *testing.B) {
	benchmarkSeriesJSerializationUsage(b, 100, 10000, 10)
}
func BenchmarkSeriesJ_010000Items_000100Points_000050Tags(b *testing.B) {
	benchmarkSeriesJSerializationUsage(b, 100, 10000, 50)
}
func BenchmarkSeriesJ_010000Items_001000Points_000010Tags(b *testing.B) {
	benchmarkSeriesJSerializationUsage(b, 1000, 10000, 10)
}
func BenchmarkSeriesJ_010000Items_001000Points_000050Tags(b *testing.B) {
	benchmarkSeriesJSerializationUsage(b, 1000, 10000, 50)
}
func BenchmarkSeriesJ_100000Items_000005Points_000010Tags(b *testing.B) {
	benchmarkSeriesJSerializationUsage(b, 5, 100000, 10)
}
func BenchmarkSeriesJ_100000Items_000005Points_000050Tags(b *testing.B) {
	benchmarkSeriesJSerializationUsage(b, 5, 100000, 50)
}
func BenchmarkSeriesJ_100000Items_000010Points_000010Tags(b *testing.B) {
	benchmarkSeriesJSerializationUsage(b, 10, 100000, 10)
}
func BenchmarkSeriesJ_100000Items_000010Points_000050Tags(b *testing.B) {
	benchmarkSeriesJSerializationUsage(b, 10, 100000, 50)
}
func BenchmarkSeriesJ_100000Items_000100Points_000010Tags(b *testing.B) {
	benchmarkSeriesJSerializationUsage(b, 100, 100000, 10)
}
func BenchmarkSeriesJ_100000Items_000100Points_000050Tags(b *testing.B) {
	benchmarkSeriesJSerializationUsage(b, 100, 100000, 50)
}
func BenchmarkSeriesJ_100000Items_001000Points_000010Tags(b *testing.B) {
	benchmarkSeriesJSerializationUsage(b, 1000, 100000, 10)
}
