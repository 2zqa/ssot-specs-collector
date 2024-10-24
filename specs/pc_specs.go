package specs

import (
	"time"

	"github.com/2zqa/ssot-specs-collector/aggregator"
	"github.com/charmbracelet/log"
	"github.com/jaypipes/ghw"
)

var (
	ghwOptions = ghw.WithAlerter(log.StandardLog())
)

type PCSpecs struct {
	Motherboard    Motherboard    `json:"motherboard"`
	CPU            CPU            `json:"cpu"`
	Disks          []Disk         `json:"disks"`
	Network        Network        `json:"network"`
	BIOS           BIOS           `json:"bios"`
	Memory         Memory         `json:"memory"`
	DIMMs          []DIMM         `json:"dimms"`
	BootTime       time.Time      `json:"boot_time"`
	Kernel         Kernel         `json:"kernel"`
	Release        Release        `json:"release"`
	OEM            OEM            `json:"oem"`
	Virtualization Virtualization `json:"virtualization"`
}

// Fetch retrieves all specs from the current system. This function may return
// errors, however the returned PCSpecs struct will still contain all specs
// that were successfully retrieved.
func Fetch() (PCSpecs, error) {
	a := aggregator.New("PC specs")
	if hasNoRoot() {
		log.Warn("Not running as root, some specs may be missing")
	}

	m, err := FetchMotherboard()
	a.CheckAndAdd(err)

	c, err := FetchCPU()
	a.CheckAndAdd(err)

	d, err := FetchDisks()
	a.CheckAndAdd(err)

	n, err := FetchNetwork()
	a.CheckAndAdd(err)

	b, err := FetchBIOS()
	a.CheckAndAdd(err)

	mem, err := FetchMemory()
	a.CheckAndAdd(err)

	dimms, err := FetchDIMMs()
	a.CheckAndAdd(err)

	k, err := FetchKernel()
	a.CheckAndAdd(err)

	r, err := FetchRelease()
	a.CheckAndAdd(err)

	bootTime, err := fetchBootTime()
	a.CheckAndAdd(err)

	o, err := FetchOEM()
	a.CheckAndAdd(err)

	v, err := FetchVirtualization()
	a.CheckAndAdd(err)

	pcSpecs := PCSpecs{
		Motherboard:    m,
		CPU:            c,
		Disks:          d,
		Network:        n,
		BIOS:           b,
		Memory:         mem,
		DIMMs:          dimms,
		BootTime:       bootTime,
		Kernel:         k,
		Release:        r,
		OEM:            o,
		Virtualization: v,
	}

	return pcSpecs, a.ErrorOrNil()
}
