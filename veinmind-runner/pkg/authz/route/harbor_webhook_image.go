package route

import (
	"encoding/json"
	"errors"
	"io/ioutil"

	"github.com/gin-gonic/gin"
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

func ParseHarborwebhookPostdata(c *gin.Context) (PullandPushData, error) {
	postData := &PullandPushData{}
	data, _ := ioutil.ReadAll(c.Request.Body)
	if err := json.Unmarshal(data, &postData); err != nil {
		return PullandPushData{}, err
	}
	return *postData, nil
}

func GetImageNames(data PullandPushData) ([]string, error) {
	resources := data.EventData.Resources
	if len(resources) < 1 {
		return []string{}, errors.New("no image choosed")
	}
	var imagenames []string
	for _, resource := range resources {
		imagenames = append(imagenames, resource.ResourceURL)
	}
	return imagenames, nil
}
