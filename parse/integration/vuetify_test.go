package integration

import (
	"testing"
)

// VuetifyTestCases contains test scenarios for Vuetify components
var VuetifyTestCases = []HTMLGoTestCase{
	{
		Name:       "v-btn basic",
		VuetifyPkg: "v",
		HTML: `
<div>
  <v-btn color="primary">Submit</v-btn>
</div>
`,
		GoCode: `package hello

var n = Body(
	Div(
		v.VBtn(
			Text("Submit"),
		).Color("primary"),
	),
)
`,
	},
	{
		Name:       "vuetify-btn basic",
		VuetifyPkg: "vuetify",
		HTML: `
<div>
  <v-btn color="primary">Submit</v-btn>
</div>
`,
		GoCode: `package hello

var n = Body(
	Div(
		vuetify.VBtn(
			Text("Submit"),
		).Color("primary"),
	),
)
`,
	},
	{
		Name:       "vvv-btn basic",
		VuetifyPkg: "vvv",
		HTML: `
<div>
  <v-btn color="primary">Submit</v-btn>
</div>
`,
		GoCode: `package hello

var n = Body(
	Div(
		vvv.VBtn(
			Text("Submit"),
		).Color("primary"),
	),
)
`,
	},
	{
		Name:       "vui-btn basic",
		VuetifyPkg: "vui",
		HTML: `
<div>
  <v-btn color="primary">Submit</v-btn>
</div>
`,
		GoCode: `package hello

var n = Body(
	Div(
		vui.VBtn(
			Text("Submit"),
		).Color("primary"),
	),
)
`,
	},
	{
		Name:       "v-card with content",
		VuetifyPkg: "v",
		HTML: `
<v-card>
  <v-card-title>Card Title</v-card-title>
  <v-card-text>This is the card content with some text.</v-card-text>
  <v-card-actions>
    <v-btn color="primary">Action</v-btn>
    <v-btn color="secondary">Cancel</v-btn>
  </v-card-actions>
</v-card>
`,
		GoCode: `package hello

var n = Body(
	v.VCard(
		v.VCardTitle(
			Text("Card Title"),
		),
		v.VCardText(
			Text("This is the card content with some text."),
		),
		v.VCardActions(
			v.VBtn(
				Text("Action"),
			).Color("primary"),
			v.VBtn(
				Text("Cancel"),
			).Color("secondary"),
		),
	),
)
`,
	},
	{
		Name:       "vuetify-card with content",
		VuetifyPkg: "vuetify",
		HTML: `
<v-card>
  <v-card-title>Card Title</v-card-title>
  <v-card-text>This is the card content with some text.</v-card-text>
  <v-card-actions>
    <v-btn color="primary">Action</v-btn>
    <v-btn color="secondary">Cancel</v-btn>
  </v-card-actions>
</v-card>
`,
		GoCode: `package hello

var n = Body(
	vuetify.VCard(
		vuetify.VCardTitle(
			Text("Card Title"),
		),
		vuetify.VCardText(
			Text("This is the card content with some text."),
		),
		vuetify.VCardActions(
			vuetify.VBtn(
				Text("Action"),
			).Color("primary"),
			vuetify.VBtn(
				Text("Cancel"),
			).Color("secondary"),
		),
	),
)
`,
	},
	{
		Name:       "v-form with inputs",
		VuetifyPkg: "v",
		HTML: `
<v-form>
  <v-text-field label="Name" required></v-text-field>
  <v-text-field label="Email" type="email"></v-text-field>
  <v-checkbox label="Subscribe to newsletter"></v-checkbox>
  <v-btn type="submit" color="success">Submit</v-btn>
</v-form>
`,
		GoCode: `package hello

var n = Body(
	v.VForm(
		v.VTextField().Label("Name").
			Attr("required", ""),
		v.VTextField().Label("Email").
			Type("email"),
		v.VCheckbox().Label("Subscribe to newsletter"),
		v.VBtn(
			Text("Submit"),
		).Attr("type", "submit").
			Color("success"),
	),
)
`,
	},
	{
		Name:       "vuetify-form with inputs",
		VuetifyPkg: "vuetify",
		HTML: `
<v-form>
  <v-text-field label="Name" required></v-text-field>
  <v-text-field label="Email" type="email"></v-text-field>
  <v-checkbox label="Subscribe to newsletter"></v-checkbox>
  <v-btn type="submit" color="success">Submit</v-btn>
</v-form>
`,
		GoCode: `package hello

var n = Body(
	vuetify.VForm(
		vuetify.VTextField().Label("Name").
			Attr("required", ""),
		vuetify.VTextField().Label("Email").
			Type("email"),
		vuetify.VCheckbox().Label("Subscribe to newsletter"),
		vuetify.VBtn(
			Text("Submit"),
		).Attr("type", "submit").
			Color("success"),
	),
)
`,
	},
	{
		Name:       "vvv-form with inputs",
		VuetifyPkg: "vvv",
		HTML: `
<v-form>
  <v-text-field label="Name" required></v-text-field>
  <v-text-field label="Email" type="email"></v-text-field>
  <v-checkbox label="Subscribe to newsletter"></v-checkbox>
  <v-btn type="submit" color="success">Submit</v-btn>
</v-form>
`,
		GoCode: `package hello

var n = Body(
	vvv.VForm(
		vvv.VTextField().Label("Name").
			Attr("required", ""),
		vvv.VTextField().Label("Email").
			Type("email"),
		vvv.VCheckbox().Label("Subscribe to newsletter"),
		vvv.VBtn(
			Text("Submit"),
		).Attr("type", "submit").
			Color("success"),
	),
)
`,
	},
}

// TestVuetify runs all Vuetify component test cases
func TestVuetify(t *testing.T) {
	RunHTMLGoTestCases(t, VuetifyTestCases)
}
