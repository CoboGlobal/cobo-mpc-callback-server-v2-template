package verifier

import (
	"encoding/json"
	"fmt"
	"github.com/CoboGlobal/cobo-mpc-callback-server-v2/internal/types"
	"github.com/CoboGlobal/cobo-mpc-callback-server-v2/pkg/log"
)

type Verifier interface {
	Verify(request *types.Request) error
}

type TssVerifier struct {
	addressWhitelist []string
}

func NewTssVerifier(addressWhitelist []string) Verifier {
	return &TssVerifier{
		addressWhitelist: addressWhitelist,
	}
}

func (v *TssVerifier) Verify(request *types.Request) error {
	if request == nil {
		return fmt.Errorf("request is nil")
	}

	switch types.RequestType(request.RequestType) {
	case types.TypePing:
		log.Debugf("Got ping request")
		return nil
	case types.TypeKeyGen:
		return v.handleKeyGen(request.RequestDetail, request.ExtraInfo)
	case types.TypeKeySign:
		return v.handleKeySign(request.RequestDetail, request.ExtraInfo)
	case types.TypeKeyReshare:
		return v.handleKeyReshare(request.RequestDetail, request.ExtraInfo)

	default:
		return fmt.Errorf("not support to process request type %v", request.RequestType)
	}

}

func (v *TssVerifier) handleKeyGen(requestDetail, extraInfo string) error {
	if requestDetail == "" || extraInfo == "" {
		return fmt.Errorf("request detail or extra info is empty")
	}

	var detail types.KeyGenDetail
	if err := json.Unmarshal([]byte(requestDetail), &detail); err != nil {
		return fmt.Errorf("failed to parse key gen detail: %w", err)
	}

	var requestInfo types.KeyGenRequestInfo
	if err := json.Unmarshal([]byte(extraInfo), &requestInfo); err != nil {
		return fmt.Errorf("failed to parse key gen request info: %w", err)
	}

	log.Debugf("key gen detail:\n%v\nrequest info:\n%v", detail, requestInfo.String())

	// key gen logic add here

	return nil
}

func (v *TssVerifier) handleKeySign(requestDetail, extraInfo string) error {
	if requestDetail == "" || extraInfo == "" {
		return fmt.Errorf("request detail or extra info is empty")
	}

	var detail types.KeySignDetail
	if err := json.Unmarshal([]byte(requestDetail), &detail); err != nil {
		return fmt.Errorf("failed to parse key sign detail: %w", err)
	}

	var requestInfo types.KeySignRequestInfo
	if err := json.Unmarshal([]byte(extraInfo), &requestInfo); err != nil {
		return fmt.Errorf("failed to parse key sign request info: %w", err)
	}

	log.Debugf("key sign detail:\n%v\nrequest info:\n%v", detail, requestInfo.String())

	// key sign logic add here

	//verify sign for example
	if err := v.verifySign(&detail, &requestInfo); err != nil {
		return fmt.Errorf("verify sign error: %w", err)
	}

	return nil
}

func (v *TssVerifier) handleKeyReshare(requestDetail, extraInfo string) error {
	if requestDetail == "" || extraInfo == "" {
		return fmt.Errorf("request detail or extra info is empty")
	}

	var detail types.KeyReshareDetail
	if err := json.Unmarshal([]byte(requestDetail), &detail); err != nil {
		return fmt.Errorf("failed to parse key reshare detail: %w", err)
	}

	var requestInfo types.KeyReshareRequestInfo
	if err := json.Unmarshal([]byte(extraInfo), &requestInfo); err != nil {
		return fmt.Errorf("failed to parse key reshare request info: %w", err)
	}

	log.Debugf("key reshare detail:\n%v\nrequest info:\n%v", detail, requestInfo.String())

	// key reshare logic add here

	return nil
}
