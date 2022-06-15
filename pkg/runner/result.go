package runner

import (
	"encoding/json"
	"fmt"
)

// Result is the output of the execution of kubepug
type Result struct {
	DeprecatedAPIs []DeprecatedAPIs `json:"DeprecatedAPIs"`
	DeletedAPIs    []DeletedAPIs    `json:"DeletedAPIs"`
}

// Items describe vulnerable items from a Kubernetes cluster
type Items struct {
	Scope      string `json:"Scope"`
	ObjectName string `json:"ObjectName"`
	Namespace  string `json:"Namespace"`
}

// DeprecatedAPIs describe APIs that were deprecated from Kubernetes but still available in the cluster
type DeprecatedAPIs struct {
	Description string  `json:"Description"`
	Group       string  `json:"Group"`
	Kind        string  `json:"Kind"`
	Version     string  `json:"Version"`
	Name        string  `json:"Name"`
	Deprecated  bool    `json:"Deprecated"`
	Items       []Items `json:"Items"`
}

// DeletedAPIs describe APIs that were deleted from Kubernetes but still available in the cluster
type DeletedAPIs struct {
	Group   string  `json:"Group"`
	Kind    string  `json:"Kind"`
	Version string  `json:"Version"`
	Name    string  `json:"Name"`
	Deleted bool    `json:"Deleted"`
	Items   []Items `json:"Items"`
}

// GetResults parses the output of a kubepug execution into a Result
func GetResult(r string) (Result, error) {
	var result Result
	err := json.Unmarshal([]byte(r), &result)
	if err != nil {
		return result, fmt.Errorf("could not unmarshal result %s: %w", r, err)
	}
	return result, nil
}
