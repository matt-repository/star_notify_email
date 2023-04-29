package main

type GraphqlRequest struct {
	Query string `json:"query"`
}

type GraphqlResponse struct {
	Data struct {
		Repository struct {
			StargazerCount int `json:"stargazerCount"`
			Stargazers     struct {
				Edges []struct {
					Node struct {
						Name      string `json:"name"`
						Email     string `json:"email"`
						URL       string `json:"url"`
						AvatarURL string `json:"avatarUrl"`
					} `json:"node"`
				} `json:"edges"`
			} `json:"stargazers"`
		} `json:"repository"`
	} `json:"data"`
}
