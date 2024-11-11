package gitdorker

// GithubResponse结构体
type GithubResponse struct {
	Total_count        float64 `json:"total_count"`
	Incomplete_results bool    `json:"incomplete_results"`
	Items              []struct {
		Html_url string `json:"html_url"`
	} `json:"items"`
}

// GithubResult结构体
type GithubResult struct {
	Status bool
	Total  float64
	Items  []string
	Link   string
}
