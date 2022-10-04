package runner

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/kubeshop/testkube/pkg/api/v1/testkube"
	"github.com/kubeshop/testkube/pkg/executor"
	"github.com/kubeshop/testkube/pkg/executor/content"
	"github.com/kubeshop/testkube/pkg/executor/output"
	"github.com/kubeshop/testkube/pkg/executor/secret"
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

// Run runs the kubepug executable and parses it's output to be Testkube-compatible
func (r *KubepugRunner) Run(execution testkube.Execution) (testkube.ExecutionResult, error) {
	path, err := r.Fetcher.Fetch(execution.Content)
	if err != nil {
		return testkube.ExecutionResult{}, fmt.Errorf("could not get content: %w", err)
	}
	output.PrintEvent("created content path", path)

	if execution.Content.IsFile() {
		output.PrintEvent("using single file", execution)
	}

	if execution.Content.IsDir() {
		output.PrintEvent("using dir", execution)
	}

	args, err := buildArgs(execution.Args, path)
	if err != nil {
		return testkube.ExecutionResult{}, fmt.Errorf("could not build up parameters: %w", err)
	}

	// add configuration files
	err = content.PlaceFiles(execution.CopyFiles)
	if err != nil {
		return testkube.ExecutionResult{}, fmt.Errorf("could not place config files: %w", err)
	}

	output.PrintEvent("running kubepug with arguments", args)
	envManager := secret.NewEnvManagerWithVars(execution.Variables)
	envManager.GetEnvs()
	for _, env := range execution.Variables {
		os.Setenv(env.Name, env.Value)
	}
	out, err := executor.Run("", "kubepug", envManager, args...)
	out = envManager.Obfuscate(out)
	if err != nil {
		return testkube.ExecutionResult{}, fmt.Errorf("could not execute kubepug: %w", err)
	}

	var kubepugResult kubepug.Result
	err = json.Unmarshal(out, &kubepugResult)
	if err != nil {
		return testkube.ExecutionResult{}, fmt.Errorf("could not unmarshal kubepug execution result: %w", err)
	}

	deprecatedAPIstep := createDeprecatedAPIsStep(kubepugResult)
	deletedAPIstep := createDeletedAPIsStep(kubepugResult)
	return testkube.ExecutionResult{
		Status: getResultStatus(kubepugResult),
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

// buildArgs builds up the arguments for
func buildArgs(args []string, inputPath string) ([]string, error) {
	for _, a := range args {
		if strings.Contains(a, "--format") {
			return []string{}, fmt.Errorf("the Testkube Kubepug executor does not accept the \"--format\" parameter: %s", a)
		}
		if strings.Contains(a, "--input-file") {
			return []string{}, fmt.Errorf("the Testkube Kubepug executor does not accept the \"--input-file\" parameter: %s", a)
		}
	}
	return append(args, "--format=json", "--input-file", inputPath), nil
}
