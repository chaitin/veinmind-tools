package ref

import (
	"github.com/distribution/distribution/reference"
	"github.com/pkg/errors"
)

func ParseReference(ref string) (repo string, tag string, err error) {
	// parse reference
	parsed, err := reference.Parse(ref)
	if err != nil {
		return
	}

	if parsed.String() != ref {
		err = errors.Errorf("basic: mismatch repo, got %q, expected %q", parsed.String(), ref)
		return
	}

	if named, ok := parsed.(reference.Named); ok {
		repo = named.Name()
	} else {
		err = errors.Errorf("basic: mismatch parse type, got nil, expect named")
		return
	}

	if tagged, ok := parsed.(reference.Tagged); ok {
		tag = tagged.Tag()
	} else {
		err = errors.Errorf("basic: mismatch parse type, got nil, expect tagged")
		return
	}

	return
}
