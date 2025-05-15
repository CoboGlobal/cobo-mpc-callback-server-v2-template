package types

import (
	"github.com/golang-jwt/jwt/v5"
)

const (
	StatusOK             = 0
	StatusInvalidRequest = 10
	StatusInvalidToken   = 20
	StatusInternalError  = 30
)

type PackageDataClaim struct {
	PackageData []byte `json:"package_data,omitempty"`
	jwt.RegisteredClaims
}
