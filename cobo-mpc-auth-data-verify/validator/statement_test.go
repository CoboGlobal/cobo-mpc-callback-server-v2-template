package validator

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/nikolalohinski/gonja/v2"
	"github.com/nikolalohinski/gonja/v2/exec"
	"github.com/test-go/testify/assert"
)

func TestRenderTemplate(t *testing.T) {
	template, err := gonja.FromString("test_var1 = {{ test_var1 }}; test_var2.test = {{ test_var2.test }}")
	if err != nil {
		panic(err)
	}

	context := exec.NewContext(map[string]interface{}{
		"test_var1": 1,
		"test_var2": map[string]interface{}{"test": 2},
	})

	result, err := template.ExecuteToString(context)
	if err != nil {
		panic(err)
	}

	fmt.Printf("template: {{ test_var1 }}\nresult: %s", result)
}

func TestBuildStatementV2(t *testing.T) {
	bizKeys := []string{
		// "mfa_create_transaction_policy",
		// "transaction",
		"withdraw_approver_approval",
		"withdraw_spender_check",
		"contract_call_approver_approval",
		"contract_call_spender_check",
		"sign_message_approver_approval",
		"sign_message_spender_check",
	}
	for _, bizKey := range bizKeys {
		data, err := getBizData(bizKey)
		assert.NoError(t, err)

		version := "1.0.0"
		templateContent, err := getTemplateContent(bizKey, version)
		assert.NoError(t, err)

		s := NewStatementBuilder(templateContent)
		message, err := s.Build(data)
		assert.NoError(t, err)
		//fmt.Printf("Data:\n %s\n", data)
		//fmt.Printf("bizKey: %s, Message:\n %s\n", bizKey, message)

		message2, err := getMessage(bizKey)
		assert.NoError(t, err)

		got, gotDiff := CompareStatementMessage(message, message2)
		if got != true {
			t.Errorf("CompareStatementMessage() got = %v, want %v", got, true)
		}
		if gotDiff != "" {
			t.Logf("gotDiff: %s", gotDiff)
		}
	}
}

func getBizData(bizKey string) (string, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error getting current directory: %v\n", err)
		return "", fmt.Errorf("error getting current directory: %w", err)
	}

	dataDir := filepath.Join(currentDir, "example_datas")
	dataFile := fmt.Sprintf("%s.json", bizKey)
	fullPath := filepath.Join(dataDir, dataFile)

	dataBytes, err := os.ReadFile(fullPath)
	if err != nil {
		fmt.Printf("Error reading data file: %v\n", err)
		return "", fmt.Errorf("error reading data file: %w", err)
	}

	return string(dataBytes), nil
}

func getMessage(bizKey string) (string, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error getting current directory: %v\n", err)
		return "", fmt.Errorf("error getting current directory: %w", err)
	}

	messageDir := filepath.Join(currentDir, "example_datas", "messages")
	messageFile := fmt.Sprintf("%s_message.json", bizKey)
	fullPath := filepath.Join(messageDir, messageFile)

	messageBytes, err := os.ReadFile(fullPath)
	if err != nil {
		fmt.Printf("Error reading message file: %v\n", err)
		return "", fmt.Errorf("error reading message file: %w", err)
	}

	return string(messageBytes), nil
}

func getTemplateContent(bizKey string, version string) (string, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("error getting current directory: %w", err)
	}

	templateDir := filepath.Join(currentDir, "json_templates")
	templateFile := fmt.Sprintf("%s_%s.json.j2", bizKey, version)
	fullPath := filepath.Join(templateDir, templateFile)

	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return "", fmt.Errorf("template file not found: %s", fullPath)
	}

	templateContent, err := os.ReadFile(fullPath)
	if err != nil {
		return "", fmt.Errorf("error reading template file: %w", err)
	}
	return string(templateContent), nil
}

func TestCompareStatementMessage(t *testing.T) {
	tests := []struct {
		name     string
		message1 string
		message2 string
		want     bool
	}{
		{
			name: "transaction message object",
			message1: `{
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
              "value": 1749715652
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
`,
			message2: `{
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
                            "value": 1749715652
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
}`,
			want: true,
		},
		{
			name:     "identical simple objects",
			message1: `{"name": "John", "age": 30}`,
			message2: `{"name": "John", "age": 30}`,
			want:     true,
		},
		{
			name:     "identical simple objects with different order",
			message1: `{"name": "John", "age": 30}`,
			message2: `{"age": 30,
				"name": "John"}`,
			want: true,
		},
		{
			name:     "different values",
			message1: `{"name": "John", "age": 30}`,
			message2: `{"name": "Jane", "age": 30}`,
			want:     false,
		},
		{
			name:     "missing field",
			message1: `{"name": "John", "age": 30}`,
			message2: `{"name": "John"}`,
			want:     false,
		},
		{
			name:     "nested objects",
			message1: `{"user": {"name": "John", "address": {"city": "Beijing"}}}`,
			message2: `{"user": {"name": "John", "address": {"city": "Shanghai"}}}`,
			want:     false,
		},
		{
			name:     "arrays",
			message1: `{"scores": [1, 2, 3]}`,
			message2: `{"scores": [1, 2, 4]}`,
			want:     false,
		},
		{
			name:     "invalid json first message",
			message1: `{"name": "John", "age": 30`,
			message2: `{"name": "John", "age": 30}`,
			want:     false,
		},
		{
			name:     "invalid json second message",
			message1: `{"name": "John", "age": 30}`,
			message2: `{"name": "John", "age": 30`,
			want:     false,
		},
		{
			name:     "different types",
			message1: `{"age": "30"}`,
			message2: `{"age": 30}`,
			want:     false,
		},
		{
			name:     "empty objects",
			message1: `{}`,
			message2: `{}`,
			want:     true,
		},
		{
			name:     "null values",
			message1: `{"value": null}`,
			message2: `{"value": null}`,
			want:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotDiff := CompareStatementMessage(tt.message1, tt.message2)
			if got != tt.want {
				t.Errorf("CompareStatementMessage() got = %v, want %v", got, tt.want)
			}
			if gotDiff != "" {
				t.Logf("gotDiff: %s", gotDiff)
			}
		})
	}
}
