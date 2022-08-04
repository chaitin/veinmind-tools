package action

var (
	PULL_ARTIFACT   string = "PULL_ARTIFACT"
	PUSH_ARTIFACT   string = "PUSH_ARTIFACT"
	DELETE_ARTIFACT string = "DELETE_ARTIFACT"
)

// only for push and pull api
// TODO: parse scan/helm/etc. api post data
type PullandPushData struct {
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
