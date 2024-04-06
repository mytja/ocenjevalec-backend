package httphandlers

// Na začetku je status vedno "ČAKA NA EVALUACIJO"
type WSNewSubmission struct {
	MessageType int    `json:"message_type"`
	Submission  string `json:"submission"`
	TeamName    string `json:"team_name"`
	TeamID      string `json:"team_id"`
	ProblemID   string `json:"problem_id"`
	ProblemName string `json:"problem_name"`
	MaxScore    int    `json:"max_score"`
}

type WSUpdateSubmissionStatus struct {
	MessageType       int    `json:"message_type"`
	Submission        string `json:"submission"`
	TeamName          string `json:"team_name"`
	TeamID            string `json:"team_id"`
	ProblemID         string `json:"problem_id"`
	ProblemName       string `json:"problem_name"`
	Verdict           string `json:"verdict"`
	Score             int    `json:"score"`
	MaxScore          int    `json:"max_score"`
	TotalScore        int    `json:"total_score"`
	SubmissionsBefore int    `json:"submissions_before"`
}

type WSChangeSubmissionID struct {
	MessageType   int    `json:"message_type"`
	OldSubmission string `json:"old_submission"`
	NewSubmission string `json:"new_submission"`
}
