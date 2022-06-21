package runner

import (
	"encoding/json"
	"fmt"

	"github.com/kubeshop/testkube/pkg/api/v1/testkube"
	"github.com/kubeshop/testkube/pkg/executor"
	"github.com/kubeshop/testkube/pkg/executor/content"
	"github.com/kubeshop/testkube/pkg/executor/output"
	kubepug "github.com/rikatz/kubepug/pkg/results"
)

func NewRunner() *KubepugRunner {
	return &KubepugRunner{
		Fetcher: content.NewFetcher(""),
	}
}

// KubepugRunner runs kubepug against cluster
type KubepugRunner struct {
	Fetcher content.ContentFetcher
}

func (r *KubepugRunner) Run(execution testkube.Execution) (result testkube.ExecutionResult, err error) {

	// use `execution.Variables` for variables passed from Test/Execution
	// variables of type "secret" will be automatically decoded

	path, err := r.Fetcher.Fetch(execution.Content)
	if err != nil {
		return result, err
	}

	output.PrintEvent("created content path", path)

	if execution.Content.IsFile() {
		output.PrintEvent("using file", execution)
		// TODO implement file based test content for string, git-file, file-uri
		//      or remove if not used
	}

	if execution.Content.IsDir() {
		output.PrintEvent("using dir", execution)
		// TODO implement file based test content for git-dir
		//      or remove if not used
	}

	out, err := executor.Run("", "kubectl", "deprecations", "--format=json", "--input-file", path) // TODO update to kubepug
	if err != nil {
		return testkube.ExecutionResult{}, fmt.Errorf("could not execute kubepug: %w", err)
	}

	var kResult kubepug.Result
	err = json.Unmarshal(out, &kResult)
	if err != nil {
		return testkube.ExecutionResult{}, fmt.Errorf("could not unmarshal kubepug execution result: %s", err)
	}

	deprecatedAPIstep := createDeprecatedAPIsStep(kResult)
	deletedAPIstep := createDeletedAPIsStep(kResult)

	// error result should be returned if something is not ok
	// return result.Err(fmt.Errorf("some test execution related error occured"))

	// TODO return ExecutionResult

	return testkube.ExecutionResult{
		Status: getResultStatus(kResult),
		Output: string(out),
		Steps: []testkube.ExecutionStepResult{
			deprecatedAPIstep,
			deletedAPIstep,
		},
	}, nil
}

// createDeprecatedAPIsStep checks the kubepug output for deprecated APIs and converts them to Testkube step result
func createDeprecatedAPIsStep(r kubepug.Result) testkube.ExecutionStepResult {
	step := testkube.ExecutionStepResult{
		Name: "Deprecated APIs",
	}

	if len(r.DeprecatedAPIs) == 0 {
		step.Status = "passed"
		return step
	}

	step.Status = "failed"
	for _, api := range r.DeletedAPIs {
		step.AssertionResults = append(step.AssertionResults, testkube.AssertionResult{
			Name:         api.Name,
			Status:       "failed",
			ErrorMessage: fmt.Sprintf("Deprecated API:\n %v", api),
		})
	}

	return step
}

// createDeletedAPISstep checks the kubepug output for deleted APIs and converts them to Testkube step result
func createDeletedAPIsStep(r kubepug.Result) testkube.ExecutionStepResult {
	step := testkube.ExecutionStepResult{
		Name: "Deleted APIs",
	}

	if len(r.DeletedAPIs) == 0 {
		step.Status = "passed"
		return step
	}

	step.Status = "failed"
	for _, api := range r.DeletedAPIs {
		step.AssertionResults = append(step.AssertionResults, testkube.AssertionResult{
			Name:         api.Name,
			Status:       "failed",
			ErrorMessage: fmt.Sprintf("Deleted API:\n %v", api),
		})
	}

	return step
}

// getResultStatus calculates the final result status
func getResultStatus(r kubepug.Result) *testkube.ExecutionStatus {
	if len(r.DeletedAPIs) == 0 && len(r.DeprecatedAPIs) == 0 {
		return testkube.ExecutionStatusPassed
	}
	return testkube.ExecutionStatusFailed
}
