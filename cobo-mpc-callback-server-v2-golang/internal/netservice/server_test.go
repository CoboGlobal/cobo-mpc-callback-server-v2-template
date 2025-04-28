package netservice

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/CoboGlobal/cobo-mpc-callback-server-v2/internal/types"
	coboWaaS2 "github.com/CoboGlobal/cobo-waas2-go-sdk/cobo_waas2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

var (
	//nolint:gosec,lll
	testServerPriKey = "-----BEGIN RSA PRIVATE KEY-----\nMIIEpQIBAAKCAQEAzLJd44qbp0N5cXKw0JNxCTQDWtC9X06rYoskRHxf5Tjy0AjW\njOU5NvpwGC4YPoK2U3Ln6YDKVF7W6wyTPePM2ECJE8y8N4B/BXpOx5hWqTmhqtZD\nfHR+KRIbkpIAjwTDJD7chINgqY53dC5fJ22F82r4e8g2kwp7C2K6SqgLr5kU/y30\nnUDgC5M27O4i1pP84t7qqMeuEAHey9zKJb640PRgUAXAWAZhL+3ylYzkSVK8ODwh\n2XJp4J1DDp20Q2M4YK98lbCpxEPncyrgRHIJid3jC0ykTPIGwGJxVXlnD6VvYl0j\n0VZ5Kq2DFIPPgmZzecQnLT9Vt5myFr5SqLOSSQIDAQABAoIBAQCpw9BsU1tuaF6D\nAVy1T2LzAAk8O1yje6pWKxHkHsalZAq1EG9oIP/HogJve2MuDNhL80N1fBPRz2ot\nPJutO417WGKXYjhDS7WNBHfrv2M4LAzxk4wa3r53L4Zgk+gUtR1mpR/cYt07ImXd\nnEvcdlAepnv4pP7mCk4sDjB0lFREx5QU5pJFbCuixAVzqXwTfe5fHbLAbbbyclFB\nMdWdVJW/t0hTcZdBUtvk1hcpgAJ6Ykl9waSOSKBvAXa9KGqPIY/f328CD7YUbXFg\nfucmlN472rsgXn7u0qP3hOox7sNw5JjEULCGg+m+JVqJ/eclDC288fI9lXOXxFTa\n38MB1TgFAoGBANyr2EKn97XDVXfGZUOCIAzq1inxkJnJCIOkBlDqwNiSfh3VFK5v\nppy2jGxGYbMdj5tca6L1Yql/DrCmWR+KJOYz3QB56YrG/pe+D2TKnSvPF3jZVQHH\nJClHW56v9PqzAbwPbyjTiStZ7oUHQ5eY5vi6ygvnUhybYNDSFAxxoDwLAoGBAO13\nz2xMSTPbmmK88UsR9S8NwshJ54a/RQl38IQJ+9TamfX/IuN2mZcs3ixr+CGv/saY\n0JMDGp0tCWdoP5ImPsq+OyjJgXXP2tzXaJrQSsCb+N4vywjvMjAmhxcdp3pKcw7F\n1NlQdaVV4InEIajbYqtyQhB2GMZyCHyq2qL5jst7AoGAJUmcV1cOklYZYQ3TGp8o\nT0Z3PcslxfakS6oxrwab43yNdvkEb51KJ/zoqXsTEzMRiw0I2xZfv4hKsSrKsHul\nVIi69VOkVODfMEDbVQqvmDF8I92FcbF2uMrn/l55JMuOpXpuLBXifcLKfQwHLdyW\nWr0lWvGRfGf86gw1ewzQKJUCgYEAolUx7aWksRehTXg+NwRaqMTub77d0CZ2ykc8\nmva8OcEKWLkGH5rW2hpo8tMIN/c44ohapPUNP38nG5KPSphsempaxMIjhucFhcyX\njKVxRIQbN8BSOpRRqcrctHeoIpg8WU/x9nDjS5gOO/9gxy7aH7um39vridUwahDe\nD2UsMXsCgYEAli9kGIRr/xe7++W4vW3RTCDDvDNIEbH6JZs1X1EVTKxWhGzgCvCO\n+lMUQP4gzlJ6DJQd4iZo7Ukj31VSHKH522BIRlKNLT7JhOdDm/Dj7ZP+EF8MO4P+\nxJ+W+aPuNbLZmvXdgbKsDSAkIFY5Wweu774/YkPp+ngRhamYoGASHX8=\n-----END RSA PRIVATE KEY-----\n"
	//nolint:lll
	testServerPubKey = "-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAzLJd44qbp0N5cXKw0JNx\nCTQDWtC9X06rYoskRHxf5Tjy0AjWjOU5NvpwGC4YPoK2U3Ln6YDKVF7W6wyTPePM\n2ECJE8y8N4B/BXpOx5hWqTmhqtZDfHR+KRIbkpIAjwTDJD7chINgqY53dC5fJ22F\n82r4e8g2kwp7C2K6SqgLr5kU/y30nUDgC5M27O4i1pP84t7qqMeuEAHey9zKJb64\n0PRgUAXAWAZhL+3ylYzkSVK8ODwh2XJp4J1DDp20Q2M4YK98lbCpxEPncyrgRHIJ\nid3jC0ykTPIGwGJxVXlnD6VvYl0j0VZ5Kq2DFIPPgmZzecQnLT9Vt5myFr5SqLOS\nSQIDAQAB\n-----END PUBLIC KEY-----\n"
	//nolint:gosec,lll
	testCheckerPriKey = "-----BEGIN RSA PRIVATE KEY-----\nMIIEowIBAAKCAQEAns81jvg4gzE+PIZgQdOqh8R0F0olnD2cqD+4lH3wqWNngOjj\ncLABRMY6kuTfxAmkBqOqt3DqVA5bfb7OXh/NwLQrY17AI2MGJ6P08eUkSd2x6x0g\nrV8EMKlCIBeLpWBff62NvgdmCJ2Iz4o9BGogBK6btNbf09NAwrl/TNBftxrD+l1J\nNxPCm6eTURWQRuIixzTu9Bs2yxpUTmyR8YnW/featPaxHhAYgvy7BanJB+dUMgrc\nhiFWBK3Gy4T7L8BTtqT4vHuUBwut/1uLZAFKE+0elaLN5c0kIGCt5r0mbPHedj0H\nFpWIOqVCO9qRX74zLlxkkzwsMLkxIMOCuj9wewIDAQABAoIBAQCVEDPahb0tz2v2\njb8OKqq3kzvQnIVe+SnxdxY/M1NQ+4AsrOzHWj1mm5ZhSTmMHex7WuakFvWsfml8\nRzwXd0y+o57SQB6jWJBvZuNEpmuAdfpJkOaaNUSOlGEAFHm8ehBJnNMd2n34ej3v\naHdLjH2PR4HZpZMklfcEj+8gX7pn99m5XoIQSydeCvcU18n5OY3qK4UOJm1iKmDE\n+C7h5yQqVqdybDFroBWqNQAZFAL688TamVcvHfGSadpOAGVo5jfEmd0jD3wKoWeq\nhwOcKxTiK62K6uqFZGk+1VsNDqEfcFvKTADLaQTN4GvFDU03hHUpTON5zvmjKkVG\nzlLP2MT5AoGBAMEMNOeGy5PULBsP7pAj/pOEmVOBh/PS7A85hfUPd/cWyvIw6ZLr\nna0vhLcGcRxsjsC3oGxEh0zqDHoi8qOrSI4rLNIbCPqc3Cxkv3dPAW+pMnAhAX4G\nekhoJ/MN0Hzt4ZitWg6LSykTVoVwzGBToZUpGuGu1LmRmz9scYr3H2GdAoGBANKY\nwtIIDZbI53iBR12vPYlkvDZZ/JBtNa+iKoSe5LVpRlPsRXdsF1VY9ZcYor+psjnN\ng3aIwHcG1OZvEn8Hzg8zlAJ4j8enWLeb8EVbzlJfo/KqaZ2cCQZhB5pSpeXReTS/\npQNGbfTfCEiFSUBwB1clQzB5JlmwW9qvssWm3qr3AoGAa9wvHwFQc2tDrWcsarrB\nvZiDtoWT+WZq4GLKds7Kv3Krt8AecSlWMvJu23gs8K2y4Ph4GKX9Vrsad49ZNJs5\n8b0r0MSsMqI73k34MGgjLElD1iSK2egyoIwZbhLU30hmGNEalS+8sdmNKQeKGXQA\nvv91do0hbAFv1XL4yaUjkn0CgYBH/WZjq9MRX14ZCIBf2x67D89y+PHoYRzADDxi\nl3pxNSqQV60rdKzJRR625voDcLv3HHS7GWZJifFPUFrPR9i9w5DuA06LHn6qTUkm\nPIrcB8ugkXaHJSbEoniZ3XTOifvX90cuRm4iDffj6oQu3dz0gk1kjZV5hVrw96yx\n+igV4wKBgFfeNRh/bFtGVLzIfu0KT47+/yjA3pbxrhlJHW1oTWUEggN+u+OGAyUK\n67JQxyUGRnaZi6lq6rec2DN9zvWjCVeDw95Ei+40bYpU3XUYCDC+R1OkIWeqhURU\nFi18jRlz6+Ru0d7jWf0nQJX4R59mWAqYoDVoJHjXE6bYTx5hf3Vl\n-----END RSA PRIVATE KEY-----\n"
	//nolint:lll
	testCheckerPubKey = "-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAns81jvg4gzE+PIZgQdOq\nh8R0F0olnD2cqD+4lH3wqWNngOjjcLABRMY6kuTfxAmkBqOqt3DqVA5bfb7OXh/N\nwLQrY17AI2MGJ6P08eUkSd2x6x0grV8EMKlCIBeLpWBff62NvgdmCJ2Iz4o9BGog\nBK6btNbf09NAwrl/TNBftxrD+l1JNxPCm6eTURWQRuIixzTu9Bs2yxpUTmyR8YnW\n/featPaxHhAYgvy7BanJB+dUMgrchiFWBK3Gy4T7L8BTtqT4vHuUBwut/1uLZAFK\nE+0elaLN5c0kIGCt5r0mbPHedj0HFpWIOqVCO9qRX74zLlxkkzwsMLkxIMOCuj9w\newIDAQAB\n-----END PUBLIC KEY-----\n"
)

