package specs

import (
	"github.com/2zqa/ssot-specs-collector/aggregator"
	"github.com/jaypipes/ghw"
)

type BIOS struct {
	Vendor  string `json:"vendor"`
	Version string `json:"version"`
	Date    string `json:"date"`
}

func FetchBIOS() (BIOS, error) {
	a := aggregator.New("bios")
	bios, err := ghw.BIOS(ghwOptions)
	if !a.CheckAndAdd(err) {
		return BIOS{}, a.ErrorOrNil()
	}

	return BIOS{
		Vendor:  bios.Vendor,
		Version: bios.Version,
		Date:    bios.Date,
	}, a.ErrorOrNil()
}
