package validator

import (
	"encoding/json"
	"fmt"

	"github.com/kluctl/kluctl/lib/go-jinja2"
)

type Statement struct {
	templateContent string
}

func NewStatement(templateContent string) *Statement {
	return &Statement{
		templateContent: templateContent,
	}
}

func (s *Statement) BuildStatementV2(bizData string) (string, error) {

	var data map[string]interface{}
	if err := json.Unmarshal([]byte(bizData), &data); err != nil {
		fmt.Printf("Error parsing JSON data for build statement: %v\n", err)
		return "", fmt.Errorf("error parsing JSON data: %w", err)
	}

	j2, err := jinja2.NewJinja2("python3", 1,
		jinja2.WithGlobals(data))
	if err != nil {
		return "", fmt.Errorf("error initializing jinja2: %w", err)
	}
	defer j2.Close()

	message, err := j2.RenderString(s.templateContent)
	if err != nil {
		return "", fmt.Errorf("error rendering template: %w", err)
	}

	return message, nil
}
