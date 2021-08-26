// +build linux_bpf

package tracer

/*
#include "../ebpf/c/tracer.h"
*/
import "C"
import (
	"strings"
	"unsafe"

	"github.com/DataDog/datadog-agent/pkg/network/ebpf/probes"
	"github.com/DataDog/datadog-agent/pkg/util/log"
	"github.com/DataDog/ebpf/manager"

	"github.com/davecgh/go-spew/spew"
)

func dumpMapsHandler(managerMap *manager.Map, manager *manager.Manager) string {
	var output strings.Builder

	mapName := managerMap.Name
	currentMap, found, err := manager.GetMap(mapName)
	if err != nil || !found {
		return ""
	}

	switch mapName {

	case "connectsock_ipv6": // maps/connectsock_ipv6 (BPF_MAP_TYPE_HASH), key C.__u64, value uintptr // C.void*
		output.WriteString("Map: '" + mapName + "', key: 'C.__u64', value: 'uintptr // C.void*'\n")
		iter := currentMap.Iterate()
		var key C.__u64
		var value uintptr // C.void*
		for iter.Next(unsafe.Pointer(&key), unsafe.Pointer(&value)) {
			output.WriteString(spew.Sdump(key, value))
		}

	case string(probes.TracerStatusMap): // maps/tracer_status (BPF_MAP_TYPE_HASH), key C.__u64, value tracerStatus
		output.WriteString("Map: '" + mapName + "', key: 'C.__u64', value: 'tracerStatus'\n")
		iter := currentMap.Iterate()
		var key C.__u64
		var value tracerStatus
		for iter.Next(unsafe.Pointer(&key), unsafe.Pointer(&value)) {
			output.WriteString(spew.Sdump(key, value))
		}

	case string(probes.ConntrackMap): // maps/conntrack (BPF_MAP_TYPE_HASH), key ConnTuple, value ConnTuple
		output.WriteString("Map: '" + mapName + "', key: 'ConnTuple', value: 'ConnTuple'\n")
		iter := currentMap.Iterate()
		var key ConnTuple
		var value ConnTuple
		for iter.Next(unsafe.Pointer(&key), unsafe.Pointer(&value)) {
			output.WriteString(spew.Sdump(key, value))
		}

	case string(probes.ConntrackTelemetryMap): // maps/conntrack_telemetry (BPF_MAP_TYPE_ARRAY), key C.u32, value conntrackTelemetry
		output.WriteString("Map: '" + mapName + "', key: 'C.u32', value: 'conntrackTelemetry'\n")
		telemetry := &conntrackTelemetry{}
		if err := currentMap.Lookup(unsafe.Pointer(&zero), unsafe.Pointer(telemetry)); err != nil {
			log.Tracef("error retrieving the contrack telemetry struct: %s", err)
		}
		output.WriteString(spew.Sdump(telemetry))

	case string(probes.SockFDLookupArgsMap): // maps/sockfd_lookup_args (BPF_MAP_TYPE_HASH), key C.__u64, value C.__u32
		output.WriteString("Map: '" + mapName + "', key: 'C.__u64', value: 'C.__u32'\n")
		iter := currentMap.Iterate()
		var key C.__u64
		var value C.__u32
		for iter.Next(unsafe.Pointer(&key), unsafe.Pointer(&value)) {
			output.WriteString(spew.Sdump(key, value))
		}

	case string(probes.SockByPidFDMap): // maps/sock_by_pid_fd (BPF_MAP_TYPE_HASH), key C.pid_fd_t, value uintptr // C.struct sock*
		output.WriteString("Map: '" + mapName + "', key: 'C.pid_fd_t', value: 'uintptr // C.struct sock*'\n")
		iter := currentMap.Iterate()
		var key C.pid_fd_t
		var value uintptr // C.struct sock*
		for iter.Next(unsafe.Pointer(&key), unsafe.Pointer(&value)) {
			output.WriteString(spew.Sdump(key, value))
		}

	case string(probes.PidFDBySockMap): // maps/pid_fd_by_sock (BPF_MAP_TYPE_HASH), key uintptr // C.struct sock*, value C.pid_fd_t
		output.WriteString("Map: '" + mapName + "', key: 'uintptr // C.struct sock*', value: 'C.pid_fd_t'\n")
		iter := currentMap.Iterate()
		var key uintptr // C.struct sock*
		var value C.pid_fd_t
		for iter.Next(unsafe.Pointer(&key), unsafe.Pointer(&value)) {
			output.WriteString(spew.Sdump(key, value))
		}

	case string(probes.ConnMap): // maps/conn_stats (BPF_MAP_TYPE_HASH), key ConnTuple, value ConnStatsWithTimestamp
		output.WriteString("Map: '" + mapName + "', key: 'ConnTuple', value: 'ConnStatsWithTimestamp'\n")
		iter := currentMap.Iterate()
		var key ConnTuple
		var value ConnStatsWithTimestamp
		for iter.Next(unsafe.Pointer(&key), unsafe.Pointer(&value)) {
			output.WriteString(spew.Sdump(key, value))
		}

	case string(probes.TcpStatsMap): // maps/tcp_stats (BPF_MAP_TYPE_HASH), key ConnTuple, value TCPStats
		output.WriteString("Map: '" + mapName + "', key: 'ConnTuple', value: 'TCPStats'\n")
		iter := currentMap.Iterate()
		var key ConnTuple
		var value TCPStats
		for iter.Next(unsafe.Pointer(&key), unsafe.Pointer(&value)) {
			output.WriteString(spew.Sdump(key, value))
		}

	case string(probes.ConnCloseBatchMap): // maps/conn_close_batch (BPF_MAP_TYPE_HASH), key C.__u32, value batch
		output.WriteString("Map: '" + mapName + "', key: 'C.__u32', value: 'batch'\n")
		iter := currentMap.Iterate()
		var key C.__u32
		var value batch
		for iter.Next(unsafe.Pointer(&key), unsafe.Pointer(&value)) {
			output.WriteString(spew.Sdump(key, value))
		}

	case "udp_recv_sock": // maps/udp_recv_sock (BPF_MAP_TYPE_HASH), key C.__u64, value C.udp_recv_sock_t
		output.WriteString("Map: '" + mapName + "', key: 'C.__u64', value: 'C.udp_recv_sock_t'\n")
		iter := currentMap.Iterate()
		var key C.__u64
		var value C.udp_recv_sock_t
		for iter.Next(unsafe.Pointer(&key), unsafe.Pointer(&value)) {
			output.WriteString(spew.Sdump(key, value))
		}

	case string(probes.PortBindingsMap): // maps/port_bindings (BPF_MAP_TYPE_HASH), key portBindingTuple, value C.__u8
		output.WriteString("Map: '" + mapName + "', key: 'portBindingTuple', value: 'C.__u8'\n")
		iter := currentMap.Iterate()
		var key portBindingTuple
		var value C.__u8
		for iter.Next(unsafe.Pointer(&key), unsafe.Pointer(&value)) {
			output.WriteString(spew.Sdump(key, value))
		}

	case string(probes.UdpPortBindingsMap): // maps/udp_port_bindings (BPF_MAP_TYPE_HASH), key portBindingTuple, value C.__u8
		output.WriteString("Map: '" + mapName + "', key: 'portBindingTuple', value: 'C.__u8'\n")
		iter := currentMap.Iterate()
		var key portBindingTuple
		var value C.__u8
		for iter.Next(unsafe.Pointer(&key), unsafe.Pointer(&value)) {
			output.WriteString(spew.Sdump(key, value))
		}

	case "pending_bind": // maps/pending_bind (BPF_MAP_TYPE_HASH), key C.__u64, value C.bind_syscall_args_t
		output.WriteString("Map: '" + mapName + "', key: 'C.__u64', value: 'C.bind_syscall_args_t'\n")
		iter := currentMap.Iterate()
		var key C.__u64
		var value C.bind_syscall_args_t
		for iter.Next(unsafe.Pointer(&key), unsafe.Pointer(&value)) {
			output.WriteString(spew.Sdump(key, value))
		}

	case string(probes.TelemetryMap): // maps/telemetry (BPF_MAP_TYPE_ARRAY), key C.u32, value kernelTelemetry
		output.WriteString("Map: '" + mapName + "', key: 'C.u32', value: 'kernelTelemetry'\n")
		telemetry := &kernelTelemetry{}
		if err := currentMap.Lookup(unsafe.Pointer(&zero), unsafe.Pointer(telemetry)); err != nil {
			// This can happen if we haven't initialized the telemetry object yet
			// so let's just use a trace log
			log.Tracef("error retrieving the telemetry struct: %s", err)
		}
		output.WriteString(spew.Sdump(telemetry))

	case "ip_route_output_flows": // maps/ip_route_output_flows (BPF_MAP_TYPE_HASH), key C.__u64, value C.ip_route_flow_t
		output.WriteString("Map: '" + mapName + "', key: 'C.__u64', value: 'C.ip_route_flow_t'\n")
		iter := currentMap.Iterate()
		var key C.__u64
		var value C.ip_route_flow_t
		for iter.Next(unsafe.Pointer(&key), unsafe.Pointer(&value)) {
			output.WriteString(spew.Sdump(key, value))
		}

	case string(probes.GatewayMap): // maps/ip_route_dest_gateways (BPF_MAP_TYPE_HASH), key ipRouteDest, value ipRouteGateway
		output.WriteString("Map: '" + mapName + "', key: 'ipRouteDest', value: 'ipRouteGateway'\n")
		iter := currentMap.Iterate()
		var key ipRouteDest
		var value ipRouteGateway
		for iter.Next(unsafe.Pointer(&key), unsafe.Pointer(&value)) {
			output.WriteString(spew.Sdump(key, value))
		}

	case string(probes.DoSendfileArgsMap): // maps/do_sendfile_args (BPF_MAP_TYPE_HASH), key C.__u64, value uintptr // C.struct sock*
		output.WriteString("Map: '" + mapName + "', key: 'C.__u64', value: 'uintptr // C.struct sock*'\n")
		iter := currentMap.Iterate()
		var key C.__u64
		var value uintptr // C.struct sock*
		for iter.Next(unsafe.Pointer(&key), unsafe.Pointer(&value)) {
			output.WriteString(spew.Sdump(key, value))
		}

	}

	return output.String()
}

func dumpPerfMapsHandler(managerMap *manager.PerfMap, manager *manager.Manager) string {
	var output strings.Builder
	mapName := managerMap.Name

	switch mapName {

	case string(probes.ConnCloseEventMap): // maps/conn_close_event (BPF_MAP_TYPE_PERF_EVENT_ARRAY), key C.__u32, value C.__u32
		output.WriteString("PerfMap: '" + mapName + "', key: 'C.__u32', value: 'C.__u32'\n")
		output.WriteString(spew.Sdump(managerMap.PerfMapStats))

	}
	return output.String()
}

func setupDumpHandler(manager *manager.Manager) {
	for _, m := range manager.Maps {
		m.DumpHandler = dumpMapsHandler
	}
	for _, m := range manager.PerfMaps {
		m.DumpHandler = dumpPerfMapsHandler
	}
}
