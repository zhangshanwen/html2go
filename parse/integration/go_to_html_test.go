package integration

import (
	"regexp"
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
			output, err := parse.GenerateGoHTML(tt.goCode, false, 2)
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

// 测试错误处理
func TestGoToHTMLErrors(t *testing.T) {
	invalidCode := `
	// 这不是有效的Go代码
	<div>Invalid HTML in Go</div>
	`

	_, err := parse.GenerateGoHTML(invalidCode, false, 2)
	if err == nil {
		t.Errorf("Expected error for invalid Go code, but got nil")
	}
}

// NormalizeFormAttrOrder 规范化表单元素属性顺序，以便在比较时忽略顺序差异
func NormalizeFormAttrOrder(html string) string {
	// 规范化input元素中type和name属性的顺序
	re := regexp.MustCompile(`<input([^>]*)type="([^"]*)"([^>]*)name="([^"]*)"([^>]*)/>`)
	html = re.ReplaceAllString(html, `<input$1name="$4"$3type="$2"$5/>`)

	// 规范化反向情况
	re = regexp.MustCompile(`<input([^>]*)name="([^"]*)"([^>]*)type="([^"]*)"([^>]*)/>`)
	html = re.ReplaceAllString(html, `<input$1name="$2"$3type="$4"$5/>`)

	// 规范化布尔属性顺序，保持调用顺序一致
	re = regexp.MustCompile(`<input([^>]*)required([^>]*)disabled([^>]*)/>`)
	html = re.ReplaceAllString(html, `<input$1disabled$2required$3/>`)

	re = regexp.MustCompile(`<input([^>]*)disabled([^>]*)required([^>]*)/>`)
	html = re.ReplaceAllString(html, `<input$1disabled$2required$3/>`)

	// 规范化vx-btn/v-xbtn属性顺序
	re = regexp.MustCompile(`<v-xbtn([^>]*)text="([^"]*)"([^>]*)color="([^"]*)"([^>]*)>`)
	html = re.ReplaceAllString(html, `<v-xbtn$1color="$4"$3text="$2"$5>`)

	re = regexp.MustCompile(`<v-xbtn([^>]*)color="([^"]*)"([^>]*)text="([^"]*)"([^>]*)>`)
	html = re.ReplaceAllString(html, `<v-xbtn$1color="$2"$3text="$4"$5>`)

	return html
}
