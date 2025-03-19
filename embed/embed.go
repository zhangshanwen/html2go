// Package embed provides embedded component data files
package embed

import (
	"embed"
	"path/filepath"
)

//go:embed data/*.json
var componentData embed.FS

// ComponentData 包含组件数据和类型信息
type ComponentData struct {
	Body []byte // 组件数据内容
	Type string // 组件类型，如 "vuetify" 或 "vuetifyx"
}

// GetComponentData returns the embedded component data files as a slice of ComponentData
func GetComponentData() ([]ComponentData, error) {
	// List all files in the data directory
	entries, err := componentData.ReadDir("data")
	if err != nil {
		return nil, err
	}

	// Read each file and add its content to the result
	result := make([]ComponentData, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		content, err := componentData.ReadFile("data/" + entry.Name())
		if err != nil {
			return nil, err
		}

		// 从文件名获取类型（去掉扩展名）
		compType := filepath.Base(entry.Name())
		compType = compType[:len(compType)-len(filepath.Ext(compType))]

		result = append(result, ComponentData{
			Body: content,
			Type: compType,
		})
	}

	return result, nil
}
