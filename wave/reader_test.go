package wave

import (
	"fmt"
	"testing"

	"github.com/DylanMeeus/GoAudio/wave"
)

func TestReadWaveFile(t *testing.T) {
	filepath := "../birds.wav"
	wav, err := wave.ReadWaveFile(filepath)
	if err != nil {
		fmt.Errorf("Could not read wave file: %w", err)
	}

	fmt.Printf("%d", len(wav.Frames))
	batches := wave.BatchSamples(wav, 1.0)
	fmt.Printf("%d", len(batches))

	waveFmt := wave.NewWaveFmt(1, 2, 44100, 16, []byte{})
	for i, batch := range batches {
		wave.WriteFrames(batch, waveFmt, fmt.Sprintf("../test/birds._pt%d.wav", i))
	}

}
