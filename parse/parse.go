package parse

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"io"
	"reflect"
	"strconv"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/theplant/htmlgo"
	"golang.org/x/net/html"
)

// 定义全局变量来存储组件定义
var componentDefinitions map[string]ComponentDefinition

func GenerateHTMLGo(pkg string, childrenMode bool, htmlCode io.Reader) string {
	// 尝试加载组件定义数据
	if componentDefinitions == nil {
		var err error
		componentDefinitions, err = ParseComponentData()
		if err != nil {
			// 记录错误但继续执行，使用默认的HTML处理
			fmt.Printf("Warning: Failed to load component data: %v\n", err)
			componentDefinitions = make(map[string]ComponentDefinition)
		}
	}

	n, err := html.Parse(htmlCode)
	if err != nil {
		panic(err)
	}
	methodNames := tagMethodNames()
	fc := &funcCall{}
	walk(n.FirstChild.FirstChild.NextSibling, fc, methodNames)

	code := string(fc.MarshalCode(methodNames, pkg, childrenMode))
	code = strings.TrimRight(code, ",\n")

	fset := token.NewFileSet()
	var f *ast.File
	f, err = parser.ParseFile(fset, "", "package hello\n var n = "+code, 0)
	if err != nil {
		hl, _ := strconv.ParseInt(strings.Split(err.Error(), ":")[0], 10, 64)
		panic(fmt.Sprintf("%s\n%s", err, codeWithLineNumber(code, hl)))
	}
	buf := bytes.NewBuffer(nil)
	err = printer.Fprint(buf, fset, f)
	if err != nil {
		panic(err)
	}
	return buf.String()
}

func codeWithLineNumber(src string, highlightLine int64) (r string) {
	lines := strings.Split(src, "\n")
	linesWithNumber := []string{}
	for i, l := range lines {
		hl := "   "
		if int64(i+1) == highlightLine {
			hl = ">> "
		}
		linesWithNumber = append(linesWithNumber, fmt.Sprintf("%s%d: %s", hl, i+1, l))
	}
	r = strings.Join(linesWithNumber, "\n")
	return
}

func pkgDot(pkg string) (r string) {
	if len(pkg) == 0 {
		return
	}
	return pkg + "."
}

type funcCall struct {
	Pkg          string
	Name         string
	Text         string
	TakeText     bool
	Children     []*funcCall
	Attrs        []html.Attribute
	IsComponent  bool                // 标记是否为组件
	ComponentDef ComponentDefinition // 存储组件定义
}

