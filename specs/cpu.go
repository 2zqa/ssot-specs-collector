package specs

import (
	"fmt"

	"github.com/2zqa/ssot-specs-collector/aggregator"
	"github.com/prometheus/procfs"
	"github.com/prometheus/procfs/sysfs"
	"github.com/shirou/gopsutil/v3/host"
)

type CPU struct {
	Name                  string   `json:"name"`
	Architecture          string   `json:"architecture"`
	CoreCount             uint64   `json:"core_count"`
	CPUCount              uint64   `json:"cpu_count"`
	MaxFrequencyMegaHertz uint64   `json:"max_frequency_megahertz"`
	Mitigations           []string `json:"mitigations"`
}

func FetchCPU() (CPU, error) {
	a := aggregator.New("cpu")
	var cpu CPU

	cpuInfo, err := getCPUInfo()
	if a.CheckAndAdd(err) {
		cpu.Name = cpuInfo.ModelName
		cpu.CoreCount = uint64(cpuInfo.CPUCores)
		cpu.CPUCount = uint64(cpuInfo.Siblings)
		cpu.Mitigations = cpuInfo.Bugs
	}

	kernelArch, err := getKernelArch()
	if a.CheckAndAdd(err) {
		cpu.Architecture = kernelArch
	}

	maxFreq, err := getMaxFrequencyMegaHertz()
	if a.CheckAndAdd(err) {
		cpu.MaxFrequencyMegaHertz = maxFreq
	}

	return cpu, a.ErrorOrNil()
}

func getCPUInfo() (procfs.CPUInfo, error) {
	proc, err := procfs.NewDefaultFS()
	if err != nil {
		return procfs.CPUInfo{}, err
	}

	cpuInfo, err := proc.CPUInfo()
	if err != nil {
		return procfs.CPUInfo{}, err
	}

	return cpuInfo[0], nil
}

func getMaxFrequencyMegaHertz() (uint64, error) {
	sys, err := sysfs.NewDefaultFS()
	if err != nil {
		return 0, err
	}

	cpuFreq, err := sys.SystemCpufreq()
	if err != nil {
		return 0, err
	}

	if len(cpuFreq) == 0 {
		return 0, fmt.Errorf("no CPU Frequency information found")
	}

	maxCPUFrequency := cpuFreq[0].CpuinfoMaximumFrequency
	if maxCPUFrequency == nil {
		return 0, fmt.Errorf("could not retrieve max CPU Frequency")
	}
	return kiloHertzToMegaHertz(*maxCPUFrequency), nil
}

func getKernelArch() (string, error) {
	host, err := host.Info()
	if err != nil {
		return "", err
	}
	return host.KernelArch, nil
}
