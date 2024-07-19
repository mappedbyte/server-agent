package sys_stats

import (
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/mappedbyte/server-agent/proto"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
	"net"
	"strconv"
	"strings"
	"time"
)

func GetSysStats() *proto.SysStats {
	percents, _ := cpu.Percent(1*time.Second, false)
	percent := 0.0
	if len(percents) > 0 {
		percent = percents[0]
	}
	ips, _ := getLocalIPs()
	ip := strings.Join(ips, ",")
	physicalCores, _ := cpu.Counts(false)
	logicalCores, _ := cpu.Counts(true)
	cpus := &proto.CpuStat{
		UsedPercent:   float32(percent),
		PhysicalCores: int32(physicalCores),
		LogicalCores:  int32(logicalCores),
	}

	memInfo, _ := mem.VirtualMemory()
	virtualMemory := &proto.MemoryStat{
		Total:       int64(memInfo.Total),
		Used:        int64(memInfo.Used),
		UsedPercent: float32(memInfo.UsedPercent),
	}
	swapMemoryInfo, _ := mem.SwapMemory()
	swapMemory := &proto.MemoryStat{
		Total:       int64(swapMemoryInfo.Total),
		Used:        int64(swapMemoryInfo.Used),
		UsedPercent: float32(swapMemoryInfo.UsedPercent),
	}
	partitions, _ := disk.Partitions(false)
	partitionStat := partitions[0]
	usage, _ := disk.Usage(partitionStat.Mountpoint)
	diskInfo := &proto.DiskStat{
		Total:       int64(usage.Total),
		Used:        int64(usage.Used),
		UsedPercent: float32(usage.UsedPercent),
	}

	now := time.Now()
	return &proto.SysStats{
		Ip:            ip,
		Cpu:           cpus,
		VirtualMemory: virtualMemory,
		SwapMemory:    swapMemory,
		Disk:          diskInfo,
		Timestamp: &timestamp.Timestamp{
			Seconds: now.Unix(),
			Nanos:   int32(now.Nanosecond()),
		},
	}
}

func GetHostInfo() *proto.HostInfo {
	cpuInfos, _ := cpu.Info()
	cpus := make([]*proto.CpuInfo, 0)
	for i := 0; i < len(cpuInfos); i++ {
		currentCpu := cpuInfos[i]
		cpuID, _ := strconv.Atoi(currentCpu.CoreID)
		c := proto.CpuInfo{
			Num:        int32(cpuID),
			VendorId:   currentCpu.VendorID,
			Family:     currentCpu.Family,
			PhysicalId: currentCpu.PhysicalID,
			Cores:      currentCpu.Cores,
			ModelName:  currentCpu.ModelName,
			Mhz:        float32(currentCpu.Mhz),
		}
		cpus = append(cpus, &c)
	}
	ips, _ := getLocalIPs()
	ip := strings.Join(ips, ",")
	hostInfo, _ := host.Info()
	now := time.Now()
	return &proto.HostInfo{
		Ip:              ip,
		HostName:        hostInfo.Hostname,
		UpTime:          strconv.FormatUint(hostInfo.Uptime, 10),
		Platform:        hostInfo.Platform,
		PlatformVersion: hostInfo.PlatformVersion,
		KernelVersion:   hostInfo.KernelVersion,
		KernelArch:      hostInfo.KernelArch,
		CpuInfos:        cpus,
		Timestamp: &timestamp.Timestamp{
			Seconds: now.Unix(),
			Nanos:   int32(now.Nanosecond()),
		},
	}
}

// getLocalIPs 返回本地机器的所有非环回 IP 地址。
func getLocalIPs() ([]string, error) {
	var ips []string

	interfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	for _, intf := range interfaces {
		addrs, err := intf.Addrs()
		if err != nil {
			return nil, err
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}

			if ip == nil || ip.IsLoopback() {
				continue
			}

			if ip.To4() != nil {
				ips = append(ips, ip.String())
			}
		}
	}

	return ips, nil
}
