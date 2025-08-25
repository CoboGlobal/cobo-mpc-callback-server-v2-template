package waas2

import (
	"context"
	"fmt"
	"strings"

	coboWaas2 "github.com/CoboGlobal/cobo-waas2-go-api/waas2"
)

type Getter interface {
	ListTransactionAndApprovalDetails(ctx context.Context, transactionIds []string) ([]*TransactionApprovalDetail, error)
}

type TransactionApprovalDetail struct {
	TransactionId  string
	Transaction    *coboWaas2.Transaction
	ApprovalDetail *coboWaas2.ApprovalDetail
	Templates      []coboWaas2.ApprovalTemplate
}

type TemplateName struct {
	TemplateKey     string
	TemplateVersion string
}

type Waas2 struct {
	client           *Client
	templateMapCache map[string]coboWaas2.ApprovalTemplate
}

func NewWaas2(client *Client) *Waas2 {
	return &Waas2{
		client: client,
	}
}

func getTemplateKey(templateName TemplateName) string {
	return templateName.TemplateKey + "_" + templateName.TemplateVersion
}

func (w *Waas2) ListTransactionAndApprovalDetails(ctx context.Context, transactionIds []string) ([]*TransactionApprovalDetail, error) {
	txApprovalDetails := make([]*TransactionApprovalDetail, len(transactionIds))

	txs, err := w.client.ListTransactions(ctx, transactionIds)
	if err != nil {
		return nil, fmt.Errorf("txs %v failed to list transactions: %w", transactionIds, err)
	}

	approvalDetails, err := w.client.ListTransactionApprovalDetails(ctx, transactionIds)
	if err != nil {
		return nil, fmt.Errorf("txs %v failed to list transaction approval details: %w", transactionIds, err)
	}

	// set tx approval details
	for i, txId := range transactionIds {
		// set transaction id
		txApprovalDetails[i].TransactionId = txId

		found := false
		for _, tx := range txs {
			if tx.TransactionId == txId {
				txApprovalDetails[i].Transaction = &tx
				found = true
				break
			}
		}
		if !found {
			return nil, fmt.Errorf("tx %s not found in txs", txId)
		}

		foundApprovalDetail := false
		for _, approvalDetail := range approvalDetails {
			if approvalDetail.TransactionId != nil && *approvalDetail.TransactionId == txId {
				txApprovalDetails[i].ApprovalDetail = &approvalDetail

				if txApprovalDetails[i].Transaction.Type == nil {
					return nil, fmt.Errorf("tx %s type is nil", txId)
				}

				transactionType := strings.ToLower(string(*txApprovalDetails[i].Transaction.Type))

				templates, err := w.getTemplateByApprovalDetail(ctx, approvalDetail, transactionType)
				if err != nil {
					return nil, fmt.Errorf("tx %s failed to get template by approval detail: %w", txId, err)
				}
				txApprovalDetails[i].Templates = templates
				foundApprovalDetail = true
				break
			}
		}
		if !foundApprovalDetail {
			return nil, fmt.Errorf("tx %s not found in approval details", txId)
		}
	}

	return txApprovalDetails, nil
}

func (w *Waas2) getTemplateByApprovalDetail(ctx context.Context, approvalDetail coboWaas2.ApprovalDetail, transactionType string) ([]coboWaas2.ApprovalTemplate, error) {
	needAllTemplateList := make([]coboWaas2.ApprovalTemplate, 0)
	needAllTemplateNameMap := make(map[string]TemplateName)
	needFetchTemplateNameMap := make(map[string]TemplateName)

	if approvalDetail.TransactionId == nil {
		return nil, fmt.Errorf("tx transaction id is nil")
	}

	handleUserDetails := func(templateKey string, userDetails []coboWaas2.ApprovalUserDetail) {
		for _, userDetail := range userDetails {
			if userDetail.TemplateVersion == nil {
				continue
			}
			templateName := TemplateName{
				TemplateKey:     templateKey,
				TemplateVersion: *userDetail.TemplateVersion,
			}
			needAllTemplateNameMap[getTemplateKey(templateName)] = templateName

			if _, exists := w.templateMapCache[getTemplateKey(templateName)]; !exists {
				needFetchTemplateNameMap[getTemplateKey(templateName)] = templateName
			}
		}
	}

	if approvalDetail.BrokerUser != nil {
		templateKey := "broker_user"
		handleUserDetails(templateKey, approvalDetail.BrokerUser.UserDetails)
	}

	if approvalDetail.Spender != nil {
		templateKey := transactionType
		handleUserDetails(templateKey, approvalDetail.Spender.UserDetails)
	}

	if approvalDetail.Approver != nil {
		templateKey := transactionType
		handleUserDetails(templateKey, approvalDetail.Approver.UserDetails)
	}

	if len(needFetchTemplateNameMap) > 0 {
		needFetchTemplateNameList := make([]TemplateName, 0)
		for templateName := range needFetchTemplateNameMap {
			needFetchTemplateNameList = append(needFetchTemplateNameList, needFetchTemplateNameMap[templateName])
		}
		templates, err := w.client.ListTransactionTemplates(ctx, needFetchTemplateNameList)
		if err != nil {
			return nil, fmt.Errorf("tx %s failed to list transaction templates: %w", *approvalDetail.TransactionId, err)
		}
		for _, template := range templates {
			if template.BusinessKey == nil || template.TemplateVersion == nil {
				return nil, fmt.Errorf("tx %s template business key or template version is nil", *approvalDetail.TransactionId)
			}
			templateName := TemplateName{
				TemplateKey:     *template.BusinessKey,
				TemplateVersion: *template.TemplateVersion,
			}
			w.templateMapCache[getTemplateKey(templateName)] = template
		}
	}

	for _, templateName := range needAllTemplateNameMap {
		needAllTemplateList = append(needAllTemplateList, w.templateMapCache[getTemplateKey(templateName)])
	}

	return needAllTemplateList, nil
}
