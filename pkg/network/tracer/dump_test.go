// +build linux_bpf

package tracer

import (
	"testing"
	//	"github.com/DataDog/ebpf"
	netebpf "github.com/DataDog/datadog-agent/pkg/network/ebpf"
)

func TestManagerDumpHandler(t *testing.T) {
	//	closedHandler *ebpf.PerfHandler, runtimeTracer bool) *manager.Manager {
	m := netebpf.NewManager(nil, false)
	setupDumpHandler(m)

	t.Log("==================ok")
	for _, emap := range m.Maps {
		t.Log(emap.DumpHandler)
	}
}
