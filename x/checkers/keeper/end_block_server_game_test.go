package keeper_test

import (
	"testing"
	"time"

	"github.com/BenWolfaardt/checkers/x/checkers/testutil"
	"github.com/BenWolfaardt/checkers/x/checkers/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

func TestForfeitUnplayed(t *testing.T) {
	_, keeper, context := setupMsgServerWithOneGameForPlayMove(t)
	ctx := sdk.UnwrapSDKContext(context)
	game1, found := keeper.GetStoredGame(ctx, "1")
	require.True(t, found)
	game1.Deadline = types.FormatDeadline(ctx.BlockTime().Add(time.Duration(-1)))
	keeper.SetStoredGame(ctx, game1)
	keeper.ForfeitExpiredGames(context)

	_, found = keeper.GetStoredGame(ctx, "1")
	require.False(t, found)

	systemInfo, found := keeper.GetSystemInfo(ctx)
	require.True(t, found)
	require.EqualValues(t, types.SystemInfo{
		NextId:        2,
		FifoHeadIndex: "-1",
		FifoTailIndex: "-1",
	}, systemInfo)
	events := sdk.StringifyEvents(ctx.EventManager().ABCIEvents())
	require.Len(t, events, 2)
	event := events[0]
	require.EqualValues(t, sdk.StringEvent{
		Type: "game-forfeited",
		Attributes: []sdk.Attribute{
			{Key: "game-index", Value: "1"},
			{Key: "winner", Value: "*"},
			{Key: "board", Value: "*b*b*b*b|b*b*b*b*|*b*b*b*b|********|********|r*r*r*r*|*r*r*r*r|r*r*r*r*"},
		},
	}, event)
}

func TestForfeitOlderUnplayed(t *testing.T) {
	msgServer, keeper, context := setupMsgServerWithOneGameForPlayMove(t)
	ctx := sdk.UnwrapSDKContext(context)
	msgServer.CreateGame(context, &types.MsgCreateGame{
		Creator: bob,
		Black:   carol,
		Red:     alice,
		Wager:   45,
		Denom:   "coin",
	})
	game1, found := keeper.GetStoredGame(ctx, "1")
	require.True(t, found)
	game1.Deadline = types.FormatDeadline(ctx.BlockTime().Add(time.Duration(-1)))
	keeper.SetStoredGame(ctx, game1)
	keeper.ForfeitExpiredGames(context)

	_, found = keeper.GetStoredGame(ctx, "1")
	require.False(t, found)

	nextGame, found := keeper.GetSystemInfo(ctx)
	require.True(t, found)
	require.EqualValues(t, types.SystemInfo{
		NextId:        3,
		FifoHeadIndex: "2",
		FifoTailIndex: "2",
	}, nextGame)
	events := sdk.StringifyEvents(ctx.EventManager().ABCIEvents())
	require.Len(t, events, 2)
	event := events[0]
	require.EqualValues(t, sdk.StringEvent{
		Type: "game-forfeited",
		Attributes: []sdk.Attribute{
			{Key: "game-index", Value: "1"},
			{Key: "winner", Value: "*"},
			{Key: "board", Value: "*b*b*b*b|b*b*b*b*|*b*b*b*b|********|********|r*r*r*r*|*r*r*r*r|r*r*r*r*"},
		},
	}, event)
}

func TestForfeit2OldestUnplayedIn1Call(t *testing.T) {
	msgServer, keeper, context := setupMsgServerWithOneGameForPlayMove(t)
	ctx := sdk.UnwrapSDKContext(context)
	msgServer.CreateGame(context, &types.MsgCreateGame{
		Creator: bob,
		Black:   carol,
		Red:     alice,
		Wager:   45,
		Denom:   "coin",
	})
	msgServer.CreateGame(context, &types.MsgCreateGame{
		Creator: carol,
		Black:   alice,
		Red:     bob,
		Wager:   45,
		Denom:   "coin",
	})
	game1, found := keeper.GetStoredGame(ctx, "1")
	require.True(t, found)
	game1.Deadline = types.FormatDeadline(ctx.BlockTime().Add(time.Duration(-1)))
	keeper.SetStoredGame(ctx, game1)
	game2, found := keeper.GetStoredGame(ctx, "2")
	require.True(t, found)
	game2.Deadline = types.FormatDeadline(ctx.BlockTime().Add(time.Duration(-1)))
	keeper.SetStoredGame(ctx, game2)
	keeper.ForfeitExpiredGames(context)

	_, found = keeper.GetStoredGame(ctx, "1")
	require.False(t, found)
	_, found = keeper.GetStoredGame(ctx, "2")
	require.False(t, found)

	systemInfo, found := keeper.GetSystemInfo(ctx)
	require.True(t, found)
	require.EqualValues(t, types.SystemInfo{
		NextId:        4,
		FifoHeadIndex: "3",
		FifoTailIndex: "3",
	}, systemInfo)
	events := sdk.StringifyEvents(ctx.EventManager().ABCIEvents())
	require.Len(t, events, 2)
	event := events[0]
	require.EqualValues(t, sdk.StringEvent{
		Type: "game-forfeited",
		Attributes: []sdk.Attribute{
			{Key: "game-index", Value: "1"},
			{Key: "winner", Value: "*"},
			{Key: "board", Value: "*b*b*b*b|b*b*b*b*|*b*b*b*b|********|********|r*r*r*r*|*r*r*r*r|r*r*r*r*"},
			{Key: "game-index", Value: "2"},
			{Key: "winner", Value: "*"},
			{Key: "board", Value: "*b*b*b*b|b*b*b*b*|*b*b*b*b|********|********|r*r*r*r*|*r*r*r*r|r*r*r*r*"},
		},
	}, event)
}

func TestForfeitOlderPlayedOnce(t *testing.T) {
	msgServer, keeper, context := setupMsgServerWithOneGameForPlayMove(t)
	ctx := sdk.UnwrapSDKContext(context)
	msgServer.PlayMove(context, &types.MsgPlayMove{
		Creator:   bob,
		GameIndex: "1",
		FromX:     1,
		FromY:     2,
		ToX:       2,
		ToY:       3,
	})
	msgServer.CreateGame(context, &types.MsgCreateGame{
		Creator: bob,
		Red:     carol,
		Black:   alice,
		Wager:   46,
		Denom:   "coin",
	})
	game1, found := keeper.GetStoredGame(ctx, "1")
	require.True(t, found)
	game1.Deadline = types.FormatDeadline(ctx.BlockTime().Add(time.Duration(-1)))
	keeper.SetStoredGame(ctx, game1)
	keeper.ForfeitExpiredGames(context)

	_, found = keeper.GetStoredGame(ctx, "1")
	require.False(t, found)

	systemInfo, found := keeper.GetSystemInfo(ctx)
	require.True(t, found)
	require.EqualValues(t, types.SystemInfo{
		NextId:        3,
		FifoHeadIndex: "2",
		FifoTailIndex: "2",
	}, systemInfo)
	events := sdk.StringifyEvents(ctx.EventManager().ABCIEvents())
	require.Len(t, events, 3)
	event := events[0]
	require.EqualValues(t, sdk.StringEvent{
		Type: "game-forfeited",
		Attributes: []sdk.Attribute{
			{Key: "game-index", Value: "1"},
			{Key: "winner", Value: "*"},
			{Key: "board", Value: "*b*b*b*b|b*b*b*b*|***b*b*b|**b*****|********|r*r*r*r*|*r*r*r*r|r*r*r*r*"},
		},
	}, event)
}
func TestForfeit2OldestPlayedOnceIn1Call(t *testing.T) {
	msgServer, keeper, context := setupMsgServerWithOneGameForPlayMove(t)
	ctx := sdk.UnwrapSDKContext(context)
	msgServer.PlayMove(context, &types.MsgPlayMove{
		Creator:   bob,
		GameIndex: "1",
		FromX:     1,
		FromY:     2,
		ToX:       2,
		ToY:       3,
	})
	msgServer.CreateGame(context, &types.MsgCreateGame{
		Creator: bob,
		Black:   carol,
		Red:     alice,
		Wager:   45,
		Denom:   "coin",
	})
	msgServer.PlayMove(context, &types.MsgPlayMove{
		Creator:   carol,
		GameIndex: "2",
		FromX:     1,
		FromY:     2,
		ToX:       2,
		ToY:       3,
	})
	msgServer.CreateGame(context, &types.MsgCreateGame{
		Creator: carol,
		Black:   alice,
		Red:     bob,
		Wager:   45,
		Denom:   "coin",
	})
	game1, found := keeper.GetStoredGame(ctx, "1")
	require.True(t, found)
	game1.Deadline = types.FormatDeadline(ctx.BlockTime().Add(time.Duration(-1)))
	keeper.SetStoredGame(ctx, game1)
	game2, found := keeper.GetStoredGame(ctx, "2")
	require.True(t, found)
	game2.Deadline = types.FormatDeadline(ctx.BlockTime().Add(time.Duration(-1)))
	keeper.SetStoredGame(ctx, game2)
	keeper.ForfeitExpiredGames(context)

	_, found = keeper.GetStoredGame(ctx, "1")
	require.False(t, found)
	_, found = keeper.GetStoredGame(ctx, "2")
	require.False(t, found)

	systemInfo, found := keeper.GetSystemInfo(ctx)
	require.True(t, found)
	require.EqualValues(t, types.SystemInfo{
		NextId:        4,
		FifoHeadIndex: "3",
		FifoTailIndex: "3",
	}, systemInfo)
	events := sdk.StringifyEvents(ctx.EventManager().ABCIEvents())
	require.Len(t, events, 3)
	event := events[0]
	require.EqualValues(t,
		sdk.StringEvent{
			Type: "game-forfeited",
			Attributes: []sdk.Attribute{
				{Key: "game-index", Value: "1"},
				{Key: "winner", Value: "*"},
				{Key: "board", Value: "*b*b*b*b|b*b*b*b*|***b*b*b|**b*****|********|r*r*r*r*|*r*r*r*r|r*r*r*r*"},
				{Key: "game-index", Value: "2"},
				{Key: "winner", Value: "*"},
				{Key: "board", Value: "*b*b*b*b|b*b*b*b*|***b*b*b|**b*****|********|r*r*r*r*|*r*r*r*r|r*r*r*r*"},
			},
		}, event)
}

func TestForfeitPlayedOnce(t *testing.T) {
	msgServer, keeper, context := setupMsgServerWithOneGameForPlayMove(t)
	ctx := sdk.UnwrapSDKContext(context)
	msgServer.PlayMove(context, &types.MsgPlayMove{
		Creator:   bob,
		GameIndex: "1",
		FromX:     1,
		FromY:     2,
		ToX:       2,
		ToY:       3,
	})
	game1, found := keeper.GetStoredGame(ctx, "1")
	require.True(t, found)
	game1.Deadline = types.FormatDeadline(ctx.BlockTime().Add(time.Duration(-1)))
	keeper.SetStoredGame(ctx, game1)
	keeper.ForfeitExpiredGames(context)

	_, found = keeper.GetStoredGame(ctx, "1")
	require.False(t, found)

	systemInfo, found := keeper.GetSystemInfo(ctx)
	require.True(t, found)
	require.EqualValues(t, types.SystemInfo{
		NextId:        2,
		FifoHeadIndex: "-1",
		FifoTailIndex: "-1",
	}, systemInfo)
	events := sdk.StringifyEvents(ctx.EventManager().ABCIEvents())
	require.Len(t, events, 3)
	event := events[0]
	require.EqualValues(t, sdk.StringEvent{
		Type: "game-forfeited",
		Attributes: []sdk.Attribute{
			{Key: "game-index", Value: "1"},
			{Key: "winner", Value: "*"},
			{Key: "board", Value: "*b*b*b*b|b*b*b*b*|***b*b*b|**b*****|********|r*r*r*r*|*r*r*r*r|r*r*r*r*"},
		},
	}, event)
}

func TestForfeitOlderPlayedTwice(t *testing.T) {
	msgServer, keeper, context := setupMsgServerWithOneGameForPlayMove(t)
	ctx := sdk.UnwrapSDKContext(context)
	msgServer.PlayMove(context, &types.MsgPlayMove{
		Creator:   bob,
		GameIndex: "1",
		FromX:     1,
		FromY:     2,
		ToX:       2,
		ToY:       3,
	})
	msgServer.PlayMove(context, &types.MsgPlayMove{
		Creator:   carol,
		GameIndex: "1",
		FromX:     0,
		FromY:     5,
		ToX:       1,
		ToY:       4,
	})
	msgServer.CreateGame(context, &types.MsgCreateGame{
		Creator: bob,
		Black:   carol,
		Red:     alice,
		Wager:   45,
		Denom:   "coin",
	})
	game1, found := keeper.GetStoredGame(ctx, "1")
	require.True(t, found)
	oldDeadline := types.FormatDeadline(ctx.BlockTime().Add(time.Duration(-1)))
	game1.Deadline = oldDeadline
	keeper.SetStoredGame(ctx, game1)
	keeper.ForfeitExpiredGames(context)

	game1, found = keeper.GetStoredGame(ctx, "1")
	require.True(t, found)
	require.EqualValues(t, types.StoredGame{
		Index:       "1",
		Board:       "",
		Turn:        "b",
		Black:       bob,
		Red:         carol,
		Winner:      "r",
		Deadline:    oldDeadline,
		MoveCount:   uint64(2),
		BeforeIndex: "-1",
		AfterIndex:  "-1",
		Wager:       45,
		Denom:       "stake",
	}, game1)

	systemInfo, found := keeper.GetSystemInfo(ctx)
	require.True(t, found)
	require.EqualValues(t, types.SystemInfo{
		NextId:        3,
		FifoHeadIndex: "2",
		FifoTailIndex: "2",
	}, systemInfo)
	events := sdk.StringifyEvents(ctx.EventManager().ABCIEvents())
	require.Len(t, events, 3)
	event := events[0]
	require.EqualValues(t, sdk.StringEvent{
		Type: "game-forfeited",
		Attributes: []sdk.Attribute{
			{Key: "game-index", Value: "1"},
			{Key: "winner", Value: "r"},
			{Key: "board", Value: "*b*b*b*b|b*b*b*b*|***b*b*b|**b*****|*r******|**r*r*r*|*r*r*r*r|r*r*r*r*"},
		},
	}, event)
}

func TestForfeit2OldestPlayedTwiceIn1Call(t *testing.T) {
	msgServer, keeper, context := setupMsgServerWithOneGameForPlayMove(t)
	ctx := sdk.UnwrapSDKContext(context)
	msgServer.PlayMove(context, &types.MsgPlayMove{
		Creator:   bob,
		GameIndex: "1",
		FromX:     1,
		FromY:     2,
		ToX:       2,
		ToY:       3,
	})
	msgServer.PlayMove(context, &types.MsgPlayMove{
		Creator:   carol,
		GameIndex: "1",
		FromX:     0,
		FromY:     5,
		ToX:       1,
		ToY:       4,
	})
	msgServer.CreateGame(context, &types.MsgCreateGame{
		Creator: bob,
		Black:   carol,
		Red:     alice,
		Wager:   45,
		Denom:   "coin",
	})
	msgServer.PlayMove(context, &types.MsgPlayMove{
		Creator:   carol,
		GameIndex: "2",
		FromX:     1,
		FromY:     2,
		ToX:       2,
		ToY:       3,
	})
	msgServer.PlayMove(context, &types.MsgPlayMove{
		Creator:   alice,
		GameIndex: "2",
		FromX:     0,
		FromY:     5,
		ToX:       1,
		ToY:       4,
	})
	msgServer.CreateGame(context, &types.MsgCreateGame{
		Creator: carol,
		Black:   alice,
		Red:     bob,
		Wager:   45,
		Denom:   "coin",
	})
	game1, found := keeper.GetStoredGame(ctx, "1")
	require.True(t, found)
	oldDeadline := types.FormatDeadline(ctx.BlockTime().Add(time.Duration(-1)))
	game1.Deadline = oldDeadline
	keeper.SetStoredGame(ctx, game1)
	game2, found := keeper.GetStoredGame(ctx, "2")
	require.True(t, found)
	game2.Deadline = oldDeadline
	keeper.SetStoredGame(ctx, game2)
	keeper.ForfeitExpiredGames(context)

	game1, found = keeper.GetStoredGame(ctx, "1")
	require.True(t, found)
	require.EqualValues(t, types.StoredGame{
		Index:       "1",
		Board:       "",
		Turn:        "b",
		Black:       bob,
		Red:         carol,
		Winner:      "r",
		Deadline:    oldDeadline,
		MoveCount:   uint64(2),
		BeforeIndex: "-1",
		AfterIndex:  "-1",
		Wager:       45,
		Denom:       "stake",
	}, game1)

	game2, found = keeper.GetStoredGame(ctx, "2")
	require.True(t, found)
	require.EqualValues(t, types.StoredGame{
		Index:       "2",
		Board:       "",
		Turn:        "b",
		Black:       carol,
		Red:         alice,
		Winner:      "r",
		Deadline:    oldDeadline,
		MoveCount:   uint64(2),
		BeforeIndex: "-1",
		AfterIndex:  "-1",
		Wager:       45,
		Denom:       "coin",
	}, game2)

	systemInfo, found := keeper.GetSystemInfo(ctx)
	require.True(t, found)
	require.EqualValues(t, types.SystemInfo{
		NextId:        4,
		FifoHeadIndex: "3",
		FifoTailIndex: "3",
	}, systemInfo)

	events := sdk.StringifyEvents(ctx.EventManager().ABCIEvents())
	require.Len(t, events, 3)
	event := events[0]
	require.EqualValues(t, sdk.StringEvent{
		Type: "game-forfeited",
		Attributes: []sdk.Attribute{
			{Key: "game-index", Value: "1"},
			{Key: "winner", Value: "r"},
			{Key: "board", Value: "*b*b*b*b|b*b*b*b*|***b*b*b|**b*****|*r******|**r*r*r*|*r*r*r*r|r*r*r*r*"},
			{Key: "game-index", Value: "2"},
			{Key: "winner", Value: "r"},
			{Key: "board", Value: "*b*b*b*b|b*b*b*b*|***b*b*b|**b*****|*r******|**r*r*r*|*r*r*r*r|r*r*r*r*"},
		},
	}, event)
}

func TestForfeitPlayedTwice(t *testing.T) {
	msgServer, keeper, context := setupMsgServerWithOneGameForPlayMove(t)
	ctx := sdk.UnwrapSDKContext(context)
	msgServer.PlayMove(context, &types.MsgPlayMove{
		Creator:   bob,
		GameIndex: "1",
		FromX:     1,
		FromY:     2,
		ToX:       2,
		ToY:       3,
	})
	msgServer.PlayMove(context, &types.MsgPlayMove{
		Creator:   carol,
		GameIndex: "1",
		FromX:     0,
		FromY:     5,
		ToX:       1,
		ToY:       4,
	})
	game1, found := keeper.GetStoredGame(ctx, "1")
	require.True(t, found)
	oldDeadline := types.FormatDeadline(ctx.BlockTime().Add(time.Duration(-1)))
	game1.Deadline = oldDeadline
	keeper.SetStoredGame(ctx, game1)
	keeper.ForfeitExpiredGames(context)

	game1, found = keeper.GetStoredGame(ctx, "1")
	require.True(t, found)
	require.EqualValues(t, types.StoredGame{
		Index:       "1",
		Board:       "",
		Turn:        "b",
		Black:       bob,
		Red:         carol,
		Winner:      "r",
		Deadline:    oldDeadline,
		MoveCount:   uint64(2),
		BeforeIndex: "-1",
		AfterIndex:  "-1",
		Wager:       45,
		Denom:       "stake",
	}, game1)

	systemInfo, found := keeper.GetSystemInfo(ctx)
	require.True(t, found)
	require.EqualValues(t, types.SystemInfo{
		NextId:        2,
		FifoHeadIndex: "-1",
		FifoTailIndex: "-1",
	}, systemInfo)
	events := sdk.StringifyEvents(ctx.EventManager().ABCIEvents())
	require.Len(t, events, 3)
	event := events[0]
	require.EqualValues(t, sdk.StringEvent{
		Type: "game-forfeited",
		Attributes: []sdk.Attribute{
			{Key: "game-index", Value: "1"},
			{Key: "winner", Value: "r"},
			{Key: "board", Value: "*b*b*b*b|b*b*b*b*|***b*b*b|**b*****|*r******|**r*r*r*|*r*r*r*r|r*r*r*r*"},
		},
	}, event)
}

func TestForfeitPlayedOnceCalledBank(t *testing.T) {
	msgServer, keeper, context, ctrl, escrow := setupMsgServerWithOneGameForPlayMoveWithMock(t)
	ctx := sdk.UnwrapSDKContext(context)
	defer ctrl.Finish()
	pay := escrow.ExpectPay(context, bob, 45).Times(1)
	escrow.ExpectRefund(context, bob, 45).Times(1).After(pay)
	msgServer.PlayMove(context, &types.MsgPlayMove{
		Creator:   bob,
		GameIndex: "1",
		FromX:     1,
		FromY:     2,
		ToX:       2,
		ToY:       3,
	})
	game1, found := keeper.GetStoredGame(ctx, "1")
	require.True(t, found)
	game1.Deadline = types.FormatDeadline(ctx.BlockTime().Add(time.Duration(-1)))
	keeper.SetStoredGame(ctx, game1)
	keeper.ForfeitExpiredGames(context)
}

func TestForfeitPlayedTwiceCalledBank(t *testing.T) {
	msgServer, keeper, context, ctrl, escrow := setupMsgServerWithOneGameForPlayMoveWithMock(t)
	ctx := sdk.UnwrapSDKContext(context)
	defer ctrl.Finish()
	payBob := escrow.ExpectPay(context, bob, 45).Times(1)
	payCarol := escrow.ExpectPay(context, carol, 45).Times(1).After(payBob)
	escrow.ExpectRefund(context, carol, 90).Times(1).After(payCarol)
	msgServer.PlayMove(context, &types.MsgPlayMove{
		Creator:   bob,
		GameIndex: "1",
		FromX:     1,
		FromY:     2,
		ToX:       2,
		ToY:       3,
	})
	msgServer.PlayMove(context, &types.MsgPlayMove{
		Creator:   carol,
		GameIndex: "1",
		FromX:     0,
		FromY:     5,
		ToX:       1,
		ToY:       4,
	})
	game1, found := keeper.GetStoredGame(ctx, "1")
	require.True(t, found)
	oldDeadline := types.FormatDeadline(ctx.BlockTime().Add(time.Duration(-1)))
	game1.Deadline = oldDeadline
	keeper.SetStoredGame(ctx, game1)
	keeper.ForfeitExpiredGames(context)
}

func TestForfeitGameAddPlayerInfo(t *testing.T) {
	msgServer, keeper, context := setupMsgServerWithOneGameForPlayMove(t)
	ctx := sdk.UnwrapSDKContext(context)

	msgServer.PlayMove(context, &types.MsgPlayMove{
		Creator:   bob,
		GameIndex: "1",
		FromX:     1,
		FromY:     2,
		ToX:       2,
		ToY:       3,
	})
	msgServer.PlayMove(context, &types.MsgPlayMove{
		Creator:   carol,
		GameIndex: "1",
		FromX:     0,
		FromY:     5,
		ToX:       1,
		ToY:       4,
	})
	game1, found := keeper.GetStoredGame(ctx, "1")
	require.True(t, found)
	oldDeadline := types.FormatDeadline(ctx.BlockTime().Add(time.Duration(-1)))
	game1.Deadline = oldDeadline
	keeper.SetStoredGame(ctx, game1)
	keeper.ForfeitExpiredGames(context)

	bobInfo, found := keeper.GetPlayerInfo(ctx, bob)
	require.True(t, found)
	require.EqualValues(t, types.PlayerInfo{
		Index:          bob,
		WonCount:       0,
		LostCount:      0,
		ForfeitedCount: 1,
	}, bobInfo)
	carolInfo, found := keeper.GetPlayerInfo(ctx, carol)
	require.True(t, found)
	require.EqualValues(t, types.PlayerInfo{
		Index:          carol,
		WonCount:       1,
		LostCount:      0,
		ForfeitedCount: 0,
	}, carolInfo)
}

func TestForfeiGameUpdatePlayerInfo(t *testing.T) {
	msgServer, keeper, context := setupMsgServerWithOneGameForPlayMove(t)
	ctx := sdk.UnwrapSDKContext(context)

	keeper.SetPlayerInfo(ctx, types.PlayerInfo{
		Index:          bob,
		WonCount:       1,
		LostCount:      2,
		ForfeitedCount: 3,
	})
	keeper.SetPlayerInfo(ctx, types.PlayerInfo{
		Index:          carol,
		WonCount:       4,
		LostCount:      5,
		ForfeitedCount: 6,
	})

	msgServer.PlayMove(context, &types.MsgPlayMove{
		Creator:   bob,
		GameIndex: "1",
		FromX:     1,
		FromY:     2,
		ToX:       2,
		ToY:       3,
	})
	msgServer.PlayMove(context, &types.MsgPlayMove{
		Creator:   carol,
		GameIndex: "1",
		FromX:     0,
		FromY:     5,
		ToX:       1,
		ToY:       4,
	})
	game1, found := keeper.GetStoredGame(ctx, "1")
	require.True(t, found)
	oldDeadline := types.FormatDeadline(ctx.BlockTime().Add(time.Duration(-1)))
	game1.Deadline = oldDeadline
	keeper.SetStoredGame(ctx, game1)
	keeper.ForfeitExpiredGames(context)

	bobInfo, found := keeper.GetPlayerInfo(ctx, bob)
	require.True(t, found)
	require.EqualValues(t, types.PlayerInfo{
		Index:          bob,
		WonCount:       1,
		LostCount:      2,
		ForfeitedCount: 4,
	}, bobInfo)
	carolInfo, found := keeper.GetPlayerInfo(ctx, carol)
	require.True(t, found)
	require.EqualValues(t, types.PlayerInfo{
		Index:          carol,
		WonCount:       5,
		LostCount:      5,
		ForfeitedCount: 6,
	}, carolInfo)
}

func TestForfeiGameCallsHook(t *testing.T) {
	msgServer, keeper, context, ctrl, hookMock := setupMsgServerWithOneGameForPlayMoveAndHooks(t)
	ctx := sdk.UnwrapSDKContext(context)
	defer ctrl.Finish()
	carolCall := hookMock.EXPECT().AfterPlayerInfoChanged(ctx, types.PlayerInfo{
		Index:          carol,
		WonCount:       5,
		LostCount:      5,
		ForfeitedCount: 6,
	}).Times(1)
	hookMock.EXPECT().AfterPlayerInfoChanged(ctx, types.PlayerInfo{
		Index:          bob,
		WonCount:       1,
		LostCount:      2,
		ForfeitedCount: 4,
	}).Times(1).After(carolCall)

	keeper.SetPlayerInfo(ctx, types.PlayerInfo{
		Index:          bob,
		WonCount:       1,
		LostCount:      2,
		ForfeitedCount: 3,
	})
	keeper.SetPlayerInfo(ctx, types.PlayerInfo{
		Index:          carol,
		WonCount:       4,
		LostCount:      5,
		ForfeitedCount: 6,
	})

	msgServer.PlayMove(context, &types.MsgPlayMove{
		Creator:   bob,
		GameIndex: "1",
		FromX:     1,
		FromY:     2,
		ToX:       2,
		ToY:       3,
	})
	msgServer.PlayMove(context, &types.MsgPlayMove{
		Creator:   carol,
		GameIndex: "1",
		FromX:     0,
		FromY:     5,
		ToX:       1,
		ToY:       4,
	})
	game1, found := keeper.GetStoredGame(ctx, "1")
	require.True(t, found)
	oldDeadline := types.FormatDeadline(ctx.BlockTime().Add(time.Duration(-1)))
	game1.Deadline = oldDeadline
	keeper.SetStoredGame(ctx, game1)
	keeper.ForfeitExpiredGames(context)
}

func TestCompleteGameCallsHook(t *testing.T) {
	msgServer, keeper, context, ctrl, hookMock := setupMsgServerWithOneGameForPlayMoveAndHooks(t)
	ctx := sdk.UnwrapSDKContext(context)
	defer ctrl.Finish()
	bobCall := hookMock.EXPECT().AfterPlayerInfoChanged(ctx, types.PlayerInfo{
		Index:          bob,
		WonCount:       2,
		LostCount:      2,
		ForfeitedCount: 3,
	}).Times(1)
	hookMock.EXPECT().AfterPlayerInfoChanged(ctx, types.PlayerInfo{
		Index:          carol,
		WonCount:       4,
		LostCount:      6,
		ForfeitedCount: 6,
	}).Times(1).After(bobCall)

	keeper.SetPlayerInfo(ctx, types.PlayerInfo{
		Index:          bob,
		WonCount:       1,
		LostCount:      2,
		ForfeitedCount: 3,
	})
	keeper.SetPlayerInfo(ctx, types.PlayerInfo{
		Index:          carol,
		WonCount:       4,
		LostCount:      5,
		ForfeitedCount: 6,
	})

	testutil.PlayAllMoves(t, msgServer, context, "1", bob, carol, testutil.Game1Moves)
}
