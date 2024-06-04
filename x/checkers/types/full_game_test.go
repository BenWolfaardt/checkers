package types_test

import (
	"testing"

	"github.com/BenWolfaardt/checkers/x/checkers/rules"
	"github.com/BenWolfaardt/checkers/x/checkers/testutil"
	"github.com/BenWolfaardt/checkers/x/checkers/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

const (
	alice = testutil.Alice
	bob   = testutil.Bob
)

func GetStoredGame1() types.StoredGame {
	return types.StoredGame{
		Black: alice,
		Red:   bob,
		Index: "1",
		Board: rules.New().String(),
		Turn:  "b",
	}
}

func TestCanGetAddressBlack(t *testing.T) {
	aliceAddress, err1 := sdk.AccAddressFromBech32(alice)
	black, err2 := GetStoredGame1().GetBlackAddress()
	require.Equal(t, aliceAddress, black)
	require.Nil(t, err2)
	require.Nil(t, err1)
}

func TestGameAddressWrongBlack(t *testing.T) {
	storedGame := GetStoredGame1()
	storedGame.Black = "cosmos1jmjfq0tplp9tmx4v9uemw72y4d2wa5nr3xn9d4"
	black, err := storedGame.GetBlackAddress()
	require.Nil(t, black)
	require.EqualError(t,
		err,
		"black address is invalid: cosmos1jmjfq0tplp9tmx4v9uemw72y4d2wa5nr3xn9d4: decoding bech32 failed: invalid checksum (expected 3xn9d3 got 3xn9d4)")
	require.EqualError(t, storedGame.Validate(), err.Error())
}

func TestCanGetAddressRed(t *testing.T) {
	bobAddress, err1 := sdk.AccAddressFromBech32(bob)
	red, err2 := GetStoredGame1().GetRedAddress()
	require.Equal(t, bobAddress, red)
	require.Nil(t, err2)
	require.Nil(t, err1)
}

func TestGameAddressWrongRed(t *testing.T) {
	storedGame := GetStoredGame1()
	storedGame.Red = "cosmos1jmjfq0tplp9tmx4v9uemw72y4d2wa5nr3xn9d4"
	red, err := storedGame.GetRedAddress()
	require.Nil(t, red)
	require.EqualError(t,
		err,
		"red address is invalid: cosmos1jmjfq0tplp9tmx4v9uemw72y4d2wa5nr3xn9d4: decoding bech32 failed: invalid checksum (expected 3xn9d3 got 3xn9d4)")
	require.EqualError(t, storedGame.Validate(), err.Error())
}

func TestGetPlayerAddressBlackCorrect(t *testing.T) {
	storedGame := GetStoredGame1()
	black, found, err := storedGame.GetPlayerAddress("b")
	require.Equal(t, alice, black.String())
	require.True(t, found)
	require.Nil(t, err)
}

func TestGetPlayerAddressBlackIncorrect(t *testing.T) {
	storedGame := GetStoredGame1()
	storedGame.Black = "notanaddress"
	black, found, err := storedGame.GetPlayerAddress("b")
	require.Nil(t, black)
	require.False(t, found)
	require.EqualError(t, err, "black address is invalid: notanaddress: decoding bech32 failed: invalid separator index -1")
}

func TestGetPlayerAddressRedCorrect(t *testing.T) {
	storedGame := GetStoredGame1()
	red, found, err := storedGame.GetPlayerAddress("r")
	require.Equal(t, bob, red.String())
	require.True(t, found)
	require.Nil(t, err)
}

func TestGetPlayerAddressRedIncorrect(t *testing.T) {
	storedGame := GetStoredGame1()
	storedGame.Red = "notanaddress"
	red, found, err := storedGame.GetPlayerAddress("r")
	require.Nil(t, red)
	require.False(t, found)
	require.EqualError(t, err, "red address is invalid: notanaddress: decoding bech32 failed: invalid separator index -1")
}

func TestGetPlayerAddressWhiteNotFound(t *testing.T) {
	storedGame := GetStoredGame1()
	white, found, err := storedGame.GetPlayerAddress("w")
	require.Nil(t, white)
	require.False(t, found)
	require.Nil(t, err)
}

func TestGetPlayerAddressAnyNotFound(t *testing.T) {
	storedGame := GetStoredGame1()
	white, found, err := storedGame.GetPlayerAddress("*")
	require.Nil(t, white)
	require.False(t, found)
	require.Nil(t, err)
}

func TestGetWinnerBlackCorrect(t *testing.T) {
	storedGame := GetStoredGame1()
	storedGame.Winner = "b"
	winner, found, err := storedGame.GetWinnerAddress()
	require.Equal(t, alice, winner.String())
	require.True(t, found)
	require.Nil(t, err)
}

func TestGetWinnerRedCorrect(t *testing.T) {
	storedGame := GetStoredGame1()
	storedGame.Winner = "r"
	winner, found, err := storedGame.GetWinnerAddress()
	require.Equal(t, bob, winner.String())
	require.True(t, found)
	require.Nil(t, err)
}

func TestGetWinnerNotYetCorrect(t *testing.T) {
	storedGame := GetStoredGame1()
	winner, found, err := storedGame.GetWinnerAddress()
	require.Nil(t, winner)
	require.False(t, found)
	require.Nil(t, err)
}
