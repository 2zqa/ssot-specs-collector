package sender

import (
	"context"
	"fmt"

	openapiclient "github.com/2zqa/ssot-specs-api-client"
	"github.com/2zqa/ssot-specs-collector/metadata"
	"github.com/2zqa/ssot-specs-collector/specs"
	"github.com/charmbracelet/log"
)

type APIClientSender struct {
}

func (p APIClientSender) Send(specs specs.PCSpecs, metadata metadata.Metadata, ctx context.Context) error {
	// Prepare payload
	uuid := metadata.UUID.String()
	openAPISpecs := convertPCSpecsToOpenAPISpecs(specs)
	devicePutInput := openapiclient.DevicePutInput{
		Specs: &openAPISpecs,
	}

	// Setup API client
	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	setAPIKey(&ctx, metadata.APIKey)

	// Send data to API
	resp, r, err := apiClient.DevicesApi.DevicesUuidPut(ctx, uuid).DevicePutInput(devicePutInput).Execute()
	if err != nil {
		return fmt.Errorf("error sending specs to API: %v", err)
	}

	log.Info("Successfully sent specs to API 🎉")
	log.Debug("Received device", "device", resp)
	log.Debug("Received full HTTP response", "response", r)

	return nil
}

func setAPIKey(ctx *context.Context, s string) {
	apiKeyMap := make(map[string]openapiclient.APIKey)
	apiKeyMap["ApiKeyAuth"] = openapiclient.APIKey{Key: s}
	*ctx = context.WithValue(*ctx, openapiclient.ContextAPIKeys, apiKeyMap)
}

func convertPCSpecsToOpenAPISpecs(specs specs.PCSpecs) openapiclient.Specs {
	coreCount := int64(specs.CPU.CoreCount)
	cpuCount := int64(specs.CPU.CPUCount)
	maxFrequencyMegaHertz := int64(specs.CPU.MaxFrequencyMegaHertz)
	memory := int64(specs.Memory.TotalSize)
	swap := int64(specs.Memory.TotalSwap)

	openapiSpecs := openapiclient.Specs{
		Motherboard: &openapiclient.Motherboard{
			Vendor:       &specs.Motherboard.Vendor,
			Name:         &specs.Motherboard.Name,
			SerialNumber: &specs.Motherboard.SerialNumber,
		},
		Cpu: &openapiclient.CPU{
			Name:                  &specs.CPU.Name,
			Architecture:          &specs.CPU.Architecture,
			CoreCount:             &coreCount,
			CpuCount:              &cpuCount,
			MaxFrequencyMegahertz: &maxFrequencyMegaHertz,
			Mitigations:           specs.CPU.Mitigations,
		},
		Disks: convertDisks(specs.Disks),
		Network: &openapiclient.Network{
			Interfaces: convertNetworkInterfaces(specs.Network.Interfaces),
			Hostname:   &specs.Network.Hostname,
		},
		Bios: &openapiclient.BIOS{
			Vendor:  &specs.BIOS.Vendor,
			Version: &specs.BIOS.Version,
			Date:    &specs.BIOS.Date,
		},
		Memory: &openapiclient.Memory{
			Memory:      &memory,
			Swap:        &swap,
			SwapDevices: convertSwapDevices(specs.Memory.SwapDevices),
		},
		BootTime: &specs.BootTime,
		Kernel: &openapiclient.Kernel{
			Name:    &specs.Kernel.Name,
			Version: &specs.Kernel.Version,
		},
		Release: &openapiclient.Release{
			Name:     &specs.Release.Name,
			Version:  &specs.Release.Version,
			Codename: &specs.Release.Codename,
		},
		Dimms: convertDIMMs(specs.DIMMs),
		Oem: &openapiclient.OEM{
			Manufacturer: &specs.OEM.Manufacturer,
			ProductName:  &specs.OEM.ProductName,
			SerialNumber: &specs.OEM.SerialNumber,
		},
		Virtualization: &openapiclient.Virtualization{
			Type: &specs.Virtualization.Type,
		},
	}

	return openapiSpecs
}

func convertDIMMs(dimms []specs.DIMM) []openapiclient.DIMMsInner {
	openapiDIMMs := make([]openapiclient.DIMMsInner, len(dimms))
	for i, dimm := range dimms {
		localDIMM := dimm
		speed := int64(localDIMM.Speed)
		sizeGigabytes := int64(localDIMM.SizeGigabytes)
		openapiDIMMs[i] = openapiclient.DIMMsInner{
			SizeGigabytes: &sizeGigabytes,
			SpeedMtS:      &speed,
			Manufacturer:  &localDIMM.Manufacturer,
			SerialNumber:  &localDIMM.SerialNumber,
			Type:          &localDIMM.Type,
			PartNumber:    &localDIMM.PartNumber,
			FormFactor:    &localDIMM.FormFactor,
			Locator:       &localDIMM.Locator,
			BankLocator:   &localDIMM.BankLocator,
		}
	}
	return openapiDIMMs
}

func convertDisks(disks []specs.Disk) []openapiclient.DisksInner {
	openapiDisks := make([]openapiclient.DisksInner, len(disks))
	for i, disk := range disks {
		// The variable gets overwritten in the loop, so we need to make a local copy to keep pointers.
		// https://stackoverflow.com/questions/45967305/copying-the-address-of-a-loop-variable-in-go
		localDisk := disk
		sizeMegaBytes := int64(localDisk.SizeMegabytes)
		openapiDisks[i] = openapiclient.DisksInner{
			Name:          &localDisk.Name,
			SizeMegabytes: &sizeMegaBytes,
			Partitions:    convertPartitions(localDisk.Partitions),
		}
	}
	return openapiDisks
}

func convertPartitions(partitions []specs.Partition) []openapiclient.Partition {
	openapiPartitions := make([]openapiclient.Partition, len(partitions))
	for i, partition := range partitions {
		localPartition := partition
		capacityMegabytes := int64(localPartition.CapacityMegabytes)

		openapiPartitions[i] = openapiclient.Partition{
			Filesystem:        &localPartition.Filesystem,
			CapacityMegabytes: &capacityMegabytes,
			Source:            &localPartition.Source,
			Target:            &localPartition.Target,
		}
	}
	return openapiPartitions
}

func convertNetworkInterfaces(interfaces []specs.NetworkInterface) []openapiclient.NetworkInterface {
	openapiInterfaces := make([]openapiclient.NetworkInterface, len(interfaces))
	for i, intf := range interfaces {
		localIntf := intf
		openapiInterfaces[i] = openapiclient.NetworkInterface{
			MacAddress: &localIntf.MACAddress,
			Driver: &openapiclient.Driver{
				Name:            &localIntf.Driver.Name,
				Version:         &localIntf.Driver.Version,
				FirmwareVersion: &localIntf.Driver.FirmwareVersion,
			},
			Ipv4Addresses: localIntf.IPv4Addresses,
			Ipv6Addresses: localIntf.IPv6Addresses,
		}
	}
	return openapiInterfaces
}

func convertSwapDevices(devices []specs.SwapDevice) []openapiclient.SwapDevice {
	openapiDevices := make([]openapiclient.SwapDevice, len(devices))
	for i, device := range devices {
		localDevice := device
		size := int64(localDevice.Size)
		openapiDevices[i] = openapiclient.SwapDevice{
			Name: &localDevice.Name,
			Size: &size,
		}
	}
	return openapiDevices
}
