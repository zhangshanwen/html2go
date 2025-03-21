package parse

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"io/ioutil"
	"regexp"
	"strconv"
	"strings"
)

// ParseGoCode 解析Go代码并返回AST
func ParseGoCode(r io.Reader) (*ast.File, error) {
	// 读取输入
	src, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, fmt.Errorf("读取输入失败: %w", err)
	}

	// 解析Go代码
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "", string(src), 0)
	if err != nil {
		return nil, fmt.Errorf("解析Go代码失败: %w", err)
	}

	return file, nil
}

// ExtractHTMLGoNodes 从Go代码中提取HTML节点
func ExtractHTMLGoNodes(goCode string) (*HTMLNode, error) {
	// 尝试将代码包装成有效的Go源代码
	goCode = preprocessGoCode(goCode)

	// 解析Go代码
	fset := token.NewFileSet()
	expr, err := parser.ParseExpr(goCode)
	if err != nil {
		// 尝试将代码包装在函数中再解析
		wrappedCode := fmt.Sprintf("func() {\n%s\n}", goCode)
		f, err := parser.ParseFile(fset, "", wrappedCode, 0)
		if err != nil {
			return nil, fmt.Errorf("解析Go代码失败: %v", err)
		}

		// 提取函数体
		fn, ok := f.Decls[0].(*ast.FuncDecl)
		if !ok || fn.Body == nil {
			return nil, fmt.Errorf("无法提取函数体")
		}

		// 最外层默认为body标签
		rootNode := NewElementNode("body")

		// 解析函数体中的语句
		for _, stmt := range fn.Body.List {
			// 处理表达式语句（可能是函数调用）
			if exprStmt, ok := stmt.(*ast.ExprStmt); ok {
				node, err := parseExpr(exprStmt.X)
				if err != nil {
					continue
				}
				if node != nil {
					rootNode.AddChild(node)
				}
			}
		}

		return rootNode, nil
	}

	// 解析直接的表达式
	node, err := parseExpr(expr)
	if err != nil {
		return nil, err
	}

	return node, nil
}

// preprocessGoCode 预处理Go代码，处理常见的不完整表达式问题
func preprocessGoCode(code string) string {
	code = strings.TrimSpace(code)

	// 如果代码为空，返回空字符串
	if code == "" {
		return ""
	}

	// 检查是否包含完整的函数调用
	if !strings.Contains(code, "(") {
		code = code + "()"
	}

	// 检查是否需要包裹在函数调用中
	if !strings.HasPrefix(code, "func") && !strings.HasPrefix(code, "package") {
		lines := strings.Split(code, "\n")
		// 检查每一行是否是独立的表达式
		for i, line := range lines {
			line = strings.TrimSpace(line)
			if line == "" || strings.HasPrefix(line, "//") {
				continue
			}

			// 修复单行表达式
			if !strings.Contains(line, "(") && !strings.HasSuffix(line, ";") {
				lines[i] = line + "()"
			}
		}
		code = strings.Join(lines, "\n")
	}

	// 检查括号是否匹配
	openCount := strings.Count(code, "(")
	closeCount := strings.Count(code, ")")
	if openCount > closeCount {
		code = code + strings.Repeat(")", openCount-closeCount)
	}

	return code
}

// parseExpr 解析Go表达式
func parseExpr(expr ast.Expr) (*HTMLNode, error) {
	switch e := expr.(type) {
	case *ast.CallExpr:
		return parseCallExpr(e)
	case *ast.SelectorExpr:
		// 处理选择器表达式，例如 pkg.Div
		return parseSelectorExpr(e)
	}
	return nil, fmt.Errorf("不支持的表达式类型: %T", expr)
}

// parseSelectorExpr 解析选择器表达式
func parseSelectorExpr(expr *ast.SelectorExpr) (*HTMLNode, error) {
	// 检查是否有效的选择器表达式
	_, ok := expr.X.(*ast.Ident)
	if !ok {
		return nil, fmt.Errorf("无效的选择器表达式")
	}

	// 创建一个假的调用表达式并解析
	callExpr := &ast.CallExpr{
		Fun: expr,
	}
	return parseCallExpr(callExpr)
}

