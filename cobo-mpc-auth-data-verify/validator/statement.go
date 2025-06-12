package validator

import (
	"encoding/json"
	"fmt"

	"github.com/google/go-cmp/cmp"
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
		jinja2.WithFilter("toString", `import json
def toString(v):
    return json.dumps(v, ensure_ascii=False)`),
		jinja2.WithFilter("toInt", `def toInt(v):
    return int(v)`),
		jinja2.WithFilter("len", `def len(v):
    return __builtins__['len'](v)`),
		jinja2.WithFilter("toList1", `import json
def toList1(v):
    return json.dumps([str(x) for x in v if x], ensure_ascii=False)`),
		jinja2.WithFilter("toList2", `import json
def toList2(v):
    return json.dumps([[str(x) for x in row if x] for row in v if row], ensure_ascii=False)`),
		jinja2.WithFilter("toRules", `import json
def toRules(v):
    return json.dumps([{str(k): str(v) for k, v in x.items()} for x in v if x], ensure_ascii=False)`),
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

func CompareStatementMessage(message1, message2 string) (bool, string) {
	var data1, data2 interface{}

	// Parse first JSON string
	if err := json.Unmarshal([]byte(message1), &data1); err != nil {
		return false, fmt.Sprintf("failed to parse first statement message: %v", err)
	}

	// Parse second JSON string
	if err := json.Unmarshal([]byte(message2), &data2); err != nil {
		return false, fmt.Sprintf("failed to parse second statement message: %v", err)
	}

	if diff := cmp.Diff(data1, data2); diff != "" {
		return false, fmt.Sprintf("statement message differences:\n%s", diff)
	}

	return true, ""
}
