package rules

import "embed"

//go:embed */*.rego
var RegoFile embed.FS
