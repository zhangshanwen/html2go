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

func GenerateHTMLGo(pkg string, vuetifyPkg string, vuetifyxPkg string, childrenMode bool, htmlCode io.Reader) string {
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

	code := string(fc.MarshalCode(methodNames, pkg, vuetifyPkg, vuetifyxPkg, childrenMode))
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

func (fc *funcCall) MarshalCode(methodNames []string, pkg string, vuetifyPkg string, vuetifyxPkg string, childrenMode bool) (r []byte) {
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
		// 根据组件类型选择包前缀
		usePkg := pkg // 默认使用pkg

		// 根据组件类型设置对应的包前缀
		if fc.ComponentDef.Type == "vuetify" {
			usePkg = vuetifyPkg
		} else if fc.ComponentDef.Type == "vuetifyx" {
			usePkg = vuetifyxPkg
		}

		// 处理已知组件
		_, _ = fmt.Fprintf(buf, "%s%s(%s", pkgDot(usePkg), fc.ComponentDef.Go, newline)

		// 处理子元素
		if !childrenMode {
			if fc.TakeText && len(fc.Children) == 1 && len(fc.Children[0].Text) > 0 {
				buf.WriteString(fmt.Sprintf("%#+v", fc.Children[0].Text))
			} else if fc.TakeText {
				buf.WriteString(`""`)
			} else {
				for _, c := range fc.Children {
					buf.Write(c.MarshalCode(methodNames, pkg, vuetifyPkg, vuetifyxPkg, childrenMode))
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
		}

		// 处理Children模式
		if childrenMode {
			buf.WriteString(".Children(")
			if fc.TakeText && len(fc.Children) == 1 && len(fc.Children[0].Text) > 0 {
				buf.WriteString(fmt.Sprintf("\n%sText(%#+v),\n", pkgDot(pkg), fc.Children[0].Text))
			} else if fc.TakeText {
				// 空文本
			} else {
				buf.WriteString("\n")
				for _, c := range fc.Children {
					buf.Write(c.MarshalCode(methodNames, pkg, vuetifyPkg, vuetifyxPkg, childrenMode))
				}
			}
			buf.WriteString(")")
		}

	} else if strings.Contains(fc.Name, "-") {
		// 处理未知的组件（使用 Tag 而不是硬编码的 h.Tag）
		_, _ = fmt.Fprintf(buf, "%sTag(%#v)", pkgDot(pkg), fc.Name)

		// 自定义组件和未知组件的处理方式区分
		if childrenMode {
			// Child mode - 先处理属性后处理子元素
			if len(fc.Attrs) > 0 {
				for i, att := range fc.Attrs {
					buf.WriteString(".")
					if i > 0 {
						buf.WriteString("\n\t")
					}
					_, _ = fmt.Fprintf(buf, "Attr(%#+v, %s)", expandAlpineKey(att.Key), normalizeGoString(att.Val))
				}
			}

			if len(fc.Children) > 0 {
				buf.WriteString(".Children(\n")
				for _, c := range fc.Children {
					buf.Write(c.MarshalCode(methodNames, pkg, vuetifyPkg, vuetifyxPkg, childrenMode))
				}
				buf.WriteString(")")
			}
		} else {
			// Normal mode - 先处理子元素后处理属性
			if len(fc.Children) > 0 {
				buf.WriteString(".Children(\n")
				for _, c := range fc.Children {
					buf.Write(c.MarshalCode(methodNames, pkg, vuetifyPkg, vuetifyxPkg, childrenMode))
				}
				buf.WriteString(")")
			}

			// 然后处理所有属性，多个属性时添加换行和缩进
			if len(fc.Attrs) > 0 {
				for i, att := range fc.Attrs {
					buf.WriteString(".")
					if i > 0 {
						buf.WriteString("\n\t")
					}
					_, _ = fmt.Fprintf(buf, "Attr(%#+v, %s)", expandAlpineKey(att.Key), normalizeGoString(att.Val))
				}
			}
		}

	} else {
		// 处理普通HTML标签
		_, _ = fmt.Fprintf(buf, "%s%s(%s", pkgDot(pkg), strcase.ToCamel(fc.Name), newline)

		// 处理子元素
		needWriteChildren := false
		if childrenMode {
			needWriteChildren = true
			if fc.TakeText && len(fc.Children) == 1 && len(fc.Children[0].Text) > 0 {
				buf.WriteString(fmt.Sprintf("%#+v", fc.Children[0].Text))
				needWriteChildren = false
			} else if fc.TakeText {
				buf.WriteString(`""`)
			}
		} else {
			if fc.TakeText && len(fc.Children) == 1 && len(fc.Children[0].Text) > 0 {
				buf.WriteString(fmt.Sprintf("%#+v", fc.Children[0].Text))
			} else if fc.TakeText {
				buf.WriteString(`""`)
				needWriteChildren = true
			} else {
				for _, c := range fc.Children {
					buf.Write(c.MarshalCode(methodNames, pkg, vuetifyPkg, vuetifyxPkg, childrenMode))
				}
			}
		}

		buf.WriteString(")")

		// 处理普通HTML标签的属性
		for i, att := range fc.Attrs {
			buf.WriteString(".")
			if i > 0 {
				buf.WriteString("\n")
			}

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

		if needWriteChildren && len(fc.Children) > 0 {
			buf.WriteString(".Children(\n")
			for _, c := range fc.Children {
				buf.Write(c.MarshalCode(methodNames, pkg, vuetifyPkg, vuetifyxPkg, childrenMode))
			}
			buf.WriteString(")")
		}
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

			// 检查是否是已映射的组件
			if componentDef, isComponent := componentDefinitions[tagName]; isComponent {
				// 是已映射的组件，设置组件相关信息
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
			} else {
				// 未映射的组件，包括普通HTML标签、自定义组件、未知的v-前缀和vx-前缀组件
				fc.Name = tagName
				fc.IsComponent = false

				// 普通HTML标签使用原始逻辑判断是否接受文本
				if !strings.Contains(tagName, "-") {
					camelName := strcase.ToCamel(tagName)
					if strings.Index(textTags, "|"+camelName+"|") >= 0 {
						fc.TakeText = true
					}
				} else {
					// 自定义组件、未知的v-前缀和vx-前缀组件默认不接受文本
					fc.TakeText = false
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
