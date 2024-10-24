package specs

import (
	"github.com/2zqa/ssot-specs-collector/aggregator"
	"github.com/jaypipes/ghw"
)

type Motherboard struct {
	Vendor       string `json:"vendor"`
	Name         string `json:"name"`
	SerialNumber string `json:"serial_number"`
}

func FetchMotherboard() (Motherboard, error) {
	a := aggregator.New("motherboard")
	baseboard, err := ghw.Baseboard(ghwOptions)
	if !a.CheckAndAdd(err) {
		return Motherboard{}, a.ErrorOrNil()
	}

	motherboard := Motherboard{
		Vendor: baseboard.Vendor,
		Name:   baseboard.Product,
	}

	if a.CheckPredicateAndAdd(baseboard.SerialNumber != "unknown", "serial number is unknown") {
		motherboard.SerialNumber = baseboard.SerialNumber
	}

	return motherboard, a.ErrorOrNil()
}
