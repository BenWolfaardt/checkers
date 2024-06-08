package leaderboard_test

import (
	"testing"

	keepertest "github.com/BenWolfaardt/checkers/testutil/keeper"
	"github.com/BenWolfaardt/checkers/testutil/nullify"
	"github.com/BenWolfaardt/checkers/x/leaderboard"
	"github.com/BenWolfaardt/checkers/x/leaderboard/types"
	"github.com/stretchr/testify/require"
)

func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		Params: types.DefaultParams(),

		Leaderboard: types.Leaderboard{
			Winners: []types.Winner{
				{
					Address: "cosmos123",
				},
				{
					Address: "cosmos456",
				},
			},
		},
		// this line is used by starport scaffolding # genesis/test/state
	}

	k, ctx := keepertest.LeaderboardKeeper(t)
	leaderboard.InitGenesis(ctx, *k, genesisState)
	got := leaderboard.ExportGenesis(ctx, *k)
	require.NotNil(t, got)

	nullify.Fill(&genesisState)
	nullify.Fill(got)

	require.Equal(t, genesisState.Leaderboard, got.Leaderboard)
	// this line is used by starport scaffolding # genesis/test/assert
}
