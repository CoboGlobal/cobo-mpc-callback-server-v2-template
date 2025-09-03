package waas2

import (
	"context"
	"fmt"
	"strings"

	coboWaas2 "github.com/CoboGlobal/cobo-waas2-go-sdk/cobo_waas2"
)

type ApproverDetailBuilder interface {
	Build(ctx context.Context, transactionIds []string) ([]*TxApprovalDetail, error)
}

type TemplateName struct {
	TemplateKey     string
	TemplateVersion string
}

type Waas2 struct {
	client           Getter
	templateMapCache map[string]coboWaas2.ApprovalTemplate
}

func NewWaas2(client *Client) *Waas2 {
	return &Waas2{
		client:           client,
		templateMapCache: make(map[string]coboWaas2.ApprovalTemplate),
	}
}

func getTemplateKey(templateName TemplateName) string {
	return templateName.TemplateKey + "_" + templateName.TemplateVersion
}

func businessKeyToTemplateKey(businessKey string) string {
	return strings.TrimPrefix(businessKey, "transaction_")
}

func (w *Waas2) Build(ctx context.Context, transactionIds []string) ([]*TxApprovalDetail, error) {
	txApprovalDetails := make([]*TxApprovalDetail, 0)

	txs, err := w.client.ListTransactions(ctx, transactionIds)
	if err != nil {
		return nil, fmt.Errorf("txs %v failed to list transactions: %w", transactionIds, err)
	}

	approvalDetails, err := w.client.ListTransactionApprovalDetails(ctx, transactionIds)
	if err != nil {
		return nil, fmt.Errorf("txs %v failed to list transaction approval details: %w", transactionIds, err)
	}

	// set tx approval details
	for _, txId := range transactionIds {
		// initialize the struct
		newApprovalDetail := &TxApprovalDetail{
			TransactionId: txId,
		}

		found := false
		for _, tx := range txs {
			if tx.TransactionId == txId {
				newApprovalDetail.Transaction = &tx
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
				newApprovalDetail.ApprovalDetail = &approvalDetail

				if newApprovalDetail.Transaction.Type == nil {
					return nil, fmt.Errorf("tx %s type is nil", txId)
				}

				transactionType := strings.ToLower(string(*newApprovalDetail.Transaction.Type))

				templates, err := w.getTemplatesByApprovalDetail(ctx, approvalDetail, transactionType)
				if err != nil {
					return nil, fmt.Errorf("tx %s failed to get template by approval detail: %w", txId, err)
				}
				newApprovalDetail.Templates = templates
				foundApprovalDetail = true
				break
			}
		}
		if !foundApprovalDetail {
			return nil, fmt.Errorf("tx %s not found in approval details", txId)
		}
		txApprovalDetails = append(txApprovalDetails, newApprovalDetail)
	}

	return txApprovalDetails, nil
}

func (w *Waas2) getTemplatesByApprovalDetail(ctx context.Context, approvalDetail coboWaas2.ApprovalDetail, transactionType string) ([]coboWaas2.ApprovalTemplate, error) {
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

	if approvalDetail.AddressOwner != nil {
		templateKey := "address_owner"
		handleUserDetails(templateKey, approvalDetail.AddressOwner.UserDetails)
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
				TemplateKey:     businessKeyToTemplateKey(*template.BusinessKey),
				TemplateVersion: *template.TemplateVersion,
			}
			w.templateMapCache[getTemplateKey(templateName)] = template
		}
	}

	for _, templateName := range needAllTemplateNameMap {
		template, exists := w.templateMapCache[getTemplateKey(templateName)]
		if !exists {
			return nil, fmt.Errorf("tx %s template %s not found in cache", *approvalDetail.TransactionId, getTemplateKey(templateName))
		}
		needAllTemplateList = append(needAllTemplateList, template)
	}

	return needAllTemplateList, nil
}
