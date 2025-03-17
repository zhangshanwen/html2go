package transform

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"strings"
)

// AttrInfo 表示组件属性的信息
type AttrInfo struct {
	Go     string `json:"go"`     // Go 方法名
	Accept string `json:"accept"` // 接受的值类型
}

// ComponentInfo 表示 Vuetify 组件的信息
type ComponentInfo struct {
	Go     string              `json:"go"`     // Go 结构体名
	Accept string              `json:"accept"` // 接受的子组件类型
	Attrs  map[string]AttrInfo `json:"attrs"`  // 属性映射
}

// ComponentMap 是主映射结构
type ComponentMap map[string]ComponentInfo

// ParseGoFiles 解析多个 Go 文件并提取组件信息
func ParseGoFiles(filePaths []string) (ComponentMap, error) {
	// 创建组件映射
	componentMap := make(ComponentMap)

	// 遍历所有文件
	for _, filePath := range filePaths {
		// 解析 Go 文件
		fset := token.NewFileSet()
		node, err := parser.ParseFile(fset, filePath, nil, parser.ParseComments)
		if err != nil {
			return nil, fmt.Errorf("error parsing %s: %v", filePath, err)
		}

		// 提取组件信息
		ExtractComponentInfo(node, componentMap)
	}

	return componentMap, nil
}

// ParseGoFile 解析单个 Go 文件并提取组件信息 (保留向后兼容性)
func ParseGoFile(filePath string) (ComponentMap, error) {
	return ParseGoFiles([]string{filePath})
}

// ParseGoDir 解析目录中的所有 Go 文件并提取组件信息
func ParseGoDir(dirPath string) (ComponentMap, error) {
	// 查找所有 Go 文件
	matches, err := filepath.Glob(filepath.Join(dirPath, "*.go"))
	if err != nil {
		return nil, fmt.Errorf("error finding Go files in %s: %v", dirPath, err)
	}

	if len(matches) == 0 {
		// If no Go files are found, return an empty component map instead of an error
		return make(ComponentMap), nil
	}

	return ParseGoFiles(matches)
}

// ExtractComponentInfo 从 AST 节点提取组件信息
func ExtractComponentInfo(node *ast.File, componentMap ComponentMap) {
	// 首先查找所有类型定义，确保记录所有构建器类型
	builderTypes := make(map[string]bool)
	for _, decl := range node.Decls {
		if genDecl, ok := decl.(*ast.GenDecl); ok {
			for _, spec := range genDecl.Specs {
				if typeSpec, ok := spec.(*ast.TypeSpec); ok {
					if strings.HasSuffix(typeSpec.Name.Name, "Builder") {
						builderTypes[typeSpec.Name.Name] = true
					}
				}
			}
		}
	}

	// 然后寻找构造函数和方法
	for _, decl := range node.Decls {
		// 查找函数声明（组件构造函数）
		if funcDecl, ok := decl.(*ast.FuncDecl); ok {
			// 处理组件构造函数
			if funcDecl.Recv == nil {
				// 检查这是否是组件构造函数（以 V 开头并返回构建器）
				funcName := funcDecl.Name.Name
				if !strings.HasPrefix(funcName, "V") {
					continue
				}

				// 从函数体中提取标签名
				tagName := ExtractTagName(funcDecl)
				if tagName == "" {
					continue
				}

				// 创建组件信息条目或使用已存在的
				info, exists := componentMap[tagName]
				if !exists {
					info = ComponentInfo{
						Go:     funcName,
						Accept: ExtractAcceptType(funcDecl),
						Attrs:  make(map[string]AttrInfo),
					}
				}

				// 将信息放回映射
				componentMap[tagName] = info
			} else if funcDecl.Recv != nil && len(funcDecl.Recv.List) > 0 {
				// 处理构建器方法
				// 获取接收者类型
				receiverType := ""
				if starExpr, ok := funcDecl.Recv.List[0].Type.(*ast.StarExpr); ok {
					if ident, ok := starExpr.X.(*ast.Ident); ok {
						receiverType = ident.Name
					}
				}

				// 检查接收者类型是否是构建器
				if strings.HasSuffix(receiverType, "Builder") {
					methodName := funcDecl.Name.Name

					// 跳过通用方法
					if isCommonMethod(methodName) {
						continue
					}

					// 将方法名转换为属性名
					attrName := camelToKebab(methodName)

					// 提取参数类型
					paramType := extractParamType(funcDecl)

					// 首先检查这个构建器是否已经有关联的组件
					found := false
					for tagName, info := range componentMap {
						// 检查构建器名称是否匹配（VBtn -> VBtnBuilder）
						builderName := info.Go + "Builder"
						if builderName == receiverType {
							// 添加到属性
							info.Attrs[attrName] = AttrInfo{
								Go:     methodName,
								Accept: paramType,
							}

							// 更新组件信息
							componentMap[tagName] = info
							found = true
							break
						}
					}

					// 如果没有找到对应的组件，但我们知道这是一个构建器，
					// 则存储该方法以备后用（可能组件构造函数在另一个文件中）
					if !found && builderTypes[receiverType] {
						// 尝试从构建器名称推导组件名
						componentName := strings.TrimSuffix(receiverType, "Builder")
						// 尝试从组件名称推导标签名
						tagName := camelToKebab(componentName)

						// 检查是否已存在该组件
						info, exists := componentMap[tagName]
						if !exists {
							// 创建新的组件条目
							info = ComponentInfo{
								Go:     componentName,
								Accept: "none", // 暂时设置为无，可能在后续文件解析中更新
								Attrs:  make(map[string]AttrInfo),
							}
						}

						// 添加属性
						info.Attrs[attrName] = AttrInfo{
							Go:     methodName,
							Accept: paramType,
						}

						// 更新组件映射
						componentMap[tagName] = info
					}
				}
			}
		}
	}

	// 查找所有构建器类型的方法
	for receiverType := range builderTypes {
		componentName := strings.TrimSuffix(receiverType, "Builder")
		findBuilderMethods(node, receiverType, componentName, componentMap)
	}
}

