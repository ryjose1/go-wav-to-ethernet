package main

import (
	"fmt"
	"log"
	"net"

	"github.com/ryjose1/go-wav-to-ethernet/config"
	"github.com/ryjose1/go-wav-to-ethernet/ethernet"
	"github.com/ryjose1/go-wav-to-ethernet/wave"

	"github.com/mdlayher/raw"
)

func main() {
	// Read file into batches
	batches := wave.GenerateBatches()

	// Process chunks to build ethernet packet payloads
	payloads := [][]byte{}
	for _, batch := range batches {
		payload := wave.SamplesToPayload(batch)
		payloads = append(payloads, payload)
	}

	// Connect to the ethernet interface
	iface, err := net.InterfaceByName(config.Iface)
	if err != nil {
		log.Printf("unable to find interface %s: %v", config.Iface, err)
	}
	conn, err := raw.ListenPacket(iface, config.EtherType, nil)
	if err != nil {
		log.Fatalf("failed to create connection")
	}

	// Send the marshaled frame to the network.
	for _, payload := range payloads {
		frame := ethernet.NewFrame(iface.HardwareAddr, payload)
		//time.Sleep(time.Second / 100000000)
		fmt.Printf("msg: %x\n\n", payload)
		ethernet.SendMessage(conn, frame)
	}
}
