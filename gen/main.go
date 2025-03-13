package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"go/ast"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/sunfmin/html2go/gen/transform" // Import the transform package
)

// AttrInfo represents information about a component attribute
type AttrInfo struct {
	Go     string `json:"go"`     // Go method name
	Accept string `json:"accept"` // Type of value accepted
}

// ComponentInfo represents information about a Vuetify component
type ComponentInfo struct {
	Go     string              `json:"go"`     // Go struct name
	Accept string              `json:"accept"` // Type of children accepted
	Attrs  map[string]AttrInfo `json:"attrs"`  // Map of attributes
}

// ComponentMap is the main mapping structure
type ComponentMap map[string]ComponentInfo

func main() {
	// Define command-line flags
	vuetifyDirFlag := flag.String("dir", "../x/ui/vuetify", "Path to the vuetify directory")
	outputPathFlag := flag.String("output", "vuetify_components_map.json", "Path for the output JSON file")
	flag.Parse()

	// Use the provided directory path and trim any whitespace
	vuetifyDir := strings.TrimSpace(*vuetifyDirFlag)
	outputPath := strings.TrimSpace(*outputPathFlag)

	// Create the component map
	componentMap := make(transform.ComponentMap)

	// Process all Go files in the vuetify directory
	err := processVuetifyDirectory(vuetifyDir, componentMap)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Convert the map to JSON
	jsonData, err := json.MarshalIndent(componentMap, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error marshaling to JSON: %v\n", err)
		os.Exit(1)
	}

	// Replace Unicode escape sequences with actual characters
	jsonString := string(jsonData)
	jsonString = strings.ReplaceAll(jsonString, "\\u003c", "<")
	jsonString = strings.ReplaceAll(jsonString, "\\u003e", ">")
	jsonString = strings.ReplaceAll(jsonString, "\\u0026", "&")

	// Write the JSON to a file
	err = ioutil.WriteFile(outputPath, []byte(jsonString), 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing to file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Component map generated successfully: %s\n", outputPath)
}

// processVuetifyDirectory processes all Go files in the Vuetify directory
func processVuetifyDirectory(dirPath string, componentMap transform.ComponentMap) error {
	// Walk through all files and directories in the Vuetify directory
	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories and non-Go files
		if info.IsDir() || !strings.HasSuffix(info.Name(), ".go") {
			return nil
		}

		// Parse the Go file and extract component information
		fileComponents, err := transform.ParseGoFile(path)
		if err != nil {
			return fmt.Errorf("error processing %s: %v", path, err)
		}

		// Merge component information
		for k, v := range fileComponents {
			componentMap[k] = v
		}

		return nil
	})

	return err
}

// extractComponentInfo extracts component information from an AST node
func extractComponentInfo(node *ast.File, componentMap ComponentMap) {
	for _, decl := range node.Decls {
		// Look for function declarations (component constructors)
		if funcDecl, ok := decl.(*ast.FuncDecl); ok {
			// Skip methods (they have receivers)
			if funcDecl.Recv != nil {
				continue
			}

			// Check if this is a component constructor (starts with V and returns a builder)
			funcName := funcDecl.Name.Name
			if !strings.HasPrefix(funcName, "V") {
				continue
			}

			// Extract the tag name from the function body
			tagName := extractTagName(funcDecl)
			if tagName == "" {
				continue
			}

			// Create a component info entry
			componentMap[tagName] = ComponentInfo{
				Go:     funcName,
				Accept: extractAcceptType(funcDecl),
				Attrs:  make(map[string]AttrInfo),
			}

			// Find the builder struct and its methods
			builderName := funcName + "Builder"
			for _, d := range node.Decls {
				if genDecl, ok := d.(*ast.GenDecl); ok {
					for _, spec := range genDecl.Specs {
						if typeSpec, ok := spec.(*ast.TypeSpec); ok && typeSpec.Name.Name == builderName {
							// 创建一个临时变量来存储组件信息
							info := componentMap[tagName]
							// 找到构建器结构，现在查找其方法
							findBuilderMethods(node, builderName, &info)
							// 将修改后的信息放回映射
							componentMap[tagName] = info
							break
						}
					}
				}
			}
		}
	}
}