// ExtractTagName 从组件构造函数中提取标签名
func ExtractTagName(funcDecl *ast.FuncDecl) string {
	// 查找标签赋值
	if funcDecl.Body == nil {
		return ""
	}

	// 遍历函数体中的所有语句
	for _, stmt := range funcDecl.Body.List {
		// 查找赋值语句
		assignStmt, ok := stmt.(*ast.AssignStmt)
		if !ok {
			continue
		}

		// 遍历赋值语句的右侧表达式
		for _, rhs := range assignStmt.Rhs {
			// 处理一元表达式 (&VAlertBuilder{...})
			if unaryExpr, ok := rhs.(*ast.UnaryExpr); ok {
				if compLit, ok := unaryExpr.X.(*ast.CompositeLit); ok {
					// 遍历复合字面量的所有元素
					for _, elt := range compLit.Elts {
						// 检查是否是键值对
						kvExpr, ok := elt.(*ast.KeyValueExpr)
						if !ok {
							continue
						}

						// 检查键是否是 "tag"
						key, ok := kvExpr.Key.(*ast.Ident)
						if !ok || key.Name != "tag" {
							continue
						}

						// 递归查找 h.Tag("tag-name") 调用
						tagName := findTagCall(kvExpr.Value)
						if tagName != "" {
							return tagName
						}
					}
				}
				continue
			}

			// 检查是否是复合字面量
			compLit, ok := rhs.(*ast.CompositeLit)
			if !ok {
				continue
			}

			// 遍历复合字面量的所有元素
			for _, elt := range compLit.Elts {
				// 检查是否是键值对
				kvExpr, ok := elt.(*ast.KeyValueExpr)
				if !ok {
					continue
				}

				// 检查键是否是 "tag"
				key, ok := kvExpr.Key.(*ast.Ident)
				if !ok || key.Name != "tag" {
					continue
				}

				// 递归查找 h.Tag("tag-name") 调用
				tagName := findTagCall(kvExpr.Value)
				if tagName != "" {
					return tagName
				}
			}
		}
	}

	return ""
}

