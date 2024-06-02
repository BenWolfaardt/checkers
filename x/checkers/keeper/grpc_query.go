package keeper

import (
	"github.com/BenWolfaardt/checkers/x/checkers/types"
)

var _ types.QueryServer = Keeper{}
