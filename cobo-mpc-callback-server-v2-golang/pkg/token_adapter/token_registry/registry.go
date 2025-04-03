package token_registry

import (
	"github.com/CoboGlobal/cobo-mpc-callback-server-v2/pkg/token_adapter"
	"github.com/CoboGlobal/cobo-mpc-callback-server-v2/pkg/token_adapter/eth_base"
	"github.com/CoboGlobal/cobo-mpc-callback-server-v2/pkg/token_adapter/solana"
	"github.com/CoboGlobal/cobo-mpc-callback-server-v2/pkg/token_adapter/tron"
)

func InitRegistry() {
	if err := token_adapter.RegisterTokenCreator("ETH", eth_base.NewToken); err != nil {
		panic(err)
	}

	if err := token_adapter.RegisterTokenCreator("ETH_USDT", eth_base.NewErc20Token); err != nil {
		panic(err)
	}

	if err := token_adapter.RegisterTokenCreator("TRON", tron.NewToken); err != nil {
		panic(err)
	}

	if err := token_adapter.RegisterTokenCreator("TRON_USDT", tron.NewTrc20Token); err != nil {
		panic(err)
	}

	if err := token_adapter.RegisterTokenCreator("SOL", solana.NewToken); err != nil {
		panic(err)
	}

	if err := token_adapter.RegisterTokenCreator("SOL_USDC", solana.NewSPLToken); err != nil {
		panic(err)
	}
}
