package transform

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"path/filepath"
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

// TestMultiFileComponent tests the ability to parse components defined across multiple files
func TestMultiFileComponent(t *testing.T) {
	// Create temporary directory for test files
	tempDir, err := ioutil.TempDir("", "transform_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create btn.go
	btnContent := `package vuetify

import (
	"context"
	"fmt"

	h "github.com/theplant/htmlgo"
)

type VBtnBuilder struct {
	tag *h.HTMLTagBuilder
}

func (b *VBtnBuilder) Symbol(v interface{}) (r *VBtnBuilder) {
	b.tag.Attr(":symbol", h.JSONString(v))
	return b
}

func (b *VBtnBuilder) Text(v string) (r *VBtnBuilder) {
	b.tag.Attr("text", v)
	return b
}

func (b *VBtnBuilder) Color(v string) (r *VBtnBuilder) {
	b.tag.Attr("color", v)
	return b
}

func (b *VBtnBuilder) SetAttr(k string, v interface{}) {
	b.tag.SetAttr(k, v)
}

func (b *VBtnBuilder) Children(children ...h.HTMLComponent) (r *VBtnBuilder) {
	b.tag.Children(children...)
	return b
}

func (b *VBtnBuilder) MarshalHTML(ctx context.Context) (r []byte, err error) {
	return b.tag.MarshalHTML(ctx)
}
`

	// Create fix-btn.go with additional methods
	fixBtnContent := `package vuetify

import (
	"github.com/qor5/web/v3"
	h "github.com/theplant/htmlgo"
)

func VBtn(text string) (r *VBtnBuilder) {
	r = &VBtnBuilder{
		tag: h.Tag("v-btn").Text(text),
	}
	return
}

func (b *VBtnBuilder) OnClick(eventFuncId string) (r *VBtnBuilder) {
	b.tag.Attr("@click", web.POST().EventFunc(eventFuncId).Go())
	return b
}

func (b *VBtnBuilder) AttrIf(key, value interface{}, add bool) (r *VBtnBuilder) {
	b.tag.AttrIf(key, value, add)
	return b
}
`

	// Write test files
	btnPath := filepath.Join(tempDir, "btn.go")
	fixBtnPath := filepath.Join(tempDir, "fix-btn.go")

	if err := ioutil.WriteFile(btnPath, []byte(btnContent), 0644); err != nil {
		t.Fatalf("Failed to write btn.go: %v", err)
	}

	if err := ioutil.WriteFile(fixBtnPath, []byte(fixBtnContent), 0644); err != nil {
		t.Fatalf("Failed to write fix-btn.go: %v", err)
	}

	// Define all non-common methods we expect to find
	// These are methods defined in both btn.go and fix-btn.go that aren't in isCommonMethod
	expectedAttrs := map[string]string{
		"symbol":   "Symbol",
		"text":     "Text",
		"color":    "Color",
		"on-click": "OnClick",
	}

	// Test 1: Parse the files in original order (btnPath then fixBtnPath)
	t.Run("Original order", func(t *testing.T) {
		componentMap, err := ParseGoFiles([]string{btnPath, fixBtnPath})
		if err != nil {
			t.Fatalf("Failed to parse files: %v", err)
		}

		// Log the actual content of the component map for debugging
		t.Logf("Component map (original order): %+v", componentMap)

		// Verify the component was correctly identified
		btn, ok := componentMap["v-btn"]
		if !ok {
			t.Fatal("v-btn component not found")
		}

		// Verify the component properties
		if btn.Go != "VBtn" {
			t.Errorf("Expected Go name 'VBtn', got %q", btn.Go)
		}

		// Log all attributes found to help debug
		t.Logf("Found attributes for v-btn: %+v", btn.Attrs)

		// Verify that each expected attribute is present
		for attrName, methodName := range expectedAttrs {
			attr, exists := btn.Attrs[attrName]
			if !exists {
				t.Errorf("Expected attribute %q not found", attrName)
				continue
			}

			if attr.Go != methodName {
				t.Errorf("For attribute %q, expected method name %q, got %q", attrName, methodName, attr.Go)
			}
		}

		// Verify that we didn't find any unexpected attributes
		// This ensures ALL methods are accounted for
		for attrName, attr := range btn.Attrs {
			_, expected := expectedAttrs[attrName]
			if !expected {
				t.Errorf("Found unexpected attribute %q with method name %q", attrName, attr.Go)
			}
		}
	})

	// Test 2: Parse the files in reverse order (fixBtnPath then btnPath)
	t.Run("Reverse order", func(t *testing.T) {
		reverseComponentMap, err := ParseGoFiles([]string{fixBtnPath, btnPath})
		if err != nil {
			t.Fatalf("Failed to parse files in reverse order: %v", err)
		}

		// Log the actual content of the component map for debugging
		t.Logf("Component map (reverse order): %+v", reverseComponentMap)

		// Verify the component was correctly identified
		btn, ok := reverseComponentMap["v-btn"]
		if !ok {
			t.Fatal("v-btn component not found when parsing in reverse order")
		}

		// Verify the component properties
		if btn.Go != "VBtn" {
			t.Errorf("Reverse order: Expected Go name 'VBtn', got %q", btn.Go)
		}

		// Verify that each expected attribute is present
		for attrName, methodName := range expectedAttrs {
			attr, exists := btn.Attrs[attrName]
			if !exists {
				t.Errorf("Reverse order: Expected attribute %q not found", attrName)
				continue
			}

			if attr.Go != methodName {
				t.Errorf("Reverse order: For attribute %q, expected method name %q, got %q", attrName, methodName, attr.Go)
			}
		}

		// Verify that we didn't find any unexpected attributes
		for attrName, attr := range btn.Attrs {
			_, expected := expectedAttrs[attrName]
			if !expected {
				t.Errorf("Reverse order: Found unexpected attribute %q with method name %q", attrName, attr.Go)
			}
		}
	})

	// Test parsing a directory
	dirComponentMap, err := ParseGoDir(tempDir)
	if err != nil {
		t.Fatalf("Failed to parse directory: %v", err)
	}

	// Verify directory parsing works the same
	btn, ok := dirComponentMap["v-btn"]
	if !ok {
		t.Fatal("v-btn component not found when parsing directory")
	}

	// Verify attributes are the same for directory parsing
	for attrName, methodName := range expectedAttrs {
		attr, exists := btn.Attrs[attrName]
		if !exists {
			t.Errorf("Directory parsing: Expected attribute %q not found", attrName)
			continue
		}

		if attr.Go != methodName {
			t.Errorf("Directory parsing: For attribute %q, expected method name %q, got %q", attrName, methodName, attr.Go)
		}
	}

	// Verify that all attributes are accounted for in directory parsing
	for attrName, attr := range btn.Attrs {
		_, expected := expectedAttrs[attrName]
		if !expected {
			t.Errorf("Directory parsing: Found unexpected attribute %q with method name %q", attrName, attr.Go)
		}
	}
}
