package integration

import (
	"testing"
)

// VuetifyXTestCases contains test scenarios for VuetifyX components
var VuetifyXTestCases = []HTMLGoTestCase{
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
}

// TestVuetifyX runs all VuetifyX component test cases
func TestVuetifyX(t *testing.T) {
	RunHTMLGoTestCases(t, VuetifyXTestCases)
}
