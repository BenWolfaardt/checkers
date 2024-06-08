package types

const (
	ConsensusVersion = uint64(3)
	// To find the ideal chunk size value, you would have to test with the real state and try different values.
	StoredGameChunkSize = 1_000
)
