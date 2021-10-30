package main

import (
	"fmt"

	"github.com/DylanMeeus/GoAudio/wave"
)

func main() {
	// Read file into batches
	filepath := "./birds.wav"
	wav, err := wave.ReadWaveFile(filepath)
	if err != nil {
		fmt.Printf("Could not read wave file: %v", err)
	}

	batches := wave.BatchSamples(wav, 1.0)

	waveFmt := wave.NewWaveFmt(1, 2, 44100, 16, []byte{})
	for i, batch := range batches {
		wave.WriteFrames(batch, waveFmt, fmt.Sprintf("./birds._pt%d.wav", i))
	}
	// Process chunks to build ethernet packets

	// Send ethernet packets
}
