package integration

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/zhangshanwen/html2go/parse"
)

func TestGoToHTML(t *testing.T) {
	tests := []struct {
		name     string
		goCode   string
		expected string
	}{
		{
			name:     "基本div元素带文本",
			goCode:   `Div(Text("Hello, world!"))`,
			expected: `<div>Hello, world!</div>`,
		},
		{
			name:     "带包前缀的元素",
			goCode:   `h.Div(h.Text("With package prefix"))`,
			expected: `<div>With package prefix</div>`,
		},
		// 添加完整包名的测试用例
		{
			name:     "使用完整包名的标准元素",
			goCode:   `htmlgo.Div(htmlgo.Text("Using full package name"))`,
			expected: `<div>Using full package name</div>`,
		},
		{
			name:     "使用完整包名的Vuetify元素",
			goCode:   `vuetify.VCardText(htmlgo.Text("Vuetify card text"))`,
			expected: `<v-card-text>Vuetify card text</v-card-text>`,
		},
		{
			name:     "使用完整包名的VuetifyX元素",
			goCode:   `vuetifyx.VXBtn().Text("VuetifyX button")`,
			expected: `<v-xbtn text="VuetifyX button"></v-xbtn>`,
		},
		{
			name:     "完整包名与简短包名混合",
			goCode:   `htmlgo.Div(v.VBtn(htmlgo.Text("Mixed packages")))`,
			expected: `<div><v-btn>Mixed packages</v-btn></div>`,
		},
		{
			name:     "完整包名的嵌套结构",
			goCode:   `htmlgo.Div(htmlgo.H1(htmlgo.Text("Title")), vuetify.VCard(vuetify.VCardTitle(htmlgo.Text("Card Title"))))`,
			expected: `<div><h1>Title</h1><v-card><v-card-title>Card Title</v-card-title></v-card></div>`,
		},
		// 添加自定义组件测试（使用Tag函数）
		{
			name:     "使用Tag函数的自定义组件",
			goCode:   `Tag("custom-component").Children(Text("Content"))`,
			expected: `<custom-component>Content</custom-component>`,
		},
		{
			name:     "带属性的自定义组件",
			goCode:   `Tag("custom-component").Attr("prop1", "value1").Children(Text("Content"))`,
			expected: `<custom-component prop1="value1">Content</custom-component>`,
		},
		{
			name:     "带包前缀的自定义组件",
			goCode:   `h.Tag("custom-component").Children(h.Text("Content"))`,
			expected: `<custom-component>Content</custom-component>`,
		},
		{
			name:     "完整包名的自定义组件",
			goCode:   `htmlgo.Tag("custom-component").Children(htmlgo.Text("Content"))`,
			expected: `<custom-component>Content</custom-component>`,
		},
		// Vuetify组件测试
		{
			name:     "简单的Vuetify组件",
			goCode:   `v.VBtn(Text("Submit"))`,
			expected: `<v-btn>Submit</v-btn>`,
		},
		{
			name:     "带属性的Vuetify组件",
			goCode:   `v.VBtn(Text("Submit")).Color("primary")`,
			expected: `<v-btn color="primary">Submit</v-btn>`,
		},
		// VuetifyX组件测试
		{
			name:     "简单的VuetifyX组件",
			goCode:   `vx.VXBtn(Text("Click me"))`,
			expected: `<v-xbtn>Click me</v-xbtn>`,
		},
		{
			name:     "带属性的VuetifyX组件",
			goCode:   `vx.VXBtn().Text("Submit").Color("primary")`,
			expected: `<v-xbtn color="primary" text="Submit"></v-xbtn>`,
		},
		// 未映射的组件测试
		{
			name:     "未映射的Vuetify组件",
			goCode:   `v.VUnknownComponent(Text("Unknown content"))`,
			expected: `<v-unknown-component>Unknown content</v-unknown-component>`,
		},
		{
			name:     "未映射的VuetifyX组件",
			goCode:   `vx.VXUnknownComponent(Text("Unknown content"))`,
			expected: `<v-xunknown-component>Unknown content</v-xunknown-component>`,
		},
		{
			name:     "完整包名的未映射组件",
			goCode:   `vuetify.VCustomComponent(htmlgo.Text("Custom content"))`,
			expected: `<v-custom-component>Custom content</v-custom-component>`,
		},
		// 复杂的混合包前缀测试
		{
			name:     "简单的混合包前缀",
			goCode:   `h.Div(v.VBtn(Text("Button")))`,
			expected: `<div><v-btn>Button</v-btn></div>`,
		},
		{
			name:     "复杂的混合包前缀",
			goCode:   `h.Div(v.VBtn(Text("Vuetify Button")), vx.VXBtn(Text("VuetifyX Button")), h.P(h.Text("Regular text")))`,
			expected: `<div><v-btn>Vuetify Button</v-btn><v-xbtn>VuetifyX Button</v-xbtn><p>Regular text</p></div>`,
		},
		{
			name:     "多级嵌套的混合包前缀",
			goCode:   `h.Div(v.VCard(v.VCardTitle(Text("Card Title")), v.VCardText(Text("Card content"))), h.Div(vx.VXBtn(Text("Action"))))`,
			expected: `<div><v-card><v-card-title>Card Title</v-card-title><v-card-text>Card content</v-card-text></v-card><div><v-xbtn>Action</v-xbtn></div></div>`,
		},
		// 布尔属性测试
		{
			name:     "布尔属性",
			goCode:   `Input().Disabled(true).Required(true)`,
			expected: `<input disabled required />`,
		},
		{
			name:     "完整包名的布尔属性",
			goCode:   `htmlgo.Input().Disabled(true).Required(true)`,
			expected: `<input disabled required />`,
		},
		// 使用Children方法
		{
			name:     "使用Children方法",
			goCode:   `Div().Children(H1(Text("Title")), P(Text("Paragraph")))`,
			expected: `<div><h1>Title</h1><p>Paragraph</p></div>`,
		},
		{
			name:     "使用包前缀和Children方法",
			goCode:   `h.Div().Children(h.H1(h.Text("Title")), h.P(h.Text("Paragraph")))`,
			expected: `<div><h1>Title</h1><p>Paragraph</p></div>`,
		},
		{
			name:     "使用完整包名和Children方法",
			goCode:   `htmlgo.Div().Children(htmlgo.H1(htmlgo.Text("Title")), htmlgo.P(htmlgo.Text("Paragraph")))`,
			expected: `<div><h1>Title</h1><p>Paragraph</p></div>`,
		},
		// 嵌套元素和复杂嵌套
		{
			name:     "嵌套元素",
			goCode:   `Div(H1(Text("Header")), P(Text("Paragraph")))`,
			expected: `<div><h1>Header</h1><p>Paragraph</p></div>`,
		},
		{
			name:     "复杂嵌套元素",
			goCode:   `Div(H1(Text("Header")), Div(P(Text("Nested paragraph")), Span(Text("Span text"))))`,
			expected: `<div><h1>Header</h1><div><p>Nested paragraph</p><span>Span text</span></div></div>`,
		},
		// 表单元素测试
		{
			name:     "表单元素",
			goCode:   `Form(Input(), Button(Text("Submit")))`,
			expected: `<form><input /><button>Submit</button></form>`,
		},
		{
			name:     "带属性的表单元素",
			goCode:   `Form(Input().Type("text").Name("username"), Button(Text("Submit")).Type("submit"))`,
			expected: `<form><input name="username" type="text" /><button type="submit">Submit</button></form>`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			goCode := tt.goCode

			// 特殊处理多行代码测试用例
			if strings.Contains(tt.name, "多行Go代码") {
				// 将多行Go代码包装成完整的Go程序
				goCode = "package hello\n\nvar n = Body(\n" + goCode + "\n)"
			}

			output, err := parse.GenerateGoHTML(goCode, false, 2)
			if err != nil {
				t.Fatalf("GenerateGoHTML error: %v", err)
			}

			// 标准化输出和预期结果，以消除空白差异
			normalizedOutput := NormalizeFormAttrOrder(parse.NormalizeHTML(output))
			normalizedExpected := NormalizeFormAttrOrder(parse.NormalizeHTML(tt.expected))

			if normalizedOutput != normalizedExpected {
				t.Errorf("Expected:\n%s\nGot:\n%s", normalizedExpected, normalizedOutput)
			}
		})
	}
}

