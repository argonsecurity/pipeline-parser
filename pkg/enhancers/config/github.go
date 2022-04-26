package config

var (
	GithubConfiguration = EnhancementConfiguration{
		LabelMapping: map[string]ObjectiveConfiguration{
			"sbom": {
				Tasks: []string{
					"argonsecurity/actions/generate-manifest",
				},
			},
			"appsec": {
				Tasks: []string{
					"argonsecurity/scanner-action",
				},
			},
		},
	}
)
