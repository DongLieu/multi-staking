package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/realio-tech/multi-staking-module/testutil"
	"github.com/realio-tech/multi-staking-module/x/multi-staking/types"
)

func (suite *KeeperTestSuite) TestEndBlocker() {
	suite.SetupTest()

	multiStakerAddress := testutil.GenAddress()
	valPubKey := testutil.GenPubKey()
	valAddr := sdk.ValAddress(valPubKey.Address())
	intDel := types.IntermediaryDelegator(multiStakerAddress)
	gasDenom := "ario"

	suite.msKeeper.SetValidatorAllowedCoin(suite.ctx, valAddr, gasDenom)
	imAccBalance := sdk.NewCoins(sdk.NewCoin(stakingtypes.DefaultParams().BondDenom, sdk.NewInt(10000)), sdk.NewCoin(gasDenom, sdk.NewInt(10000)))
	err := suite.app.BankKeeper.MintCoins(suite.ctx, minttypes.ModuleName, imAccBalance)
	suite.NoError(err)
	err = suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, minttypes.ModuleName, intDel, imAccBalance)
	suite.NoError(err)

	rate, err := sdk.NewDecFromStr("1")
	suite.NoError(err)

	unlockEntry := types.UnlockEntry{
		CreationHeight: suite.ctx.BlockHeight(),
		UnlockingCoin: types.MultiStakingCoin{
			Denom:      gasDenom,
			Amount:     sdk.NewInt(1000),
			BondWeight: rate,
		},
	}
	newUbd := types.MultiStakingUnlock{
		UnlockID: &types.UnlockID{
			MultiStakerAddr: multiStakerAddress.String(),
			ValAddr:         valAddr.String(),
		},
		Entries: []types.UnlockEntry{unlockEntry},
	}

	suite.msKeeper.SetMultiStakingUnlock(suite.ctx, newUbd)
	matureUnbondingDelegations := suite.msKeeper.GetMatureUnbondingDelegations(suite.ctx)
	completionTime := suite.ctx.BlockHeader().Time

	unDT := stakingtypes.UnbondingDelegationEntry{
		CreationHeight: suite.ctx.BlockHeight(),
		CompletionTime: completionTime,
		InitialBalance: sdk.NewInt(1000),
		Balance:        sdk.NewInt(1000),
	}

	unD := stakingtypes.UnbondingDelegation{
		DelegatorAddress: intDel.String(),
		ValidatorAddress: valAddr.String(),
		Entries:          []stakingtypes.UnbondingDelegationEntry{unDT},
	}

	matureUnbondingDelegations = append(matureUnbondingDelegations, unD)

	suite.msKeeper.EndBlocker(suite.ctx, matureUnbondingDelegations)

}