// findTagCall 递归查找 h.Tag("tag-name") 调用
func findTagCall(expr ast.Expr) string {
	switch e := expr.(type) {
	case *ast.CallExpr:
		// 检查是否是 h.Tag("tag-name") 调用
		sel, ok := e.Fun.(*ast.SelectorExpr)
		if ok {
			ident, ok := sel.X.(*ast.Ident)
			if ok && ident.Name == "h" && sel.Sel.Name == "Tag" {
				if len(e.Args) > 0 {
					lit, ok := e.Args[0].(*ast.BasicLit)
					if ok && lit.Kind == token.STRING {
						// 移除标签名中的引号
						return strings.Trim(lit.Value, "\"'")
					}
				}
			}
		}

		// 如果不是 h.Tag 调用，则检查函数表达式
		// 这可能是链式调用的一部分，如 h.Tag("v-alert").Children(...)
		if sel, ok := e.Fun.(*ast.SelectorExpr); ok {
			if tagName := findTagCall(sel.X); tagName != "" {
				return tagName
			}
		}

	case *ast.SelectorExpr:
		// 检查选择器的左侧
		return findTagCall(e.X)
	}

	return ""
}

// ExtractAcceptType 提取组件接受的子组件类型
func ExtractAcceptType(funcDecl *ast.FuncDecl) string {
	if funcDecl.Type.Params != nil && len(funcDecl.Type.Params.List) > 0 {
		param := funcDecl.Type.Params.List[0]

		if ellipsis, ok := param.Type.(*ast.Ellipsis); ok {
			if sel, ok := ellipsis.Elt.(*ast.SelectorExpr); ok {
				if x, ok := sel.X.(*ast.Ident); ok && x.Name == "h" && sel.Sel.Name == "HTMLComponent" {
					return "...h.HTMLComponent"
				}
			}
		}
	}
	return "none"
}

// findBuilderMethods 查找构建器结构的所有方法
func findBuilderMethods(node *ast.File, builderName string, componentName string, componentMap ComponentMap) {
	// 寻找对应的组件
	var componentInfo *ComponentInfo
	var tagName string

	// 首先尝试通过组件名查找
	for t, info := range componentMap {
		if info.Go == componentName {
			componentInfo = &info
			tagName = t
			break
		}
	}

	// 如果没有找到，尝试通过标签名推导
	if componentInfo == nil {
		tagName = camelToKebab(componentName)
		// 默认创建一个新组件
		componentInfo = &ComponentInfo{
			Go:     componentName,
			Accept: "none", // 可能会在后续文件解析中更新
			Attrs:  make(map[string]AttrInfo),
		}
		componentMap[tagName] = *componentInfo
	}

	// 查找该构建器的所有方法
	for _, decl := range node.Decls {
		if funcDecl, ok := decl.(*ast.FuncDecl); ok {
			// 检查这是否是构建器结构的方法
			if funcDecl.Recv != nil && len(funcDecl.Recv.List) > 0 {
				if starExpr, ok := funcDecl.Recv.List[0].Type.(*ast.StarExpr); ok {
					if ident, ok := starExpr.X.(*ast.Ident); ok && ident.Name == builderName {
						// 这是构建器结构的方法
						methodName := funcDecl.Name.Name

						// 跳过通用方法
						if isCommonMethod(methodName) {
							continue
						}

						// 将方法名转换为属性名（camelCase 到 kebab-case）
						attrName := camelToKebab(methodName)

						// 提取参数类型
						paramType := extractParamType(funcDecl)

						// 添加到属性
						info := componentMap[tagName]
						info.Attrs[attrName] = AttrInfo{
							Go:     methodName,
							Accept: paramType,
						}
						componentMap[tagName] = info
					}
				}
			}
		}
	}
}

// isCommonMethod 检查方法是否是所有构建器通用的
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
		"AttrIf":          true, // Added AttrIf as a common method
	}

	return commonMethods[name]
}

// extractParamType 提取方法的第一个参数的类型
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

		// 如果无法确定类型，默认为字符串
		return "string"
	}

	return "string"
}

// camelToKebab 将 camelCase 转换为 kebab-case
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
