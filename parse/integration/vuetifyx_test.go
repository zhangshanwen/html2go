package integration

import (
	"testing"
)

// VuetifyXTestCases contains test scenarios for VuetifyX components
var VuetifyXTestCases = []HTMLGoTestCase{
	{
		Name:        "empty-pkg with empty VuetifyxPkg",
		Pkg:         "",
		VuetifyxPkg: "",
		HTML: `
<div>
  <vx-btn text="Submit" color="primary" on-click="handleSubmit"></vx-btn>
</div>
`,
		GoCode: `package hello

var n = Body(
	Div(
		VXBtn().Text("Submit").
			Color("primary").
			OnClick("handleSubmit"),
	),
)
`,
	},
	{
		Name:        "h-pkg with empty VuetifyxPkg",
		Pkg:         "h",
		VuetifyxPkg: "",
		HTML: `
<div>
  <vx-btn text="Submit" color="primary" on-click="handleSubmit"></vx-btn>
</div>
`,
		GoCode: `package hello

var n = h.Body(
	h.Div(
		VXBtn().Text("Submit").
			Color("primary").
			OnClick("handleSubmit"),
	),
)
`,
	},
	{
		Name:        "vx-btn basic",
		VuetifyxPkg: "vx",
		HTML: `
<div>
  <vx-btn text="Submit" color="primary" on-click="handleSubmit"></vx-btn>
</div>
`,
		GoCode: `package hello

var n = Body(
	Div(
		vx.VXBtn().Text("Submit").
			Color("primary").
			OnClick("handleSubmit"),
	),
)
`,
	},
	{
		Name:        "vuetifyx-btn basic",
		VuetifyxPkg: "vuetifyx",
		HTML: `
<div>
  <vx-btn text="Submit" color="primary" on-click="handleSubmit"></vx-btn>
</div>
`,
		GoCode: `package hello

var n = Body(
	Div(
		vuetifyx.VXBtn().Text("Submit").
			Color("primary").
			OnClick("handleSubmit"),
	),
)
`,
	},
	{
		Name:        "vvvx-btn basic",
		VuetifyxPkg: "vvvx",
		HTML: `
<div>
  <vx-btn text="Submit" color="primary" on-click="handleSubmit"></vx-btn>
</div>
`,
		GoCode: `package hello

var n = Body(
	Div(
		vvvx.VXBtn().Text("Submit").
			Color("primary").
			OnClick("handleSubmit"),
	),
)
`,
	},
	{
		Name:        "vuix-btn basic",
		VuetifyxPkg: "vuix",
		HTML: `
<div>
  <vx-btn text="Submit" color="primary" on-click="handleSubmit"></vx-btn>
</div>
`,
		GoCode: `package hello

var n = Body(
	Div(
		vuix.VXBtn().Text("Submit").
			Color("primary").
			OnClick("handleSubmit"),
	),
)
`,
	},
	{
		Name:        "vx-btn with boolean attributes",
		VuetifyxPkg: "vx",
		HTML: `
<div>
  <vx-btn text="Submit" disabled flat block></vx-btn>
</div>
`,
		GoCode: `package hello

var n = Body(
	Div(
		vx.VXBtn().Text("Submit").
			Disabled(true).
			Flat(true).
			Block(true),
	),
)
`,
	},
	{
		Name:        "vuetifyx-btn with boolean attributes",
		VuetifyxPkg: "vuetifyx",
		HTML: `
<div>
  <vx-btn text="Submit" disabled flat block></vx-btn>
</div>
`,
		GoCode: `package hello

var n = Body(
	Div(
		vuetifyx.VXBtn().Text("Submit").
			Disabled(true).
			Flat(true).
			Block(true),
	),
)
`,
	},
	{
		Name:        "vx-dialog with nested content",
		VuetifyxPkg: "vx",
		HTML: `
<vx-dialog title="Confirmation" persistent width="400" max-width="500">
  <p>Are you sure you want to delete this item?</p>
  <vx-btn text="Cancel" color="grey"></vx-btn>
  <vx-btn text="Delete" color="red"></vx-btn>
</vx-dialog>
`,
		GoCode: `package hello

var n = Body(
	vx.VXDialog(
		P(
			Text("Are you sure you want to delete this item?"),
		),
		vx.VXBtn().Text("Cancel").
			Color("grey"),
		vx.VXBtn().Text("Delete").
			Color("red"),
	).Title("Confirmation").
		Persistent(true).
		Width(400).
		MaxWidth(500),
)
`,
	},
	{
		Name:        "vuetifyx-dialog with nested content",
		VuetifyxPkg: "vuetifyx",
		HTML: `
<vx-dialog title="Confirmation" persistent width="400" max-width="500">
  <p>Are you sure you want to delete this item?</p>
  <vx-btn text="Cancel" color="grey"></vx-btn>
  <vx-btn text="Delete" color="red"></vx-btn>
</vx-dialog>
`,
		GoCode: `package hello

var n = Body(
	vuetifyx.VXDialog(
		P(
			Text("Are you sure you want to delete this item?"),
		),
		vuetifyx.VXBtn().Text("Cancel").
			Color("grey"),
		vuetifyx.VXBtn().Text("Delete").
			Color("red"),
	).Title("Confirmation").
		Persistent(true).
		Width(400).
		MaxWidth(500),
)
`,
	},
	{
		Name:         "vx-dialog with children mode",
		VuetifyxPkg:  "vx",
		ChildrenMode: true,
		HTML: `
<vx-dialog title="Confirmation" persistent>
  <p>Are you sure you want to delete this item?</p>
</vx-dialog>
`,
		GoCode: `package hello

var n = Body().Children(
	vx.VXDialog().Title("Confirmation").
		Persistent(true).Children(
		P().Children(
			Text("Are you sure you want to delete this item?"),
		),
	),
)
`,
	},
	{
		Name:         "vvvx-dialog with children mode",
		VuetifyxPkg:  "vvvx",
		ChildrenMode: true,
		HTML: `
<vx-dialog title="Confirmation" persistent>
  <p>Are you sure you want to delete this item?</p>
</vx-dialog>
`,
		GoCode: `package hello

var n = Body().Children(
	vvvx.VXDialog().Title("Confirmation").
		Persistent(true).Children(
		P().Children(
			Text("Are you sure you want to delete this item?"),
		),
	),
)
`,
	},
	{
		Name:        "vx-select with items",
		VuetifyxPkg: "vx",
		HTML: `
<vx-select label="Choose a country" clearable multiple required>
  <option value="us">United States</option>
  <option value="ca">Canada</option>
  <option value="mx">Mexico</option>
</vx-select>
`,
		GoCode: `package hello

var n = Body(
	vx.VXSelect(
		Option("United States").Value("us"),
		Option("Canada").Value("ca"),
		Option("Mexico").Value("mx"),
	).Label("Choose a country").
		Clearable(true).
		Multiple(true).
		Required(true),
)
`,
	},
	{
		Name:        "vuetifyx-select with items",
		VuetifyxPkg: "vuetifyx",
		HTML: `
<vx-select label="Choose a country" clearable multiple required>
  <option value="us">United States</option>
  <option value="ca">Canada</option>
  <option value="mx">Mexico</option>
</vx-select>
`,
		GoCode: `package hello

var n = Body(
	vuetifyx.VXSelect(
		Option("United States").Value("us"),
		Option("Canada").Value("ca"),
		Option("Mexico").Value("mx"),
	).Label("Choose a country").
		Clearable(true).
		Multiple(true).
		Required(true),
)
`,
	},
	{
		Name:        "vx-checkbox with attributes",
		VuetifyxPkg: "vx",
		HTML: `
<vx-checkbox label="Accept terms" model-value="true" error-messages="You must accept the terms"></vx-checkbox>
`,
		GoCode: `package hello

var n = Body(
	vx.VXCheckbox().Label("Accept terms").
		ModelValue("true").
		ErrorMessages("You must accept the terms"),
)
`,
	},
	{
		Name:        "vuetifyx-checkbox with attributes",
		VuetifyxPkg: "vuetifyx",
		HTML: `
<vx-checkbox label="Accept terms" model-value="true" error-messages="You must accept the terms"></vx-checkbox>
`,
		GoCode: `package hello

var n = Body(
	vuetifyx.VXCheckbox().Label("Accept terms").
		ModelValue("true").
		ErrorMessages("You must accept the terms"),
)
`,
	},
	{
		Name:        "vvvx-checkbox with attributes",
		VuetifyxPkg: "vvvx",
		HTML: `
<vx-checkbox label="Accept terms" model-value="true" error-messages="You must accept the terms"></vx-checkbox>
`,
		GoCode: `package hello

var n = Body(
	vvvx.VXCheckbox().Label("Accept terms").
		ModelValue("true").
		ErrorMessages("You must accept the terms"),
)
`,
	},
	{
		Name:        "vuetifyx-datepicker example",
		VuetifyxPkg: "vuetifyx",
		HTML: `
<vx-date-picker label="Select date" model-value="2023-01-15" range first-day-of-week="1"></vx-date-picker>
`,
		GoCode: `package hello

var n = Body(
	vuetifyx.VXDatepicker().Attr("label", "Select date").
		Attr("model-value", "2023-01-15").
		Attr("range", "").
		Attr("first-day-of-week", "1"),
)
`,
	},
	{
		Name:        "vvvx-components example",
		VuetifyxPkg: "vvvx",
		HTML: `
<vx-tabs>
  <vx-tab title="Dashboard">
    <vx-chart type="bar" :data="chartData" height="300"></vx-chart>
  </vx-tab>
  <vx-tab title="Settings">
    <vx-settings-panel user-id="123"></vx-settings-panel>
  </vx-tab>
</vx-tabs>
`,
		GoCode: `package hello

var n = Body(
	vvvx.VXTabs(
		Tag("vx-tab").Children(
			vvvx.VXChart().Attr("type", "bar").
				Attr("x-bind:data", "chartData").
				Attr("height", "300"),
		).Attr("title", "Dashboard"),
		Tag("vx-tab").Children(
			Tag("vx-settings-panel").Attr("user-id", "123"),
		).Attr("title", "Settings"),
	),
)
`,
	},
	{
		Name:        "vuix-component complex",
		VuetifyxPkg: "vuix",
		HTML: `
<vx-data-grid :columns="columns" :rows="rows" sortable pagination>
  <template v-slot:toolbar>
    <vx-search-field placeholder="Search data..."></vx-search-field>
    <vx-export-button formats="csv,excel"></vx-export-button>
  </template>
</vx-data-grid>
`,
		GoCode: `package hello

var n = Body(
	Tag("vx-data-grid").Children(
		Template(
			Tag("vx-search-field").Attr("placeholder", "Search data..."),
			Tag("vx-export-button").Attr("formats", "csv,excel"),
		).Attr("v-slot:toolbar", ""),
	).Attr("x-bind:columns", "columns").
		Attr("x-bind:rows", "rows").
		Attr("sortable", "").
		Attr("pagination", ""),
)
`,
	},
}

// TestVuetifyX runs all VuetifyX component test cases
func TestVuetifyX(t *testing.T) {
	RunHTMLGoTestCases(t, VuetifyXTestCases)
}
