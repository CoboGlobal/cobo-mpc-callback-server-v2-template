package validator

import (
	"fmt"
)

//all data
// 	result    int
// 	pubkey    string
// 	signature string
// 	algorithm string
//  templateContent string
//  bizData string
//  callbackData string

type Validator interface {
	Verify() error
}

type AuthData struct {
	Result    int    `json:"result"`
	Pubkey    string `json:"pubkey"`
	Signature string `json:"signature"`
	Algorithm string `json:"algorithm"`
	Template  string `json:"template"`
	BizData   string `json:"bizData"`
}

type AuthValidator struct {
	authData *AuthData
}

func NewAuthValidator(authData *AuthData) *AuthValidator {
	return &AuthValidator{
		authData: authData,
	}
}

func (v *AuthValidator) Verify() error {
	statement := NewStatement(v.authData.Template)

	// step 1: build statement from biz data and template content
	message, err := statement.BuildStatementV2(v.authData.BizData)
	if err != nil {
		return fmt.Errorf("error building statement: %w", err)
	}

	// step 2: verify signature of message and result
	sv := NewSignatureValidator(message, v.authData.Pubkey, v.authData.Signature, v.authData.Algorithm, v.authData.Result)
	err = sv.Verify()
	if err != nil {
		return fmt.Errorf("error verifying message: %w", err)
	}

	return nil
}
