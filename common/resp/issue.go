package resp

// Issue repo's issues
type Issue struct {
	URL               string  `json:"url"`
	RepositoryURL     string  `json:"repository_url"`
	LabelsURL         string  `json:"labels_url"`
	CommentsURL       string  `json:"comments_url"`
	EventsURL         string  `json:"events_url"`
	HTMLURL           string  `json:"html_url"`
	ID                uint64  `json:"id"`
	Number            uint64  `json:"number"`
	Title             string  `json:"title"`
	User              Owner   `json:"user"`
	Labels            []Label `json:"labels"`
	State             string  `json:"state"`
	Locked            bool    `json:"locked"`
	Assignee          *Owner  `json:"assignee"`
	Assignees         []Owner `json:"assignees"`
	Commits           uint64  `json:"commits"`
	CreatedAt         string  `json:"created_at"`
	UpdatedAt         string  `json:"updated_at"`
	ClosedAt          *string `json:"closed_at"`
	AuthorAssociation string  `json:"author_association"`
	Body              []byte  `json:"body"`
}

// Label issue's label
type Label struct {
	ID      uint   `json:"id"`
	URL     string `json:"url"`
	Name    string `json:"name"`
	Color   string `json:"color"`
	Default bool   `json:"default"`
}

// Milestone issue's milestone
type Milestone struct {
	URL          string  `json:"url"`
	HTMLURL      string  `json:"html_url"`
	LabelsURL    string  `json:"labels_url"`
	ID           uint    `json:"id"`
	Title        string  `json:"title"`
	Description  string  `json:"description"`
	Creator      Owner   `json:"creator"`
	OpenIssues   uint    `json:"open_issues"`
	ClosedIssues uint    `json:"closed_issues"`
	State        string  `json:"state"`
	CreatedAt    string  `json:"created_at"`
	UpdatedAt    string  `json:"updated_at"`
	DueOn        string  `json:"due_on"`
	ClosedAt     *string `json:"closed_at"`
}
