package service

import (
	"encoding/json"
	"fmt"

	"github.com/CoboGlobal/cobo-mpc-callback-server-v2/internal/config"
	netService "github.com/CoboGlobal/cobo-mpc-callback-server-v2/internal/netservice"
	"github.com/CoboGlobal/cobo-mpc-callback-server-v2/internal/types"
	"github.com/CoboGlobal/cobo-mpc-callback-server-v2/internal/verifier"
	coboWaaS2 "github.com/CoboGlobal/cobo-waas2-go-sdk/cobo_waas2"
)

type CallbackService interface {
	Start()
	HandleRequest(rawRequest []byte) (*coboWaaS2.TSSCallbackResponse, error)
}

type Service struct {
	callbackSrv *netService.Service

	vfr verifier.Verifier
}

func New(cfg *config.Config, vfr verifier.Verifier) *Service {
	s := &Service{
		vfr: vfr,
	}
	s.callbackSrv = netService.New(cfg.CallbackServer, s.HandleRequest)
	return s
}

func (s *Service) Start() {
	s.callbackSrv.Start()
}

func (s *Service) HandleRequest(rawRequest []byte) (*coboWaaS2.TSSCallbackResponse, error) {
	req := &coboWaaS2.TSSCallbackRequest{}
	if err := json.Unmarshal(rawRequest, req); err != nil {
		status := int32(types.StatusInternalError)
		errStr := fmt.Sprintf("failed to parse raw request %v", err.Error())
		return &coboWaaS2.TSSCallbackResponse{
			Status:    &status,
			Error:     &errStr,
			RequestId: req.RequestId,
		}, nil
	}

	//reqJSON, _ := req.MarshalJSON()
	//log.Debugf("Callback request: %v", string(reqJSON))
	if err := s.vfr.Verify(req); err != nil {
		status := int32(types.StatusInternalError)
		errStr := fmt.Sprintf("reject sign request: %v", err)
		return &coboWaaS2.TSSCallbackResponse{
			Status:    &status,
			Error:     &errStr,
			RequestId: req.RequestId,
		}, err
	}

	status := int32(types.StatusOK)
	action := coboWaaS2.TSSCALLBACKACTIONTYPE_APPROVE
	return &coboWaaS2.TSSCallbackResponse{
		Action:    &action,
		Status:    &status,
		RequestId: req.RequestId,
	}, nil
}
