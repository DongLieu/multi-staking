package types_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/require"

	"github.com/realio-tech/multi-staking-module/testutil"
	"github.com/realio-tech/multi-staking-module/x/multi-staking/types"
)

func TestMsgBeginRedelegate_ValidateBasic(t *testing.T) {
	mulStakeAddr := testutil.GenAddress()
	valSrcAddr := testutil.GenValAddress()
	valDstAddr := testutil.GenValAddress()
	denom := "ario"
	coin := sdk.NewCoin(denom, sdk.NewInt(10000))

	tests := []struct {
		name string
		msg  types.MsgBeginRedelegate
		err  error
	}{
		{
			name: "happy path",
			msg: types.MsgBeginRedelegate{
				MultiStakerAddress:  mulStakeAddr.String(),
				ValidatorSrcAddress: valSrcAddr.String(),
				ValidatorDstAddress: valDstAddr.String(),
				Amount:              coin,
			},
		},
		{
			name: "invalid address",
			msg: types.MsgBeginRedelegate{
				MultiStakerAddress:  "",
				ValidatorSrcAddress: "",
				ValidatorDstAddress: "",
				Amount:              coin,
			},
			err: sdkerrors.ErrInvalidAddress,
		},
		{
			name: "invalid shares amount",
			msg: types.MsgBeginRedelegate{
				MultiStakerAddress:  mulStakeAddr.String(),
				ValidatorSrcAddress: valSrcAddr.String(),
				ValidatorDstAddress: valDstAddr.String(),
				Amount:              sdk.NewCoin(denom, sdk.ZeroInt()),
			},
			err: sdkerrors.ErrInvalidRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.ValidateBasic()
			if tt.err != nil {
				require.ErrorIs(t, err, tt.err)
				return
			}
			require.NoError(t, err)
		})
	}
}
