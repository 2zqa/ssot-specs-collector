package specs

import (
	"github.com/2zqa/ssot-specs-collector/aggregator"
	"github.com/prometheus/procfs/sysfs"
)

type OEM struct {
	Manufacturer string `json:"manufacturer"`
	ProductName  string `json:"product_name"`
	SerialNumber string `json:"serial_number"`
}

func FetchOEM() (OEM, error) {
	a := aggregator.New("oem")
	sys, err := sysfs.NewDefaultFS()
	if !a.CheckAndAdd(err) {
		return OEM{}, a.ErrorOrNil()
	}

	dmi, err := sys.DMIClass()
	if !a.CheckAndAdd(err) {
		return OEM{}, a.ErrorOrNil()
	}

	oem := OEM{}
	if a.CheckPredicateAndAdd(dmi.SystemVendor != nil, "system vendor is nil") {
		oem.Manufacturer = *dmi.SystemVendor
	}
	if a.CheckPredicateAndAdd(dmi.ProductName != nil, "product name is nil") {
		oem.ProductName = *dmi.ProductName
	}
	if a.CheckPredicateAndAdd(dmi.ProductSerial != nil, "product serial is nil") {
		oem.SerialNumber = *dmi.ProductSerial
	}

	return oem, a.ErrorOrNil()
}
