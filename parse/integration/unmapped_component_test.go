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
		Tag("custom-component").Children(
			Span("Child content"),
		).Attr("prop1", "value1").
			Attr("prop2", "value2"),
	),
)
`,
	},
	{
		Name: "Custom component with pkg prefix",
		Pkg:  "h",
		HTML: `
<div>
  <custom-component prop="value">Child content</custom-component>
</div>
`,
		GoCode: `package hello

var n = h.Body(
	h.Div(
		h.Tag("custom-component").Children(
			h.Text("Child content"),
		).Attr("prop", "value"),
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
	Tag("custom-grid").Attr("row-gap", "10px").
		Attr("column-gap", "15px").Children(
		Tag("custom-cell").Attr("span", "2").Children(
			Text("Cell 1"),
		),
		Tag("custom-cell").Attr("span", "1").Children(
			Text("Cell 2"),
		),
	),
)
`,
	},
	{
		Name:        "Mixed known and unknown components",
		VuetifyPkg:  "v",
		VuetifyxPkg: "vx",
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
		Tag("custom-header").Children(
			v.VBtn(
				Text("Action"),
			).Color("primary"),
		).Attr("title", "Page Title"),
		vx.VXBtn().Text("VuetifyX Button"),
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
		Tag("v-unknown").Children(
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
		Tag("vx-unknown").Children(
			Span("Child content"),
		).Attr("prop1", "value1").
			Attr("prop2", "value2"),
	),
)
`,
	},
	{
		Name: "Unknown v- and vx- components with pkg prefix",
		Pkg:  "h",
		HTML: `
<div>
  <v-unknown prop1="value1">Unknown Vuetify</v-unknown>
  <vx-unknown prop2="value2">Unknown VuetifyX</vx-unknown>
</div>
`,
		GoCode: `package hello

var n = h.Body(
	h.Div(
		h.Tag("v-unknown").Children(
			h.Text("Unknown Vuetify"),
		).Attr("prop1", "value1"),
		h.Tag("vx-unknown").Children(
			h.Text("Unknown VuetifyX"),
		).Attr("prop2", "value2"),
	),
)
`,
	},
	{
		Name:        "Mixed known and all types of unknown components",
		VuetifyPkg:  "v",
		VuetifyxPkg: "vx",
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
		Tag("custom-component").Children(
			Text("Custom content"),
		).Attr("prop", "value"),
		Tag("v-unknown").Children(
			Text("Unknown Vuetify"),
		).Attr("color", "primary"),
		Tag("vx-unknown").Children(
			Text("Unknown VuetifyX"),
		).Attr("text", "Test"),
		v.VBtn(
			Text("Known component"),
		).Color("success"),
	),
)
`,
	},
	{
		Name:        "All component types with all package prefixes",
		Pkg:         "h",
		VuetifyPkg:  "v",
		VuetifyxPkg: "vx",
		HTML: `
<div>
  <span>Standard HTML element</span>
  <custom-component>Custom element</custom-component>
  <v-btn>Known Vuetify component</v-btn>
  <vx-btn>Known VuetifyX component</vx-btn>
  <v-unknown>Unknown Vuetify component</v-unknown>
  <vx-unknown>Unknown VuetifyX component</vx-unknown>
</div>
`,
		GoCode: `package hello

var n = h.Body(
	h.Div(
		h.Span("Standard HTML element"),
		h.Tag("custom-component").Children(
			h.Text("Custom element"),
		),
		v.VBtn(
			h.Text("Known Vuetify component"),
		),
		vx.VXBtn(
			h.Text("Known VuetifyX component"),
		),
		h.Tag("v-unknown").Children(
			h.Text("Unknown Vuetify component"),
		),
		h.Tag("vx-unknown").Children(
			h.Text("Unknown VuetifyX component"),
		),
	),
)
`,
	},
}

// TestUnmappedComponents runs tests for components not in the mapping
func TestUnmappedComponents(t *testing.T) {
	RunHTMLGoTestCases(t, UnmappedComponentTestCases)
}
