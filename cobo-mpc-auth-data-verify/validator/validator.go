package validator

import (
	"fmt"
)

//all data
// 	result    int
// 	pubkey    string
// 	signature string
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
	Message   string `json:"message"`
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
	if v.authData == nil {
		return fmt.Errorf("auth data is nil")
	}

	statement := NewStatementBuilder(v.authData.Template)

	// step 1: build statement message from biz data and template
	buildMsg, err := statement.Build(v.authData.BizData)
	if err != nil {
		return fmt.Errorf("error building statement: %w", err)
	}

	// step 2: verify statement message and build message are equal
	equal, diff := CompareStatementMessage(buildMsg, v.authData.Message)
	if !equal {
		return fmt.Errorf("source message and build message are not equal: %s", diff)
	}

	// step 3: verify signature of message and result
	sv := NewSignatureValidator(v.authData.Message, v.authData.Pubkey, v.authData.Signature, v.authData.Result)
	err = sv.Verify()
	if err != nil {
		return fmt.Errorf("error verifying message: %w", err)
	}

	// step 3: verify result is approved
	if v.authData.Result != 2 {
		return fmt.Errorf("result is not approved(2): %d", v.authData.Result)
	}

	// step 4: verify biz data and tss callback data are match

	return nil
}
