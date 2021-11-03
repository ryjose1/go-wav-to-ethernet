package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"math"
	"net"

	"github.com/DylanMeeus/GoAudio/wave"
	"github.com/mdlayher/ethernet"
	"github.com/mdlayher/raw"
)

const (
	// Equivalent to 1518 bytes of audio data
	// 1400 bytes/packet * 8 bits/byte * 1 channel/16 bits * 1 sample/2 channels *1 second/44100 samples
	MaxSecondsPerPacket = .0079365079
	WaveFile            = "./birds.wav"
	TestOutputFormat    = "./test/birds._pt%d.wav"
	EtherType           = 0xcccc
	Iface               = "enp0s8"
)

type intsToBytesFunc func(i int) []byte

var (
	// intsToBytesFm to map X-bit int to byte functions
	intsToBytesFm = map[int]intsToBytesFunc{
		16: int16ToBytes,
		32: int32ToBytes,
	}
	// max value depending on the bit size
	maxValues = map[int]int{
		8:  math.MaxInt8,
		16: math.MaxInt16,
		32: math.MaxInt32,
		64: math.MaxInt64,
	}
)

func int16ToBytes(i int) []byte {
	b := make([]byte, 2)
	in := uint16(i)
	binary.LittleEndian.PutUint16(b, in)
	return b
}

func int32ToBytes(i int) []byte {
	b := make([]byte, 4)
	in := uint32(i)
	binary.LittleEndian.PutUint32(b, in)
	return b
}

// sendMessages continuously sends a message over a connection at regular intervals,
// sourced from specified hardware address.
func sendMessage(c net.PacketConn, f ethernet.Frame) {
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

func NewFrame(source net.HardwareAddr, msg []byte) *ethernet.Frame {
	// Message is broadcast to all machines in same network segment.
	return &ethernet.Frame{
		Destination: ethernet.Broadcast,
		Source:      source,
		EtherType:   EtherType,
		Payload:     msg,
	}
}

// Turn the samples into raw data...
func samplesToRawData(samples []wave.Frame, props wave.WaveFmt) []byte {
	raw := []byte{}
	for _, s := range samples {
		// the samples are scaled - rescale them?
		rescaled := rescaleFrame(s, props.BitsPerSample)
		bits := intsToBytesFm[props.BitsPerSample](rescaled)
		raw = append(raw, bits...)
	}
	return raw
}

// rescale frames back to the original values..
func rescaleFrame(s wave.Frame, bits int) int {
	rescaled := float64(s) * float64(maxValues[bits])
	return int(rescaled)
}

func main() {
	// Read file into batches
	filepath := WaveFile
	wav, err := wave.ReadWaveFile(filepath)
	if err != nil {
		fmt.Printf("Could not read wave file: %v", err)
	}

	batches := wave.BatchSamples(wav, MaxSecondsPerPacket)

	iface, err := net.InterfaceByName(Iface)
	if err != nil {
		log.Printf("unable to find interface %s: %v", Iface, err)
	}

	ethFrames := []ethernet.Frame{}
	// Process chunks to build ethernet packets
	for i, batch := range batches {
		waveFmt := wave.NewWaveFmt(1, 2, 44100, 16, []byte{})
		if i == 0 {
			wave.WriteFrames(batch, waveFmt, fmt.Sprintf(TestOutputFormat, i))
			// 2 bytes of data per frame in 16bits/sample rate, should have 759 packets
			fmt.Printf("%d\n", len(batch))

		}

		msg := samplesToRawData(batch, waveFmt)
		ethFrame := NewFrame(iface.HardwareAddr, msg)
		ethFrames = append(ethFrames, *ethFrame)
	}

	conn, err := raw.ListenPacket(iface, EtherType, nil)
	if err != nil {
		log.Fatalf("failed to create connection")
	}

	// Send the marshaled frame to the network.
	for _, frame := range ethFrames {

		//time.Sleep(time.Second / 100000000)
		fmt.Printf("msg: %x\n\n", frame.Payload)
		sendMessage(conn, frame)
	}
}
