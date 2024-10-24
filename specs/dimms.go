package specs

import (
	"strconv"
	"strings"

	"github.com/2zqa/ssot-specs-collector/aggregator"
	"github.com/7Linternational/dmidecode"
	"github.com/charmbracelet/log"
	"github.com/inhies/go-bytesize"
)

// emptyDIMMArray can be used when no disks are detected. This is preferred
// over using nil because json.Marshal returns "null" for nil slices, but we
// want to return an empty array instead ("[]").
var emptyDIMMArray = make([]DIMM, 0)

const TypeMemoryDevice int = 17

type DIMM struct {
	SizeGigabytes uint64 `json:"size_gigabytes"`
	Speed         uint64 `json:"speed_mt_s"`
	Manufacturer  string `json:"manufacturer"`
	SerialNumber  string `json:"serial_number"`
	Type          string `json:"type"`
	PartNumber    string `json:"part_number"`
	FormFactor    string `json:"form_factor"`
	Locator       string `json:"locator"`
	BankLocator   string `json:"bank_locator"`
}

func FetchDIMMs() ([]DIMM, error) {
	a := aggregator.New("DIMM")
	dmi := dmidecode.New()
	if err := dmi.Run(); !a.CheckAndAdd(err) {
		return emptyDIMMArray, a
	}

	rawDIMMs, err := dmi.SearchByType(TypeMemoryDevice)
	if !a.CheckAndAdd(err) {
		return emptyDIMMArray, a
	}

	dimmCount := len(rawDIMMs)
	log.Debug("Total dimms detected", "count", dimmCount)
	if dimmCount == 0 {
		return emptyDIMMArray, a.ErrorOrNil()
	}

	// Convert dmi DIMMs into DIMMs
	dimms := make([]DIMM, 0, dimmCount)
	for i := 0; i < dimmCount; i++ {
		if dimm := rawDIMMs[i]; shouldNotFilterDIMM(dimm) {
			dimm, err := NewDIMM(dimm)
			a.CheckAndAdd(err)
			// Append the DIMM even if there was an error retrieving it
			// There can still be useful info in it
			dimms = append(dimms, dimm)
		} else {
			log.Debug("Skipping dimm", "dimm", dimm)
		}
	}
	log.Info("Retrieved dimms", "count", len(dimms))
	return dimms, a.ErrorOrNil()
}

func NewDIMM(dmiDIMM dmidecode.Record) (DIMM, error) {
	a := aggregator.New("DIMM")
	dimm := DIMM{
		Manufacturer: strings.TrimSpace(dmiDIMM["Manufacturer"]),
		SerialNumber: strings.TrimSpace(dmiDIMM["Serial Number"]),
		Type:         strings.TrimSpace(dmiDIMM["Type"]),
		PartNumber:   strings.TrimSpace(dmiDIMM["Part Number"]),
		FormFactor:   strings.TrimSpace(dmiDIMM["Form Factor"]),
		Locator:      strings.TrimSpace(dmiDIMM["Locator"]),
		BankLocator:  strings.TrimSpace(dmiDIMM["Bank Locator"]),
	}

	size, err := parseDIMMSize(dmiDIMM["Size"])
	if a.CheckAndAdd(err) {
		dimm.SizeGigabytes = size
	}

	speed, err := parseDIMMSpeed(dmiDIMM["Speed"])
	if a.CheckAndAdd(err) {
		dimm.Speed = speed
	}

	return dimm, a.ErrorOrNil()
}

// parseDIMMSpeed parses the speed of a DIMM into a uint64 representing the speed in MT/s.
func parseDIMMSpeed(speedString string) (uint64, error) {
	speedWithoutUnit := speedString[:len(speedString)-5]
	return strconv.ParseUint(speedWithoutUnit, 10, 64)
}

// Parses a DIMM size string into a uint64 representing the size in gigabytes.
func parseDIMMSize(sizeString string) (uint64, error) {
	// First, parse the size string into a bytesize.Size
	parsedSize, err := bytesize.Parse(sizeString)
	if err != nil {
		return 0, err
	}

	// Convert said bytesize.Size into a string representing the size in gigabytes
	sizeInGB := parsedSize.Format("%.0f", "GB", false)

	// Remove the "GB" suffix from the string and parse it into a uint64
	sizeUint, err := strconv.ParseUint(sizeInGB[:len(sizeInGB)-2], 10, 64)
	if err != nil {
		return 0, err
	}
	return sizeUint, nil
}

func shouldNotFilterDIMM(dimm dmidecode.Record) bool {
	return dimm["Size"] != "No Module Installed"
}
