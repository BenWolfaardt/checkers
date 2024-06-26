package keeper_test

import (
	"strconv"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	keepertest "github.com/BenWolfaardt/checkers/testutil/keeper"
	"github.com/BenWolfaardt/checkers/testutil/nullify"
	"github.com/BenWolfaardt/checkers/x/checkers/types"
)

// Prevent strconv unused error
var _ = strconv.IntSize

func TestPlayerInfoQuerySingle(t *testing.T) {
	keeper, ctx := keepertest.CheckersKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	msgs := createNPlayerInfo(keeper, ctx, 2)
	for _, tc := range []struct {
		desc     string
		request  *types.QueryGetPlayerInfoRequest
		response *types.QueryGetPlayerInfoResponse
		err      error
	}{
		{
			desc: "First",
			request: &types.QueryGetPlayerInfoRequest{
				Index: msgs[0].Index,
			},
			response: &types.QueryGetPlayerInfoResponse{PlayerInfo: msgs[0]},
		},
		{
			desc: "Second",
			request: &types.QueryGetPlayerInfoRequest{
				Index: msgs[1].Index,
			},
			response: &types.QueryGetPlayerInfoResponse{PlayerInfo: msgs[1]},
		},
		{
			desc: "KeyNotFound",
			request: &types.QueryGetPlayerInfoRequest{
				Index: strconv.Itoa(100000),
			},
			err: status.Error(codes.NotFound, "not found"),
		},
		{
			desc: "InvalidRequest",
			err:  status.Error(codes.InvalidArgument, "invalid request"),
		},
	} {
		t.Run(tc.desc, func(t *testing.T) {
			response, err := keeper.PlayerInfo(wctx, tc.request)
			if tc.err != nil {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				require.Equal(t,
					nullify.Fill(tc.response),
					nullify.Fill(response),
				)
			}
		})
	}
}

func TestPlayerInfoQueryPaginated(t *testing.T) {
	keeper, ctx := keepertest.CheckersKeeper(t)
	wctx := sdk.WrapSDKContext(ctx)
	msgs := createNPlayerInfo(keeper, ctx, 5)

	request := func(next []byte, offset, limit uint64, total bool) *types.QueryAllPlayerInfoRequest {
		return &types.QueryAllPlayerInfoRequest{
			Pagination: &query.PageRequest{
				Key:        next,
				Offset:     offset,
				Limit:      limit,
				CountTotal: total,
			},
		}
	}
	t.Run("ByOffset", func(t *testing.T) {
		step := 2
		for i := 0; i < len(msgs); i += step {
			resp, err := keeper.PlayerInfoAll(wctx, request(nil, uint64(i), uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.PlayerInfo), step)
			require.Subset(t,
				nullify.Fill(msgs),
				nullify.Fill(resp.PlayerInfo),
			)
		}
	})
	t.Run("ByKey", func(t *testing.T) {
		step := 2
		var next []byte
		for i := 0; i < len(msgs); i += step {
			resp, err := keeper.PlayerInfoAll(wctx, request(next, 0, uint64(step), false))
			require.NoError(t, err)
			require.LessOrEqual(t, len(resp.PlayerInfo), step)
			require.Subset(t,
				nullify.Fill(msgs),
				nullify.Fill(resp.PlayerInfo),
			)
			next = resp.Pagination.NextKey
		}
	})
	t.Run("Total", func(t *testing.T) {
		resp, err := keeper.PlayerInfoAll(wctx, request(nil, 0, 0, true))
		require.NoError(t, err)
		require.Equal(t, len(msgs), int(resp.Pagination.Total))
		require.ElementsMatch(t,
			nullify.Fill(msgs),
			nullify.Fill(resp.PlayerInfo),
		)
	})
	t.Run("InvalidRequest", func(t *testing.T) {
		_, err := keeper.PlayerInfoAll(wctx, nil)
		require.ErrorIs(t, err, status.Error(codes.InvalidArgument, "invalid request"))
	})
}
