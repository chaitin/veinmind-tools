package request

import (
	"encoding/base64"
	"fmt"
	"github.com/google/go-containerregistry/pkg/authn"
	"net/http"
)

type basicTransport struct {
	inner  http.RoundTripper
	auth   authn.Authenticator
	target string
}

var _ http.RoundTripper = (*basicTransport)(nil)

// RoundTrip implements http.RoundTripper
func (bt *basicTransport) RoundTrip(in *http.Request) (*http.Response, error) {
	if bt.auth != authn.Anonymous {
		auth, err := bt.auth.Authorization()
		if err != nil {
			return nil, err
		}

		// http.Client handles redirects at a layer above the http.RoundTripper
		// abstraction, so to avoid forwarding Authorization headers to places
		// we are redirected, only set it when the authorization header matches
		// the host with which we are interacting.
		// In case of redirect http.Client can use an empty Host, check URL too.
		if in.Host == bt.target || in.URL.Host == bt.target {
			if bearer := auth.RegistryToken; bearer != "" {
				hdr := fmt.Sprintf("Bearer %s", bearer)
				in.Header.Set("Authorization", hdr)
			} else if user, pass := auth.Username, auth.Password; user != "" && pass != "" {
				delimited := fmt.Sprintf("%s:%s", user, pass)
				encoded := base64.StdEncoding.EncodeToString([]byte(delimited))
				hdr := fmt.Sprintf("Basic %s", encoded)
				in.Header.Set("Authorization", hdr)
			} else if token := auth.Auth; token != "" {
				hdr := fmt.Sprintf("Basic %s", token)
				in.Header.Set("Authorization", hdr)
			}
		}
	}
	return bt.inner.RoundTrip(in)
}
