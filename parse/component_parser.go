package parse

import (
	"encoding/json"
	"fmt"

	"github.com/sunfmin/html2go/embed"
)

// AttributeDefinition 表示组件属性的定义
type AttributeDefinition struct {
	Go     string `json:"go"`     // Go属性名称
	Accept string `json:"accept"` // 该属性接受的类型
}

// ComponentDefinition 表示一个组件的定义
type ComponentDefinition struct {
	Go     string                         `json:"go"`     // Go组件名称
	Accept string                         `json:"accept"` // 该组件接受的内容
	Attrs  map[string]AttributeDefinition `json:"attrs"`  // 该组件的属性
}

// ParseComponentData 解析组件数据并返回组件定义的映射
func ParseComponentData() (map[string]ComponentDefinition, error) {
	// 获取嵌入的组件数据
	dataBytes, err := embed.GetComponentData()
	if err != nil {
		return nil, fmt.Errorf("获取组件数据失败: %w", err)
	}

	// 创建返回的映射
	result := make(map[string]ComponentDefinition)

	// 处理每个JSON文件
	for _, data := range dataBytes {
		// 解析JSON数据到临时映射
		var tempComponents map[string]ComponentDefinition
		if err := json.Unmarshal(data, &tempComponents); err != nil {
			return nil, fmt.Errorf("解析组件数据失败: %w", err)
		}

		// 合并到结果映射中
		for k, v := range tempComponents {
			result[k] = v
		}
	}

	return result, nil
}
