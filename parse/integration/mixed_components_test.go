package integration

import (
	"testing"
)

// MixedComponentsTestCases contains test scenarios for mixed component usage
var MixedComponentsTestCases = []HTMLGoTestCase{
	{
		Name:        "Standard and Vuetify components with empty pkg",
		Pkg:         "",
		VuetifyPkg:  "v",
		VuetifyxPkg: "",
		HTML: `
<div class="container">
  <h1>Mixed Components Example</h1>
  <v-card elevation="2">
    <v-card-title>Card Title</v-card-title>
    <v-card-text>
      <p>Standard HTML inside Vuetify component</p>
      <span class="highlight">Important information</span>
    </v-card-text>
    <v-card-actions>
      <v-btn color="primary">Action</v-btn>
    </v-card-actions>
  </v-card>
</div>
`,
		GoCode: `package hello

var n = Body(
	Div(
		H1("Mixed Components Example"),
		v.VCard(
			v.VCardTitle(
				Text("Card Title"),
			),
			v.VCardText(
				P(
					Text("Standard HTML inside Vuetify component"),
				),
				Span("Important information").Class("highlight"),
			),
			v.VCardActions(
				v.VBtn(
					Text("Action"),
				).Color("primary"),
			),
		).Elevation("2"),
	).Class("container"),
)
`,
	},
	{
		Name:        "HTML, Vuetify and VuetifyX with h pkg",
		Pkg:         "h",
		VuetifyPkg:  "v",
		VuetifyxPkg: "vx",
		HTML: `
<div>
  <h2>Component Showcase</h2>
  <v-card>
    <v-card-title>Card with VuetifyX component</v-card-title>
    <v-card-text>
      <vx-btn color="success" text="Save Data"></vx-btn>
      <p>Standard paragraph with <strong>bold text</strong></p>
    </v-card-text>
  </v-card>
</div>
`,
		GoCode: `package hello

var n = h.Body(
	h.Div(
		h.H2("Component Showcase"),
		v.VCard(
			v.VCardTitle(
				h.Text("Card with VuetifyX component"),
			),
			v.VCardText(
				vx.VXBtn().Color("success").
					Text("Save Data"),
				h.P(
					h.Text("Standard paragraph with"),
					h.Strong("bold text"),
				),
			),
		),
	),
)
`,
	},
	{
		Name:        "Standard, Vuetify, VuetifyX, and unmapped components",
		Pkg:         "html",
		VuetifyPkg:  "vuetify",
		VuetifyxPkg: "vuetifyx",
		HTML: `
<div>
  <h3>Complex Layout</h3>
  <custom-component title="My Custom Component">
    <v-btn color="primary" x-large>Vuetify Button</v-btn>
    <vx-select label="Options" :items="['Option 1', 'Option 2']"></vx-select>
    <span>Standard HTML element</span>
    <v-unknown custom-prop="test">Unknown component</v-unknown>
  </custom-component>
</div>
`,
		GoCode: `package hello

var n = html.Body(
	html.Div(
		html.H3("Complex Layout"),
		html.Tag("custom-component").Children(
			vuetify.VBtn(
				html.Text("Vuetify Button"),
			).Color("primary").
				Attr("x-large", ""),
			vuetifyx.VXSelect().Label("Options").
				Attr("x-bind:items", |backquote|["Option 1", "Option 2"]|backquote|),
			html.Span("Standard HTML element"),
			html.Tag("v-unknown").Children(
				html.Text("Unknown component"),
			).Attr("custom-prop", "test"),
		).Attr("title", "My Custom Component"),
	),
)
`,
	},
	{
		Name:        "Empty and custom package prefixes",
		Pkg:         "",
		VuetifyPkg:  "",
		VuetifyxPkg: "",
		HTML: `
<div>
  <h4>No Package Prefixes</h4>
  <v-btn>Vuetify Button</v-btn>
  <vx-btn>VuetifyX Button</vx-btn>
  <custom-element>Custom Element</custom-element>
</div>
`,
		GoCode: `package hello

var n = Body(
	Div(
		H4("No Package Prefixes"),
		VBtn(
			Text("Vuetify Button"),
		),
		VXBtn(
			Text("VuetifyX Button"),
		),
		Tag("custom-element").Children(
			Text("Custom Element"),
		),
	),
)
`,
	},
	{
		Name:        "HTML with pkg and empty Vuetify Packages",
		Pkg:         "h",
		VuetifyPkg:  "",
		VuetifyxPkg: "",
		HTML: `
<div>
  <h1>HTML with h package</h1>
  <v-btn>Vuetify with no prefix</v-btn>
  <vx-btn>VuetifyX with no prefix</vx-btn>
  <custom-element>Custom Element with h</custom-element>
</div>
`,
		GoCode: `package hello

var n = h.Body(
	h.Div(
		h.H1("HTML with h package"),
		VBtn(
			h.Text("Vuetify with no prefix"),
		),
		VXBtn(
			h.Text("VuetifyX with no prefix"),
		),
		h.Tag("custom-element").Children(
			h.Text("Custom Element with h"),
		),
	),
)
`,
	},
	{
		Name:        "Nested mixed components with different packages",
		Pkg:         "html",
		VuetifyPkg:  "vui",
		VuetifyxPkg: "vuix",
		HTML: `
<div>
  <v-container>
    <v-row>
      <v-col cols="6">
        <vx-card title="Card Title">
          <custom-component class="my-custom">
            <h4>Title inside custom component</h4>
            <p>Text content</p>
          </custom-component>
          <v-unknown prop="value"></v-unknown>
        </vx-card>
      </v-col>
      <v-col cols="6">
        <v-card>
          <v-card-title>Standard Card</v-card-title>
          <v-card-text>Content text</v-card-text>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</div>
`,
		GoCode: `package hello

var n = html.Body(
	html.Div(
		vui.VContainer(
			vui.VRow(
				vui.VCol(
					html.Tag("vx-card").Children(
						html.Tag("custom-component").Children(
							html.H4("Title inside custom component"),
							html.P(
								html.Text("Text content"),
							),
						).Attr("class", "my-custom"),
						html.Tag("v-unknown").Attr("prop", "value"),
					).Attr("title", "Card Title"),
				).Cols("6"),
				vui.VCol(
					vui.VCard(
						vui.VCardTitle(
							html.Text("Standard Card"),
						),
						vui.VCardText(
							html.Text("Content text"),
						),
					),
				).Cols("6"),
			),
		),
	),
)
`,
	},
	{
		Name:        "Form with mixed components and multiple Vuetify elements",
		Pkg:         "h",
		VuetifyPkg:  "v",
		VuetifyxPkg: "vx",
		HTML: `
<form>
  <v-text-field label="Username" required></v-text-field>
  <v-text-field label="Password" type="password" required></v-text-field>
  <vx-checkbox label="Remember me"></vx-checkbox>
  <custom-captcha verify="true"></custom-captcha>
  <div class="actions">
    <v-btn type="submit" color="primary">Login</v-btn>
    <v-btn type="button" text>Cancel</v-btn>
  </div>
</form>
`,
		GoCode: `package hello

var n = h.Body(
	h.Form(
		v.VTextField().Label("Username").
			Attr("required", ""),
		v.VTextField().Label("Password").
			Type("password").
			Attr("required", ""),
		vx.VXCheckbox().Label("Remember me"),
		h.Tag("custom-captcha").Attr("verify", "true"),
		h.Div(
			v.VBtn(
				h.Text("Login"),
			).Attr("type", "submit").
				Color("primary"),
			v.VBtn(
				h.Text("Cancel"),
			).Attr("type", "button").
				Text(""),
		).Class("actions"),
	),
)
`,
	},
	{
		Name:        "All types of components with boolean attributes",
		Pkg:         "web",
		VuetifyPkg:  "v",
		VuetifyxPkg: "vx",
		HTML: `
<div>
  <input type="checkbox" checked disabled />
  <v-checkbox label="Vuetify Checkbox" input-value="true" readonly></v-checkbox>
  <vx-switch label="VuetifyX Switch" active></vx-switch>
  <custom-toggle enabled active></custom-toggle>
  <vx-unknown visible hidden></vx-unknown>
</div>
`,
		GoCode: `package hello

var n = web.Body(
	web.Div(
		web.Input("").Type("checkbox").
			Checked(true).
			Disabled(true),
		v.VCheckbox().Label("Vuetify Checkbox").
			Attr("input-value", "true").
			Readonly(true),
		web.Tag("vx-switch").Attr("label", "VuetifyX Switch").
			Attr("active", ""),
		web.Tag("custom-toggle").Attr("enabled", "").
			Attr("active", ""),
		web.Tag("vx-unknown").Attr("visible", "").
			Attr("hidden", ""),
	),
)
`,
	},
	{
		Name:         "Children mode with mixed components",
		ChildrenMode: true,
		Pkg:          "h",
		VuetifyPkg:   "v",
		VuetifyxPkg:  "vx",
		HTML: `
<div>
  <v-container>
    <v-row>
      <v-col>
        <vx-card>
          <h3>Card Title</h3>
          <p>Card content</p>
          <custom-element></custom-element>
        </vx-card>
      </v-col>
    </v-row>
  </v-container>
</div>
`,
		GoCode: `package hello

var n = h.Body().Children(
	h.Div().Children(
		v.VContainer().Children(
			v.VRow().Children(
				v.VCol().Children(
					h.Tag("vx-card").Children(
						h.H3("Card Title"),
						h.P().Children(
							h.Text("Card content"),
						),
						h.Tag("custom-element"),
					),
				),
			),
		),
	),
)
`,
	},
	{
		Name:        "Unknown components that resemble Vuetify and VuetifyX",
		Pkg:         "h",
		VuetifyPkg:  "v",
		VuetifyxPkg: "vx",
		HTML: `
<div>
  <v-something title="Similar to Vuetify">Content</v-something>
  <vx-something title="Similar to VuetifyX">Content</vx-something>
  <v-btn color="primary">Real Vuetify</v-btn>
  <vx-btn text="Submit">Real VuetifyX</vx-btn>
</div>
`,
		GoCode: `package hello

var n = h.Body(
	h.Div(
		h.Tag("v-something").Children(
			h.Text("Content"),
		).Attr("title", "Similar to Vuetify"),
		h.Tag("vx-something").Children(
			h.Text("Content"),
		).Attr("title", "Similar to VuetifyX"),
		v.VBtn(
			h.Text("Real Vuetify"),
		).Color("primary"),
		vx.VXBtn(
			h.Text("Real VuetifyX"),
		).Text("Submit"),
	),
)
`,
	},
	{
		Name:        "Deeply nested mixed components",
		Pkg:         "h",
		VuetifyPkg:  "v",
		VuetifyxPkg: "vx",
		HTML: `
<v-container>
  <v-row>
    <v-col cols="12">
      <v-card>
        <v-card-title>
          <span class="headline">Complex Form</span>
          <v-spacer></v-spacer>
          <v-btn icon><v-icon>close</v-icon></v-btn>
        </v-card-title>
        <v-card-text>
          <vx-form>
            <div class="form-group">
              <vx-text-field label="Name"></vx-text-field>
              <v-text-field label="Email" type="email"></v-text-field>
              <custom-field label="Custom Input"></custom-field>
            </div>
          </vx-form>
        </v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn color="primary" text>Cancel</v-btn>
          <vx-btn primary>Submit</vx-btn>
        </v-card-actions>
      </v-card>
    </v-col>
  </v-row>
</v-container>
`,
		GoCode: `package hello

var n = h.Body(
	v.VContainer(
		v.VRow(
			v.VCol(
				v.VCard(
					v.VCardTitle(
						h.Span("Complex Form").Class("headline"),
						v.VSpacer(),
						v.VBtn(
							v.VIcon(
								h.Text("close"),
							),
						).Icon(""),
					),
					v.VCardText(
						h.Tag("vx-form").Children(
							h.Div(
								h.Tag("vx-text-field").Attr("label", "Name"),
								v.VTextField().Label("Email").
									Type("email"),
								h.Tag("custom-field").Attr("label", "Custom Input"),
							).Class("form-group"),
						),
					),
					v.VCardActions(
						v.VSpacer(),
						v.VBtn(
							h.Text("Cancel"),
						).Color("primary").
							Text(""),
						vx.VXBtn(
							h.Text("Submit"),
						).Attr("primary", ""),
					),
				),
			).Cols("12"),
		),
	),
)
`,
	},
}

// TestMixedComponents runs tests for different combinations of components
func TestMixedComponents(t *testing.T) {
	RunHTMLGoTestCases(t, MixedComponentsTestCases)
}
