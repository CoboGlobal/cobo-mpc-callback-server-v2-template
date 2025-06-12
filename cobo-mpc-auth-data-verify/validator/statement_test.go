package validator

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/kluctl/kluctl/lib/go-jinja2"
	"github.com/test-go/testify/assert"
)

func TestRenderTemplate(t *testing.T) {
	j2, err := jinja2.NewJinja2("example", 1,
		jinja2.WithGlobal("test_var1", 1),
		jinja2.WithGlobal("test_var2", map[string]any{"test": 2}))
	if err != nil {
		panic(err)
	}
	defer j2.Close()

	template := "{{ test_var1 }}"

	s, err := j2.RenderString(template)
	if err != nil {
		panic(err)
	}

	fmt.Printf("template: %s\nresult: %s", template, s)
}

func TestBuildStatementV2(t *testing.T) {
	bizKeys := []string{
		//"mfa_create_transaction_policy",
		//"mfa_delete_transaction_policy",
		//"mfa_edit_transaction_policy",
		//"mfa_adjust_priorities",
		"transaction",
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
		fmt.Printf("bizKey: %s, Message:\n %s\n", bizKey, message)
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
			name:     "identical simple objects",
			message1: `{"name": "John", "age": 30}`,
			message2: `{"name": "John", "age": 30}`,
			want:     true,
		},
		{
			name:     "identical simple objects with different order",
			message1: `{"name": "John", "age": 30}`,
			message2: `{"age": 30, "name": "John"}`,
			want:     true,
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
