package dto

type struct Subversion {
	HN_BTC_VERSION string `json:"hn_btc_version"`
	HN_ETH_VERSION string `json:"hn_eth_version"`
	HN_QTUM_VERSION string `json:"hn_qtum_version"`
}

type struct PBlock {
	HN_BTC_HEIGHT int64 `json:"hn_btc_height"`
	HN_ETH_HEIGHT int64 `json:"hn_eth_height"`
	HN_QTUM_HEIGHT int64 `json:"hn_qtum_height"`
}

type struct Verification {
	HN_ID string `json:"hn_id"`
	HN_VERSION string `json:"hn_version"`
	HN_SUBVERSION subversion `json:"hn_subversion"`
	HN_SERVICE string `json:"hn_service"`
	HN_TIME int64 `json:"hn_time"`
	HN_THIS_ADDR string `json:"hn_addr"`
	HN_BLOCK_HEIGHT PBlock `json:"hn_height"`
}