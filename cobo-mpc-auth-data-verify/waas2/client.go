package waas2

import (
	"context"
	"fmt"
	"os"
	"strings"

	coboWaas2 "github.com/CoboGlobal/cobo-waas2-go-sdk/cobo_waas2"
	"github.com/CoboGlobal/cobo-waas2-go-sdk/cobo_waas2/crypto"
)

type Getter interface {
	ListTransactions(ctx context.Context, transactionIds []string) ([]coboWaas2.Transaction, error)
	ListTransactionApprovalDetails(ctx context.Context, transactionIds []string) ([]coboWaas2.ApprovalDetail, error)
	ListTransactionTemplates(ctx context.Context, templateNames []TemplateName) ([]coboWaas2.ApprovalTemplate, error)
}

type Client struct {
	env    int
	signer crypto.ApiSigner
	client *coboWaas2.APIClient
}

func NewClient(apiSecret string, env int) *Client {
	configuration := coboWaas2.NewConfiguration()
	client := coboWaas2.NewAPIClient(configuration)

	return &Client{
		env: env,
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

func (c *Client) ListTransactions(ctx context.Context, transactionIds []string) ([]coboWaas2.Transaction, error) {
	ctx = c.createContext(ctx)

	req := c.client.TransactionsAPI.ListTransactions(ctx)
	req = req.TransactionIds(strings.Join(transactionIds, ","))

	resp, r, err := req.Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `.TransactionsAPI.ListTransactions``: %v\n", err)
		if apiErr, ok := err.(*coboWaas2.GenericOpenAPIError); ok {
			fmt.Fprintf(os.Stderr, "Error response: %s\n", string(apiErr.Body()))
		}
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
		return nil, err
	}
	// todo: handle pagination

	return resp.Data, nil
}

func (c *Client) ListTransactionApprovalDetails(ctx context.Context, transactionIds []string) ([]coboWaas2.ApprovalDetail, error) {
	ctx = c.createContext(ctx)

	req := c.client.TransactionsAPI.ListApprovalDetails(ctx)
	req = req.TransactionIds(strings.Join(transactionIds, ","))
	resp, r, err := req.Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `.TransactionsAPI.ListApprovalDetails``: %v\n", err)
		if apiErr, ok := err.(*coboWaas2.GenericOpenAPIError); ok {
			fmt.Fprintf(os.Stderr, "Error response: %s\n", string(apiErr.Body()))
		}
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
		return nil, err
	}
	return resp, nil
}

func (c *Client) ListTransactionTemplates(ctx context.Context, templateNames []TemplateName) ([]coboWaas2.ApprovalTemplate, error) {
	ctx = c.createContext(ctx)
	templates := make([]coboWaas2.ApprovalTemplate, 0)
	for _, templateName := range templateNames {
		req := c.client.TransactionsAPI.ListTransactionTemplates(ctx)
		req = req.TemplateKey(templateName.TemplateKey)
		req = req.TemplateVersion(templateName.TemplateVersion)
		resp, r, err := req.Execute()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error when calling `TransactionsAPI.ListTransactionTemplates``: %v\n", err)
			if apiErr, ok := err.(*coboWaas2.GenericOpenAPIError); ok {
				fmt.Fprintf(os.Stderr, "Error response: %s\n", string(apiErr.Body()))
			}
			fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
			return nil, err
		}
		templates = append(templates, resp...)
	}

	return templates, nil
}
