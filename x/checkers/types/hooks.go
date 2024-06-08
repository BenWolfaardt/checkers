package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Take inspiration from: https://github.com/cosmos/cosmos-sdk/blob/v0.45.4/x/gov/types/hooks.go#L7-L20
var _ CheckersHooks = MultiCheckersHooks{}

type MultiCheckersHooks []CheckersHooks

func NewMultiCheckersHooks(hooks ...CheckersHooks) MultiCheckersHooks {
	return hooks
}

func (h MultiCheckersHooks) AfterPlayerInfoChanged(ctx sdk.Context, playerInfo PlayerInfo) {
	for i := range h {
		h[i].AfterPlayerInfoChanged(ctx, playerInfo)
	}
}
