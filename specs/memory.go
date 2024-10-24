package specs

import (
	"github.com/2zqa/ssot-specs-collector/aggregator"
	"github.com/shirou/gopsutil/mem"
)

type Memory struct {
	TotalSize   uint64       `json:"memory"`
	TotalSwap   uint64       `json:"swap"`
	SwapDevices []SwapDevice `json:"swap_devices"`
}

type SwapDevice struct {
	Name string `json:"name"`
	Size uint64 `json:"size"`
}

func FetchMemory() (Memory, error) {
	a := aggregator.New("memory")
	var memory Memory

	memorySize, err := getTotalMemorySize()
	if a.CheckAndAdd(err) {
		memory.TotalSize = memorySize
	}

	swapSize, err := getTotalSwapSize()
	if a.CheckAndAdd(err) {
		memory.TotalSwap = swapSize
	}

	sd, err := FetchSwapDevices()
	a.CheckAndAdd(err)
	// Add swap devices even if there was an error
	// because there might still be useful information in the response
	memory.SwapDevices = sd

	return memory, a.ErrorOrNil()
}

func FetchSwapDevices() ([]SwapDevice, error) {
	a := aggregator.New("swap devices")

	swapDevices, swapDeviceErr := mem.SwapDevices()
	if !a.CheckAndAdd(swapDeviceErr) {
		return nil, a
	}

	newSwapDevices := make([]SwapDevice, len(swapDevices))
	for i := range swapDevices {
		newSwapDevices[i] = NewSwapDevice(swapDevices[i])
	}

	return newSwapDevices, a.ErrorOrNil()
}

// NewSwapDevice converts a gopsutil SwapDevice into a SwapDevice
func NewSwapDevice(swapDevice *mem.SwapDevice) SwapDevice {
	totalMemory := swapDevice.FreeBytes + swapDevice.UsedBytes
	return SwapDevice{
		Name: swapDevice.Name,
		Size: byteToMegabyte(totalMemory),
	}
}

// getTotalSwapSize returns the total swap size in megabytes
func getTotalSwapSize() (uint64, error) {
	swap, err := mem.SwapMemory()
	if err != nil {
		return 0, err
	}

	return byteToMegabyte(swap.Total), nil

}

// getTotalMemorySize returns the total memory size in megabytes
func getTotalMemorySize() (uint64, error) {
	memory, err := mem.VirtualMemory()
	if err != nil {
		return 0, err
	}

	return byteToMegabyte(memory.Total), nil
}
