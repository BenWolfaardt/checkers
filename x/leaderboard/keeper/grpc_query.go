package keeper

import (
	"github.com/BenWolfaardt/checkers/x/leaderboard/types"
)

var _ types.QueryServer = Keeper{}
