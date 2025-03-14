package integration

import (
	"testing"
)

// UnmappedComponentTestCases contains test scenarios for components not in mapping
var UnmappedComponentTestCases = []HTMLGoTestCase{
	{
		Name: "Custom component",
		HTML: `
<div>
  <custom-component prop1="value1" prop2="value2">
    <span>Child content</span>
  </custom-component>
</div>
`,
		GoCode: `package hello

var n = Body(
	Div(
		CustomComponent(
			Span(
				Text("Child content"),
			),
		).Attr("prop1", "value1").
			Attr("prop2", "value2"),
	),
)
`,
	},
	{
		Name:         "Child mode custom component",
		ChildrenMode: true,
		HTML: `
<custom-grid row-gap="10px" column-gap="15px">
  <custom-cell span="2">Cell 1</custom-cell>
  <custom-cell span="1">Cell 2</custom-cell>
</custom-grid>
`,
		GoCode: `package hello

var n = Body().Children(
	CustomGrid().Attr("row-gap", "10px").
		Attr("column-gap", "15px").Children(
		CustomCell().Attr("span", "2").Children(
			Text("Cell 1"),
		),
		CustomCell().Attr("span", "1").Children(
			Text("Cell 2"),
		),
	),
)
`,
	},
	{
		Name: "Mixed known and unknown components",
		HTML: `
<div>
  <custom-header title="Page Title">
    <v-btn color="primary">Action</v-btn>
  </custom-header>
  <vx-btn text="VuetifyX Button"></vx-btn>
</div>
`,
		GoCode: `package hello

var n = Body(
	Div(
		CustomHeader(
			VBtn(
				Text("Action"),
			).Color("primary"),
		).Attr("title", "Page Title"),
		VXBtn().Text("VuetifyX Button"),
	),
)
`,
	},
}

// TestUnmappedComponents runs tests for components not in the mapping
func TestUnmappedComponents(t *testing.T) {
	RunHTMLGoTestCases(t, UnmappedComponentTestCases)
}
