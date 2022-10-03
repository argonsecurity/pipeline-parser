package models

type Matrix map[string]any

type JobStrategy struct {
	Matrix      *Matrix `yaml:"matrix,omitempty"`
	MaxParallel string  `yaml:"maxParallel,omitempty"`
	Parallel    string  `yaml:"parallel,omitempty"`
}

type DeploymentHook struct {
	Steps *Steps `yaml:"steps,omitempty"`
	Pool  *Pool  `yaml:"pool,omitempty"`
}

type DeploymentStrategy struct {
	RunOnce *BaseDeploymentStrategy `yaml:"runOnce,omitempty"`
	Rolling *RollingStrategy        `yaml:"rolling,omitempty"`
	Canary  *CanaryStrategy         `yaml:"canary,omitempty"`
}

type BaseDeploymentStrategy struct {
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
	MaxParallel            string `yaml:"maxParallel,omitempty"`
	BaseDeploymentStrategy `yaml:",inline"`
}

type CanaryStrategy struct {
	Increments             []string `yaml:"increments,omitempty"`
	BaseDeploymentStrategy `yaml:",inline"`
}
