package repository

// StarterDocument mirrors the Elasticsearch document structure optimized for search.
type StarterDocument struct {
	ID               int64  `json:"id"`
	Domain           string `json:"domain"`
	Name             string `json:"name"`
	DepartmentName   string `json:"department_name,omitempty"`
	BusinessUnitName string `json:"business_unit_name,omitempty"`

	FullText     string   `json:"full_text"`
	SearchTokens []string `json:"search_tokens"`
}

// StarterDocumentBuilder ==========================Builder==========================
type StarterDocumentBuilder struct {
	doc StarterDocument
}

// NewStarterDocumentBuilder creates a new builder instance.
func NewStarterDocumentBuilder() *StarterDocumentBuilder {
	return &StarterDocumentBuilder{
		doc: StarterDocument{},
	}
}

// ID sets the ID field.
func (b *StarterDocumentBuilder) ID(id int64) *StarterDocumentBuilder {
	b.doc.ID = id
	return b
}

// Domain sets the Domain field.
func (b *StarterDocumentBuilder) Domain(domain string) *StarterDocumentBuilder {
	b.doc.Domain = domain
	return b
}

// Name sets the Name field.
func (b *StarterDocumentBuilder) Name(name string) *StarterDocumentBuilder {
	b.doc.Name = name
	return b
}

func (b *StarterDocumentBuilder) DepartmentName(departmentName string) *StarterDocumentBuilder {
	b.doc.DepartmentName = departmentName
	return b
}

// BusinessUnitName sets the BusinessUnitName field.
func (b *StarterDocumentBuilder) BusinessUnitName(businessUnitName string) *StarterDocumentBuilder {
	b.doc.BusinessUnitName = businessUnitName
	return b
}

// FullText sets the FullText field.
func (b *StarterDocumentBuilder) FullText(fullText string) *StarterDocumentBuilder {
	b.doc.FullText = fullText
	return b
}

// SearchTokens sets the SearchTokens field.
func (b *StarterDocumentBuilder) SearchTokens(tokens []string) *StarterDocumentBuilder {
	b.doc.SearchTokens = tokens
	return b
}

// AddSearchToken adds a single token to SearchTokens.
func (b *StarterDocumentBuilder) AddSearchToken(token string) *StarterDocumentBuilder {
	b.doc.SearchTokens = append(b.doc.SearchTokens, token)
	return b
}

// Build returns the constructed StarterDocument.
func (b *StarterDocumentBuilder) Build() StarterDocument {
	return b.doc
}

// BuildPtr returns a pointer to the constructed StarterDocument.
func (b *StarterDocumentBuilder) BuildPtr() *StarterDocument {
	return &b.doc
}