// 辅助函数：处理方法链
func processMethodChain(node *HTMLNode, methodName string, args []ast.Expr) {
	// 处理Children方法
	if methodName == "Children" {
		for _, arg := range args {
			childNode, err := parseExpr(arg)
			if err == nil && childNode != nil {
				node.AddChild(childNode)
			}
		}
		return
	}

	// 处理Attr方法
	if methodName == "Attr" && len(args) >= 2 {
		attrName, ok1 := extractStringLiteral(args[0])
		attrValue, ok2 := extractStringLiteral(args[1])
		if ok1 && ok2 {
			node.SetAttr(attrName, attrValue)
			return
		}
	}

	// 处理其他方法作为属性
	attrName := strings.ToLower(methodName)
	attrValue := extractArgs(args)

	// 处理布尔属性（无参数或布尔参数）
	if len(args) == 0 || (len(args) == 1 && isBooleanExpr(args[0])) {
		// 布尔值属性
		value := true
		if len(args) == 1 {
			if boolValue, ok := extractValue(args[0]).(bool); ok {
				value = boolValue
			}
		}

		if value {
			node.SetAttr(attrName, "")
		}
		return
	}

	// 普通属性
	node.SetAttr(attrName, fmt.Sprintf("%v", attrValue))
}

// 判断表达式是否为布尔表达式
func isBooleanExpr(expr ast.Expr) bool {
	if ident, ok := expr.(*ast.Ident); ok {
		return ident.Name == "true" || ident.Name == "false"
	}
	return false
}

// parseCallExpr 解析函数调用表达式
func parseCallExpr(callExpr *ast.CallExpr) (*HTMLNode, error) {
	// 获取函数名
	var funcName string

	switch fun := callExpr.Fun.(type) {
	case *ast.Ident:
		// 直接函数调用，如 Div(...)
		funcName = fun.Name
	case *ast.SelectorExpr:
		// 带包名的函数调用，如 pkg.Div(...)
		funcName = fun.Sel.Name

		// 处理链式调用，例如 Div().Children(...)
		if callExprInner, ok := fun.X.(*ast.CallExpr); ok {
			innerNode, err := parseCallExpr(callExprInner)
			if err != nil {
				return nil, err
			}

			if innerNode != nil {
				// 处理方法调用
				processMethodChain(innerNode, funcName, callExpr.Args)
				return innerNode, nil
			}
		}
	case *ast.CallExpr:
		// 处理链式调用的内层函数，如 Div().Class("...")
		innerNode, err := parseCallExpr(fun)
		if err != nil {
			return nil, err
		}

		// 处理当前层的方法调用
		if innerNode != nil {
			methodName := ""
			if selExpr, ok := fun.Fun.(*ast.SelectorExpr); ok {
				methodName = selExpr.Sel.Name
			}

			if methodName != "" {
				processMethodChain(innerNode, methodName, callExpr.Args)
			}
		}

		return innerNode, nil
	default:
		return nil, fmt.Errorf("不支持的函数类型: %T", fun)
	}

	// 处理Text函数特殊情况
	if funcName == "Text" && len(callExpr.Args) > 0 {
		text, _ := extractStringLiteral(callExpr.Args[0])
		return NewTextNode(text), nil
	}

	// 处理Tag函数特殊情况 - 用于创建自定义组件
	if funcName == "Tag" && len(callExpr.Args) > 0 {
		tagName, ok := extractStringLiteral(callExpr.Args[0])
		if ok {
			return NewElementNode(tagName), nil
		}
	}

	// 使用反向映射解析HTML标签
	rm := NewReverseMapping()
	tagName := rm.ResolveHTMLTag(funcName)

	// 处理特殊的未映射组件命名
	if strings.HasPrefix(funcName, "V") && !strings.HasPrefix(tagName, "v-") {
		// 处理Vuetify组件
		tagName = formatVuetifyComponentName(funcName)
	} else if strings.HasPrefix(funcName, "VX") && !strings.HasPrefix(tagName, "vx-") {
		// 处理VuetifyX组件
		tagName = formatVuetifyXComponentName(funcName)
	}

	// 创建HTML节点
	node := NewElementNode(tagName)

	// 处理参数
	for _, arg := range callExpr.Args {
		// 检查是否为文本内容
		if strLit, ok := arg.(*ast.BasicLit); ok && strLit.Kind == token.STRING {
			textValue, _ := strconv.Unquote(strLit.Value)
			node.AddChild(NewTextNode(textValue))
			continue
		}

		// 检查是否为子元素
		if childCallExpr, ok := arg.(*ast.CallExpr); ok {
			childNode, err := parseCallExpr(childCallExpr)
			if err == nil && childNode != nil {
				node.AddChild(childNode)
			}
			continue
		}

		// 检查是否为selExpr
		if childSelExpr, ok := arg.(*ast.SelectorExpr); ok {
			childCallExpr := &ast.CallExpr{
				Fun: childSelExpr,
			}
			childNode, err := parseCallExpr(childCallExpr)
			if err == nil && childNode != nil {
				node.AddChild(childNode)
			}
			continue
		}
	}

	return node, nil
}

