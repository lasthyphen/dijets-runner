package api

import (
	"github.com/lasthyphen/dijetsgo/api/admin"
	"github.com/lasthyphen/dijetsgo/api/health"
	"github.com/lasthyphen/dijetsgo/api/info"
	"github.com/lasthyphen/dijetsgo/api/ipcs"
	"github.com/lasthyphen/dijetsgo/api/keystore"
	"github.com/lasthyphen/dijetsgo/indexer"
	"github.com/lasthyphen/dijetsgo/vms/avm"
	"github.com/lasthyphen/dijetsgo/vms/platformvm"
	"github.com/lasthyphen/dijeth/plugin/evm"
)

// Issues API calls to a node
// TODO: byzantine api. check if appropiate. improve implementation.
type Client interface {
	PChainAPI() platformvm.Client
	XChainAPI() avm.Client
	XChainWalletAPI() avm.WalletClient
	CChainAPI() evm.Client
	CChainEthAPI() EthClient // ethclient websocket wrapper that adds mutexed calls, and lazy conn init (on first call)
	InfoAPI() info.Client
	HealthAPI() health.Client
	IpcsAPI() ipcs.Client
	KeystoreAPI() keystore.Client
	AdminAPI() admin.Client
	PChainIndexAPI() indexer.Client
	CChainIndexAPI() indexer.Client
	// TODO add methods
}
