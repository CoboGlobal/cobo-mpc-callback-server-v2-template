package verifier

import (
	"encoding/json"
	"fmt"

	"github.com/CoboGlobal/cobo-mpc-callback-server-v2/pkg/log"
	coboWaaS2 "github.com/CoboGlobal/cobo-waas2-go-sdk/cobo_waas2"
)

type Verifier interface {
	Verify(request *coboWaaS2.TSSCallbackRequest) error
}

type TssVerifier struct {
	addressWhitelist []string
}

func NewTssVerifier(addressWhitelist []string) Verifier {
	return &TssVerifier{
		addressWhitelist: addressWhitelist,
	}
}

func (v *TssVerifier) Verify(request *coboWaaS2.TSSCallbackRequest) error {
	if request == nil || request.RequestType == nil {
		return fmt.Errorf("request or request type is nil")
	}

	switch *request.RequestType {
	case coboWaaS2.TSSCALLBACKREQUESTTYPE__0:
		log.Debugf("Got ping request")
		return nil
	case coboWaaS2.TSSCALLBACKREQUESTTYPE__1:
		return v.handleKeyGen(*request.RequestDetail, *request.ExtraInfo)
	case coboWaaS2.TSSCALLBACKREQUESTTYPE__2:
		return v.handleKeySign(*request.RequestDetail, *request.ExtraInfo)
	case coboWaaS2.TSSCALLBACKREQUESTTYPE__3:
		return v.handleKeyReshare(*request.RequestDetail, *request.ExtraInfo)

	default:
		return fmt.Errorf("not support to process request type %v", *request.RequestType)
	}

}

func (v *TssVerifier) handleKeyGen(requestDetail, extraInfo string) error {
	if requestDetail == "" || extraInfo == "" {
		return fmt.Errorf("request detail or extra info is empty")
	}

	var detail coboWaaS2.TSSKeyGenRequest
	if err := json.Unmarshal([]byte(requestDetail), &detail); err != nil {
		return fmt.Errorf("failed to parse key gen detail: %w", err)
	}

	var extra coboWaaS2.TSSKeyGenExtra
	if err := json.Unmarshal([]byte(extraInfo), &extra); err != nil {
		return fmt.Errorf("failed to parse key gen extra: %w", err)
	}

	log.Debugf("key gen detail:\n%v\n extra:\n%v", detail, extra)

	// key gen logic add here

	return nil
}

func (v *TssVerifier) handleKeySign(requestDetail, extraInfo string) error {
	if requestDetail == "" || extraInfo == "" {
		return fmt.Errorf("request detail or extra info is empty")
	}

	var detail coboWaaS2.TSSKeySignRequest
	if err := json.Unmarshal([]byte(requestDetail), &detail); err != nil {
		return fmt.Errorf("failed to parse key sign detail: %w", err)
	}

	var extra coboWaaS2.TSSKeySignExtra
	if err := json.Unmarshal([]byte(extraInfo), &extra); err != nil {
		return fmt.Errorf("failed to parse key sign extra: %w", err)
	}

	log.Debugf("key sign detail:\n%v\n extra:\n%v", detail, extra)

	// key sign logic add here

	//verify sign for example
	if err := v.verifySign(&detail, &extra); err != nil {
		return fmt.Errorf("verify sign error: %w", err)
	}

	return nil
}

func (v *TssVerifier) handleKeyReshare(requestDetail, extraInfo string) error {
	if requestDetail == "" || extraInfo == "" {
		return fmt.Errorf("request detail or extra info is empty")
	}

	var detail coboWaaS2.TSSKeyReshareRequest
	if err := json.Unmarshal([]byte(requestDetail), &detail); err != nil {
		return fmt.Errorf("failed to parse key reshare detail: %w", err)
	}

	var extra coboWaaS2.TSSKeyReshareExtra
	if err := json.Unmarshal([]byte(extraInfo), &extra); err != nil {
		return fmt.Errorf("failed to parse key reshare extra: %w", err)
	}

	log.Debugf("key reshare detail:\n%v\n extra:\n%v", detail, extra)

	// key reshare logic add here

	return nil
}
