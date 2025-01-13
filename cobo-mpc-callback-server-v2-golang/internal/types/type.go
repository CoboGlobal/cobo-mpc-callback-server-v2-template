package types

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
)

const (
	ActionApprove = "APPROVE"
	ActionReject  = "REJECT"

	StatusOK             = 0
	StatusInvalidRequest = 10
	StatusInvalidToken   = 20
	StatusInternalError  = 30
)

type PackageDataClaim struct {
	PackageData []byte `json:"package_data,omitempty"`
	jwt.RegisteredClaims
}

type Request struct {
	RequestID     string `json:"request_id,omitempty"`
	RequestType   int    `json:"request_type,omitempty"`
	RequestDetail string `json:"request_detail,omitempty"`
	ExtraInfo     string `json:"extra_info,omitempty"`
}

type Response struct {
	Status    int    `json:"status,omitempty"`
	RequestID string `json:"request_id,omitempty"`
	Action    string `json:"action,omitempty"` //[APPROVE, REJECT]
	ErrStr    string `json:"error,omitempty"`
}

func (r *Response) String() string {
	return fmt.Sprintf("Status: %d, RequestID: %s, Action: %s, ErrStr: %s", r.Status, r.RequestID, r.Action, r.ErrStr)
}
