package specs

import (
	"net"
	"os"
	"strings"

	"github.com/2zqa/ssot-specs-collector/aggregator"
	"github.com/charmbracelet/log"
	"github.com/safchain/ethtool"
)

type Network struct {
	Interfaces []NetworkInterface `json:"interfaces"`
	Hostname   string             `json:"hostname"`
}

type NetworkInterface struct {
	MACAddress    string   `json:"mac_address"`
	Driver        Driver   `json:"driver"`
	IPv4Addresses []string `json:"ipv4_addresses"`
	IPv6Addresses []string `json:"ipv6_addresses"`
}

type Driver struct {
	Name            string `json:"name"`
	Version         string `json:"version"`
	FirmwareVersion string `json:"firmware_version"`
}

func FetchNetwork() (Network, error) {
	a := aggregator.New("network")
	var network Network

	h, err := os.Hostname()
	if a.CheckAndAdd(err) {
		network.Hostname = h
	}

	ni, err := FetchNetworkInterfaces()
	a.CheckAndAdd(err)
	// Add network interface even if there was an error
	// because there might still be useful information in the response
	network.Interfaces = ni

	return network, a.ErrorOrNil()
}

func FetchNetworkInterfaces() ([]NetworkInterface, error) {
	a := aggregator.New("network interfaces")

	ifaces, err := net.Interfaces()
	if !a.CheckAndAdd(err) {
		return nil, a
	}

	ethHandle, err := ethtool.NewEthtool()
	if a.CheckAndAdd(err) {
		defer ethHandle.Close()
	}

	ifaceCount := len(ifaces)
	log.Debug("Total interfaces detected", "count", ifaceCount)
	networkInterfaces := make([]NetworkInterface, 0, ifaceCount)
	for _, iface := range ifaces {
		if shouldNotUseInterface(iface) {
			log.Debug("Skipping interface", "interface", iface)
			continue
		}
		iface, err := FetchNetworkInterface(iface, ethHandle)
		a.CheckAndAdd(err)
		// Add interface even if there was an error,
		// there may still be usable data in it
		networkInterfaces = append(networkInterfaces, iface)
	}
	log.Info("Retrieved network interfaces", "count", len(networkInterfaces))
	return networkInterfaces, a.ErrorOrNil()
}

func FetchNetworkInterface(netIface net.Interface, ethHandle *ethtool.Ethtool) (NetworkInterface, error) {
	a := aggregator.New("network interface")
	specsIface := NetworkInterface{
		IPv4Addresses: make([]string, 0),
		IPv6Addresses: make([]string, 0),
	}
	specsIface.MACAddress = netIface.HardwareAddr.String()

	if ethHandle != nil {
		specsIface.Driver = NewDriver(netIface.Name, ethHandle)
	}

	addresses, err := netIface.Addrs()
	if !a.CheckAndAdd(err) {
		return specsIface, a
	}

	for _, address := range addresses {
		ipString := address.String()
		if shouldNotUseIP(ipString) {
			log.Debug("Skipping IP", "ip", ipString)
			continue
		}

		if isIPv4(ipString) {
			specsIface.IPv4Addresses = append(specsIface.IPv4Addresses, ipString)
		}
		if isIPv6(ipString) {
			specsIface.IPv6Addresses = append(specsIface.IPv6Addresses, ipString)
		}
	}

	return specsIface, a.ErrorOrNil()
}

func NewDriver(interfaceName string, ethHandle *ethtool.Ethtool) Driver {
	if ethHandle == nil {
		log.Error("Ethtool handle is nil")
		return Driver{}
	}

	driverInfo, err := ethHandle.DriverInfo(interfaceName)
	if err != nil {
		log.Error(err)
		return Driver{}
	}

	driver := Driver{
		Name:    driverInfo.Driver,
		Version: driverInfo.Version,
	}

	if version := driverInfo.FwVersion; version != "N/A" {
		driver.FirmwareVersion = version
	}

	return driver
}

func shouldNotUseInterface(iface net.Interface) bool {
	// Don't use interface if it's down
	if iface.Flags&net.FlagUp == 0 {
		return true
	}
	// Don't use if it's a loopback interface
	if iface.Flags&net.FlagLoopback != 0 {
		return true
	}
	return false
}

func shouldNotUseIP(ipString string) bool {
	ip, _, err := net.ParseCIDR(ipString)
	if err != nil {
		log.Error(err)
		return true
	}
	return ip.IsLoopback()
}

func isIPv4(address string) bool {
	return strings.Count(address, ":") < 2
}

func isIPv6(address string) bool {
	return strings.Count(address, ":") >= 2
}
