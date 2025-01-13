package service

import (
	"encoding/json"
	"fmt"

	"github.com/CoboGlobal/cobo-mpc-callback-server-v2/internal/config"
	netService "github.com/CoboGlobal/cobo-mpc-callback-server-v2/internal/netservice"
	"github.com/CoboGlobal/cobo-mpc-callback-server-v2/internal/types"
	"github.com/CoboGlobal/cobo-mpc-callback-server-v2/internal/verifier"
)

type CallbackService interface {
	Start()
	HandleRequest(rawRequest []byte) (*types.Response, error)
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

func (s *Service) HandleRequest(rawRequest []byte) (*types.Response, error) {
	req := &types.Request{}
	if err := json.Unmarshal(rawRequest, req); err != nil {
		return &types.Response{
			Status:    types.StatusInternalError,
			ErrStr:    fmt.Sprintf("failed to parse raw request %v", err.Error()),
			RequestID: req.RequestID,
		}, nil
	}

	if err := s.vfr.Verify(req); err != nil {
		return &types.Response{
			Status:    types.StatusInternalError,
			ErrStr:    fmt.Sprintf("reject sign request: %v", err),
			RequestID: req.RequestID,
		}, err
	}

	return &types.Response{
		Action:    types.ActionApprove,
		Status:    types.StatusOK,
		RequestID: req.RequestID,
	}, nil
}
