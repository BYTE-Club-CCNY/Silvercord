package leetcode_graphql

type AcSubmission struct {
	Difficulty  string `json:"difficulty"`
	Count       int    `json:"count"`
	Submissions int    `json:"submissions"`
}

type SubmitStats struct {
	AcSubmissionNum []AcSubmission `json:"acSubmissionNum"`
}

type MatchedUser struct {
	Username    string      `json:"username"`
	SubmitStats SubmitStats `json:"submitStats"`
}

type GetOnlineUsernameResponse struct {
	Data struct {
		MatchedUser *MatchedUser `json:"matchedUser"`
	} `json:"data"`
}

type Question struct {
	Title      string `json:"title"`
	Difficulty string `json:"difficulty"`
}

type GetDifficultyResponse struct {
	Data struct {
		Question Question `json:"question"`
	} `json:"data"`
}

type ExtractProblemResponse struct {
	Problem string `json:"problem"`
}

