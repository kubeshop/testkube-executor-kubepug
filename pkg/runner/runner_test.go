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
	t.Run("runner should return success on valid yaml gist", func(t *testing.T) {
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
}

func TestRunGitFile(t *testing.T) {
	t.Skip()
}

func TestRunGitDirectory(t *testing.T) {
	t.Skip()
}

func TestRunDirectConnection(t *testing.T) {
	// Will likely not be implemented
	t.Skip()
}

func TestRunWithSpecificK8sVersion(t *testing.T) {
	t.Skip()
}
