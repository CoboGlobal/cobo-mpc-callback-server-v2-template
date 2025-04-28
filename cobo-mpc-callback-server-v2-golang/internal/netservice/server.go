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

	"github.com/CoboGlobal/cobo-mpc-callback-server-v2/internal/types"
	"github.com/CoboGlobal/cobo-mpc-callback-server-v2/pkg/log"
	coboWaaS2 "github.com/CoboGlobal/cobo-waas2-go-sdk/cobo_waas2"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type Config struct {
	ServiceName        string `mapstructure:"service_name"`
	Endpoint           string `mapstructure:"endpoint"`
	TokenExpireMinutes uint64 `mapstructure:"token_expire_minutes"`
	ClientPubKeyPath   string `mapstructure:"client_public_key_path"`
	ServicePriKeyPath  string `mapstructure:"service_private_key_path"`
	EnableDebug        bool   `mapstructure:"enable_debug"`
}

type RequestHandler func(rawRequest []byte) (*coboWaaS2.TSSCallbackResponse, error)

type Service struct {
	config            Config
	clientPublicKey   *rsa.PublicKey
	servicePrivateKey *rsa.PrivateKey
	tokenExpireTime   time.Duration
	handler           RequestHandler
}

func New(cfg Config, handler RequestHandler) *Service {
	if cfg.Endpoint == "" || cfg.TokenExpireMinutes == 0 || handler == nil {
		log.Fatal("Callback service config have empty values")
	}

	cPubKey, sPriKey, err := parserKey(cfg)
	if err != nil {
		log.Fatalf("Failed to parse callback public and private key: %v", err)
	}

	return &Service{
		clientPublicKey:   cPubKey,
		servicePrivateKey: sPriKey,
		tokenExpireTime:   time.Duration(cfg.TokenExpireMinutes) * time.Minute,
		handler:           handler,
		config:            cfg,
	}
}

func parserKey(config Config) (*rsa.PublicKey, *rsa.PrivateKey, error) {
	clientPubKeyByte, err := os.ReadFile(config.ClientPubKeyPath)
	if err != nil {
		return nil, nil, err
	}
	clientPubKey, err := jwt.ParseRSAPublicKeyFromPEM(clientPubKeyByte)
	if err != nil {
		return nil, nil, err
	}

	servicePriKeyByte, err := os.ReadFile(config.ServicePriKeyPath)
	if err != nil {
		return nil, nil, err
	}
	servicePriKey, err := jwt.ParseRSAPrivateKeyFromPEM(servicePriKeyByte)
	if err != nil {
		return nil, nil, err
	}

	return clientPubKey, servicePriKey, nil
}

func (s *Service) Start() {
	if !s.config.EnableDebug {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()

	r.GET("/ping", s.Ping)
	api := r.Group("/v2")
	api.Use(s.jwtAuthMiddleware())
	api.POST("/check", s.RiskControl)

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

func (s *Service) RiskControl(c *gin.Context) {
	rawRequest, err := s.GetRawRequest(c)
	if err != nil {
		status := int32(types.StatusInvalidRequest)
		errStr := err.Error()
		rsp := &coboWaaS2.TSSCallbackResponse{
			Status: &status,
			Error:  &errStr,
		}
		s.SendResponse(c, rsp, http.StatusOK)
		return
	}
	rsp, err := s.Process(rawRequest)
	if err != nil {
		log.Errorf("Callback process request error: %v", err)
		s.SendResponse(c, rsp, http.StatusBadRequest)
		return
	}
	s.SendResponse(c, rsp, http.StatusOK)
}

func (s *Service) Process(rawRequest []byte) (*coboWaaS2.TSSCallbackResponse, error) {
	//log.Debugf("Callback process request: %v", string(rawRequest))
	if s.handler == nil {
		log.Errorf("callback service no handler registered")
		return nil, fmt.Errorf("callback service no handler registered")
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

	var claim types.PackageDataClaim
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
			status := int32(types.StatusInvalidToken)
			errStr := err.Error()
			rsp := &coboWaaS2.TSSCallbackResponse{
				Status: &status,
				Error:  &errStr,
			}
			s.SendResponse(c, rsp, http.StatusOK)
			return
		}
		c.Next()
	}
}

func (s *Service) SendResponse(c *gin.Context, rsp *coboWaaS2.TSSCallbackResponse, httpStatusCode int) {
	if rsp == nil {
		log.Errorf("callback server response nil")
		c.JSON(http.StatusInternalServerError, fmt.Errorf("callback server response nil"))
		c.Abort()
		return
	}

	rspJSON, _ := rsp.MarshalJSON()
	if *rsp.Status == types.StatusOK {
		log.WithField("request_id", *rsp.RequestId).Infof("Callback server http code %v, response: %v", httpStatusCode, string(rspJSON))
	} else {
		log.WithField("request_id", *rsp.RequestId).Errorf("Callback server http code %v, response: %v", httpStatusCode, string(rspJSON))
	}

	data, err := rsp.MarshalJSON()
	if err != nil {
		log.WithField("request_id", *rsp.RequestId).Errorf("callback server response marshal error: %v", err)
		c.JSON(http.StatusInternalServerError, err.Error())
		c.Abort()
		return
	}

	token, err := s.CreateToken(data)
	if err != nil {
		log.WithField("request_id", *rsp.RequestId).Errorf("callback server response create token error: %v", err)
		c.JSON(http.StatusInternalServerError, err.Error())
		c.Abort()
		return
	}
	c.JSON(httpStatusCode, token)
	c.Abort()
}

func (s *Service) CreateToken(data []byte) (string, error) {
	expirationTime := time.Now().Add(s.tokenExpireTime)

	claims := &types.PackageDataClaim{
		PackageData: data,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			Issuer:    s.config.ServiceName,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(s.servicePrivateKey)
}
