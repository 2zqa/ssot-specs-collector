package specs

import (
	"github.com/2zqa/ssot-specs-collector/aggregator"
	"github.com/charmbracelet/log"
	"github.com/jaypipes/ghw"
	"github.com/jaypipes/ghw/pkg/block"
)

// emptyDiskArray can be used when no disks are detected. This is preferred
// over using nil because json.Marshal returns "null" for nil slices, but we
// want to return an empty array instead ("[]").
var emptyDiskArray = make([]Disk, 0)

type Disk struct {
	Name          string      `json:"name"`
	SizeMegabytes uint64      `json:"size_megabytes"`
	Partitions    []Partition `json:"partitions"`
}

type Partition struct {
	Filesystem        string `json:"filesystem"`
	CapacityMegabytes uint64 `json:"capacity_megabytes"`
	Source            string `json:"source"`
	Target            string `json:"target"`
}

func FetchDisks() ([]Disk, error) {
	a := aggregator.New("disks")
	block, err := ghw.Block(ghwOptions)
	if !a.CheckAndAdd(err) {
		return emptyDiskArray, a.ErrorOrNil()
	}

	diskCount := len(block.Disks)
	log.Debug("Total disks detected", "count", diskCount)
	if diskCount == 0 {
		return emptyDiskArray, a.ErrorOrNil()
	}

	// Convert ghw Disks into Disks
	disks := make([]Disk, 0, diskCount)
	for i := 0; i < diskCount; i++ {
		if disk := block.Disks[i]; shouldNotFilterDisk(disk) {
			disks = append(disks, NewDisk(disk))
		} else {
			log.Debug("Skipping disk", "disk", disk)
		}
	}
	log.Info("Retrieved disks", "count", len(disks))
	return disks, a.ErrorOrNil()
}

func shouldNotFilterDisk(disk *block.Disk) bool {
	switch disk.StorageController {
	case block.STORAGE_CONTROLLER_UNKNOWN, block.STORAGE_CONTROLLER_LOOP:
		return false
	default:
		return true
	}
}

func NewDisk(disk *block.Disk) Disk {
	model := disk.Model
	return Disk{
		Name:          model,
		SizeMegabytes: byteToMegabyte(disk.SizeBytes),
		Partitions:    NewPartitions(disk),
	}
}

func NewPartitions(disk *block.Disk) []Partition {
	specPartitions := make([]Partition, len(disk.Partitions))
	for i, p := range disk.Partitions {
		specPartitions[i] = NewPartition(p)
	}
	return specPartitions
}

func NewPartition(partition *block.Partition) Partition {
	return Partition{
		Filesystem:        partition.Type,
		CapacityMegabytes: byteToMegabyte(partition.SizeBytes),
		Source:            partition.Name,
		Target:            partition.MountPoint,
	}
}
