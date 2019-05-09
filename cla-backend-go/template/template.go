package template

type CLATemplate struct {
	Name        string
	Description string

	MetaFields []MetaField

	IclaFields []Field
	CclaFields []Field

	HtmlBody string
}

type MetaField struct {
	Name             string
	Description      string
	TemplateVariable string
}

type Field struct {
	AnchorString string
	Type         string
	IsOptional   bool
	IsEditable   bool
	Width        int
	Height       int
	OffsetX      int
	OffsetY      int
}
