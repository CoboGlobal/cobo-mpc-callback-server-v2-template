package types

import (
	"encoding/json"
	coboWaaS2 "github.com/CoboGlobal/cobo-waas2-go-sdk/cobo_waas2"
)

// =================== request_type ===================.

type RequestType int

const (
	TypePing       RequestType = 0
	TypeKeyGen     RequestType = 1
	TypeKeySign    RequestType = 2
	TypeKeyReshare RequestType = 3
)

func (t RequestType) String() string {
	switch t {
	case TypePing:
		return "Ping"
	case TypeKeyGen:
		return "KeyGen"
	case TypeKeySign:
		return "KeySign"
	case TypeKeyReshare:
		return "KeyReshare"
	default:
		return "Unknown"
	}
}

// =================== request_detail ===================.

type CurveType int

const (
	SECP256K1 CurveType = 0
	ED25519   CurveType = 2
)

type SignatureType int32

const (
	UNKNOWN_TYPE SignatureType = 0
	Ecdsa        SignatureType = 1
	Eddsa        SignatureType = 2
	Schnorr      SignatureType = 3
)

type TssProtocol int32

const (
	UNKNOWN_PROTOCOL TssProtocol = 0
	GG18             TssProtocol = 1
	Lindell          TssProtocol = 2
	EddsaTSS         TssProtocol = 3
)

type KeyGenDetail struct {
	Threshold int       `json:"threshold"`
	Curve     CurveType `json:"curve"`
	NodeIDs   []string  `json:"node_ids"`
	TaskID    string    `json:"task_id"`
	BizTaskID string    `json:"biz_task_id"`
}

type KeySignDetail struct {
	GroupID       string        `json:"group_id"`
	RootPubKey    string        `json:"root_pub_key"`
	UsedNodeIDs   []string      `json:"used_node_ids"`
	Bip32PathList []string      `json:"bip32_path_list"`
	MsgHashList   []string      `json:"msg_hash_list"`
	TweakList     []string      `json:"tweak_list"`
	SignatureType SignatureType `json:"signature_type"`
	TssProtocol   TssProtocol   `json:"tss_protocol"`
	TaskID        string        `json:"task_id"`
	BizTaskID     string        `json:"biz_task_id"`
}

type KeyReshareDetail struct {
	OldGroupID   string    `json:"old_group_id"`
	RootPubKey   string    `json:"root_pub_key"`
	Curve        CurveType `json:"curve"`
	UsedNodeIDs  []string  `json:"used_node_ids"`
	OldThreshold int       `json:"old_threshold"`
	NewThreshold int       `json:"new_threshold"`
	NewNodeIDs   []string  `json:"new_node_ids"`
	TaskID       string    `json:"task_id"`
	BizTaskID    string    `json:"biz_task_id"`
}

// =================== extra_info ===================.

type KeyGenRequestInfo struct {
	Org                       *coboWaaS2.OrgInfo             `json:"org,omitempty"`
	Project                   *coboWaaS2.MPCProject          `json:"project,omitempty"`
	Vault                     *coboWaaS2.MPCVault            `json:"vault,omitempty"`
	TargetKeyShareHolderGroup *coboWaaS2.KeyShareHolderGroup `json:"target_key_share_holder_group,omitempty"`
	TSSRequest                *coboWaaS2.TSSRequest          `json:"tss_request,omitempty"`
}

func (r *KeyGenRequestInfo) String() string {
	if r == nil {
		return ""
	}
	info, _ := json.Marshal(r)
	return string(info)
}

type KeySignRequestInfo struct {
	Org                       *coboWaaS2.OrgInfo             `json:"org,omitempty"`
	Project                   *coboWaaS2.MPCProject          `json:"project,omitempty"`
	Vault                     *coboWaaS2.MPCVault            `json:"vault,omitempty"`
	Wallet                    *coboWaaS2.WalletInfo          `json:"wallet,omitempty"`
	SignerKeyShareHolderGroup *coboWaaS2.KeyShareHolderGroup `json:"signer_key_share_holder_group,omitempty"`
	SourceAddresses           []*coboWaaS2.AddressInfo       `json:"source_addresses,omitempty"`
	Transaction               *coboWaaS2.Transaction         `json:"transaction,omitempty"`
	StakingActivity           *coboWaaS2.Activity            `json:"staking_activity,omitempty"`
}

func (r *KeySignRequestInfo) String() string {
	if r == nil {
		return ""
	}
	info, _ := json.Marshal(r)
	return string(info)
}

type KeyReshareRequestInfo struct {
	Org                       *coboWaaS2.OrgInfo             `json:"org,omitempty"`
	Project                   *coboWaaS2.MPCProject          `json:"project,omitempty"`
	Vault                     *coboWaaS2.MPCVault            `json:"vault,omitempty"`
	SourceKeyShareHolderGroup *coboWaaS2.KeyShareHolderGroup `json:"source_key_share_holder_group,omitempty"`
	TargetKeyShareHolderGroup *coboWaaS2.KeyShareHolderGroup `json:"target_key_share_holder_group,omitempty"`
	TSSRequest                *coboWaaS2.TSSRequest          `json:"tss_request,omitempty"`
}

func (r *KeyReshareRequestInfo) String() string {
	if r == nil {
		return ""
	}
	info, _ := json.Marshal(r)
	return string(info)
}