// TestGoToHTMLErrors tests error handling
func TestGoToHTMLErrors(t *testing.T) {
	invalidCode := `
	// This is not valid Go code
	<div>Invalid HTML in Go</div>
	`

	_, err := parse.GenerateGoHTML(invalidCode, false, 2)
	if err == nil {
		t.Errorf("Expected error for invalid Go code, but got nil")
	}
}

// TestMultilineGoHTML demonstrates how to test multi-line Go code
func TestMultilineGoHTML(t *testing.T) {
	// Skip this test, it's only for README example
	t.Skip("Skipping multi-line code test, just for README example")

	tests := []struct {
		name     string
		goCode   string
		expected string
	}{
		{
			name: "Multiline Go Code - Complex Nested Structure",
			goCode: `Div(
	Class("container"),
	H1(
		Class("title"),
		Text("Multi-line Go Code Test")
	),
	Div(
		Class("content"),
		P(Text("This is the first paragraph")),
		P(Text("This is the second paragraph")),
		Ul(
			Li(Text("List item 1")),
			Li(Text("List item 2")),
			Li(Text("List item 3"))
		)
	),
	Button(
		Type("button"),
		Class("btn btn-primary"),
		Text("Submit")
	)
)`,
			expected: `<div class="container"><h1 class="title">Multi-line Go Code Test</h1><div class="content"><p>This is the first paragraph</p><p>This is the second paragraph</p><ul><li>List item 1</li><li>List item 2</li><li>List item 3</li></ul></div><button class="btn btn-primary" type="button">Submit</button></div>`,
		},
		{
			name: "Multiline Go Code - Vuetify Component Structure",
			goCode: `v.VApp(
	v.VMain(
		v.VContainer(
			v.VRow(
				v.VCol(
					v.VCard(
						v.VCardTitle(Text("Card Title")),
						v.VCardText(
							Text("This is card content with some descriptive text."),
							v.VDivider(),
							P(Text("Paragraph below divider"))
						),
						v.VCardActions(
							v.VBtn(Text("Cancel")).Color("error"),
							v.VBtn(Text("Confirm")).Color("primary")
						)
					)
				).Cols("12").Md("6")
			).Justify("center")
		).Fluid(true)
	)
)`,
			expected: `<v-app><v-main><v-container fluid><v-row justify="center"><v-col cols="12" md="6"><v-card><v-card-title>Card Title</v-card-title><v-card-text>This is card content with some descriptive text.<v-divider></v-divider><p>Paragraph below divider</p></v-card-text><v-card-actions><v-btn color="error">Cancel</v-btn><v-btn color="primary">Confirm</v-btn></v-card-actions></v-card></v-col></v-row></v-container></v-main></v-app>`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Wrap multi-line Go code in a complete Go program
			wrappedCode := fmt.Sprintf("package demo\n\nvar n = Body(\n%s\n)", tt.goCode)

			// Parse directly from string
			output, err := parse.GenerateGoHTML(wrappedCode, false, 2)
			if err != nil {
				t.Fatalf("GenerateGoHTML error: %v", err)
			}

			// Normalize output and expected result to eliminate whitespace differences
			normalizedOutput := NormalizeFormAttrOrder(parse.NormalizeHTML(output))
			normalizedExpected := NormalizeFormAttrOrder(parse.NormalizeHTML(tt.expected))

			if normalizedOutput != normalizedExpected {
				t.Errorf("Expected:\n%s\nGot:\n%s", normalizedExpected, normalizedOutput)
			}
		})
	}
}

