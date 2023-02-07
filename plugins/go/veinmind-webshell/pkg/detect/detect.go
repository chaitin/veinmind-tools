package detect

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"io"
	"io/fs"
	"io/ioutil"
	"mime/multipart"
	"net/http"

	"github.com/chaitin/veinmind-tools/plugins/go/veinmind-webshell/pkg/filter"
)

type FileInfo struct {
	Path        string
	Reader      io.Reader
	RawFileInfo fs.FileInfo
	ScriptType  filter.ScriptType
}

type Result struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    struct {
		RiskLevel int    `json:"risk_level"`
		ID        string `json:"id"`
		Type      string `json:"type"`
		Reason    string `json:"reason"`
		Engine    string `json:"engine"`
	} `json:"data"`
}

type Kit struct {
	ctx    context.Context
	token  string
	client *http.Client
}

type KitOption func(kit *Kit)

func WithToken(token string) KitOption {
	return func(kit *Kit) {
		kit.token = token
	}
}

func WithDefaultToken() KitOption {
	return func(kit *Kit) {
		kit.token = token
	}
}

func WithClient(client *http.Client) KitOption {
	return func(kit *Kit) {
		kit.client = client
	}
}

func WithDefaultClient() KitOption {
	return func(kit *Kit) {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		kit.client = &http.Client{Transport: tr}
	}
}

func NewKit(ctx context.Context, opts ...KitOption) (*Kit, error) {
	k := new(Kit)
	k.ctx = ctx

	for _, opt := range opts {
		opt(k)
	}
	return k, nil
}

func (k *Kit) Detect(info FileInfo) (*Result, error) {
	buf := new(bytes.Buffer)
	writer := multipart.NewWriter(buf)

	part, err := writer.CreateFormFile("file", info.Path)
	if err != nil {
		return nil, err
	}

	_, err = io.Copy(part, info.Reader)
	if err != nil {
		return nil, err
	}

	_ = writer.WriteField("tag", "veinmind-webshell")
	_ = writer.WriteField("type", info.ScriptType.String())

	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, buf)
	if err != nil {
		return nil, err
	}

	req.Header.Add("X-Ca-Token", k.token)
	req.Header.Add("Content-Type", writer.FormDataContentType())

	resp, err := k.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	res := &Result{}
	err = json.Unmarshal(body, res)
	if err != nil {
		return nil, err
	}

	return res, nil
}
