package scan

import (
	"context"

	"github.com/chaitin/veinmind-tools/veinmind-runner/pkg/target"
)

// DispatchTask declare func that how to scan a target object
type DispatchTask func(ctx context.Context, targets []*target.Target) error
