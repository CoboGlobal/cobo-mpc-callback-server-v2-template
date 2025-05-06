package verifier

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/kluctl/kluctl/lib/go-jinja2"
)

func RenderTemplate() {
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

type Statement struct{}

// PushContent 结构体表示推送内容
type PushContent struct {
	TitleMap map[string]string
	BodyMap  map[string]string
}

// PushMsgManager 模拟Python代码中的PushMsgManager
type PushMsgManager struct{}

// GetPushContent 获取推送内容
func (pm *PushMsgManager) GetPushContent(businessKey string, data map[string]interface{}) (map[string]string, map[string]string) {
	// 这里实现获取推送内容的逻辑
	// 这是一个示例实现，实际应根据您的业务逻辑来实现
	titleMap := map[string]string{"default": "通知"}
	bodyMap := map[string]string{"default": "您有一条新消息"}

	return titleMap, bodyMap
}

// NowTs 返回当前的Unix时间戳
func NowTs() int64 {
	return time.Now().Unix()
}

// FillData 自动填充数据
func (s *Statement) FillData(data map[string]interface{}, version string, expiredMinutes int) map[string]interface{} {
	// 复制原始数据，避免修改原始数据
	filledData := make(map[string]interface{})
	for k, v := range data {
		filledData[k] = v
	}

	// 添加statement_uuid，如果不存在
	if _, exists := filledData["statement_uuid"]; !exists {
		filledData["statement_uuid"] = uuid.New().String()
	}

	// 添加created_time
	filledData["created_time"] = NowTs()

	// 添加expired_time
	if expiredMinutes > 0 {
		filledData["expired_time"] = NowTs() + int64(expiredMinutes*60)
	} else {
		filledData["expired_time"] = 0
	}

	// 添加template_version
	filledData["template_version"] = version

	return filledData
}

// BuildStatementV2 构建语句
func (s *Statement) BuildStatementV2(businessKey string, data map[string]interface{}, version string, expiredMinutes int) (string, map[string]string, map[string]string, map[string]interface{}, error) {
	// Step 0: 自动填充字段
	filledData := s.FillData(data, version, expiredMinutes)

	// Step 1: 找到模板路径
	currentDir, err := os.Getwd()
	if err != nil {
		return "", nil, nil, nil, fmt.Errorf("error getting current directory: %w", err)
	}

	// 构建模板路径，这里假设与Python代码中的路径结构相似
	templateDir := filepath.Join(currentDir, "json_templates")
	templateFile := fmt.Sprintf("%s_%s.json.j2", businessKey, version)
	fullPath := filepath.Join(templateDir, templateFile)

	// 检查模板文件是否存在
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return "", nil, nil, nil, fmt.Errorf("template file not found: %s", fullPath)
	}

	// Step 2: 初始化Jinja2环境并渲染模板
	// 创建Jinja2实例
	j2, err := jinja2.NewJinja2("python3", 1,
		jinja2.WithGlobals(filledData))
	if err != nil {
		return "", nil, nil, nil, fmt.Errorf("error initializing jinja2: %w", err)
	}
	defer j2.Close()

	// 读取模板文件内容
	templateContent, err := os.ReadFile(fullPath)
	if err != nil {
		return "", nil, nil, nil, fmt.Errorf("error reading template file: %w", err)
	}

	// 转换数据为JSON，以便传递给Jinja2
	// dataJSON, err := json.Marshal(filledData)
	// if err != nil {
	// 	return "", nil, nil, nil, fmt.Errorf("error marshaling data to JSON: %w", err)
	// }

	// Step 3: 渲染模板
	message, err := j2.RenderString(string(templateContent))
	if err != nil {
		return "", nil, nil, nil, fmt.Errorf("error rendering template: %w", err)
	}

	// 获取推送内容
	pushMgr := &PushMsgManager{}
	titleMap, bodyMap := pushMgr.GetPushContent(businessKey, filledData)

	return message, titleMap, bodyMap, filledData, nil
}

// 使用示例
func ExampleUsage() {
	// 创建Statement实例
	s := &Statement{}

	// 获取当前目录
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error getting current directory: %v\n", err)
		return
	}

	// 构建数据文件路径
	dataDir := filepath.Join(currentDir, "example_datas")
	bizKey := "mfa_create_transaction_policy"
	dataFile := fmt.Sprintf("%s.json", bizKey)
	fullPath := filepath.Join(dataDir, dataFile)

	// 读取JSON数据
	dataBytes, err := os.ReadFile(fullPath)
	if err != nil {
		fmt.Printf("Error reading data file: %v\n", err)
		return
	}

	// 解析JSON数据
	var data map[string]interface{}
	if err := json.Unmarshal(dataBytes, &data); err != nil {
		fmt.Printf("Error parsing JSON data: %v\n", err)
		return
	}

	// 调用BuildStatementV2
	message, titleMap, bodyMap, filledData, err := s.BuildStatementV2(bizKey, data, "1.0.0", 30)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	// 使用结果
	fmt.Printf("Message: %s\n", message)
	fmt.Printf("Title Map: %v\n", titleMap)
	fmt.Printf("Body Map: %v\n", bodyMap)
	fmt.Printf("Filled Data: %v\n", filledData)
}
