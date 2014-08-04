package seapi

type Owner struct {
	User_id       int
	Display_name  string
	Reputation    int
	User_type     string //one of unregistered, registered, moderator, or does_not_exist
	Profile_image string
	Link          string
}

type SearchResult struct {
	SearchResult_id      int
	Last_edit_date       int64
	Creation_date        int64
	Last_activity_date   int64
	Locked_date          int64
	Community_owned_date int64
	Score                int
	Answer_count         int
	Accepted_answer_id   int
	Bounty_closes_date   int64
	Bounty_amount        int
	Closed_date          int64
	Protected_date       int64
	Title                string
	Tags                 []string
	Closed_reason        string
	Up_vote_count        int
	Down_vote_count      int
	Favorite_count       int
	View_count           int
	Owner                Owner
	Link                 string
	Is_answered          bool
}

type SearchResultCollection struct {
	Items         []SearchResult
	Error_id      int
	Error_name    string
	Error_message string

	Has_more        bool
	Quota_max       int
	Quota_remaining int
}
