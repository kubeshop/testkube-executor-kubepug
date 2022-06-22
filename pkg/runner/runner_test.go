package runner

import (
	"testing"

	"github.com/kubeshop/testkube/pkg/api/v1/testkube"
	"github.com/stretchr/testify/assert"
)

func TestRunString(t *testing.T) {
	t.Run("runner should return success and empty result on empty string", func(t *testing.T) {
		runner := NewRunner()
		execution := testkube.NewQueuedExecution()
		execution.Content = testkube.NewStringTestContent("")

		result, err := runner.Run(*execution)

		assert.NoError(t, err)
		assert.Equal(t, testkube.ExecutionStatusPassed, result.Status)
		assert.Equal(t, 2, len(result.Steps))
		assert.Equal(t, "passed", result.Steps[0].Status)
		assert.Equal(t, "passed", result.Steps[1].Status)
	})

	t.Run("runner should return success and empty result on passing yaml", func(t *testing.T) {
		runner := NewRunner()
		execution := testkube.NewQueuedExecution()
		execution.Content = testkube.NewStringTestContent(`
apiVersion: v1
items:
- apiVersion: v1
  kind: ConfigMap
  metadata:
    annotations:
      control-plane.alpha.kubernetes.io/leader: '{"holderIdentity":"ingress-nginx-controller-646d5d4d67-7nx7r","leaseDurationSeconds":30,"acquireTime":"2022-05-31T23:08:52Z","renewTime":"2022-06-20T16:17:51Z","leaderTransitions":12}'
    creationTimestamp: "2021-10-07T13:44:37Z"
    name: ingress-controller-leader
    namespace: default
    resourceVersion: "170745168"
    uid: 9bb57467-b5c4-41fe-83a8-9513ae86fbff

`)

		result, err := runner.Run(*execution)

		assert.NoError(t, err)
		assert.Equal(t, testkube.ExecutionStatusPassed, result.Status)
		assert.Equal(t, 2, len(result.Steps))
		assert.Equal(t, "passed", result.Steps[0].Status)
		assert.Equal(t, "passed", result.Steps[1].Status)
	})
	t.Run("runner should return failure and list of deprecated APIs result on yaml containing deprecated API", func(t *testing.T) {
		runner := NewRunner()
		execution := testkube.NewQueuedExecution()
		execution.Content = testkube.NewStringTestContent(`
apiVersion: v1
items:
- apiVersion: v1
  conditions:
  - message: '{"health":"true"}'
    status: "True"
    type: Healthy
  kind: ComponentStatus
  metadata:
    creationTimestamp: null
    name: etcd-1
`)

		result, err := runner.Run(*execution)

		assert.NoError(t, err)
		assert.Equal(t, testkube.ExecutionStatusFailed, result.Status)
		assert.Equal(t, 2, len(result.Steps))
		assert.Equal(t, "failed", result.Steps[0].Status)
		assert.Equal(t, "passed", result.Steps[1].Status)
	})
}

func TestRunFileURI(t *testing.T) {
	t.Run("runner should return success on valid yaml gist file URI", func(t *testing.T) {
		runner := NewRunner()
		execution := testkube.NewQueuedExecution()
		execution.Content = &testkube.TestContent{
			Type_: string(testkube.TestContentTypeFileURI),
			Uri:   "https://gist.githubusercontent.com/vLia/b3df9e43f55fd43d1bca93cdfd5ae27c/raw/535e8db46f33693a793c616fc1e2b4d77c4b06d2/example-k8s-pod-yaml",
		}

		result, err := runner.Run(*execution)

		assert.NoError(t, err)
		assert.Equal(t, testkube.ExecutionStatusPassed, result.Status)
		assert.Equal(t, 2, len(result.Steps))
		assert.Equal(t, "passed", result.Steps[0].Status)
		assert.Equal(t, "passed", result.Steps[1].Status)
	})
	t.Run("runner should return failure on yaml gist file URI with deprecated/deleted APIs", func(t *testing.T) {
		runner := NewRunner()
		execution := testkube.NewQueuedExecution()
		execution.Content = &testkube.TestContent{
			Type_: string(testkube.TestContentTypeFileURI),
			Uri:   "https://gist.githubusercontent.com/vLia/91289de9cc8b6953be5f90b0a52fa8d3/raw/a8ed0b07361b84873c6b71fb8be6e334224062d4/example-k8s-pod-yaml-deprecated",
		}

		result, err := runner.Run(*execution)

		assert.NoError(t, err)
		assert.Equal(t, testkube.ExecutionStatusFailed, result.Status)
		assert.Equal(t, 2, len(result.Steps))
		assert.Equal(t, "failed", result.Steps[0].Status)
		assert.Equal(t, "passed", result.Steps[1].Status)
	})
}

