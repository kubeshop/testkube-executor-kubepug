package runner

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResultParser(t *testing.T) {
	t.Run("GetResults should return empty result when there is no JSON output", func(t *testing.T) {
		t.Parallel()
		output := `{"DeprecatedAPIs":null,"DeletedAPIs":null}`
		result, err := GetResult(output)
		assert.NoError(t, err)
		assert.Equal(t, Result{
			DeprecatedAPIs: nil,
			DeletedAPIs:    nil,
		}, result)
	})
	t.Run("GetResult should return error for invalid JSON", func(t *testing.T) {
		t.Parallel()
		output := `invalid JSON`
		_, err := GetResult(output)
		assert.Error(t, err)
	})
	t.Run("GetResult should return populated DeprecatedAPIs when there's a DeprecatedAPI finding", func(t *testing.T) {
		t.Parallel()
		output := `{
			"DeprecatedAPIs": [
			  {
				"Description": "ComponentStatus (and ComponentStatusList) holds the cluster validation info. Deprecated: This API is deprecated in v1.19+",
				"Group": "",
				"Kind": "ComponentStatus",
				"Version": "v1",
				"Name": "",
				"Deprecated": true,
				"Items": [
				  {
					"Scope": "GLOBAL",
					"ObjectName": "scheduler",
					"Namespace": ""
				  },
				  {
					"Scope": "GLOBAL",
					"ObjectName": "etcd-0",
					"Namespace": ""
				  },
				  {
					"Scope": "GLOBAL",
					"ObjectName": "etcd-1",
					"Namespace": ""
				  },
				  {
					"Scope": "GLOBAL",
					"ObjectName": "controller-manager",
					"Namespace": ""
				  }
				]
			  }
			],
			"DeletedAPIs": null
			}`
		expected := Result{
			DeprecatedAPIs: []DeprecatedAPIs{
				{
					Description: "ComponentStatus (and ComponentStatusList) holds the cluster validation info. Deprecated: This API is deprecated in v1.19+",
					Kind:        "ComponentStatus",
					Version:     "v1",
					Deprecated:  true,
					Items: []Items{
						{
							Scope:      "GLOBAL",
							ObjectName: "scheduler",
						},
						{
							Scope:      "GLOBAL",
							ObjectName: "etcd-0",
						},
						{
							Scope:      "GLOBAL",
							ObjectName: "etcd-1",
						},
						{
							Scope:      "GLOBAL",
							ObjectName: "controller-manager",
						},
					},
				},
			},
		}
		result, err := GetResult(output)
		assert.NoError(t, err)
		assert.Equal(t, expected, result)
	})
	t.Run("GetResult should return populated DeletedAPIs when there's a DeletedAPIs finding", func(t *testing.T) {
		t.Parallel()
		output := `{
			"DeprecatedAPIs": null,
			"DeletedAPIs": [
			  {
				"Group": "extensions",
				"Kind": "Ingress",
				"Version": "v1beta1",
				"Name": "ingresses",
				"Deleted": true,
				"Items": [
				  {
					"Scope": "OBJECT",
					"ObjectName": "cli-testkube-api-server-testkube",
					"Namespace": "testkube"
				  },
				  {
					"Scope": "OBJECT",
					"ObjectName": "oauth2-proxy",
					"Namespace": "testkube"
				  },
				  {
					"Scope": "OBJECT",
					"ObjectName": "testapi",
					"Namespace": "testkube"
				  },
				  {
					"Scope": "OBJECT",
					"ObjectName": "testdash",
					"Namespace": "testkube"
				  },
				  {
					"Scope": "OBJECT",
					"ObjectName": "testkube-dashboard-testkube",
					"Namespace": "testkube"
				  },
				  {
					"Scope": "OBJECT",
					"ObjectName": "ui-testkube-api-server-testkube",
					"Namespace": "testkube"
				  }
				]
			  },
			  {
				"Group": "policy",
				"Kind": "PodSecurityPolicy",
				"Version": "v1beta1",
				"Name": "podsecuritypolicies",
				"Deleted": true,
				"Items": [
				  {
					"Scope": "GLOBAL",
					"ObjectName": "gce.gke-metrics-agent",
					"Namespace": ""
				  }
				]
			  }
			]
		  }
		`
		expected := Result{
			DeletedAPIs: []DeletedAPIs{
				{
					Group:   "extensions",
					Kind:    "Ingress",
					Version: "v1beta1",
					Name:    "ingresses",
					Deleted: true,
					Items: []Items{
						{
							Scope:      "OBJECT",
							ObjectName: "cli-testkube-api-server-testkube",
							Namespace:  "testkube",
						},
						{
							Scope:      "OBJECT",
							ObjectName: "oauth2-proxy",
							Namespace:  "testkube",
						},
						{
							Scope:      "OBJECT",
							ObjectName: "testapi",
							Namespace:  "testkube",
						},
						{
							Scope:      "OBJECT",
							ObjectName: "testdash",
							Namespace:  "testkube",
						},
						{
							Scope:      "OBJECT",
							ObjectName: "testkube-dashboard-testkube",
							Namespace:  "testkube",
						},
						{
							Scope:      "OBJECT",
							ObjectName: "ui-testkube-api-server-testkube",
							Namespace:  "testkube",
						},
					},
				},
				{
					Group:   "policy",
					Kind:    "PodSecurityPolicy",
					Version: "v1beta1",
					Name:    "podsecuritypolicies",
					Deleted: true,
					Items: []Items{
						{
							Scope:      "GLOBAL",
							ObjectName: "gce.gke-metrics-agent",
							Namespace:  "",
						},
					},
				},
			},
		}
		result, err := GetResult(output)
		assert.NoError(t, err)
		assert.Equal(t, expected, result)
	})
	t.Run("GetResult should return populated DeprecatedAPIs and DeletedAPIs when there's both finding", func(t *testing.T) {
		t.Parallel()
		output := `{
			"DeprecatedAPIs": [
			  {
				"Description": "ComponentStatus (and ComponentStatusList) holds the cluster validation info. Deprecated: This API is deprecated in v1.19+",
				"Group": "",
				"Kind": "ComponentStatus",
				"Version": "v1",
				"Name": "",
				"Deprecated": true,
				"Items": [
				  {
					"Scope": "GLOBAL",
					"ObjectName": "scheduler",
					"Namespace": ""
				  },
				  {
					"Scope": "GLOBAL",
					"ObjectName": "etcd-0",
					"Namespace": ""
				  },
				  {
					"Scope": "GLOBAL",
					"ObjectName": "etcd-1",
					"Namespace": ""
				  },
				  {
					"Scope": "GLOBAL",
					"ObjectName": "controller-manager",
					"Namespace": ""
				  }
				]
			  }
			],
			"DeletedAPIs": [
			  {
				"Group": "extensions",
				"Kind": "Ingress",
				"Version": "v1beta1",
				"Name": "ingresses",
				"Deleted": true,
				"Items": [
				  {
					"Scope": "OBJECT",
					"ObjectName": "cli-testkube-api-server-testkube",
					"Namespace": "testkube"
				  },
				  {
					"Scope": "OBJECT",
					"ObjectName": "oauth2-proxy",
					"Namespace": "testkube"
				  },
				  {
					"Scope": "OBJECT",
					"ObjectName": "testapi",
					"Namespace": "testkube"
				  },
				  {
					"Scope": "OBJECT",
					"ObjectName": "testdash",
					"Namespace": "testkube"
				  },
				  {
					"Scope": "OBJECT",
					"ObjectName": "testkube-dashboard-testkube",
					"Namespace": "testkube"
				  },
				  {
					"Scope": "OBJECT",
					"ObjectName": "ui-testkube-api-server-testkube",
					"Namespace": "testkube"
				  }
				]
			  },
			  {
				"Group": "policy",
				"Kind": "PodSecurityPolicy",
				"Version": "v1beta1",
				"Name": "podsecuritypolicies",
				"Deleted": true,
				"Items": [
				  {
					"Scope": "GLOBAL",
					"ObjectName": "gce.gke-metrics-agent",
					"Namespace": ""
				  }
				]
			  }
			]
		  }
		`
		expected := Result{
			DeprecatedAPIs: []DeprecatedAPIs{
				{
					Description: "ComponentStatus (and ComponentStatusList) holds the cluster validation info. Deprecated: This API is deprecated in v1.19+",
					Kind:        "ComponentStatus",
					Version:     "v1",
					Deprecated:  true,
					Items: []Items{
						{
							Scope:      "GLOBAL",
							ObjectName: "scheduler",
						},
						{
							Scope:      "GLOBAL",
							ObjectName: "etcd-0",
						},
						{
							Scope:      "GLOBAL",
							ObjectName: "etcd-1",
						},
						{
							Scope:      "GLOBAL",
							ObjectName: "controller-manager",
						},
					},
				},
			},
			DeletedAPIs: []DeletedAPIs{
				{
					Group:   "extensions",
					Kind:    "Ingress",
					Version: "v1beta1",
					Name:    "ingresses",
					Deleted: true,
					Items: []Items{
						{
							Scope:      "OBJECT",
							ObjectName: "cli-testkube-api-server-testkube",
							Namespace:  "testkube",
						},
						{
							Scope:      "OBJECT",
							ObjectName: "oauth2-proxy",
							Namespace:  "testkube",
						},
						{
							Scope:      "OBJECT",
							ObjectName: "testapi",
							Namespace:  "testkube",
						},
						{
							Scope:      "OBJECT",
							ObjectName: "testdash",
							Namespace:  "testkube",
						},
						{
							Scope:      "OBJECT",
							ObjectName: "testkube-dashboard-testkube",
							Namespace:  "testkube",
						},
						{
							Scope:      "OBJECT",
							ObjectName: "ui-testkube-api-server-testkube",
							Namespace:  "testkube",
						},
					},
				},
				{
					Group:   "policy",
					Kind:    "PodSecurityPolicy",
					Version: "v1beta1",
					Name:    "podsecuritypolicies",
					Deleted: true,
					Items: []Items{
						{
							Scope:      "GLOBAL",
							ObjectName: "gce.gke-metrics-agent",
							Namespace:  "",
						},
					},
				},
			},
		}
		result, err := GetResult(output)
		assert.NoError(t, err)
		assert.Equal(t, expected, result)
	})
}
