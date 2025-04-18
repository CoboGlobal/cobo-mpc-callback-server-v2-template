package netservice

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

var (
	//nolint:gosec,lll
	testClientPriKey = "-----BEGIN RSA PRIVATE KEY-----\nMIIEowIBAAKCAQEAns81jvg4gzE+PIZgQdOqh8R0F0olnD2cqD+4lH3wqWNngOjj\ncLABRMY6kuTfxAmkBqOqt3DqVA5bfb7OXh/NwLQrY17AI2MGJ6P08eUkSd2x6x0g\nrV8EMKlCIBeLpWBff62NvgdmCJ2Iz4o9BGogBK6btNbf09NAwrl/TNBftxrD+l1J\nNxPCm6eTURWQRuIixzTu9Bs2yxpUTmyR8YnW/featPaxHhAYgvy7BanJB+dUMgrc\nhiFWBK3Gy4T7L8BTtqT4vHuUBwut/1uLZAFKE+0elaLN5c0kIGCt5r0mbPHedj0H\nFpWIOqVCO9qRX74zLlxkkzwsMLkxIMOCuj9wewIDAQABAoIBAQCVEDPahb0tz2v2\njb8OKqq3kzvQnIVe+SnxdxY/M1NQ+4AsrOzHWj1mm5ZhSTmMHex7WuakFvWsfml8\nRzwXd0y+o57SQB6jWJBvZuNEpmuAdfpJkOaaNUSOlGEAFHm8ehBJnNMd2n34ej3v\naHdLjH2PR4HZpZMklfcEj+8gX7pn99m5XoIQSydeCvcU18n5OY3qK4UOJm1iKmDE\n+C7h5yQqVqdybDFroBWqNQAZFAL688TamVcvHfGSadpOAGVo5jfEmd0jD3wKoWeq\nhwOcKxTiK62K6uqFZGk+1VsNDqEfcFvKTADLaQTN4GvFDU03hHUpTON5zvmjKkVG\nzlLP2MT5AoGBAMEMNOeGy5PULBsP7pAj/pOEmVOBh/PS7A85hfUPd/cWyvIw6ZLr\nna0vhLcGcRxsjsC3oGxEh0zqDHoi8qOrSI4rLNIbCPqc3Cxkv3dPAW+pMnAhAX4G\nekhoJ/MN0Hzt4ZitWg6LSykTVoVwzGBToZUpGuGu1LmRmz9scYr3H2GdAoGBANKY\nwtIIDZbI53iBR12vPYlkvDZZ/JBtNa+iKoSe5LVpRlPsRXdsF1VY9ZcYor+psjnN\ng3aIwHcG1OZvEn8Hzg8zlAJ4j8enWLeb8EVbzlJfo/KqaZ2cCQZhB5pSpeXReTS/\npQNGbfTfCEiFSUBwB1clQzB5JlmwW9qvssWm3qr3AoGAa9wvHwFQc2tDrWcsarrB\nvZiDtoWT+WZq4GLKds7Kv3Krt8AecSlWMvJu23gs8K2y4Ph4GKX9Vrsad49ZNJs5\n8b0r0MSsMqI73k34MGgjLElD1iSK2egyoIwZbhLU30hmGNEalS+8sdmNKQeKGXQA\nvv91do0hbAFv1XL4yaUjkn0CgYBH/WZjq9MRX14ZCIBf2x67D89y+PHoYRzADDxi\nl3pxNSqQV60rdKzJRR625voDcLv3HHS7GWZJifFPUFrPR9i9w5DuA06LHn6qTUkm\nPIrcB8ugkXaHJSbEoniZ3XTOifvX90cuRm4iDffj6oQu3dz0gk1kjZV5hVrw96yx\n+igV4wKBgFfeNRh/bFtGVLzIfu0KT47+/yjA3pbxrhlJHW1oTWUEggN+u+OGAyUK\n67JQxyUGRnaZi6lq6rec2DN9zvWjCVeDw95Ei+40bYpU3XUYCDC+R1OkIWeqhURU\nFi18jRlz6+Ru0d7jWf0nQJX4R59mWAqYoDVoJHjXE6bYTx5hf3Vl\n-----END RSA PRIVATE KEY-----\n"
	//nolint:lll
	testClientPubKey = "-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAns81jvg4gzE+PIZgQdOq\nh8R0F0olnD2cqD+4lH3wqWNngOjjcLABRMY6kuTfxAmkBqOqt3DqVA5bfb7OXh/N\nwLQrY17AI2MGJ6P08eUkSd2x6x0grV8EMKlCIBeLpWBff62NvgdmCJ2Iz4o9BGog\nBK6btNbf09NAwrl/TNBftxrD+l1JNxPCm6eTURWQRuIixzTu9Bs2yxpUTmyR8YnW\n/featPaxHhAYgvy7BanJB+dUMgrchiFWBK3Gy4T7L8BTtqT4vHuUBwut/1uLZAFK\nE+0elaLN5c0kIGCt5r0mbPHedj0HFpWIOqVCO9qRX74zLlxkkzwsMLkxIMOCuj9w\newIDAQAB\n-----END PUBLIC KEY-----\n"
)

type testCase struct {
	name               string
	srvError           error
	responseHttpStatus int
}

var (
	errRsp error
)

var testCases = []testCase{
	{
		name:               "http status 200",
		srvError:           nil,
		responseHttpStatus: http.StatusOK,
	},
	{
		name:               "bad http status",
		srvError:           fmt.Errorf("test error"),
		responseHttpStatus: http.StatusBadRequest,
	},
}

func createRequestJWT(t *testing.T, data []byte) string {
	t.Helper()

	checkerPriKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(testClientPriKey))
	assert.NoError(t, err)

	expirationTime := time.Now().Add(10 * time.Second)

	claims := &PackageDataClaim{
		PackageData: data,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			Issuer:    "test issuer",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	str, err := token.SignedString(checkerPriKey)
	assert.NoError(t, err)

	return str
}

func TestService(t *testing.T) {
	serverCfg := Config{
		ServiceName: "",
		Endpoint:    "localhost:9999",
		EnableDebug: false,
	}
	checkerPubKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(testClientPubKey))
	assert.NoError(t, err)

	service := &Service{
		clientPublicKey: checkerPubKey,
		config:          serverCfg,
		handler: func(rawEvent []byte) error {
			return errRsp
		},
	}
	go service.Start()
	time.Sleep(500 * time.Millisecond)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			errRsp = tc.srvError

			// build request jwt
			reqData := []byte("test request")
			reqJWT := createRequestJWT(t, reqData)
			values := url.Values{"TSS_JWT_MSG": {reqJWT}}

			ctx := context.Background()
			req, err := http.NewRequestWithContext(ctx, http.MethodPost, "http://localhost:9999/v2/event", strings.NewReader(values.Encode()))
			assert.NoError(t, err)
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			client := &http.Client{}

			// send request
			resp, err := client.Do(req)
			assert.NoError(t, err)
			defer resp.Body.Close()

			// get response
			assert.Equal(t, tc.responseHttpStatus, resp.StatusCode)

			respBody, err := io.ReadAll(resp.Body)
			assert.NoError(t, err)
			fmt.Println("respBody ", string(respBody))
		})
	}
}
