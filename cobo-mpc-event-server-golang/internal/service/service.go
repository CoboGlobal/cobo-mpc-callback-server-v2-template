package service

import (
	"github.com/CoboGlobal/cobo-mpc-event-server/internal/config"
	netService "github.com/CoboGlobal/cobo-mpc-event-server/internal/netservice"
	"github.com/CoboGlobal/cobo-mpc-event-server/pkg/log"
)

type EventService interface {
	Start()
	HandleEvent(rawEvent []byte) error
}

type Service struct {
	eventSrv *netService.Service
}

func New(cfg *config.Config) *Service {
	s := &Service{}
	s.eventSrv = netService.New(cfg.EventServer, s.HandleEvent)
	return s
}

func (s *Service) Start() {
	s.eventSrv.Start()
}

func (s *Service) HandleEvent(rawEvent []byte) error {

	log.Debugf("Get event: %s", string(rawEvent))
	// TODO: handle event
	return nil
}
