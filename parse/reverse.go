package parse

import (
	"fmt"
	"regexp"
	"strings"
)

// ReverseMapping 存储Go函数名到HTML标签和属性的反向映射
type ReverseMapping struct {
	// Go组件名到HTML标签的映射
	GoComponentToHTMLTag map[string]string
	// Go组件名到属性映射的映射
	GoMethodToHTMLAttr map[string]map[string]AttributeDefinition
	// 特殊标签和属性映射
	commonAttrMap map[string]string
}

// NewReverseMapping 创建并初始化反向映射结构
func NewReverseMapping() *ReverseMapping {
	rm := &ReverseMapping{
		GoComponentToHTMLTag: make(map[string]string),
		GoMethodToHTMLAttr:   make(map[string]map[string]AttributeDefinition),
		commonAttrMap:        make(map[string]string),
	}

	// 初始化常用属性映射
	rm.initCommonAttrMap()

	// 从组件定义初始化反向映射（包括Vuetify和VuetifyX组件）
	rm.initFromComponentDefinitions()

	// 添加任何特殊情况的映射规则
	rm.addSpecialMappings()

	return rm
}

// 初始化常用属性映射
func (rm *ReverseMapping) initCommonAttrMap() {
	// HTML常用属性的映射
	rm.commonAttrMap = map[string]string{
		"Class":       "class",
		"Id":          "id",
		"Style":       "style",
		"Type":        "type",
		"Href":        "href",
		"Src":         "src",
		"Alt":         "alt",
		"Value":       "value",
		"Placeholder": "placeholder",
		"Name":        "name",
		"Title":       "title",
		"Width":       "width",
		"Height":      "height",
		"TabIndex":    "tabindex",
		"Target":      "target",
		"Rel":         "rel",
		"Method":      "method",
		"Action":      "action",
		"ForAttr":     "for",
		"For":         "for",
		"Role":        "role",
		"Disabled":    "disabled",
		"Readonly":    "readonly",
		"Required":    "required",
		"Checked":     "checked",
		"Selected":    "selected",
		"Multiple":    "multiple",
		"Color":       "color",    // Vuetify 常用属性
		"Text":        "text",     // VuetifyX 常用属性
		"OnClick":     "on-click", // 事件处理
	}
}

// 从组件定义初始化反向映射
func (rm *ReverseMapping) initFromComponentDefinitions() {
	// 确保组件定义已加载
	if componentDefinitions == nil {
		var err error
		componentDefinitions, err = ParseComponentData()
		if err != nil {
			fmt.Printf("Warning: Failed to load component data: %v\n", err)
			componentDefinitions = make(map[string]ComponentDefinition)
		}
	}

	// 从组件定义构建反向映射
	for tag, compDef := range componentDefinitions {
		// 映射组件名
		rm.GoComponentToHTMLTag[compDef.Go] = tag

		// 映射组件属性
		attrMap := make(map[string]AttributeDefinition)
		rm.GoMethodToHTMLAttr[compDef.Go] = attrMap

		for attrName, attrDef := range compDef.Attrs {
			attrMap[attrDef.Go] = AttributeDefinition{
				Go:     attrName,
				Accept: attrDef.Accept,
			}
		}
	}
}

// addSpecialMappings 添加任何特殊处理的映射规则
func (rm *ReverseMapping) addSpecialMappings() {
	// 所有组件映射都应该从componentDefinitions获取
	// 这个函数保留为扩展点，以便将来可能需要添加特殊处理规则
	// 但不再硬编码Vuetify和VuetifyX组件的映射
}

// ResolveHTMLTag 将Go函数名解析为HTML标签名
func (rm *ReverseMapping) ResolveHTMLTag(funcName string) string {
	// 1. 检查是否是Text函数
	if funcName == "Text" {
		return "#text" // 特殊标记，表示文本节点
	}

	// 2. 处理 Body 函数 - 在我们的处理中忽略 body 标签
	if funcName == "Body" {
		return "body"
	}

	// 3. 检查特殊组件映射
	if tag, ok := rm.GoComponentToHTMLTag[funcName]; ok {
		return tag
	}

	// 4. 检查是否为原生HTML标签（转小写）
	lowerFuncName := strings.ToLower(funcName)
	if isHTMLTag(lowerFuncName) {
		return lowerFuncName
	}

	// 5. 使用驼峰转连字符作为后备选项
	return kebabCase(funcName)
}

// kebabCase 将驼峰命名转换为连字符命名
func kebabCase(s string) string {
	// 在大写字母前添加连字符，然后转小写
	s = regexp.MustCompile(`([a-z0-9])([A-Z])`).ReplaceAllString(s, "$1-$2")
	return strings.ToLower(s)
}

