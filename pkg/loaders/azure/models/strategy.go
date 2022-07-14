package models

type Matrix map[string]map[string]string

type JobStrategy struct {
	Matrix      *Matrix `yaml:"matrix,omitempty"`
	MaxParallel int     `yaml:"maxParallel,omitempty"`
	Parallel    string  `yaml:"parallel,omitempty"`
}

type DeploymentHook struct {
	Steps *Steps `yaml:"steps,omitempty"`
	Pool  *Pool  `yaml:"pool,omitempty"`
}

type DeploymentStrategy struct {
	PreDeploy        *DeploymentHook `yaml:"preDeploy,omitempty"`
	Deploy           *DeploymentHook `yaml:"deploy,omitempty"`
	RouteTraffic     *DeploymentHook `yaml:"routeTraffic,omitempty"`
	PostRouteTraffic *DeploymentHook `yaml:"postRouteTraffic,omitempty"`
	On               struct {
		Failure *DeploymentHook `yaml:"failure,omitempty"`
		Success *DeploymentHook `yaml:"success,omitempty"`
	} `yaml:"on,omitempty"`
}

type RollingStrategy struct {
	MaxParallel        string `yaml:"maxParallel,omitempty"`
	DeploymentStrategy `yaml:",inline"`
}

type CanaryDeploymentStrategy struct {
	Increments         []string `yaml:"increments,omitempty"`
	DeploymentStrategy `yaml:",inline"`
}
