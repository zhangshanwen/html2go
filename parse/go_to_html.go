package parse

import (
	"fmt"
	"strings"
)

// GenerateGoHTML 将Go代码转换为HTML
func GenerateGoHTML(goCode string, formatOutput bool, indentSize int) (string, error) {
	// 解析Go代码，生成HTML节点树
	rootNode, err := ExtractHTMLGoNodes(goCode)
	if err != nil {
		return "", fmt.Errorf("解析Go代码失败: %v", err)
	}

	// 生成HTML输出
	html := GenerateHTML(rootNode, formatOutput, indentSize)
	return html, nil
}

// NormalizeHTML 标准化HTML字符串，用于测试比较
func NormalizeHTML(html string) string {
	// 删除所有空白字符
	html = strings.ReplaceAll(html, "\n", "")
	html = strings.ReplaceAll(html, "\t", "")
	html = strings.ReplaceAll(html, "\r", "")

	// 标准化空格 - 删除连续空格，替换为单个空格
	html = strings.Join(strings.Fields(html), " ")

	// 标准化标签之间的空格
	html = strings.ReplaceAll(html, "> <", "><")

	// 删除开始和结束处的空格
	html = strings.TrimSpace(html)

	return html
}

// ExtractPackageName 从Go代码中提取包名
func ExtractPackageName(goCode string) string {
	// 查找 package 关键字
	packageIndex := strings.Index(goCode, "package ")
	if packageIndex == -1 {
		return ""
	}

	// 从 package 关键字之后截取一行
	startPos := packageIndex + len("package ")
	endPos := strings.IndexAny(goCode[startPos:], "\n\r;")
	if endPos == -1 {
		// 如果没有找到换行符，使用整个剩余部分
		return strings.TrimSpace(goCode[startPos:])
	}

	// 提取包名并去除空格
	return strings.TrimSpace(goCode[startPos : startPos+endPos])
}

// ParseSingleElement 解析单个HTML元素的Go代码
func ParseSingleElement(goCode string) (*HTMLNode, error) {
	// 如果代码不是表达式，尝试包装为表达式
	if !strings.Contains(goCode, "return") && !strings.HasPrefix(strings.TrimSpace(goCode), "func") {
		if !strings.HasSuffix(strings.TrimSpace(goCode), ")") && !strings.Contains(goCode, "}") {
			goCode = goCode + "()"
		}
	}

	// 使用通用解析器提取节点
	return ExtractHTMLGoNodes(goCode)
}

// ProcessElementAttributes 处理元素的属性
func ProcessElementAttributes(node *HTMLNode, attrs map[string]string) {
	for name, value := range attrs {
		node.SetAttr(name, value)
	}
}

// FixIncompleteGoCode 修复不完整的Go代码
func FixIncompleteGoCode(goCode string) string {
	goCode = strings.TrimSpace(goCode)

	// 检查是否是完整的函数或表达式
	if strings.HasPrefix(goCode, "func") || strings.HasPrefix(goCode, "return") {
		return goCode
	}

	// 检查是否缺少分号或换行符
	lastChar := ""
	if len(goCode) > 0 {
		lastChar = goCode[len(goCode)-1:]
	}

	if lastChar != ";" && lastChar != "}" && lastChar != "\n" {
		goCode += "\n"
	}

	return goCode
}