func (fc *funcCall) MarshalCode(methodNames []string, pkg string, childrenMode bool) (r []byte) {
	buf := bytes.NewBuffer(nil)

	if len(fc.Text) > 0 {
		buf.WriteString(fmt.Sprintf("%sText(%#+v),\n", pkgDot(pkg), fc.Text))
		return buf.Bytes()
	}

	newline := "\n"
	if fc.TakeText {
		newline = ""
	}

	// 处理组件或普通HTML标签
	if fc.IsComponent {
		// 处理组件
		_, _ = fmt.Fprintf(buf, "%s%s(%s", pkgDot(pkg), fc.ComponentDef.Go, newline)
	} else if strings.Contains(fc.Name, "-") {
		// 处理未知的自定义组件，转换为驼峰命名
		componentName := convertToCamelCase(fc.Name)
		_, _ = fmt.Fprintf(buf, "%s%s(%s", pkgDot(pkg), componentName, newline)
	} else {
		// 处理普通HTML标签，使用原始逻辑
		_, _ = fmt.Fprintf(buf, "%s%s(%s", pkgDot(pkg), strcase.ToCamel(fc.Name), newline)
	}

	needWriteChilren := false
	if childrenMode {
		needWriteChilren = true
		if fc.TakeText && len(fc.Children) == 1 && len(fc.Children[0].Text) > 0 {
			buf.WriteString(fmt.Sprintf("%#+v", fc.Children[0].Text))
			needWriteChilren = false
		} else if fc.TakeText {
			buf.WriteString(`""`)
		}
	} else {
		if fc.TakeText && len(fc.Children) == 1 && len(fc.Children[0].Text) > 0 {
			buf.WriteString(fmt.Sprintf("%#+v", fc.Children[0].Text))
		} else if fc.TakeText {
			buf.WriteString(`""`)
			needWriteChilren = true
		} else {
			for _, c := range fc.Children {
				buf.Write(c.MarshalCode(methodNames, pkg, childrenMode))
			}
		}
	}

	buf.WriteString(")")

	// 处理属性
	for i, att := range fc.Attrs {
		buf.WriteString(".")
		if i > 0 {
			buf.WriteString("\n")
		}

		if fc.IsComponent {
			// 检查属性是否在组件定义中存在
			if attrDef, ok := fc.ComponentDef.Attrs[att.Key]; ok {
				var val interface{} = att.Val

				// 根据接受类型处理值
				if attrDef.Accept == "bool" {
					val = true
					if att.Val == "false" {
						val = false
					}
				} else if attrDef.Accept == "int" {
					var err error
					val, err = strconv.ParseInt(att.Val, 10, 64)
					if err != nil {
						// 如果解析失败，保持原始值
						val = att.Val
					}
				}

				_, _ = fmt.Fprintf(buf, "%s(%s)", attrDef.Go, normalizeGoString(val))
			} else {
				// 如果属性未在定义中找到，使用Attr
				_, _ = fmt.Fprintf(buf, "Attr(%#+v, %s)", expandAlpineKey(att.Key), normalizeGoString(att.Val))
			}
		} else if strings.HasPrefix(fc.Name, "v-") || strings.HasPrefix(fc.Name, "V") {
			// Vuetify组件的特殊属性处理
			attName := strcase.ToCamel(att.Key)
			var val interface{} = att.Val

			// 处理布尔属性
			if att.Val == "" || att.Val == "true" {
				if att.Key == "required" || att.Key == "disabled" || att.Key == "readonly" || att.Key == "multiple" {
					val = true
				}
			}

			// 处理数字属性
			if att.Key == "width" || att.Key == "height" || att.Key == "max-width" || att.Key == "max-height" {
				if i, err := strconv.ParseInt(att.Val, 10, 64); err == nil {
					val = i
				}
			}

			_, _ = fmt.Fprintf(buf, "%s(%s)", attName, normalizeGoString(val))
		} else if strings.Contains(fc.Name, "-") {
			// 未知组件的属性都使用Attr
			_, _ = fmt.Fprintf(buf, "Attr(%#+v, %s)", expandAlpineKey(att.Key), normalizeGoString(att.Val))
		} else {
			// 原始HTML属性处理逻辑
			attFuncName := getFuncName(att.Key, methodNames)
			if len(attFuncName) > 0 {
				var val interface{} = att.Val
				if strings.Index(boolAttr, "|"+attFuncName+"|") >= 0 {
					val = true
				}
				if strings.Index(intAttr, "|"+attFuncName+"|") >= 0 {
					var err error
					val, err = strconv.ParseInt(att.Val, 10, 64)
					if err != nil {
						panic(err)
					}
				}
				_, _ = fmt.Fprintf(buf, "%s(%s)", attFuncName, normalizeGoString(val))
			} else {
				_, _ = fmt.Fprintf(buf, "Attr(%#+v, %s)", expandAlpineKey(att.Key), normalizeGoString(att.Val))
			}
		}
	}

	if needWriteChilren && len(fc.Children) > 0 {
		buf.WriteString(".Children(\n")
		for _, c := range fc.Children {
			buf.Write(c.MarshalCode(methodNames, pkg, childrenMode))
		}
		buf.WriteString(")")
	}

	buf.WriteString(",\n")

	return buf.Bytes()
}

func expandAlpineKey(key string) (r string) {
	if strings.Index(key, ":") == 0 {
		return fmt.Sprintf("x-bind%s", key)
	}

	if strings.Index(key, "@") == 0 {
		return fmt.Sprintf("x-on%s", key)
	}
	return key
}

