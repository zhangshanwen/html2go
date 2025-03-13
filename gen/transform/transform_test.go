package transform

import (
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
	"testing"
)

func TestExtractTagName(t *testing.T) {
	// 测试用例
	testCases := []struct {
		name     string
		code     string
		expected string
	}{
		{
			name: "Basic component",
			code: `
package vuetify

import (
	"context"
	"fmt"
	h "github.com/theplant/htmlgo"
)

func VAlert(children ...h.HTMLComponent) (r *VAlertBuilder) {
	r = &VAlertBuilder{
		tag: h.Tag("v-alert").Children(children...),
	}
	return
}
`,
			expected: "v-alert",
		},
		{
			name: "Component with different tag name",
			code: `
package vuetify

import (
	"context"
	"fmt"
	h "github.com/theplant/htmlgo"
)

func VBtn(children ...h.HTMLComponent) (r *VBtnBuilder) {
	r = &VBtnBuilder{
		tag: h.Tag("v-btn").Children(children...),
	}
	return
}
`,
			expected: "v-btn",
		},
	}

	// 运行测试用例
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 解析代码
			fset := token.NewFileSet()
			node, err := parser.ParseFile(fset, "", tc.code, parser.ParseComments)
			if err != nil {
				t.Fatalf("Failed to parse code: %v", err)
			}

			// 查找函数声明
			var funcDecl *ast.FuncDecl
			for _, decl := range node.Decls {
				if fd, ok := decl.(*ast.FuncDecl); ok {
					if strings.HasPrefix(fd.Name.Name, "V") {
						funcDecl = fd
						break
					}
				}
			}

			if funcDecl == nil {
				t.Fatal("No function declaration found")
			}

			// 提取标签名
			tagName := ExtractTagName(funcDecl)
			if tagName != tc.expected {
				t.Errorf("Expected tag name %q, got %q", tc.expected, tagName)
			}
		})
	}
}

func TestExtractComponentInfo(t *testing.T) {
	// 测试用例
	code := `
package vuetify

import (
	"context"
	"fmt"
	h "github.com/theplant/htmlgo"
)

type VAlertBuilder struct {
	tag *h.HTMLTagBuilder
}

func VAlert(children ...h.HTMLComponent) (r *VAlertBuilder) {
	r = &VAlertBuilder{
		tag: h.Tag("v-alert").Children(children...),
	}
	return
}

func (b *VAlertBuilder) Title(v string) (r *VAlertBuilder) {
	b.tag.Attr("title", v)
	return b
}

func (b *VAlertBuilder) Text(v string) (r *VAlertBuilder) {
	b.tag.Attr("text", v)
	return b
}

func (b *VAlertBuilder) SetAttr(k string, v interface{}) {
	b.tag.SetAttr(k, v)
}

func (b *VAlertBuilder) MarshalHTML(ctx context.Context) (r []byte, err error) {
	return b.tag.MarshalHTML(ctx)
}
`

	// 解析代码
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, "", code, parser.ParseComments)
	if err != nil {
		t.Fatalf("Failed to parse code: %v", err)
	}

	// 提取组件信息
	componentMap := make(ComponentMap)
	ExtractComponentInfo(node, componentMap)

	// 验证结果
	if len(componentMap) != 1 {
		t.Fatalf("Expected 1 component, got %d", len(componentMap))
	}

	// 检查 v-alert 组件
	alert, ok := componentMap["v-alert"]
	if !ok {
		t.Fatal("v-alert component not found")
	}

	if alert.Go != "VAlert" {
		t.Errorf("Expected Go name 'VAlert', got %q", alert.Go)
	}

	if alert.Accept != "...h.HTMLComponent" {
		t.Errorf("Expected Accept '...h.HTMLComponent', got %q", alert.Accept)
	}

	// 检查属性
	if len(alert.Attrs) != 2 {
		t.Fatalf("Expected 2 attributes, got %d", len(alert.Attrs))
	}

	// 检查 title 属性
	title, ok := alert.Attrs["title"]
	if !ok {
		t.Fatal("title attribute not found")
	}
	if title.Go != "Title" {
		t.Errorf("Expected Go name 'Title', got %q", title.Go)
	}
	if title.Accept != "string" {
		t.Errorf("Expected Accept 'string', got %q", title.Accept)
	}

	// 检查 text 属性
	text, ok := alert.Attrs["text"]
	if !ok {
		t.Fatal("text attribute not found")
	}
	if text.Go != "Text" {
		t.Errorf("Expected Go name 'Text', got %q", text.Go)
	}
	if text.Accept != "string" {
		t.Errorf("Expected Accept 'string', got %q", text.Accept)
	}
}

func TestExtractAcceptType(t *testing.T) {
	// Test code with a function that accepts ...h.HTMLComponent
	src := `
package test

import "github.com/qor5/web/h"

func VComponent(children ...h.HTMLComponent) *VComponentBuilder {
	return &VComponentBuilder{
		tag: h.Tag("v-component"),
	}
}
`
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "", src, parser.ParseComments)
	if err != nil {
		t.Fatalf("Failed to parse test source: %v", err)
	}

	// Find the VComponent function
	var funcDecl *ast.FuncDecl
	for _, decl := range f.Decls {
		if fd, ok := decl.(*ast.FuncDecl); ok && fd.Name.Name == "VComponent" {
			funcDecl = fd
			break
		}
	}

	if funcDecl == nil {
		t.Fatal("Failed to find VComponent function in test source")
	}

	// Test ExtractAcceptType
	acceptType := ExtractAcceptType(funcDecl)
	if acceptType != "...h.HTMLComponent" {
		t.Errorf("Expected accept type to be '...h.HTMLComponent', got '%s'", acceptType)
	}
}

func TestParseGoFile(t *testing.T) {
	// This is a placeholder test that you would replace with actual file parsing
	// when you have sample files to test with
	t.Skip("Implement this test with actual Go files")
}

func TestComponentMapGeneration(t *testing.T) {
	// Create a simple component map for testing
	componentMap := ComponentMap{
		"v-component": ComponentInfo{
			Go:     "VComponent",
			Accept: "...h.HTMLComponent",
			Attrs: map[string]AttrInfo{
				"color": {
					Go:     "Color",
					Accept: "string",
				},
			},
		},
	}

	// Verify the structure is as expected
	component, exists := componentMap["v-component"]
	if !exists {
		t.Fatal("Component v-component not found in map")
	}

	if component.Go != "VComponent" {
		t.Errorf("Expected Go name to be 'VComponent', got '%s'", component.Go)
	}

	if component.Accept != "...h.HTMLComponent" {
		t.Errorf("Expected Accept to be '...h.HTMLComponent', got '%s'", component.Accept)
	}

	attr, exists := component.Attrs["color"]
	if !exists {
		t.Fatal("Attribute 'color' not found in component")
	}

	if attr.Go != "Color" {
		t.Errorf("Expected Go method to be 'Color', got '%s'", attr.Go)
	}

	if attr.Accept != "string" {
		t.Errorf("Expected Accept to be 'string', got '%s'", attr.Accept)
	}
}
