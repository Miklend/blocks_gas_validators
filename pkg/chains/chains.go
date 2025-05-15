package chains

type ChainInfo struct {
	Name        string
	Network     string
	URL         string
	BlockTime   float64
	EtherscanId string
}

var AlchemyChains = map[string]ChainInfo{
	"ethereum": {
		Name:        "ethereum",
		Network:     "Mainnet",
		URL:         "://eth-mainnet.g.alchemy.com/v2/",
		BlockTime:   12.0,
		EtherscanId: "1",
	},
	"polygon": {
		Name:        "polygon",
		Network:     "Mainnet",
		URL:         "://polygon-mainnet.g.alchemy.com/v2/",
		BlockTime:   2.1,
		EtherscanId: "137",
	},
	"bnb": {
		Name:        "bnb",
		Network:     "Mainnet",
		URL:         "://bnb-mainnet.g.alchemy.com/v2/",
		BlockTime:   3.0,
		EtherscanId: "56",
	},
	"avalanche": {
		Name:        "avalanche",
		Network:     "Mainnet",
		URL:         "://avax-mainnet.g.alchemy.com/v2/",
		BlockTime:   2.0,
		EtherscanId: "43114",
	},
	"optimism": {
		Name:        "optimism",
		Network:     "Mainnet",
		URL:         "://opt-mainnet.g.alchemy.com/v2/",
		BlockTime:   2.0,
		EtherscanId: "10",
	},
	"base": {
		Name:        "base",
		Network:     "Mainnet",
		URL:         "://base-mainnet.g.alchemy.com/v2/",
		BlockTime:   2.0,
		EtherscanId: "2",
	},
}
