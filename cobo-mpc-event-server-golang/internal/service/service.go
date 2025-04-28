package service

import (
	"encoding/json"

	"github.com/CoboGlobal/cobo-mpc-event-server/internal/config"
	netService "github.com/CoboGlobal/cobo-mpc-event-server/internal/netservice"
	"github.com/CoboGlobal/cobo-mpc-event-server/pkg/log"
	coboWaaS2 "github.com/CoboGlobal/cobo-waas2-go-sdk/cobo_waas2"
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
	//
	var event coboWaaS2.TSSEvent
	err := event.UnmarshalJSON(rawEvent)
	if err != nil {
		log.Errorf("Unmarshal event failed: %v", err)
		return err
	}

	jsonEvent, err := event.MarshalJSON()
	if err != nil {
		log.Errorf("Marshal event failed: %v", err)
		return err
	}
	log.Debugf("Event: %+v", string(jsonEvent))

	eventData := event.Data.GetActualInstance()

	if data, ok := eventData.(*coboWaaS2.TSSKeyGenEventData); ok {

		if data.ExtraInfo != nil && *data.ExtraInfo != "" {
			var extraInfo coboWaaS2.TSSKeyGenExtra
			err = json.Unmarshal([]byte(*data.ExtraInfo), &extraInfo)
			if err != nil {
				log.Errorf("Unmarshal extra info failed: %v", err)
				return err
			}

			extraInfoJSON, err := extraInfo.MarshalJSON()
			if err != nil {
				log.Errorf("Marshal extra info failed: %v", err)
				return err
			}
			log.Debugf("ExtraInfo: %+v", string(extraInfoJSON))
		}

	} else if data, ok := eventData.(*coboWaaS2.TSSKeyReshareEventData); ok {

		if data.ExtraInfo != nil && *data.ExtraInfo != "" {
			var extraInfo coboWaaS2.TSSKeyReshareExtra
			err = json.Unmarshal([]byte(*data.ExtraInfo), &extraInfo)
			if err != nil {
				log.Errorf("Unmarshal extra info failed: %v", err)
				return err
			}

			extraInfoJSON, err := extraInfo.MarshalJSON()
			if err != nil {
				log.Errorf("Marshal extra info failed: %v", err)
				return err
			}
			log.Debugf("ExtraInfo: %+v", string(extraInfoJSON))
		}

	} else if data, ok := eventData.(*coboWaaS2.TSSKeySignEventData); ok {

		if data.ExtraInfo != nil && *data.ExtraInfo != "" {
			var extraInfo coboWaaS2.TSSKeySignExtra
			err = json.Unmarshal([]byte(*data.ExtraInfo), &extraInfo)
			if err != nil {
				log.Errorf("Unmarshal extra info failed: %v", err)
				return err
			}

			extraInfoJSON, err := extraInfo.MarshalJSON()
			if err != nil {
				log.Errorf("Marshal extra info failed: %v", err)
				return err
			}
			log.Debugf("ExtraInfo: %+v", string(extraInfoJSON))
		}

	} else if data, ok := eventData.(*coboWaaS2.TSSKeyShareSignEventData); ok {

		if data.ExtraInfo != nil && *data.ExtraInfo != "" {
			var extraInfo coboWaaS2.TSSKeyShareSignExtra
			err = json.Unmarshal([]byte(*data.ExtraInfo), &extraInfo)
			if err != nil {
				log.Errorf("Unmarshal extra info failed: %v", err)
				return err
			}

			extraInfoJSON, err := extraInfo.MarshalJSON()
			if err != nil {
				log.Errorf("Marshal extra info failed: %v", err)
				return err
			}
			log.Debugf("ExtraInfo: %+v", string(extraInfoJSON))
		}

	}

	return nil
}
