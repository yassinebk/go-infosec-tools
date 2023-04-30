package main

import (
	"fmt"
	"log"

	"github.com/google/gopacket"
	"github.com/google/gopacket/pcap"
)

var (
	iface    = "wlan0"
	snaplen  = int32(1600)
	promisc  = false
	timeout  = pcap.BlockForever
	filter   = "tcp and port 80"
	devFound = false
)

func getAllDevices() {
	devices, err := pcap.FindAllDevs()
	if err != nil {
		log.Panicln(err)
	}

	for _, device := range devices {
		fmt.Println("Name: ", device.Name)
		fmt.Println("Description: ", device.Description)
		fmt.Println("Devices addresses: ", device.Addresses)
		fmt.Println("Devices flags: ", device.Flags)
		for _, address := range device.Addresses {
			fmt.Println("IP address: ", address.IP)
			fmt.Println("Subnet mask: ", address.Netmask)
		}
		fmt.Println("---------------------------")
	}
}

func findDevice(iface string) *pcap.Interface {
	devices, err := pcap.FindAllDevs()
	if err != nil {
		log.Panicln(err)
	}
	for _, device := range devices {
		if device.Name == iface {
			devFound = true
			return &device

		}
	}

	log.Panicln("Device not found", iface)
	return nil
}

func main() {

	device := findDevice(iface)
	fmt.Println("Device found: ", device)

	handle, err := pcap.OpenLive(device.Name, snaplen, promisc, timeout)
	if err != nil {
		log.Panicln(err)
	}

	defer handle.Close()

	if err := handle.SetBPFFilter(filter); err != nil {
		log.Panicln(err)
	}

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	for packet := range packetSource.Packets() {
		fmt.Println(packet)
	}

}
