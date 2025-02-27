package v2_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/testutil"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"

	"github.com/evmos/ethermint/encoding"

	"github.com/evmos/evmos/v10/app"
	erc20types "github.com/evmos/evmos/v10/x/erc20/types"
	v2 "github.com/haqq-network/haqq/x/erc20/migrations/v2"
)

func TestUpdateParams(t *testing.T) {
	encCfg := encoding.MakeConfig(app.ModuleBasics)
	erc20Key := sdk.NewKVStoreKey(erc20types.StoreKey)
	tErc20Key := sdk.NewTransientStoreKey(fmt.Sprintf("%s_test", erc20types.StoreKey))
	ctx := testutil.DefaultContext(erc20Key, tErc20Key)
	paramstore := paramtypes.NewSubspace(
		encCfg.Codec, encCfg.Amino, erc20Key, tErc20Key, "erc20",
	)
	paramstore = paramstore.WithKeyTable(erc20types.ParamKeyTable())
	require.True(t, paramstore.HasKeyTable())

	// check no params
	require.False(t, paramstore.Has(ctx, erc20types.ParamStoreKeyEnableErc20))
	require.False(t, paramstore.Has(ctx, erc20types.ParamStoreKeyEnableEVMHook))

	// Run migrations
	err := v2.UpdateParams(ctx, &paramstore)
	require.NoError(t, err)

	// Make sure the params are set
	require.True(t, paramstore.Has(ctx, erc20types.ParamStoreKeyEnableErc20))
	require.True(t, paramstore.Has(ctx, erc20types.ParamStoreKeyEnableEVMHook))

	var enableERC20, enableEVMHook bool

	// Make sure the new params are set
	require.NotPanics(t, func() {
		paramstore.Get(ctx, erc20types.ParamStoreKeyEnableErc20, &enableERC20)
		paramstore.Get(ctx, erc20types.ParamStoreKeyEnableEVMHook, &enableEVMHook)
	})

	// check the params are updated
	require.True(t, enableERC20)
	require.True(t, enableEVMHook)
}
