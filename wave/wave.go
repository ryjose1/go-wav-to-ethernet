package wave

import (
	"encoding/binary"
	"fmt"
	"math"

	"github.com/ryjose1/go-wav-to-ethernet/config"

	"github.com/DylanMeeus/GoAudio/wave"
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

func GenerateBatches() [][]wave.Frame {
	filepath := config.WaveFile
	wav, err := wave.ReadWaveFile(filepath)
	if err != nil {
		fmt.Printf("Could not read wave file: %v", err)
	}

	return wave.BatchSamples(wav, config.MaxSecondsPerPacket)
}

func SamplesToPayload(batch []wave.Frame) []byte {
	waveFmt := wave.NewWaveFmt(1, 2, 44100, 16, []byte{})
	/*
		if i == 0 {
			TestOutputFormat := "./test/birds._pt%d.wav"
			wave.WriteFrames(batch, waveFmt, fmt.Sprintf(TestOutputFormat, i))
			// 2 bytes of data per frame in 16bits/sample rate, should have 759 packets
			fmt.Printf("%d\n", len(batch))

		}
	*/
	return samplesToRawData(batch, waveFmt)
}
