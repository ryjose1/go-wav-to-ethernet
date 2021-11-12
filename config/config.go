package config

const (
	// Equivalent to 1518 bytes of audio data
	// 1400 bytes/packet * 8 bits/byte * 1 channel/16 bits * 1 sample/2 channels *1 second/44100 samples
	MaxSecondsPerPacket = .0079365079
	OneSecond           = 1.0
	WaveFile            = "birds"
	EtherType           = 0xcccc
	Iface               = "enp0s8"
)