func TestRunGitFile(t *testing.T) {
	t.Run("runner should return error on non-existent Git path", func(t *testing.T) {
		runner := NewRunner()
		execution := testkube.NewQueuedExecution()
		execution.Content = &testkube.TestContent{
			Type_: string(testkube.TestContentTypeGitFile),
			Repository: &testkube.Repository{
				Uri:    "https://github.com/kubeshop/testkube-dashboard/",
				Branch: "main",
				Path:   "manifests/fake-deployment.yaml",
			},
		}

		_, err := runner.Run(*execution)

		assert.Error(t, err)
	})
	t.Run("runner should return deprecated and deleted APIs on Git file containing deprecated and delete API definitions", func(t *testing.T) {
		runner := NewRunner()
		execution := testkube.NewQueuedExecution()
		execution.Content = &testkube.TestContent{
			Type_: string(testkube.TestContentTypeGitFile),
			Repository: &testkube.Repository{
				Uri:    "https://github.com/kubeshop/testkube-dashboard/",
				Branch: "main",
				Path:   "manifests/deployment.yaml",
			},
		}

		result, err := runner.Run(*execution)

		assert.NoError(t, err)
		assert.Equal(t, testkube.ExecutionStatusPassed, result.Status)
		assert.Equal(t, 2, len(result.Steps))
		assert.Equal(t, "passed", result.Steps[0].Status)
		assert.Equal(t, "passed", result.Steps[1].Status)
	})
}

func TestRunGitDirectory(t *testing.T) {
	t.Run("runner should return success on manifests from Git directory", func(t *testing.T) {
		runner := NewRunner()
		execution := testkube.NewQueuedExecution()
		execution.Content = &testkube.TestContent{
			Type_: string(testkube.TestContentTypeGitDir),
			Repository: &testkube.Repository{
				Uri:    "https://github.com/kubeshop/testkube-dashboard/",
				Branch: "main",
				Path:   "manifests",
			},
		}

		result, err := runner.Run(*execution)

		assert.NoError(t, err)
		assert.Equal(t, testkube.ExecutionStatusPassed, result.Status)
		assert.Equal(t, 2, len(result.Steps))
		assert.Equal(t, "passed", result.Steps[0].Status)
		assert.Equal(t, "passed", result.Steps[1].Status)
	})
}

func TestRunDirectConnection(t *testing.T) {
	// Will likely not be implemented
	t.Skip()
}

func TestRunWithSpecificK8sVersion(t *testing.T) {
	// To be implemented
	t.Run("runner should return failure and list of deprecated APIs result on yaml containing deprecated API with current K8s version", func(t *testing.T) {
		runner := NewRunner()
		execution := testkube.NewQueuedExecution()
		execution.Content = testkube.NewStringTestContent(`
apiVersion: v1
items:
- apiVersion: v1
  conditions:
  - message: '{"health":"true"}'
    status: "True"
    type: Healthy
  kind: ComponentStatus
  metadata:
    creationTimestamp: null
    name: etcd-1
`)

		result, err := runner.Run(*execution)

		assert.NoError(t, err)
		assert.Equal(t, testkube.ExecutionStatusFailed, result.Status)
		assert.Equal(t, 2, len(result.Steps))
		assert.Equal(t, "failed", result.Steps[0].Status)
		assert.Equal(t, "passed", result.Steps[1].Status)
	})
	t.Run("runner should return success on yaml containing deprecated API with old K8s version", func(t *testing.T) {
		runner := NewRunner()
		execution := testkube.NewQueuedExecution()
		execution.Args = []string{
			"--k8s-version=v1.18.0", // last version v1/ComponentStatus was valid
		}
		execution.Content = testkube.NewStringTestContent(`
apiVersion: v1
items:
- apiVersion: v1
  conditions:
  - message: '{"health":"true"}'
    status: "True"
    type: Healthy
  kind: ComponentStatus
  metadata:
    creationTimestamp: null
    name: etcd-1
`)

		result, err := runner.Run(*execution)

		assert.NoError(t, err)
		assert.Equal(t, testkube.ExecutionStatusPassed, result.Status)
		assert.Equal(t, 2, len(result.Steps))
		assert.Equal(t, "passed", result.Steps[0].Status)
		assert.Equal(t, "passed", result.Steps[1].Status)
	})
}
