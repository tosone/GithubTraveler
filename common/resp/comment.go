package resp

// Comment ..
type Comment struct {
	URL               string
	HTMLURL           string
	IssueURL          string
	ID                string
	User              Owner
	CreatedAt         string
	UpdatedAt         string
	AuthorAssociation string
	Body              []byte
}