// NormalizeFormAttrOrder normalizes form element attribute order to ignore order differences in comparisons
func NormalizeFormAttrOrder(html string) string {
	// Normalize type and name attributes order in input elements
	re := regexp.MustCompile(`<input([^>]*)type="([^"]*)"([^>]*)name="([^"]*)"([^>]*)/>`)
	html = re.ReplaceAllString(html, `<input$1name="$4"$3type="$2"$5/>`)

	// Normalize reverse case
	re = regexp.MustCompile(`<input([^>]*)name="([^"]*)"([^>]*)type="([^"]*)"([^>]*)/>`)
	html = re.ReplaceAllString(html, `<input$1name="$2"$3type="$4"$5/>`)

	// Normalize boolean attribute order, maintain consistent calling order
	re = regexp.MustCompile(`<input([^>]*)required([^>]*)disabled([^>]*)/>`)
	html = re.ReplaceAllString(html, `<input$1disabled$2required$3/>`)

	re = regexp.MustCompile(`<input([^>]*)disabled([^>]*)required([^>]*)/>`)
	html = re.ReplaceAllString(html, `<input$1disabled$2required$3/>`)

	// Normalize vx-btn/v-xbtn attribute order
	re = regexp.MustCompile(`<v-xbtn([^>]*)text="([^"]*)"([^>]*)color="([^"]*)"([^>]*)>`)
	html = re.ReplaceAllString(html, `<v-xbtn$1color="$4"$3text="$2"$5>`)

	re = regexp.MustCompile(`<v-xbtn([^>]*)color="([^"]*)"([^>]*)text="([^"]*)"([^>]*)>`)
	html = re.ReplaceAllString(html, `<v-xbtn$1color="$2"$3text="$4"$5>`)

	return html
}
