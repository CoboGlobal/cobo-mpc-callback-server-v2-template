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
  {% if is_for_sign %}
  "_post_actions": [
    {
        "type": "sign_message",
        "data": {
            "transaction_id": {{ transaction_id | toString }},
            "transaction_type": {{ type | toString }}
        }
    }
  ],
  {% endif %}
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
              "value": {{ show_info.org_name | toString }},
              "label": {{ show_info.environment | toString }}
            }
          },
          {
            "_component_type": "text",
            "_is_in_list": true,
            "key": "Source",
            "data": {
              "value": {{ show_info.wallet_name | toString }},
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
            "key": "Chain",
            "data": {
              "value": {{ chain_id | toString }}
            }
          },
          {
            "_component_type": "text",
            "_style": {
              "is_highlighted": true,
              "is_large_size": false,
              "is_bold": false
            },
            "_actions": [
              "copy"
            ],
            "key": "Transaction ID",
            "data": {
              "value": {{ transaction_id | toString }}
            }
          },
          {
            "_component_type": "text",
            "_style": {
              "is_highlighted": true,
              "is_large_size": false,
              "is_bold": false
            },
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
            "_style": {
              "is_highlighted": true,
              "is_large_size": false,
              "is_bold": false
            },
            "_actions": [
              "copy"
            ],
            "key": "From",
            "data": {
              {% if source.source_type in ["Org-Controlled","User-Controlled","Web3","Safe{Wallet}"] %}
              "label": {{ show_info.from_address_label | toString }},
              "value": {{ source.address | toString }}
              {% elif source.source_type in ["Asset","Main","Sub"] %}
              "value": {{ show_info.wallet_name | toString }}
              {% endif %}
            }
          },
          {% set has_account_output = destination.get("account_output") %}
          {% set has_single_utxo = destination.get("utxo_outputs") and (destination.utxo_outputs | len) == 1 %}
          {% if has_account_output or has_single_utxo %}
          {
            "_component_type": "text",
            "_style": {
              "is_highlighted": true,
              "is_large_size": false,
              "is_bold": false
            },
            "_actions": [
              "copy"
            ],
            "key": "To",
            "data": {
              "label": {{ show_info.to_address_label | toString }},
              {% if destination.get("account_output") and destination.account_output.get("memo") %}
              "value": {{ (destination.account_output.address ~ "|" ~ destination.account_output.memo) | toString }}
              {% elif destination.get("account_output") %}
              "value": {{ destination.account_output.address | toString }}
              {% else %}
              "value": {{ destination.utxo_outputs[0].address | toString }}
              {% endif %}
            }
          }
          {% endif %}
          {% if destination.get("utxo_outputs") and destination.utxo_outputs | len > 1 %}
          {% for output in destination.utxo_outputs %}
          {
            "_component_type": "text",
            "_style": {
              "is_highlighted": true,
              "is_large_size": false,
              "is_bold": false
            },
            "_actions": [
              "copy"
            ],
            {% if loop.index == 1 %}
            "key": "To",
            {% endif %}
            "data": {
              "value": {{ output.address | toString }},
              "sub_value": {{ (output.amount ~ " " ~ token_id) | toString }}
            }
          }{% if not loop.last %},{% endif %}
          {% endfor %}
          {% endif %}
        ]
      },
      {
        "_component_type": "section",
        "components": [
          {% set has_account_output = destination.get("account_output") %}
          {% set has_single_utxo = destination.get("utxo_outputs") and (destination.utxo_outputs | len) == 1 %}
          {% if has_account_output or has_single_utxo %}
          {
            "_component_type": "text",
            "_style": {
              "is_highlighted": false,
              "is_large_size": true,
              "is_bold": true
            },
            "key": "Value",
            "data": {
              {% if destination.get("account_output") %}
              "value": {{ (destination.account_output.amount ~ " " ~ token_id) | toString }}
              {% else %}
              "value": {{ (destination.utxo_outputs[0].amount ~ " " ~ token_id) | toString }}
              {% endif %}
            }
          }
          {% endif %}
          {% if destination.get("utxo_outputs") and destination.utxo_outputs | len > 1 %}
          {
            "_component_type": "text",
            "_style": {
              "is_highlighted": false,
              "is_large_size": true,
              "is_bold": true
            },
            "key": "Value",
            "data": {
              "value": {{ (show_info.total_amount ~ " " ~ token_id) | toString }}
            }
          }
          {% endif %},
          {
            "_component_type": "text",
            "key": "Estimated Fee",
            "data": {
              "value": {% if fee.estimated_fee_used %}{{ (fee.estimated_fee_used ~ " " ~ fee.token_id) | toString }}{% else %}{{ (fee.fee_used ~ " " ~ fee.token_id) | toString }}{% endif %}
            }
          }
          {% if fee.get("max_fee_amount") %},
          {
            "_component_type": "text",
            "key": "Max Fee",
            "data": {
              "value": {{ (fee.max_fee_amount ~ " " ~ fee.token_id) | toString }}
            }
          }
          {% endif %}
        ]
      }
      {% if description %},
      {
        "_component_type": "section",
        "components": [
          {
            "_component_type": "paragraph",
            "key": "Note",
            "_actions": [
              "copy"
            ],
            "data": [
               {{ description | toString }}
            ]
          }
        ]
      }
      {% endif %}
    ]
  }
}
	`

	testBizData = `
{
      "block_info": {
        "block_hash": "0x17a5cf8ea7e177cee3ccbfe6d9bc395b7073468e1c55d59d4d55177faefcd2c7",
        "block_number": 9128859,
        "block_timestamp": 1756954536000
      },
      "category": [],
      "chain_id": "SETH",
      "cobo_category": [],
      "cobo_id": "20250904104620000107060000000331",
      "confirmed_num": 64,
      "confirming_threshold": 64,
      "created_timestamp": 1756953980334,
      "description": "",
      "destination": {
        "account_output": {
          "address": "0x44deca05623d5d85faa426c972064a944591561f",
          "amount": "0.00001"
        },
        "destination_type": "Address",
        "force_external": false,
        "force_internal": false
      },
      "fee": {
        "estimated_fee_used": "0.000031559277456",
        "fee_type": "EVM_EIP_1559",
        "fee_used": "0.000031559277456",
        "gas_limit": "21000",
        "gas_used": "21000",
        "max_fee_per_gas": "1502822736",
        "max_priority_fee_per_gas": "1500000000",
        "token_id": "SETH"
      },
      "initiator": "zhaozhe",
      "initiator_type": "Web",
      "is_loop": false,
      "raw_tx_info": {
        "raw_tx": "0x02f87383aa36a7168459682f0084599341508252089444deca05623d5d85faa426c972064a944591561f8609184e72a00080c080a040189e1e8761df7832e2d80566d0ee86d55630d3184d8ee1b48a491849dd00b8a006e353280710b816e39536ef7d6ee3e486b5f8d841524b2f289de322c9d3c6a0",
        "selected_utxos": [],
        "used_nonce": 22
      },
      "request_id": "web_send_cee0d0a8-0eba-49a4-9d57-961b4aea46f6",
      "source": {
        "address": "0x13c475d9dae8058a8d9a8a724a83f1755db0165c",
        "signer_key_share_holder_group_id": "54274576-6c89-49ca-9924-84988767ee12",
        "source_type": "Org-Controlled",
        "wallet_id": "e499fcb7-4cde-4586-bd14-17d7ea49c25f"
      },
      "status": "Completed",
      "token_id": "SETH",
      "transaction_hash": "0xb44ecf73e8c4dd344398dcdc5106d8e799e6a8c59523901957b154d2f2307559",
      "transaction_id": "7bbe0fe1-187a-436a-80be-f30c586a3914",
      "type": "Withdrawal",
      "updated_timestamp": 1756974676030,
      "wallet_id": "e499fcb7-4cde-4586-bd14-17d7ea49c25f",
      "created_time": 1756953981,
      "header_title": "Transaction: Spender Approval",
      "is_for_sign": false,
      "pubkey": "4136f05547ecd8ab37a8f56908d9137ea45476f88e733c283c3c1caed17c56c855c4cb01fd3f114ea4f82f03e76b79f05ca947c5333697e2ba84337b43c5d39a",
      "result": 2,
      "show_info": {
        "environment": "DEVELOP",
        "from_address_label": "",
        "org_name": "Portal内部测试_Dev_1",
        "to_address_label": "seth测试地址anthony111",
        "wallet_name": "yangming钱包1"
      },
      "signature": "a9679e83726e9e9beca67e12a359b7c6f60a891ae7deaa6fedf636a00f8865f9c81c92bfc45dec7d753d95ad15b6883a75b6da42016d70488625b1f54db93bb0",
      "statement_uuid": "92609b58-f685-458c-be43-388f28f4ec00",
      "template_version": "1.0.0",
      "user_email": "zhaozhe@cobo.com"
    }
	`
	testMessage = `{"_theme":"structured","_biz_version":"1.0.0","header":{"title":"Transaction: Spender Approval","title_icon":""},"body":{"components":[{"_component_type":"section","components":[{"_component_type":"text","_is_in_list":true,"key":"Organization","data":{"value":"Portal内部测试_Dev_1","label":"DEVELOP"}},{"_component_type":"text","_is_in_list":true,"key":"Source","data":{"value":"yangming钱包1","label":"Org-Controlled"}},{"_component_type":"text","_is_in_list":true,"key":"Initiator","data":{"value":"zhaozhe"}},{"_component_type":"date_time","key":"Created Time","data":{"value":1756953981}},{"_component_type":"text","key":"Message ID","data":{"value":"92609b58-f685-458c-be43-388f28f4ec00"}}]},{"_component_type":"section","components":[{"_component_type":"text","key":"Transaction Type","data":{"value":"Withdrawal"}},{"_component_type":"text","key":"Chain","data":{"value":"SETH"}},{"_component_type":"text","_style":{"is_highlighted":true,"is_large_size":false,"is_bold":false},"_actions":["copy"],"key":"Transaction ID","data":{"value":"7bbe0fe1-187a-436a-80be-f30c586a3914"}},{"_component_type":"text","_style":{"is_highlighted":true,"is_large_size":false,"is_bold":false},"_actions":["copy"],"key":"Request ID","data":{"value":"web_send_cee0d0a8-0eba-49a4-9d57-961b4aea46f6"}}]},{"_component_type":"section","components":[{"_component_type":"text","_style":{"is_highlighted":true,"is_large_size":false,"is_bold":false},"_actions":["copy"],"key":"From","data":{"label":"","value":"0x13c475d9dae8058a8d9a8a724a83f1755db0165c"}},{"_component_type":"text","_style":{"is_highlighted":true,"is_large_size":false,"is_bold":false},"_actions":["copy"],"key":"To","data":{"label":"seth测试地址anthony111","value":"0x44deca05623d5d85faa426c972064a944591561f"}}]},{"_component_type":"section","components":[{"_component_type":"text","_style":{"is_highlighted":false,"is_large_size":true,"is_bold":true},"key":"Value","data":{"value":"0.00001 SETH"}},{"_component_type":"text","key":"Estimated Fee","data":{"value":"0.000031559277456 SETH"}}]}]}}`

	testPubkey    = "4136f05547ecd8ab37a8f56908d9137ea45476f88e733c283c3c1caed17c56c855c4cb01fd3f114ea4f82f03e76b79f05ca947c5333697e2ba84337b43c5d39a"
	testSignature = "a9679e83726e9e9beca67e12a359b7c6f60a891ae7deaa6fedf636a00f8865f9c81c92bfc45dec7d753d95ad15b6883a75b6da42016d70488625b1f54db93bb0"
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
				// Message:   testMessage,
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
				//Message:   testMessage,
			},
			wantErr: true,
			errMsg:  "compact rendered template failed",
		},
		{
			name: "invalid biz data",
			authData: &AuthData{
				Result:    testResult,
				Pubkey:    testPubkey,
				Signature: testSignature,
				Template:  testTemplate,
				BizData:   "invalid_biz_data",
				// Message:   testMessage,
			},
			wantErr: true,
			errMsg:  "error parsing JSON data",
		},
		{
			name: "invalid pubkey",
			authData: &AuthData{
				Result:    testResult,
				Pubkey:    "invalid_pubkey",
				Signature: testSignature,
				Template:  testTemplate,
				BizData:   testBizData,
				// Message:   testMessage,
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
				// Message:   testMessage,
			},
			wantErr: true,
			errMsg:  "error verifying message",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := NewAuthValidator(tt.authData)
			err := v.VerifyAuthDataAndResult()

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
