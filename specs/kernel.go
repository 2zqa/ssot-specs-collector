package specs

import (
	"runtime"

	"github.com/2zqa/ssot-specs-collector/aggregator"
	"github.com/shirou/gopsutil/v3/host"
)

type Kernel struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

func FetchKernel() (Kernel, error) {
	a := aggregator.New("kernel")
	kernel := Kernel{
		Name: runtime.GOOS,
	}

	version, err := host.KernelVersion()
	if a.CheckAndAdd(err) {
		kernel.Version = version
	}

	return kernel, a.ErrorOrNil()
}
