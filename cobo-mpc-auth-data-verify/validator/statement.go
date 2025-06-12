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

// getJinja2Filters returns a list of filters to be used with Jinja2
func getJinja2Filters() []jinja2.Jinja2Opt {
	return []jinja2.Jinja2Opt{
		jinja2.WithFilter("toString", "lambda v: json.dumps(v, ensure_ascii=False)"),
		jinja2.WithFilter("toInt", "lambda v: int(v)"),
		jinja2.WithFilter("len", "lambda v: len(v)"),
		jinja2.WithFilter("toList1", "lambda v: json.dumps([str(x) for x in v if x], ensure_ascii=False)"),
		jinja2.WithFilter("toList2", "lambda v: json.dumps([[str(x) for x in row if x] for row in v if row], ensure_ascii=False)"),
		jinja2.WithFilter("toRules", "lambda v: json.dumps([{str(k): str(v) for k, v in x.items()} for x in v if x], ensure_ascii=False)"),
	}
}

func (s *StatementBuilder) Build(bizData string) (string, error) {
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(bizData), &data); err != nil {
		fmt.Printf("Error parsing JSON data for build statement: %v\n", err)
		return "", fmt.Errorf("error parsing JSON data: %w", err)
	}

	options := append([]jinja2.Jinja2Opt{
		jinja2.WithGlobals(data),
	}, getJinja2Filters()...)

	j2, err := jinja2.NewJinja2("python3", 1, options...)
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
