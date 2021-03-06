// This software is Copyright (c) 2019 e-Money A/S. It is not offered under an open source license.
//
// Please contact partners@e-money.com for licensing related questions.

package keeper

import (
	"fmt"

	"github.com/e-money/em-ledger/x/issuer/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

func NewQuerier(k Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case types.QueryIssuers:
			return listIssuers(ctx, k)
		default:
			return nil, sdk.ErrUnknownRequest(fmt.Sprintf("unknown issuer query endpoint: %s", path[0]))
		}

		return []byte{}, nil
	}
}

func listIssuers(ctx sdk.Context, k Keeper) ([]byte, sdk.Error) {
	issuers := k.GetIssuers(ctx)
	res, err := types.ModuleCdc.MarshalJSON(issuers)
	if err != nil {
		return nil, sdk.ErrInternal(sdk.AppendMsgToErr("failed to marshal JSON", err.Error()))
	}

	return res, nil
}
