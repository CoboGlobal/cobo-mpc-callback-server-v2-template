package netservice

import (
	"crypto/rsa"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/CoboGlobal/cobo-mpc-event-server/pkg/log"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type Config struct {
	ServiceName      string `mapstructure:"service_name"`
	Endpoint         string `mapstructure:"endpoint"`
	ClientPubKeyPath string `mapstructure:"client_public_key_path"`
	EnableDebug      bool   `mapstructure:"enable_debug"`
}

type EventHandler func(rawEvent []byte) error

type Service struct {
	config          Config
	clientPublicKey *rsa.PublicKey
	handler         EventHandler
}

func New(cfg Config, handler EventHandler) *Service {
	if cfg.Endpoint == "" || handler == nil {
		log.Fatal("Event service config have empty values")
	}

	cPubKey, err := parserKey(cfg)
	if err != nil {
		log.Fatalf("Failed to parse event server key: %v", err)
	}

	return &Service{
		clientPublicKey: cPubKey,
		handler:         handler,
		config:          cfg,
	}
}

func parserKey(config Config) (*rsa.PublicKey, error) {
	clientPubKeyByte, err := os.ReadFile(config.ClientPubKeyPath)
	if err != nil {
		return nil, err
	}
	clientPubKey, err := jwt.ParseRSAPublicKeyFromPEM(clientPubKeyByte)
	if err != nil {
		return nil, err
	}

	return clientPubKey, nil
}

func (s *Service) Start() {
	if !s.config.EnableDebug {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()

	r.GET("/ping", s.Ping)
	api := r.Group("/v2")
	api.Use(s.jwtAuthMiddleware())
	api.POST("/event", s.handleRequest)

	log.Infof("%v %v is running", s.config.ServiceName, s.config.Endpoint)
	log.Fatal(r.Run(s.config.Endpoint))
}

func (s *Service) Ping(c *gin.Context) {
	log.Debugf("Got ping")
	pong := make(map[string]string, 2)
	pong["server"] = s.config.ServiceName
	pong["timestamp"] = strconv.FormatInt(time.Now().UnixMilli(), 10)
	c.JSON(http.StatusOK, pong)
}

func (s *Service) handleRequest(c *gin.Context) {
	rawRequest, err := s.GetRawRequest(c)
	if err != nil {
		log.Errorf("event server get raw request error: %v", err)
		s.SendResponse(c, http.StatusBadRequest)
		return
	}
	err = s.Process(rawRequest)
	if err != nil {
		log.Errorf("event server process request error: %v", err)
		s.SendResponse(c, http.StatusBadRequest)
		return
	}
	s.SendResponse(c, http.StatusOK)
}

func (s *Service) Process(rawRequest []byte) error {
	if s.handler == nil {
		log.Errorf("event service no handler registered")
		return fmt.Errorf("event service no handler registered")
	}
	return s.handler(rawRequest)
}

func (s *Service) GetRawRequest(c *gin.Context) ([]byte, error) {
	data, exist := c.Get("request")
	if !exist {
		return nil, errors.New("request field not exist")
	}

	byteData, ok := data.([]byte)
	if !ok {
		return nil, errors.New("request data is not []byte type")
	}

	var claim PackageDataClaim
	if err := json.Unmarshal(byteData, &claim); err != nil {
		return nil, err
	}

	return claim.PackageData, nil
}

// ============================== HELP FUNCTION ===========================

func (s *Service) ExtractToken(c *gin.Context) string {
	bearToken := c.PostForm("TSS_JWT_MSG")
	return strings.Trim(bearToken, " ")
}

func (s *Service) VerifyToken(c *gin.Context) (*jwt.Token, error) {
	tokenString := s.ExtractToken(c)

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Make sure that the token method conform to "SigningMethodRSA"
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return s.clientPublicKey, nil
	})
	if err != nil {
		return nil, err
	}

	return token, nil
}

func (s *Service) TokenValid(c *gin.Context) error {
	token, err := s.VerifyToken(c)
	if err != nil {
		return err
	}

	if token.Valid {
		data, err := json.Marshal(token.Claims)
		if err != nil {
			return errors.New("fail to convert token.Claims to []byte")
		}

		c.Set("request", data)
		return nil
	}

	return errors.New("invalid token")
}

func (s *Service) jwtAuthMiddleware() func(c *gin.Context) {
	return func(c *gin.Context) {
		err := s.TokenValid(c)
		if err != nil {
			log.Errorf("event server jwt auth middleware error: %v", err)
			s.SendResponse(c, http.StatusOK)
			return
		}
		c.Next()
	}
}

func (s *Service) SendResponse(c *gin.Context, httpStatusCode int) {
	c.JSON(httpStatusCode, "")
	c.Abort()
}