type testCase struct {
	name               string
	srvR               coboWaaS2.TSSCallbackResponse
	srvError           error
	responseHttpStatus int
	responseAction     string
	responseStatus     int
}

var (
	rsp                  coboWaaS2.TSSCallbackResponse
	errRsp               error
	actionApprove        = coboWaaS2.TSSCALLBACKACTIONTYPE_APPROVE
	actionReject         = coboWaaS2.TSSCALLBACKACTIONTYPE_REJECT
	statusOK             = int32(types.StatusOK)
	statusInternalError  = int32(types.StatusInternalError)
)

var testCases = []testCase{
	{
		name:               "approve",
		srvR:               coboWaaS2.TSSCallbackResponse{Action: &actionApprove, Status: &statusOK},
		srvError:           nil,
		responseHttpStatus: http.StatusOK,
		responseAction:     string(actionApprove),
		responseStatus:     types.StatusOK,
	},
	{
		name:               "bad http status",
		srvR:               coboWaaS2.TSSCallbackResponse{Status: &statusInternalError},
		srvError:           fmt.Errorf("test error"),
		responseHttpStatus: http.StatusBadRequest,
		responseAction:     "",
		responseStatus:     types.StatusInternalError,
	},
	{
		name:               "response error",
		srvR:               coboWaaS2.TSSCallbackResponse{Status: &statusInternalError},
		srvError:           nil,
		responseHttpStatus: http.StatusOK,
		responseAction:     "",
		responseStatus:     types.StatusInternalError,
	},
	{
		name:               "response reject",
		srvR:               coboWaaS2.TSSCallbackResponse{Action: &actionReject, Status: &statusOK},
		srvError:           nil,
		responseHttpStatus: http.StatusOK,
		responseAction:     string(actionReject),
		responseStatus:     types.StatusOK,
	},
}

