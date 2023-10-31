package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

func GetIntermediaryAccount(delAddr string, valAddr string) sdk.AccAddress {
	// TODO: Make this better namespaced.
	// Following Osmosis Superfluid in the future to resolve this comment
	return authtypes.NewModuleAddress(delAddr + valAddr)
}
