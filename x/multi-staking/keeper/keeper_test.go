package keeper_test

import (
	"testing"

	"github.com/realio-tech/multi-staking-module/testing/simapp"
	multistakingkeeper "github.com/realio-tech/multi-staking-module/x/multi-staking/keeper"
	"github.com/stretchr/testify/suite"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
)

type KeeperTestSuite struct {
	suite.Suite

	ctx           sdk.Context
	app           simapp.SimApp
	msKeeper      *multistakingkeeper.Keeper
	stakingKeeper *stakingkeeper.Keeper
}

func (suite *KeeperTestSuite) SetupTest() {
	app := simapp.Setup(false)
	ctx := app.BaseApp.NewContext(false, tmproto.Header{})

	suite.app, suite.ctx, suite.msKeeper, suite.stakingKeeper = *app, ctx, &app.MultiStakingKeeper, &app.StakingKeeper
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}
