package main

import (
	"fmt"
	"log"
	"net"

	"github.com/DylanMeeus/GoAudio/wave"
	"github.com/mdlayher/ethernet"
)

const (
	// Equivalent to 1518 bytes of audio data
	// 1518 bytes/packet * 8 bits/byte * 1 channel/16 bits * 1 sample/2 channels *1 second/44100 samples
	MaxSecondsPerPacket = .0086054422
)

func main() {
	// Read file into batches
	filepath := "./birds.wav"
	wav, err := wave.ReadWaveFile(filepath)
	if err != nil {
		fmt.Printf("Could not read wave file: %v", err)
	}

	batches := wave.BatchSamples(wav, MaxSecondsPerPacket)

	// Process chunks to build ethernet packets
	for i, batch := range batches {
		if i == 0 {
			waveFmt := wave.NewWaveFmt(1, 2, 44100, 16, []byte{})
			wave.WriteFrames(batch, waveFmt, fmt.Sprintf("./test/birds._pt%d.wav", i))
			// 2 bytes of data per frame in 16bits/sample rate, should have 759 packets
			fmt.Printf("%d\n", len(batch))
		}
	}

	// Send ethernet packets
	f := &ethernet.Frame{
		// Broadcast frame to all machines on same network segment.
		Destination: ethernet.Broadcast,
		// Identify our machine as the sender.
		Source: net.HardwareAddr{0xde, 0xad, 0xbe, 0xef, 0xde, 0xad},
		// Identify frame with an unused EtherType.
		EtherType: 0xcccc,
		// Send a simple message.
		Payload: []byte("hello world"),
	}
	// Marshal the Go representation of a frame to
	// the Ethernet frame format.
	_, err = f.MarshalBinary()
	if err != nil {
		log.Fatalf("failed to marshal frame: %v", err)
	}
	// Send the marshaled frame to the network.
	//sendEthernetFrame(b)

}
