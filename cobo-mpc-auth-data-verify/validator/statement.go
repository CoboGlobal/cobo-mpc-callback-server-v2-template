package validator

import (
	"encoding/json"
	"fmt"

	"github.com/kluctl/kluctl/lib/go-jinja2"
)

type StatementBuilder struct {
	template string
}

func NewStatementBuilder(template string) *StatementBuilder {
	return &StatementBuilder{
		template: template,
	}
}

func (s *StatementBuilder) Build(bizData string) (string, error) {

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

	message, err := j2.RenderString(s.template)
	if err != nil {
		return "", fmt.Errorf("error rendering template: %w", err)
	}

	return message, nil
}
