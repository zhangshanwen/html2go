package integration

import (
	"testing"
)

// VuetifyTestCases contains test scenarios for Vuetify components
var VuetifyTestCases = []HTMLGoTestCase{
	{
		Name: "v-btn basic",
		HTML: `
<div>
  <v-btn color="primary">Submit</v-btn>
</div>
`,
		GoCode: `package hello

var n = Body(
	Div(
		VBtn(
			Text("Submit"),
		).Color("primary"),
	),
)
`,
	},
	{
		Name: "v-card with content",
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
	VCard(
		VCardTitle(
			Text("Card Title"),
		),
		VCardText(
			Text("This is the card content with some text."),
		),
		VCardActions(
			VBtn(
				Text("Action"),
			).Color("primary"),
			VBtn(
				Text("Cancel"),
			).Color("secondary"),
		),
	),
)
`,
	},
	{
		Name: "v-form with inputs",
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
	VForm(
		VTextField().Label("Name").
			Required(true),
		VTextField().Label("Email").
			Type("email"),
		VCheckbox().Label("Subscribe to newsletter"),
		VBtn(
			Text("Submit"),
		).Type("submit").
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
