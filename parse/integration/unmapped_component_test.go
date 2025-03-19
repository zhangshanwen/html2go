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
		h.Tag("custom-component").Children(
			Span("Child content"),
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
	h.Tag("custom-grid").Attr("row-gap", "10px").
		Attr("column-gap", "15px").Children(
		h.Tag("custom-cell").Attr("span", "2").Children(
			Text("Cell 1"),
		),
		h.Tag("custom-cell").Attr("span", "1").Children(
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
		h.Tag("custom-header").Children(
			VBtn(
				Text("Action"),
			).Color("primary"),
		).Attr("title", "Page Title"),
		VXBtn().Text("VuetifyX Button"),
	),
)
`,
	},
	{
		Name: "Unknown vuetify component",
		HTML: `
<div>
  <v-unknown prop1="value1" prop2="value2">
    <span>Child content</span>
  </v-unknown>
</div>
`,
		GoCode: `package hello

var n = Body(
	Div(
		h.Tag("v-unknown").Children(
			Span("Child content"),
		).Attr("prop1", "value1").
			Attr("prop2", "value2"),
	),
)
`,
	},
	{
		Name: "Unknown vuetifyx component",
		HTML: `
<div>
  <vx-unknown prop1="value1" prop2="value2">
    <span>Child content</span>
  </vx-unknown>
</div>
`,
		GoCode: `package hello

var n = Body(
	Div(
		h.Tag("vx-unknown").Children(
			Span("Child content"),
		).Attr("prop1", "value1").
			Attr("prop2", "value2"),
	),
)
`,
	},
	{
		Name: "Mixed known and all types of unknown components",
		HTML: `
<div>
  <custom-component prop="value">Custom content</custom-component>
  <v-unknown color="primary">Unknown Vuetify</v-unknown>
  <vx-unknown text="Test">Unknown VuetifyX</vx-unknown>
  <v-btn color="success">Known component</v-btn>
</div>
`,
		GoCode: `package hello

var n = Body(
	Div(
		h.Tag("custom-component").Children(
			Text("Custom content"),
		).Attr("prop", "value"),
		h.Tag("v-unknown").Children(
			Text("Unknown Vuetify"),
		).Attr("color", "primary"),
		h.Tag("vx-unknown").Children(
			Text("Unknown VuetifyX"),
		).Attr("text", "Test"),
		VBtn(
			Text("Known component"),
		).Color("success"),
	),
)
`,
	},
}

// TestUnmappedComponents runs tests for components not in the mapping
func TestUnmappedComponents(t *testing.T) {
	RunHTMLGoTestCases(t, UnmappedComponentTestCases)
}