func normalizeGoString(val interface{}) (r interface{}) {
	strval, ok := val.(string)
	if !ok {
		return fmt.Sprintf("%#+v", val)
	}

	if strings.Contains(strval, "'") {
		strval = strings.ReplaceAll(strval, "'", "\"")
	}
	if strings.ContainsAny(strval, "\n\t\"") && !strings.Contains(strval, "`") {
		strval = fmt.Sprintf("`%s`", strval)
		return strval
	}

	return fmt.Sprintf("%#+v", strval)
}

const intAttr = "|TabIndex|"
const boolAttr = "|Required|Readonly|Disabled|Checked|"

const textTags = "|Abbr|B|Bdi|Bdo|Button|Caption|Code|Del|Dfn|Em|Figcaption|H1|H2|H3|H4|H5|" +
	"H6|I|Img|Input|Kbd|Label|Legend|Link|Mark|Object|Option|Param|Pre|Q|Rp|Rt|S|" +
	"Script|Small|Source|Span|Strong|Style|Sub|Sup|Textarea|Th|Time|Title|Track|U|Var|Wbr|"

func walk(n *html.Node, fc *funcCall, methodNames []string) {
	switch n.Type {
	case html.ElementNode:
		if len(strings.TrimSpace(n.Data)) > 0 {
			tagName := strings.TrimSpace(n.Data)

			// 检查是否是组件
			if componentDef, isComponent := componentDefinitions[tagName]; isComponent {
				// 是组件，设置组件相关信息
				fc.IsComponent = true
				fc.ComponentDef = componentDef
				fc.Name = tagName

				// 根据Accept属性判断是否接受文本
				if componentDef.Accept == "none" {
					fc.TakeText = false
				} else if strings.Contains(componentDef.Accept, "...h.HTMLComponent") {
					// 允许子组件，不设置TakeText为true
					fc.TakeText = false
				}
			} else if strings.HasPrefix(tagName, "v-") {
				// Vuetify组件 (v-btn, v-card 等)
				// 转换为对应的驼峰命名格式 (VBtn, VCard 等)
				fc.Name = strings.ToUpper(tagName[2:3]) + tagName[3:]
				// Vuetify组件通常使用子元素而非TakeText
				fc.TakeText = false
			} else if strings.Contains(tagName, "-") {
				// 对于自定义组件（包含连字符的标签），保持原始名称
				fc.Name = tagName
			} else {
				// 不是组件，使用原始逻辑
				fc.Name = strcase.ToCamel(tagName)
				// 检查是否是允许文本的HTML标签
				if strings.Index(textTags, "|"+fc.Name+"|") >= 0 {
					fc.TakeText = true
				}
			}
		}
	case html.TextNode:
		if len(strings.TrimSpace(n.Data)) > 0 {
			fc.Text = strings.TrimSpace(n.Data)
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if c.Type == html.TextNode && len(strings.TrimSpace(c.Data)) == 0 {
			continue
		}
		if c.Type == html.CommentNode {
			continue
		}

		ch := &funcCall{Attrs: c.Attr}
		fc.Children = append(fc.Children, ch)
		walk(c, ch, methodNames)
	}
}

func getFuncName(name string, methodNames []string) (r string) {
	for _, m := range methodNames {
		if strings.ToLower(name) == strings.ToLower(m) {
			return m
		}
	}
	return ""
}

func tagMethodNames() (r []string) {
	tag := htmlgo.Tag("")
	tagType := reflect.TypeOf(tag)
	for i := 0; i < tagType.NumMethod(); i++ {
		r = append(r, tagType.Method(i).Name)
	}
	return
}

// 将连字符格式的标签名转换为驼峰命名
func convertToCamelCase(name string) string {
	parts := strings.Split(name, "-")
	for i := range parts {
		parts[i] = strcase.ToCamel(parts[i])
	}
	return strings.Join(parts, "")
}
