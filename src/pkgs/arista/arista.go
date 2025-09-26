package arista

// Top-level JSON-RPC response
type BGPEvpnSummaryResponse struct {
	JSONRPC string                 `json:"jsonrpc"`
	ID      int                    `json:"id"`
	Result  []BGPEvpnSummaryResult `json:"result"`
}

// Each element in "result" has a "vrfs" map (keys like "default")
type BGPEvpnSummaryResult struct {
	Vrfs map[string]VRF `json:"vrfs"`
}

// A VRF contains router info and a map of BGP peers keyed by neighbor IP
type VRF struct {
	VRF      string          `json:"vrf"`
	RouterID string          `json:"routerId"`
	ASN      string          `json:"asn"`
	Peers    map[string]Peer `json:"peers"`
}

// Peer is one neighbor entry as reported by EOS
type Peer struct {
	Version          int     `json:"version"`
	MsgReceived      int     `json:"msgReceived"`
	MsgSent          int     `json:"msgSent"`
	InMsgQueue       int     `json:"inMsgQueue"`
	OutMsgQueue      int     `json:"outMsgQueue"`
	ASN              string  `json:"asn"`
	PrefixAccepted   int     `json:"prefixAccepted"`
	PrefixReceived   int     `json:"prefixReceived"`
	UpDownTime       float64 `json:"upDownTime"`
	UnderMaintenance bool    `json:"underMaintenance"`
	PeerState        string  `json:"peerState"`
	PrefixAdvertised int     `json:"prefixAdvertised"`
}
