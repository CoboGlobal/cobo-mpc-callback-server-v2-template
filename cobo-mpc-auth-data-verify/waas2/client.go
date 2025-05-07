package waas2

import (
	"context"
	"fmt"
	"os"

	coboWaas2 "github.com/CoboGlobal/cobo-waas2-go-sdk/cobo_waas2"
	"github.com/CoboGlobal/cobo-waas2-go-sdk/cobo_waas2/crypto"
)

type Client struct {
	env    int
	signer crypto.ApiSigner
	client *coboWaas2.APIClient
}

func NewClient(apiSecret string) *Client {
	configuration := coboWaas2.NewConfiguration()
	client := coboWaas2.NewAPIClient(configuration)

	return &Client{
		env: coboWaas2.DevEnv,
		signer: crypto.Ed25519Signer{
			Secret: apiSecret,
		},
		client: client,
	}
}

func (c *Client) createContext(ctx context.Context) context.Context {
	// Select the environment that you use and comment out the other line of code
	ctx = context.WithValue(ctx, coboWaas2.ContextEnv, c.env)
	ctx = context.WithValue(ctx, coboWaas2.ContextPortalSigner, c.signer)
	return ctx
}

func (c *Client) GetTransactionApprovalDetail(ctx context.Context, transactionId string) (*coboWaas2.TransactionApprovalDetail, error) {
	ctx = c.createContext(ctx)

	resp, r, err := c.client.TransactionsAPI.GetTransactionApprovalDetail(ctx, transactionId).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `WalletsAPI.CreateWallet``: %v\n", err)
		if apiErr, ok := err.(*coboWaas2.GenericOpenAPIError); ok {
			fmt.Fprintf(os.Stderr, "Error response: %s\n", string(apiErr.Body()))
		}
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
		return nil, err
	}

	return resp, nil
}
