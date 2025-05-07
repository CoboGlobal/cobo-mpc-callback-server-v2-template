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

type AuthValidator interface {
	Verify() error
}

type Validator struct {
	result          int
	pubkey          string
	signature       string
	algorithm       string
	templateContent string
	bizData         string
}

func NewValidator(result int, pubkey string, signature string, algorithm string, templateContent string, bizData string) *Verifier {
	return &Validator{
		result:          result,
		pubkey:          pubkey,
		signature:       signature,
		algorithm:       algorithm,
		templateContent: templateContent,
		bizData:         bizData,
	}
}

func (v *Validator) Verify() error {
	statement := NewStatement(v.templateContent)

	// step 1: build statement from biz data and template content
	message, err := statement.BuildStatementV2(v.bizData)
	if err != nil {
		return fmt.Errorf("error building statement: %w", err)
	}

	// step 2: verify signature of message and result
	validator := NewValidator(message, v.pubkey, v.signature, v.algorithm, v.result)
	err = validator.Verify()
	if err != nil {
		return fmt.Errorf("error verifying message: %w", err)
	}

	return nil
}
