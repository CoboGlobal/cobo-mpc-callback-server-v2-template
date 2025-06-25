package validator

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"sort"

	"github.com/google/go-cmp/cmp"
	"github.com/nikolalohinski/gonja/v2"
	"github.com/nikolalohinski/gonja/v2/exec"
	"github.com/nikolalohinski/gonja/v2/loaders"
)

type StatementBuilder struct {
	template string
}

func NewStatementBuilder(template string) *StatementBuilder {
	return &StatementBuilder{
		template: template,
	}
}

// getGonjaFilters returns a list of custom filters to be used with Gonja
func getGonjaFilters() map[string]exec.FilterFunction {
	return map[string]exec.FilterFunction{
		"toString": func(e *exec.Evaluator, in *exec.Value, params *exec.VarArgs) *exec.Value {
			bytes, _ := json.Marshal(in.Interface())
			return exec.AsValue(string(bytes))
		},
		"toInt": func(e *exec.Evaluator, in *exec.Value, params *exec.VarArgs) *exec.Value {
			switch val := in.Interface().(type) {
			case float64:
				return exec.AsValue(int(val))
			case int:
				return exec.AsValue(val)
			case string:
				var result int
				fmt.Sscanf(val, "%d", &result)
				return exec.AsValue(result)
			default:
				return exec.AsValue(0)
			}
		},
		"len": func(e *exec.Evaluator, in *exec.Value, params *exec.VarArgs) *exec.Value {
			switch val := in.Interface().(type) {
			case []interface{}:
				return exec.AsValue(len(val))
			case map[string]interface{}:
				return exec.AsValue(len(val))
			case string:
				return exec.AsValue(len(val))
			default:
				return exec.AsValue(0)
			}
		},
		"toList1": func(e *exec.Evaluator, in *exec.Value, params *exec.VarArgs) *exec.Value {
			if slice, ok := in.Interface().([]interface{}); ok {
				var result []string
				for _, item := range slice {
					if item != nil {
						result = append(result, fmt.Sprintf("%v", item))
					}
				}
				bytes, _ := json.Marshal(result)
				return exec.AsValue(string(bytes))
			}
			return exec.AsValue("[]")
		},
		"toList2": func(e *exec.Evaluator, in *exec.Value, params *exec.VarArgs) *exec.Value {
			if slice, ok := in.Interface().([]interface{}); ok {
				var result [][]string
				for _, row := range slice {
					if rowSlice, ok := row.([]interface{}); ok {
						var rowResult []string
						for _, item := range rowSlice {
							if item != nil {
								rowResult = append(rowResult, fmt.Sprintf("%v", item))
							}
						}
						if len(rowResult) > 0 {
							result = append(result, rowResult)
						}
					}
				}
				bytes, _ := json.Marshal(result)
				return exec.AsValue(string(bytes))
			}
			return exec.AsValue("[]")
		},
		"toRules": func(e *exec.Evaluator, in *exec.Value, params *exec.VarArgs) *exec.Value {
			if slice, ok := in.Interface().([]interface{}); ok {
				var result []map[string]string
				for _, item := range slice {
					if mapItem, ok := item.(map[string]interface{}); ok {
						ruleMap := make(map[string]string)
						for k, v := range mapItem {
							if v != nil {
								ruleMap[k] = fmt.Sprintf("%v", v)
							}
						}
						if len(ruleMap) > 0 {
							result = append(result, ruleMap)
						}
					}
				}
				bytes, _ := json.Marshal(result)
				return exec.AsValue(string(bytes))
			}
			return exec.AsValue("[]")
		},
	}
}

func getGonjaDictMethods() *exec.MethodSet[map[string]interface{}] {
	return exec.NewMethodSet[map[string]interface{}](map[string]exec.Method[map[string]interface{}]{
		// keys method is builtins method defined in https://github.com/NikolaLohinski/gonja/blob/master/builtins/methods/dict.go#L9
		"keys": func(self map[string]interface{}, selfValue *exec.Value, arguments *exec.VarArgs) (interface{}, error) {
			if err := arguments.Take(); err != nil {
				return nil, exec.ErrInvalidCall(err)
			}
			keys := make([]string, 0)
			for key := range self {
				keys = append(keys, key)
			}
			sort.Strings(keys)
			return keys, nil
		},
		"get": func(self map[string]interface{}, selfValue *exec.Value, arguments *exec.VarArgs) (interface{}, error) {
			if len(arguments.Args) != 1 {
				return nil, exec.ErrInvalidCall(fmt.Errorf("get method expects exactly 1 argument, got %d", len(arguments.Args)))
			}
			key, ok := arguments.Args[0].Interface().(string)
			if !ok {
				return nil, exec.ErrInvalidCall(fmt.Errorf("get method expects string key"))
			}
			return self[key], nil
		},
	})
}

func (s *StatementBuilder) Build(bizData string) (string, error) {
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(bizData), &data); err != nil {
		fmt.Printf("Error parsing JSON data for build statement: %v\n", err)
		return "", fmt.Errorf("error parsing JSON data: %w", err)
	}

	template, err := getGonjaTemplate(s.template)
	if err != nil {
		return "", fmt.Errorf("error parsing template: %w", err)
	}

	// Create context with data
	context := exec.NewContext(data)

	// Render the template
	message, err := template.ExecuteToString(context)
	if err != nil {
		return "", fmt.Errorf("error rendering template: %w", err)
	}

	return message, nil
}

func getGonjaTemplate(source string) (*exec.Template, error) {
	// Get custom filters
	customFilters := getGonjaFilters()

	// Create a new filter set with custom filters
	customFilterSet := exec.NewFilterSet(customFilters)

	// Merge with default filters using Update method
	filterSet := gonja.DefaultEnvironment.Filters.Update(customFilterSet)

	methods := gonja.DefaultEnvironment.Methods

	// Get custom dict methods to the methods
	methods.Dict = getGonjaDictMethods()

	// Create environment with merged filters and methods
	env := &exec.Environment{
		Context:           gonja.DefaultEnvironment.Context,
		Filters:           filterSet,
		ControlStructures: gonja.DefaultEnvironment.ControlStructures,
		Tests:             gonja.DefaultEnvironment.Tests,
		Methods:           methods,
	}

	sourceBytes := []byte(source)
	rootID := fmt.Sprintf("root-%s", string(sha256.New().Sum(sourceBytes)))

	loader, err := loaders.NewFileSystemLoader("")
	if err != nil {
		return nil, err
	}
	shiftedLoader, err := loaders.NewShiftedLoader(rootID, bytes.NewReader(sourceBytes), loader)
	if err != nil {
		return nil, err
	}

	return exec.NewTemplate(rootID, gonja.DefaultConfig, shiftedLoader, env)
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
