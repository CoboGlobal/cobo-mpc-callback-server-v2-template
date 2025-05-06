package verifier

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/kluctl/kluctl/lib/go-jinja2"
)

type Statement struct {
	templateContent string
	version         string
}

func NewStatement(templateContent string, version string) *Statement {
	return &Statement{
		templateContent: templateContent,
		version:         version,
	}
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

func (s *Statement) BuildStatementV2(bizKey string, bizData string, expiredMinutes int) (string, map[string]string, map[string]string, map[string]interface{}, error) {

	var data map[string]interface{}
	if err := json.Unmarshal([]byte(bizData), &data); err != nil {
		fmt.Printf("Error parsing JSON data: %v\n", err)
		return "", nil, nil, nil, fmt.Errorf("error parsing JSON data: %w", err)
	}

	filledData := s.FillData(data, s.version, expiredMinutes)

	j2, err := jinja2.NewJinja2("python3", 1,
		jinja2.WithGlobals(filledData))
	if err != nil {
		return "", nil, nil, nil, fmt.Errorf("error initializing jinja2: %w", err)
	}
	defer j2.Close()

	message, err := j2.RenderString(s.templateContent)
	if err != nil {
		return "", nil, nil, nil, fmt.Errorf("error rendering template: %w", err)
	}

	// 获取推送内容
	pushMgr := &PushMsgManager{}
	titleMap, bodyMap := pushMgr.GetPushContent(bizKey, filledData)

	return message, titleMap, bodyMap, filledData, nil
}

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