// 格式化Vuetify组件名称
func formatVuetifyComponentName(name string) string {
	if !strings.HasPrefix(name, "V") {
		return name
	}

	// 移除开头的"V"并转换为kebab-case
	name = "v-" + camelToKebab(name[1:])
	return name
}

// 格式化VuetifyX组件名称
func formatVuetifyXComponentName(name string) string {
	if !strings.HasPrefix(name, "VX") {
		return name
	}

	// VuetifyX组件应该是v-xbtn这样的格式
	return "v-x" + camelToKebab(name[2:])
}

// 将驼峰命名转换为kebab-case
func camelToKebab(s string) string {
	re := regexp.MustCompile(`([a-z0-9])([A-Z])`)
	kebab := re.ReplaceAllString(s, "${1}-${2}")
	return strings.ToLower(kebab)
}

// extractArgs 提取参数值
func extractArgs(args []ast.Expr) interface{} {
	if len(args) == 0 {
		return nil
	}

	if len(args) == 1 {
		return extractValue(args[0])
	}

	// 多个参数组成的参数列表
	values := make([]interface{}, len(args))
	for i, arg := range args {
		values[i] = extractValue(arg)
	}
	return values
}

// extractValue 提取表达式的值
func extractValue(expr ast.Expr) interface{} {
	switch e := expr.(type) {
	case *ast.BasicLit:
		// 基本字面量
		switch e.Kind {
		case token.STRING:
			value, _ := strconv.Unquote(e.Value)
			return value
		case token.INT:
			value, _ := strconv.Atoi(e.Value)
			return value
		case token.FLOAT:
			value, _ := strconv.ParseFloat(e.Value, 64)
			return value
		}
	case *ast.Ident:
		// 标识符，如变量名或布尔值
		if e.Name == "true" {
			return true
		} else if e.Name == "false" {
			return false
		}
		return e.Name
	case *ast.CompositeLit:
		// 复合字面量，如数组或结构体
		values := make([]interface{}, len(e.Elts))
		for i, elt := range e.Elts {
			values[i] = extractValue(elt)
		}
		return values
	case *ast.CallExpr:
		// 函数调用表达式
		node, err := parseCallExpr(e)
		if err == nil && node != nil {
			// 如果是文本节点，返回文本内容
			if node.Type == "text" {
				return node.TextValue
			}
			// 否则返回节点的HTML表示
			return nodeToString(node)
		}
	}
	return nil
}

// nodeToString 将节点转换为字符串表示
func nodeToString(node *HTMLNode) string {
	if node == nil {
		return ""
	}

	// 简单实现，实际上应该递归处理子节点
	if node.Type == "text" {
		return node.TextValue
	}

	// 返回元素节点的标签名称
	return "<" + node.TagName + ">"
}

// extractStringLiteral 提取字符串字面量
func extractStringLiteral(expr ast.Expr) (string, bool) {
	switch e := expr.(type) {
	case *ast.BasicLit:
		if e.Kind == token.STRING {
			value, err := strconv.Unquote(e.Value)
			if err != nil {
				return "", false
			}
			return value, true
		}
	case *ast.CallExpr:
		// 如果是Text函数调用
		if ident, ok := e.Fun.(*ast.Ident); ok && ident.Name == "Text" && len(e.Args) > 0 {
			return extractStringLiteral(e.Args[0])
		} else if selExpr, ok := e.Fun.(*ast.SelectorExpr); ok {
			if selExpr.Sel.Name == "Text" && len(e.Args) > 0 {
				return extractStringLiteral(e.Args[0])
			}
		}

		// 尝试解析为节点
		node, err := parseCallExpr(e)
		if err == nil && node != nil && node.Type == "text" {
			return node.TextValue, true
		}
	}
	return "", false
}
