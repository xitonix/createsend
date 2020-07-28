package clients

// Template represents a template.
type Template struct {
	// ID template id.
	ID string `json:"TemplateID"`
	// Name template name.
	Name string
	// PreviewURL the HTML preview URL of the template.
	PreviewURL string
	// ScreenshotURL the template's screenshot URL.
	ScreenshotURL string
}
