package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func AccAddrAndValAddrFromStrings(accAddrString string, valAddrString string) (sdk.AccAddress, sdk.ValAddress, error) {
	accAddr, err := sdk.AccAddressFromBech32(accAddrString)
	if err != nil {
		return sdk.AccAddress{}, sdk.ValAddress{}, err
	}
	valAcc, err := sdk.ValAddressFromBech32(valAddrString)
	if err != nil {
		return sdk.AccAddress{}, sdk.ValAddress{}, err
	}

	return accAddr, valAcc, nil
}
