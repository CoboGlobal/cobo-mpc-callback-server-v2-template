package validator

import (
	"strings"
	"testing"
)

const (
	testTemplate = `
{
  "_theme": "structured",
  "_biz_version": {{ template_version | toString }},
  "header": {
    "title": {{ header_title | toString }},
    "title_icon": ""
  },
  "body": {
    "components": [
      {
        "_component_type": "section",
        "components": [
          {
            "_component_type": "text",
            "_is_in_list": true,
            "key": "Organization",
            "data": {
              "value": {{ org_name | toString }},
              "label": {{ environment | toString }}
            }
          },
          {
            "_component_type": "text",
            "_is_in_list": true,
            "key": "Source",
            "data": {
              "value": {{ wallet_name | toString }},
              "label": {{ source.source_type | toString }}
            }
          },
          {
            "_component_type": "text",
            "_is_in_list": true,
            "key": "Initiator",
            "data": {
              "value": {{ initiator | toString }}
            }
          },
          {
            "_component_type": "date_time",
            "key": "Created Time",
            "data": {
              "value": {{ created_time | toInt }}
            }
          },
          {% if expired_time %}
          {
            "_component_type": "date_time",
            "key": "Expired Time",
            "data": {
              "value": {{ expired_time | toInt }}
            }
          },
          {% endif %}
          {
            "_component_type": "text",
            "key": "Message ID",
            "data": {
              "value": {{ statement_uuid | toString }}
            }
          }
        ]
      },
      {
        "_component_type": "section",
        "components": [
          {
            "_component_type": "text",
            "key": "Transaction Type",
            "data": {
              "value": {{ type | toString }}
            }
          },
          {
            "_component_type": "text",
            "_actions": [
              "copy"
            ],
            "key": "Request ID",
            "data": {
              "value": {{ request_id | toString }}
            }
          }
        ]
      },
      {
        "_component_type": "section",
        "components": [
          {
            "_component_type": "text",
            "_actions": [
              "copy"
            ],
            "key": "From Address",
            "data": {
              {% if source.source_type == "Org-Controlled" %}
              "value": {{ source.address | toString }}
              {% elif source.source_type == "Asset" %}
              "value": {{ wallet_name | toString }}
              {% endif %}
            }
          },
          {
            "_component_type": "text",
            "_actions": [
              "copy"
            ],
            "key": "To Address",
            "data": {
              {% if destination.account_output.memo %}
              "value": {{ destination.account_output.address ~ "|" ~ destination.account_output.memo | toString }}
              {% else %}
              "value": {{ destination.account_output.address | toString }}
              {% endif %}
            }
          }
        ]
      },
      {
        "_component_type": "section",
        "components": [
          {
            "_component_type": "text",
            "key": "Amount",
            "data": {
              "value": {{ (destination.account_output.amount ~ " " ~ token_id) | toString }}
            }
          },
          {
            "_component_type": "text",
            "key": "Fee",
            "data": {
              "value": {{ (fee.fee_used ~ " " ~ fee.token_id) | toString }}
            }
          }
        ]
      }
    ]
  }
}
	`

	testBizData = `
{
    "transaction_id": "2ed8c173-3d1b-49d1-9621-46e35d20ee92",
    "wallet_id": "104fbf43-27b1-47f0-ae5c-dbe141bcfe27",
    "type": "Withdrawal",
    "status": "Completed",
    "initiator_type": "App",
    "source": {
        "source_type": "Asset",
        "wallet_id": "104fbf43-27b1-47f0-ae5c-dbe141bcfe27"
    },
    "destination": {
        "destination_type": "Address",
        "account_output": {
            "address": "0xd6d15f37737b6f67d15388bf77269a4387e41fe8",
            "amount": "30",
            "memo": null
        }
    },
    "request_id": "Swap-Broker-Custodial-Payback-Transfer-eede2929-82de-4ec3-8ca2-c79b50715654",
    "fee": {
        "fee_type": "Fixed",
        "token_id": "BSC_BNB",
        "fee_used": "0",
        "estimated_fee_used": null,
        "max_fee_amount": null,
        "token_symbol": "BNB"
    },
    "initiator": "Token Swap",
    "transaction_hash": "L9912a18198baf148a03232d9b2ea412",
    "org_id": "f926907c-8141-4cc2-9d12-7808590e2167",
    "wallet_name": "zc-bitget-test",
    "org_name": "Cobo",
    "environment": "Prod",
    "created_time": 1749712477,
    "expired_time": 1749714277,
    "statement_uuid": "eede2929-82de-4ec3-8ca2-c79b50715654",
    "header_title": "Transaction: Approver Approval",
    "token_id": "BSC_USDT",
    "template_version": "1.0.0"
}
	`
	testMessage = `
{
  "_theme": "structured",
  "_biz_version": "1.0.0",
  "header": {
    "title": "Transaction: Approver Approval",
    "title_icon": ""
  },
  "body": {
    "components": [
      {
        "_component_type": "section",
        "components": [
          {
            "_component_type": "text",
            "_is_in_list": true,
            "key": "Organization",
            "data": {
              "value": "Cobo",
              "label": "Prod"
            }
          },
          {
            "_component_type": "text",
            "_is_in_list": true,
            "key": "Source",
            "data": {
              "value": "zc-bitget-test",
              "label": "Asset"
            }
          },
          {
            "_component_type": "text",
            "_is_in_list": true,
            "key": "Initiator",
            "data": {
              "value": "Token Swap"
            }
          },
          {
            "_component_type": "date_time",
            "key": "Created Time",
            "data": {
              "value": 1749712477
            }
          },
          {
            "_component_type": "date_time",
            "key": "Expired Time",
            "data": {
              "value": 1749714277
            }
          },
          {
            "_component_type": "text",
            "key": "Message ID",
            "data": {
              "value": "eede2929-82de-4ec3-8ca2-c79b50715654"
            }
          }
        ]
      },
      {
        "_component_type": "section",
        "components": [
          {
            "_component_type": "text",
            "key": "Transaction Type",
            "data": {
              "value": "Withdrawal"
            }
          },
          {
            "_component_type": "text",
            "_actions": [
              "copy"
            ],
            "key": "Request ID",
            "data": {
              "value": "Swap-Broker-Custodial-Payback-Transfer-eede2929-82de-4ec3-8ca2-c79b50715654"
            }
          }
        ]
      },
      {
        "_component_type": "section",
        "components": [
          {
            "_component_type": "text",
            "_actions": [
              "copy"
            ],
            "key": "From Address",
            "data": {
              "value": "zc-bitget-test"
            }
          },
          {
            "_component_type": "text",
            "_actions": [
              "copy"
            ],
            "key": "To Address",
            "data": {
              "value": "0xd6d15f37737b6f67d15388bf77269a4387e41fe8"
            }
          }
        ]
      },
      {
        "_component_type": "section",
        "components": [
          {
            "_component_type": "text",
            "key": "Amount",
            "data": {
              "value": "30 BSC_USDT"
            }
          },
          {
            "_component_type": "text",
            "key": "Fee",
            "data": {
              "value": "0 BSC_BNB"
            }
          }
        ]
      }
    ]
  }
}
	`

	testPubkey    = "dce5743d58cd0fbd5dcbca1faa2ee184c9c0b0642f97160a9aa063bbba5ba726634abd1571d1c704256d083fbe4e800bce90f069ccf42a4123b67d5f2b164d09"
	testSignature = "4b5290dcd03f3efee3baf32c6f27b75b7d016b29a58ab3edaaf02cede5dfa0f1a6ac8bd15cc3ccaddee274ea43e550eb85be8c48ed7e795c69981fae62f4e410"
	testResult    = 2
)

