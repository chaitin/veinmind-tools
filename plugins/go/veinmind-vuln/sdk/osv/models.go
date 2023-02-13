package osv

import "time"

// Package represents a package identifier for OSV.
type Package struct {
	Index     int
	Path      string
	PURL      string `json:"purl,omitempty"`
	Name      string `json:"name,omitempty"`
	Ecosystem string `json:"ecosystem,omitempty"`
}

type SourceInfo struct {
	Path string `json:"path"`
	Type string `json:"type"`
}

func (s SourceInfo) String() string {
	return s.Type + ":" + s.Path
}

// Query represents a query to OSV.
type Query struct {
	Type    string
	Commit  string     `json:"commit,omitempty"`
	Package Package    `json:"package,omitempty"`
	Version string     `json:"version,omitempty"`
	Source  SourceInfo `json:"omit"`
}

// BatchedQuery represents a batched query to OSV.
type BatchedQuery struct {
	Queries []*Query `json:"queries"`
}

// MinimalVulnerability represents an unhydrated vulnerability entry from OSV.
type MinimalVulnerability struct {
	ID string `json:"id"`
}

// Response represents a full response from OSV.
type Response struct {
	Vulns []Vulnerability `json:"vulns"`
}

// MinimalResponse represents an unhydrated response from OSV.
type MinimalResponse struct {
	Vulns []MinimalVulnerability `json:"vulns"`
}

// BatchedResponse represents an unhydrated batched response from OSV.
type BatchedResponse struct {
	Results []MinimalResponse `json:"results"`
}

// HydratedBatchedResponse represents a hydrated batched response from OSV.
type HydratedBatchedResponse struct {
	Results []Response `json:"results"`
}

type Vulnerability struct {
	SchemaVersion string    `json:"schema_version,omitempty"`
	ID            string    `json:"id,omitempty"`
	Modified      time.Time `json:"modified,omitempty"`
	Published     time.Time `json:"published,omitempty"`
	Aliases       []string  `json:"aliases,omitempty"`
	Summary       string    `json:"summary,omitempty"`
	Details       string    `json:"details,omitempty"`
	Affected      []struct {
		Package struct {
			Ecosystem string `json:"ecosystem,omitempty"`
			Name      string `json:"name,omitempty"`
			Purl      string `json:"purl,omitempty"`
		} `json:"package,omitempty"`
		Ranges []struct {
			Type   string `json:"type,omitempty"`
			Events []struct {
				Introduced   string `json:"introduced,omitempty"`
				Fixed        string `json:"fixed,omitempty"`
				LastAffected string `json:"last_affected,omitempty"`
				Limit        string `json:"limit,omitempty"`
			} `json:"events,omitempty"`
			DatabaseSpecific map[string]interface{} `json:"database_specific,omitempty"`
		} `json:"ranges,omitempty"`
		Versions          []string               `json:"versions,omitempty"`
		DatabaseSpecific  map[string]interface{} `json:"database_specific,omitempty"`
		EcosystemSpecific map[string]interface{} `json:"ecosystem_specific,omitempty"`
	} `json:"affected,omitempty"`
	References []struct {
		Type string `json:"type,omitempty"`
		URL  string `json:"url,omitempty"`
	} `json:"references,omitempty"`
	DatabaseSpecific map[string]interface{} `json:"database_specific,omitempty"`
}

func (v *Vulnerability) GetAliases() []string {
	if len(v.Aliases) > 0 {
		return v.Aliases
	}
	return []string{v.ID}
}
