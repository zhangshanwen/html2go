package parse

import (
	"regexp"
	"strings"
)

// HTMLNode 表示HTML节点结构
type HTMLNode struct {
	Type      string            // 节点类型 "element" 或 "text"
	TagName   string            // 元素标签名
	Attrs     map[string]string // 属性映射
	Children  []*HTMLNode       // 子节点
	TextValue string            // 文本节点的值
}

// NewElementNode 创建新的元素节点
func NewElementNode(tagName string) *HTMLNode {
	return &HTMLNode{
		Type:     "element",
		TagName:  tagName,
		Attrs:    make(map[string]string),
		Children: make([]*HTMLNode, 0),
	}
}

// NewTextNode 创建新的文本节点
func NewTextNode(text string) *HTMLNode {
	return &HTMLNode{
		Type:      "text",
		TextValue: text,
	}
}

// AddChild 添加子节点
func (n *HTMLNode) AddChild(child *HTMLNode) {
	if n.Type == "element" {
		n.Children = append(n.Children, child)
	}
}

// SetAttr 设置属性
func (n *HTMLNode) SetAttr(name, value string) {
	if n.Type == "element" {
		n.Attrs[name] = value
	}
}

// ToHTML 将节点转换为HTML字符串
func (n *HTMLNode) ToHTML(formatOutput bool, indentSize int) string {
	var sb strings.Builder
	n.writeHTML(&sb, formatOutput, indentSize, 0)
	return sb.String()
}

// writeHTML 递归生成HTML
func (n *HTMLNode) writeHTML(sb *strings.Builder, formatOutput bool, indentSize int, level int) {
	// 处理文本节点
	if n.Type == "text" {
		sb.WriteString(n.TextValue)
		return
	}

	// 如果是特殊标记的文本节点，直接输出文本值并返回
	if n.TagName == "#text" {
		sb.WriteString(n.TextValue)
		return
	}

	// body 标签特殊处理 - 只输出其子元素
	if n.TagName == "body" && len(n.Children) > 0 {
		for _, child := range n.Children {
			child.writeHTML(sb, formatOutput, indentSize, level)
		}
		return
	}

	indent := ""
	newline := ""
	if formatOutput {
		indent = strings.Repeat(" ", level*indentSize)
		newline = "\n"
	}

	// 开始标签
	sb.WriteString(indent)
	sb.WriteString("<")
	sb.WriteString(n.TagName)

	// 添加属性
	for name, value := range n.Attrs {
		sb.WriteString(" ")
		sb.WriteString(name)
		// 布尔属性特殊处理
		if value == "" {
			continue
		}
		sb.WriteString("=\"")
		sb.WriteString(value)
		sb.WriteString("\"")
	}

	// 自闭合或常规闭合标签
	if len(n.Children) == 0 {
		// 自闭合元素列表
		selfClosingTags := map[string]bool{
			"area": true, "base": true, "br": true, "col": true, "embed": true,
			"hr": true, "img": true, "input": true, "link": true, "meta": true,
			"param": true, "source": true, "track": true, "wbr": true,
		}

		if selfClosingTags[n.TagName] {
			sb.WriteString(" />")
			sb.WriteString(newline)
		} else {
			sb.WriteString(">")
			sb.WriteString("</")
			sb.WriteString(n.TagName)
			sb.WriteString(">")
			sb.WriteString(newline)
		}
	} else {
		sb.WriteString(">")

		// 输出子节点
		childrenAreTextOnly := true
		for _, child := range n.Children {
			if child.Type != "text" && child.TagName != "#text" {
				childrenAreTextOnly = false
				break
			}
		}

		if !childrenAreTextOnly && formatOutput {
			sb.WriteString(newline)
		}

		for _, child := range n.Children {
			if !childrenAreTextOnly {
				child.writeHTML(sb, formatOutput, indentSize, level+1)
			} else {
				child.writeHTML(sb, false, indentSize, 0)
			}
		}

		// 闭合标签
		if !childrenAreTextOnly && formatOutput {
			sb.WriteString(indent)
		}
		sb.WriteString("</")
		sb.WriteString(n.TagName)
		sb.WriteString(">")
		sb.WriteString(newline)
	}
}

// GenerateHTML 从节点树生成格式化的HTML
func GenerateHTML(root *HTMLNode, formatOutput bool, indentSize int) string {
	// 如果根节点是body类型，直接处理其子节点
	if root.TagName == "body" {
		var sb strings.Builder
		for _, child := range root.Children {
			sb.WriteString(child.ToHTML(formatOutput, indentSize))
		}
		return sb.String()
	}
	return root.ToHTML(formatOutput, indentSize)
}

// EscapeHTML 对HTML特殊字符进行转义
func EscapeHTML(s string) string {
	var sb strings.Builder
	for i := 0; i < len(s); i++ {
		c := s[i]
		switch c {
		case '&':
			sb.WriteString("&amp;")
		case '<':
			sb.WriteString("&lt;")
		case '>':
			sb.WriteString("&gt;")
		case '"':
			sb.WriteString("&quot;")
		case '\'':
			sb.WriteString("&#39;")
		default:
			sb.WriteByte(c)
		}
	}
	return sb.String()
}

// CleanHTML 标准化HTML输出，用于测试比较
func CleanHTML(html string) string {
	// 删除所有空白字符
	html = strings.ReplaceAll(html, "\n", "")
	html = strings.ReplaceAll(html, "\t", "")
	// 删除连续空格
	html = regexp.MustCompile(`\s+`).ReplaceAllString(html, " ")
	// 删除标签之间的空格
	html = regexp.MustCompile(`>\s+<`).ReplaceAllString(html, "><")
	// 删除开始和结束处的空格
	html = strings.TrimSpace(html)
	return html
}