func TestAuthValidator_Verify(t *testing.T) {
	tests := []struct {
		name     string
		authData *AuthData
		wantErr  bool
		errMsg   string
	}{
		{
			name: "valid auth data",
			authData: &AuthData{
				Result:    testResult,
				Pubkey:    testPubkey,
				Signature: testSignature,
				Template:  testTemplate,
				BizData:   testBizData,
				Message:   testMessage,
			},
			wantErr: false,
		},
		{
			name:     "nil auth data",
			authData: nil,
			wantErr:  true,
			errMsg:   "auth data is nil",
		},
		{
			name: "invalid template",
			authData: &AuthData{
				Result:    testResult,
				Pubkey:    testPubkey,
				Signature: testSignature,
				Template:  "invalid_template",
				BizData:   testBizData,
				Message:   testMessage,
			},
			wantErr: true,
			errMsg:  "failed to parse first statement message",
		},
		{
			name: "invalid biz data",
			authData: &AuthData{
				Result:    testResult,
				Pubkey:    testPubkey,
				Signature: testSignature,
				Template:  testTemplate,
				BizData:   "invalid_biz_data",
				Message:   testMessage,
			},
			wantErr: true,
			errMsg:  "error parsing JSON data",
		},
		{
			name: "invalid message",
			authData: &AuthData{
				Result:    testResult,
				Pubkey:    testPubkey,
				Signature: testSignature,
				Template:  testTemplate,
				BizData:   testBizData,
				Message:   "invalid_message",
			},
			wantErr: true,
			errMsg:  "source message and build message are not equal",
		},
		{
			name: "invalid pubkey",
			authData: &AuthData{
				Result:    testResult,
				Pubkey:    "invalid_pubkey",
				Signature: testSignature,
				Template:  testTemplate,
				BizData:   testBizData,
				Message:   testMessage,
			},
			wantErr: true,
			errMsg:  "error decoding pubkey",
		},
		{
			name: "invalid signature",
			authData: &AuthData{
				Result:    testResult,
				Pubkey:    testPubkey,
				Signature: "invalid_signature",
				Template:  testTemplate,
				BizData:   testBizData,
				Message:   testMessage,
			},
			wantErr: true,
			errMsg:  "error verifying message",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := NewAuthValidator(tt.authData)
			err := v.Verify()

			if tt.wantErr {
				if err == nil {
					t.Errorf("Verify() error = nil, want error containing %q", tt.errMsg)
				} else if !contains(err.Error(), tt.errMsg) {
					t.Errorf("Verify() error = %v, want error containing %q", err, tt.errMsg)
				}
			} else {
				if err != nil {
					t.Errorf("Verify() unexpected error = %v", err)
				}
			}
		})
	}
}

// contains checks if a string contains another string, with optional case-insensitive comparison
func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}
