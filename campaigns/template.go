package campaigns

type EditableField struct {
	Content string
	ALT     string
	HREF    string
}

type RepeaterItem struct {
	Layout      []string
	SingleLines []EditableField
	MultiLines  []EditableField
	Images      []EditableField
}

type Repeater struct {
	Items []RepeaterItem
}

type TemplateContent struct {
	SingleLines []EditableField
	MultiLines  []EditableField
	Images      []EditableField
	Repeaters   []Repeater
}

type Template struct {
	BasicDetails
	ID              string `json:"TemplateID"`
	TemplateContent TemplateContent
}
