package types

import "github.com/BenWolfaardt/checkers/x/leaderboard/types"

const (
	PlayerInfoChunkSize = types.DefaultLength * uint64(2)
	ConsensusVersion    = uint64(2)
)