func createRequestJWT(t *testing.T, data []byte) string {
	t.Helper()

	checkerPriKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(testCheckerPriKey))
	assert.NoError(t, err)

	expirationTime := time.Now().Add(10 * time.Second)

	claims := &types.PackageDataClaim{
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

func parserResponseJWT(t *testing.T, tokenStr string) *coboWaaS2.TSSCallbackResponse {
	t.Helper()

	serverPublicKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(testServerPubKey))
	assert.NoError(t, err)

	token, err := jwt.ParseWithClaims(tokenStr, &types.PackageDataClaim{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return serverPublicKey, nil
	})
	assert.NoError(t, err)

	if !token.Valid {
		assert.Fail(t, "unexpected token: %v", token)
	}

	rspClaim, ok := token.Claims.(*types.PackageDataClaim)
	assert.Equal(t, true, ok)

	rsp := &coboWaaS2.TSSCallbackResponse{}
	err = json.Unmarshal(rspClaim.PackageData, rsp)
	assert.NoError(t, err)

	return rsp
}

func TestService(t *testing.T) {
	serverCfg := Config{
		ServiceName: "",
		Endpoint:    "localhost:9999",
		EnableDebug: false,
	}
	checkerPubKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(testCheckerPubKey))
	assert.NoError(t, err)
	serverPriKey, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(testServerPriKey))
	assert.NoError(t, err)

	service := &Service{
		clientPublicKey:   checkerPubKey,
		servicePrivateKey: serverPriKey,
		tokenExpireTime:   10 * time.Second,
		config:            serverCfg,
		handler: func(rawRequest []byte) (*coboWaaS2.TSSCallbackResponse, error) {
			return &rsp, errRsp
		},
	}
	go service.Start()
	time.Sleep(500 * time.Millisecond)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rsp = tc.srvR
			errRsp = tc.srvError

			// build request jwt
			reqData := []byte("test request")
			reqJWT := createRequestJWT(t, reqData)
			values := url.Values{"TSS_JWT_MSG": {reqJWT}}

			ctx := context.Background()
			req, err := http.NewRequestWithContext(ctx, http.MethodPost, "http://localhost:9999/v2/check", strings.NewReader(values.Encode()))
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
			respJWT := string(respBody)
			respJWT = strings.Trim(respJWT, "\"")
			fmt.Println("respBody ", string(respBody))

			// parse response
			respData := parserResponseJWT(t, respJWT)
			if respData.Action != nil {
				assert.Equal(t, tc.responseAction, string(*respData.Action))
			}
			if respData.Status != nil {
				assert.Equal(t, tc.responseStatus, int(*respData.Status))
			}
		})
	}
}