// extractTagName extracts the tag name from a component constructor function
func extractTagName(funcDecl *ast.FuncDecl) string {
	// Look for the tag assignment in the function body
	if funcDecl.Body != nil {
		for _, stmt := range funcDecl.Body.List {
			if assignStmt, ok := stmt.(*ast.AssignStmt); ok {
				for _, rhs := range assignStmt.Rhs {
					if callExpr, ok := rhs.(*ast.CallExpr); ok {
						if selectorExpr, ok := callExpr.Fun.(*ast.SelectorExpr); ok {
							if ident, ok := selectorExpr.X.(*ast.Ident); ok && ident.Name == "h" && selectorExpr.Sel.Name == "Tag" {
								if len(callExpr.Args) > 0 {
									if basicLit, ok := callExpr.Args[0].(*ast.BasicLit); ok {
										// Remove quotes from the tag name
										tagName := strings.Trim(basicLit.Value, "\"")
										return tagName
									}
								}
							}
						}
					}
				}
			}
		}
	}
	return ""
}

// extractAcceptType extracts the type of children accepted by a component
func extractAcceptType(funcDecl *ast.FuncDecl) string {
	if funcDecl.Type.Params != nil && len(funcDecl.Type.Params.List) > 0 {
		param := funcDecl.Type.Params.List[0]

		if ellipsis, ok := param.Type.(*ast.Ellipsis); ok {
			if sel, ok := ellipsis.Elt.(*ast.SelectorExpr); ok {
				if x, ok := sel.X.(*ast.Ident); ok && x.Name == "h" && sel.Sel.Name == "HTMLComponent" {
					return "array<HTMLComponent>"
				}
			}
		}
	}
	return "none"
}

// findBuilderMethods finds all methods of a builder struct
func findBuilderMethods(node *ast.File, builderName string, componentInfo *ComponentInfo) {
	for _, decl := range node.Decls {
		if funcDecl, ok := decl.(*ast.FuncDecl); ok {
			// Check if this is a method of the builder struct
			if funcDecl.Recv != nil && len(funcDecl.Recv.List) > 0 {
				if starExpr, ok := funcDecl.Recv.List[0].Type.(*ast.StarExpr); ok {
					if ident, ok := starExpr.X.(*ast.Ident); ok && ident.Name == builderName {
						// This is a method of the builder struct
						methodName := funcDecl.Name.Name

						// Skip common methods
						if isCommonMethod(methodName) {
							continue
						}

						// Convert method name to attribute name (camelCase to kebab-case)
						attrName := camelToKebab(methodName)

						// Extract parameter type
						paramType := extractParamType(funcDecl)

						// Add to attributes
						componentInfo.Attrs[attrName] = AttrInfo{
							Go:     methodName,
							Accept: paramType,
						}
					}
				}
			}
		}
	}
}

// isCommonMethod checks if a method is common to all builders
func isCommonMethod(name string) bool {
	commonMethods := map[string]bool{
		"SetAttr":         true,
		"Attr":            true,
		"Children":        true,
		"AppendChildren":  true,
		"PrependChildren": true,
		"Class":           true,
		"ClassIf":         true,
		"On":              true,
		"Bind":            true,
		"MarshalHTML":     true,
	}

	return commonMethods[name]
}

// extractParamType extracts the type of the first parameter of a method
func extractParamType(funcDecl *ast.FuncDecl) string {
	if funcDecl.Type.Params != nil && len(funcDecl.Type.Params.List) > 0 {
		param := funcDecl.Type.Params.List[0]

		switch t := param.Type.(type) {
		case *ast.Ident:
			return t.Name
		case *ast.SelectorExpr:
			if x, ok := t.X.(*ast.Ident); ok {
				return x.Name + "." + t.Sel.Name
			}
		case *ast.ArrayType:
			if elt, ok := t.Elt.(*ast.SelectorExpr); ok {
				if x, ok := elt.X.(*ast.Ident); ok {
					return "[]" + x.Name + "." + elt.Sel.Name
				}
			} else if elt, ok := t.Elt.(*ast.Ident); ok {
				return "[]" + elt.Name
			}
		case *ast.InterfaceType:
			return "interface{}"
		}

		// Default to string if we can't determine the type
		return "string"
	}

	return "string"
}

// camelToKebab converts camelCase to kebab-case
func camelToKebab(s string) string {
	var result strings.Builder

	for i, r := range s {
		if i > 0 && 'A' <= r && r <= 'Z' {
			result.WriteRune('-')
			result.WriteRune(r - 'A' + 'a')
		} else {
			result.WriteRune(r)
		}
	}

	return strings.ToLower(result.String())
}