// ResolveHTMLAttr 将Go方法名解析为HTML属性名和值
func (rm *ReverseMapping) ResolveHTMLAttr(componentName string, methodName string, value interface{}) (string, string) {
	// 1. 检查组件特定属性映射
	if attrMap, ok := rm.GoMethodToHTMLAttr[componentName]; ok {
		if attr, ok := attrMap[methodName]; ok {
			// 如果是布尔属性且值为true，返回无值的属性
			if attr.Accept == "bool" && isTrueValue(value) {
				return attr.Go, ""
			}
			return attr.Go, formatAttrValue(attr.Accept, value)
		}
	}

	// 2. 检查通用属性映射
	if attrName, ok := rm.commonAttrMap[methodName]; ok {
		// 特殊处理布尔属性
		if isBooleanAttr(attrName) && isTrueValue(value) {
			return attrName, ""
		}
		return attrName, formatValue(value)
	}

	// 3. 处理Attr方法调用
	if methodName == "Attr" {
		if pair, ok := value.([]interface{}); ok && len(pair) == 2 {
			if attrName, ok := pair[0].(string); ok {
				// 检查是否是布尔属性
				if isBooleanAttr(attrName) && isTrueValue(pair[1]) {
					return attrName, ""
				}
				return attrName, formatValue(pair[1])
			}
		}
	}

	// 4. 使用驼峰转连字符作为后备选项
	attrName := kebabCase(methodName)

	// 特殊处理布尔属性
	if isBooleanAttr(attrName) && isTrueValue(value) {
		return attrName, ""
	}

	return attrName, formatValue(value)
}

// 检查是否为HTML标准标签
func isHTMLTag(tag string) bool {
	htmlTags := map[string]bool{
		"a": true, "abbr": true, "address": true, "area": true, "article": true,
		"aside": true, "audio": true, "b": true, "base": true, "bdi": true,
		"bdo": true, "blockquote": true, "body": true, "br": true, "button": true,
		"canvas": true, "caption": true, "cite": true, "code": true, "col": true,
		"colgroup": true, "data": true, "datalist": true, "dd": true, "del": true,
		"details": true, "dfn": true, "dialog": true, "div": true, "dl": true,
		"dt": true, "em": true, "embed": true, "fieldset": true, "figcaption": true,
		"figure": true, "footer": true, "form": true, "h1": true, "h2": true,
		"h3": true, "h4": true, "h5": true, "h6": true, "head": true, "header": true,
		"hr": true, "html": true, "i": true, "iframe": true, "img": true,
		"input": true, "ins": true, "kbd": true, "label": true, "legend": true,
		"li": true, "link": true, "main": true, "map": true, "mark": true,
		"meta": true, "meter": true, "nav": true, "noscript": true, "object": true,
		"ol": true, "optgroup": true, "option": true, "output": true, "p": true,
		"param": true, "picture": true, "pre": true, "progress": true, "q": true,
		"rp": true, "rt": true, "ruby": true, "s": true, "samp": true,
		"script": true, "section": true, "select": true, "small": true, "source": true,
		"span": true, "strong": true, "style": true, "sub": true, "summary": true,
		"sup": true, "svg": true, "table": true, "tbody": true, "td": true,
		"template": true, "textarea": true, "tfoot": true, "th": true, "thead": true,
		"time": true, "title": true, "tr": true, "track": true, "u": true,
		"ul": true, "var": true, "video": true, "wbr": true,
	}
	return htmlTags[tag]
}

// 检查是否为布尔属性
func isBooleanAttr(attrName string) bool {
	// 首先尝试从组件定义中查找属性类型
	if componentDefinitions != nil {
		for _, compDef := range componentDefinitions {
			for attr, attrDef := range compDef.Attrs {
				if attrName == attr && attrDef.Accept == "bool" {
					return true
				}
			}
		}
	}

	// 作为后备，使用常见的布尔属性列表
	booleanAttrs := map[string]bool{
		"disabled": true, "checked": true, "selected": true, "readonly": true,
		"required": true, "multiple": true, "autofocus": true, "autoplay": true,
		"controls": true, "default": true, "defer": true, "download": true,
		"hidden": true, "ismap": true, "loop": true, "muted": true, "novalidate": true,
		"open": true, "reversed": true, "scoped": true, "seamless": true,
		"typemustmatch": true, "formnovalidate": true, "async": true, "compact": true,
		"declare": true, "inert": true, "truespeed": true,
	}
	return booleanAttrs[attrName]
}

// 检查值是否表示"true"
func isTrueValue(value interface{}) bool {
	if boolVal, ok := value.(bool); ok {
		return boolVal
	}
	if strVal, ok := value.(string); ok {
		return strVal == "true" || strVal == ""
	}
	return false
}

// 检查值是否为字符串对
func isStringPair(value interface{}) bool {
	if pair, ok := value.([]interface{}); ok && len(pair) == 2 {
		_, ok1 := pair[0].(string)
		return ok1
	}
	return false
}

// 根据属性类型格式化属性值
func formatAttrValue(acceptType string, value interface{}) string {
	switch acceptType {
	case "bool":
		// 布尔属性的处理已在 ResolveHTMLAttr 中完成
		// 这里只处理非空值的情况
		if val, ok := value.(bool); ok {
			if val {
				return "" // HTML布尔属性为true时不需要值
			}
			return "false"
		}
		return formatValue(value)
	case "int", "float":
		return fmt.Sprintf("%v", value)
	default:
		return formatValue(value)
	}
}

// 格式化值为字符串
func formatValue(value interface{}) string {
	if value == nil {
		return ""
	}

	switch v := value.(type) {
	case string:
		return v
	case bool:
		if v {
			return "true"
		}
		return "false"
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
		return fmt.Sprintf("%d", v)
	case float32, float64:
		return fmt.Sprintf("%g", v)
	default:
		return fmt.Sprintf("%v", v)
	}
}
