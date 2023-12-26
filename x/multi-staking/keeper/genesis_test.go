package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/realio-tech/multi-staking-module/testutil"
	"github.com/realio-tech/multi-staking-module/x/multi-staking/types"
)

func (suite *KeeperTestSuite) TestInitGenesis() {
	suite.SetupTest()
	multiStakerAddress := testutil.GenAddress()
	valPubKey := testutil.GenPubKey()
	valAddr := sdk.ValAddress(valPubKey.Address())
	denom := "ario"

	multiStakingLock := types.MultiStakingLock{
		LockID: &types.LockID{
			MultiStakerAddr: multiStakerAddress.String(),
			ValAddr:         valAddr.String(),
		},
		LockedCoin: types.MultiStakingCoin{
			Denom:      denom,
			Amount:     sdk.NewInt(1000),
			BondWeight: sdk.NewDec(1),
		},
	}
	validatorAllowedCoin := types.ValidatorAllowedCoin{
		ValAddr:    valAddr.String(),
		TokenDenom: denom,
	}

	var delegations []stakingtypes.Delegation
	genesisDelegations := suite.stakingKeeper.GetAllDelegations(suite.ctx)
	delegations = append(delegations, genesisDelegations...)

	validators := suite.stakingKeeper.GetAllValidators(suite.ctx)

	params := suite.stakingKeeper.GetParams(suite.ctx)

	stakingGenesisState := stakingtypes.NewGenesisState(params, validators, delegations)

	expectedGenesisState := types.GenesisState{
		MultiStakingLocks:     []types.MultiStakingLock{multiStakingLock},
		ValidatorAllowedToken: []types.ValidatorAllowedCoin{validatorAllowedCoin},
		StakingGenesisState:   stakingGenesisState,
	}

	suite.msKeeper.InitGenesis(suite.ctx, expectedGenesisState)

	actualGenesisState := suite.msKeeper.ExportGenesis(suite.ctx)
	suite.Require().NotNil(actualGenesisState)
	suite.Require().Equal(expectedGenesisState.ValidatorAllowedToken, actualGenesisState.ValidatorAllowedToken)
	suite.Require().Equal(expectedGenesisState.StakingGenesisState.Delegations, actualGenesisState.StakingGenesisState.Delegations)
	suite.Require().Equal(expectedGenesisState.StakingGenesisState.Validators, actualGenesisState.StakingGenesisState.Validators)
}
