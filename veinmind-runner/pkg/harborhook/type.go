package harborhook

type PushImage struct {
	Type      string `json:"type"`
	OccurAt   int    `json:"occur_at"`
	Operator  string `json:"operator"`
	EventData struct {
		Resources []struct {
			Digest      string `json:"digest"`
			Tag         string `json:"tag"`
			ResourceURL string `json:"resource_url"`
		} `json:"resources"`
		Repository struct {
			DateCreated  int    `json:"date_created"`
			Name         string `json:"name"`
			Namespace    string `json:"namespace"`
			RepoFullName string `json:"repo_full_name"`
			RepoType     string `json:"repo_type"`
		} `json:"repository"`
	} `json:"event_data"`
}
type PullImage struct {
	Type      string `json:"type"`
	OccurAt   int    `json:"occur_at"`
	Operator  string `json:"operator"`
	EventData struct {
		Resources []struct {
			Digest      string `json:"digest"`
			Tag         string `json:"tag"`
			ResourceURL string `json:"resource_url"`
		} `json:"resources"`
		Repository struct {
			DateCreated  int    `json:"date_created"`
			Name         string `json:"name"`
			Namespace    string `json:"namespace"`
			RepoFullName string `json:"repo_full_name"`
			RepoType     string `json:"repo_type"`
		} `json:"repository"`
	} `json:"event_data"`
}
