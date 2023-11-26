package keeper

import (
	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/realio-tech/multi-staking-module/x/multi-staking/types"
)

func (k Keeper) LockMultiStakingTokenAndMintBondToken(
	ctx sdk.Context, delAcc sdk.AccAddress, valAcc sdk.ValAddress,
	multiStakingToken sdk.Coin,
) (mintedBondToken sdk.Coin, err error) {
	intermediaryAcc := k.GetIntermediaryAccountDelegator(ctx, delAcc)

	// get bond denom weight
	bondDenomWeight, isBondToken := k.GetBondTokenWeight(ctx, multiStakingToken.Denom)
	if !isBondToken {
		return sdk.Coin{}, errors.Wrapf(
			sdkerrors.ErrInvalidRequest, "invalid coin denomination: got %s", multiStakingToken.Denom,
		)
	}

	// lock coin in intermediary account
	err = k.bankKeeper.SendCoins(ctx, delAcc, intermediaryAcc, sdk.NewCoins(multiStakingToken))
	if err != nil {
		return sdk.Coin{}, err
	}

	// update multistaking lock
	multiStakingLockKey := types.MultiStakingLockID(delAcc, valAcc)
	multiStakingLock, found := k.GetMultiStakingLock(ctx, multiStakingLockKey)
	if !found {
		multiStakingLock = types.NewMultiStakingLock(multiStakingToken.Amount, multiStakingLock.ConversionRatio, intermediaryAcc.String())
	} else {
		multiStakingLock = multiStakingLock.AddTokenToMultiStakingLock(multiStakingToken.Amount, bondDenomWeight)
	}
	k.SetMultiStakingLock(ctx, multiStakingLockKey, multiStakingLock)

	// Calculate the amount of bond denom to be minted
	// minted bond amount = multistaking token * bond token weight
	mintedBondAmount := bondDenomWeight.MulInt(multiStakingToken.Amount).RoundInt()
	mintedBondToken = sdk.NewCoin(k.stakingKeeper.BondDenom(ctx), mintedBondAmount)

	// mint bond token to intermediary account
	k.bankKeeper.MintCoins(ctx, types.ModuleName, sdk.NewCoins(mintedBondToken))
	k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, intermediaryAcc, sdk.NewCoins(mintedBondToken))

	return mintedBondToken, nil
}