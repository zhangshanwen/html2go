// Package embed provides embedded component data files
package embed

import (
	"embed"
)

//go:embed data/*.json
var componentData embed.FS

// GetComponentData returns the embedded component data files as [][]byte
func GetComponentData() ([][]byte, error) {
	// List all files in the data directory
	entries, err := componentData.ReadDir("data")
	if err != nil {
		return nil, err
	}

	// Read each file and add its content to the result
	result := make([][]byte, 0, len(entries))
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		content, err := componentData.ReadFile("data/" + entry.Name())
		if err != nil {
			return nil, err
		}

		result = append(result, content)
	}

	return result, nil
}
