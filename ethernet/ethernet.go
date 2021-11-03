package ethernet

import (
	"log"
	"net"

	"github.com/mdlayher/ethernet"
	"github.com/mdlayher/raw"
	"github.com/ryjose1/go-wav-to-ethernet/config"
)

func NewFrame(source net.HardwareAddr, msg []byte) *ethernet.Frame {
	// Message is broadcast to all machines in same network segment.
	return &ethernet.Frame{
		Destination: ethernet.Broadcast,
		Source:      source,
		EtherType:   config.EtherType,
		Payload:     msg,
	}
}

// sendMessages continuously sends a message over a connection at regular intervals,
// sourced from specified hardware address.
func SendMessage(c net.PacketConn, f *ethernet.Frame) {
	b, err := f.MarshalBinary()
	if err != nil {
		log.Fatalf("failed to marshal ethernet frame: %v", err)
	}

	// Required by Linux, even though the Ethernet frame has a destination.
	// Unused by BSD.
	addr := &raw.Addr{
		HardwareAddr: ethernet.Broadcast,
	}

	if _, err := c.WriteTo(b, addr); err != nil {
		log.Fatalf("failed to send message: %v", err)
	}
	return
}
