package validator

import (
	"fmt"
)

type Validator interface {
	VerifyAuthData() error
	VerifyAuthDataAndResult() error
}

type AuthData struct {
	Result    int    `json:"result"`
	Pubkey    string `json:"pubkey"`
	Signature string `json:"signature"`
	Message   string `json:"message"`
	Template  string `json:"template"`
	BizData   string `json:"biz_data"`
}

type AuthValidator struct {
	authData *AuthData
}

func NewAuthValidator(authData *AuthData) *AuthValidator {
	return &AuthValidator{
		authData: authData,
	}
}

func (v *AuthValidator) VerifyAuthDataAndResult() error {
	err := v.VerifyAuthData()
	if err != nil {
		return fmt.Errorf("error verifying auth data: %w", err)
	}

	// step 4: verify result is approved
	if v.authData.Result != 2 {
		return fmt.Errorf("result is not approved(2): %d", v.authData.Result)
	}

	return nil
}

func (v *AuthValidator) VerifyAuthData() error {
	if v.authData == nil {
		return fmt.Errorf("auth data is nil")
	}

	statement := NewStatementBuilder(v.authData.Template)

	// step 1: build statement message from biz data and template
	buildMsg, err := statement.Build(v.authData.BizData)
	if err != nil {
		return fmt.Errorf("error building statement: %w", err)
	}

	originalMsg := v.authData.Message
	if originalMsg == "" {
		originalMsg = buildMsg
	} else {
		// step 2: verify statement message and build message are equal
		equal, diff := CompareStatementMessage(buildMsg, originalMsg)
		if !equal {
			return fmt.Errorf("source message and build message are not equal: %s", diff)
		}
	}

	// step 3: verify signature of message and result
	sv := NewSignatureValidator(originalMsg, v.authData.Pubkey, v.authData.Signature, v.authData.Result)
	err = sv.Verify()
	if err != nil {
		return fmt.Errorf("error verifying message: %w", err)
	}

	return nil
}
