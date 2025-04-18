package netservice

import (
	"github.com/golang-jwt/jwt/v5"
)

type PackageDataClaim struct {
	PackageData []byte `json:"package_data,omitempty"`
	jwt.RegisteredClaims
}
