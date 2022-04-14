package models

type Jobs struct {
	NormalJobs               map[string]*NormalJob
	ReusableWorkflowCallJobs map[string]*ReusableWorkflowCallJob
}

type Concurrency struct {
	CancelInProgress *bool   `mapstructure:"cancel-in-progress,omitempty" yaml:"cancel-in-progress,omitempty"`
	Group            *string `mapstructure:"group" yaml:"group"`
}

func (c *Concurrency) UnmarshalText(text []byte) error {
	s := string(text)
	c.Group = &s
	return nil
}

// NormalJob Each job must have an id to associate with the job. The key job_id is a string and its value is a map of the job's configuration data. You must replace <job_id> with a string that is unique to the jobs object. The <job_id> must start with a letter or _ and contain only alphanumeric characters, -, or _.
type NormalJob struct {

	// Concurrency ensures that only a single job or workflow using the same concurrency group will run at a time. A concurrency group can be any string or expression. The expression can use any context except for the secrets context.
	// You can also specify concurrency at the workflow level.
	// When a concurrent job or workflow is queued, if another job or workflow using the same concurrency group in the repository is in progress, the queued job or workflow will be pending. Any previously pending job or workflow in the concurrency group will be canceled. To also cancel any currently running job or workflow in the same concurrency group, specify cancel-in-progress: true.
	Concurrency *Concurrency `mapstructure:"concurrency,omitempty" yaml:"concurrency,omitempty" json:"concurrency,omitempty"`

	// A container to run any steps in a job that don't already specify a container. If you have steps that use both script and container actions, the container actions will run as sibling containers on the same network with the same volume mounts.
	// If you do not set a container, all steps will run directly on the host specified by runs-on unless a step refers to an action configured to run in a container.
	Container interface{} `mapstructure:"container,omitempty" yaml:"container,omitempty" json:"container,omitempty"`

	// Prevents a workflow run from failing when a job fails. Set to true to allow a workflow run to pass when this job fails.
	ContinueOnError bool `mapstructure:"continue-on-error,omitempty" yaml:"continue-on-error,omitempty" json:"continue-on-error,omitempty"`

	// A map of default settings that will apply to all steps in the job.
	Defaults *Defaults `mapstructure:"defaults,omitempty" yaml:"defaults,omitempty" json:"defaults,omitempty"`

	// A map of environment variables that are available to all steps in the job.
	Env interface{} `mapstructure:"env,omitempty" yaml:"env,omitempty" json:"env,omitempty"`

	// The environment that the job references.
	Environment interface{} `mapstructure:"environment,omitempty" yaml:"environment,omitempty" json:"environment,omitempty"`

	// You can use the if conditional to prevent a job from running unless a condition is met. You can use any supported context and expression to create a conditional.
	// Expressions in an if conditional do not require the ${{ }} syntax. For more information, see https://help.github.com/en/articles/contexts-and-expression-syntax-for-github-actions.
	If string `mapstructure:"if,omitempty" yaml:"if,omitempty" json:"if,omitempty"`

	// The name of the job displayed on GitHub.
	Name  string      `mapstructure:"name,omitempty" yaml:"name,omitempty" json:"name,omitempty"`
	Needs interface{} `mapstructure:"needs,omitempty" yaml:"needs,omitempty" json:"needs,omitempty"`

	// A map of outputs for a job. Job outputs are available to all downstream jobs that depend on this job.
	Outputs     map[string]string `mapstructure:"outputs,omitempty" yaml:"outputs,omitempty" json:"outputs,omitempty"`
	Permissions *PermissionsEvent `mapstructure:"permissions,omitempty" yaml:"permissions,omitempty" json:"permissions,omitempty"`

	// The type of machine to run the job on. The machine can be either a GitHub-hosted runner, or a self-hosted runner.
	RunsOn interface{} `mapstructure:"runs-on" yaml:"runs-on" json:"runs-on"`

	// Additional containers to host services for a job in a workflow. These are useful for creating databases or cache services like redis. The runner on the virtual machine will automatically create a network and manage the life cycle of the service containers.
	// When you use a service container for a job or your step uses container actions, you don't need to set port information to access the service. Docker automatically exposes all ports between containers on the same network.
	// When both the job and the action run in a container, you can directly reference the container by its hostname. The hostname is automatically mapped to the service name.
	// When a step does not use a container action, you must access the service using localhost and bind the ports.
	Services map[string]*Container `mapstructure:"services,omitempty" yaml:"services,omitempty" json:"services,omitempty"`

	// A job contains a sequence of tasks called steps. Steps can run commands, run setup tasks, or run an action in your repository, a public repository, or an action published in a Docker registry. Not all steps run actions, but all actions run as a step. Each step runs in its own process in the virtual environment and has access to the workspace and filesystem. Because steps run in their own process, changes to environment variables are not preserved between steps. GitHub provides built-in steps to set up and complete a job.
	//
	Steps *[]Step `mapstructure:"steps,omitempty" yaml:"steps,omitempty" json:"steps,omitempty"`

	// A strategy creates a build matrix for your jobs. You can define different variations of an environment to run each job in.
	Strategy *Strategy `mapstructure:"strategy,omitempty" yaml:"strategy,omitempty" json:"strategy,omitempty"`

	// The maximum number of minutes to let a workflow run before GitHub automatically cancels it. Default: 360
	TimeoutMinutes *float64 `mapstructure:"timeout-minutes,omitempty" yaml:"timeout-minutes,omitempty" json:"timeout-minutes,omitempty"`
}

// ReusableWorkflowCallJob Each job must have an id to associate with the job. The key job_id is a string and its value is a map of the job's configuration data. You must replace <job_id> with a string that is unique to the jobs object. The <job_id> must start with a letter or _ and contain only alphanumeric characters, -, or _.
type ReusableWorkflowCallJob struct {

	// You can use the if conditional to prevent a job from running unless a condition is met. You can use any supported context and expression to create a conditional.
	// Expressions in an if conditional do not require the ${{ }} syntax. For more information, see https://help.github.com/en/articles/contexts-and-expression-syntax-for-github-actions.
	If string `mapstructure:"if,omitempty" yaml:"if,omitempty" json:"if,omitempty"`

	// The name of the job displayed on GitHub.
	Name        string            `mapstructure:"name,omitempty" yaml:"name,omitempty" json:"name,omitempty"`
	Needs       interface{}       `mapstructure:"needs,omitempty" yaml:"needs,omitempty" json:"needs,omitempty"`
	Permissions *PermissionsEvent `mapstructure:"permissions,omitempty" yaml:"permissions,omitempty" json:"permissions,omitempty"`

	// When a job is used to call a reusable workflow, you can use 'secrets' to provide a map of secrets that are passed to the called workflow. Any secrets that you pass must match the names defined in the called workflow.
	Secrets interface{} `mapstructure:"secrets,omitempty" yaml:"secrets,omitempty" json:"secrets,omitempty"`

	// The location and version of a reusable workflow file to run as a job, of the form './{path/to}/{localfile}.yml' or '{owner}/{repo}/{path}/{filename}@{ref}'. {ref} can be a SHA, a release tag, or a branch name. Using the commit SHA is the safest for stability and security.
	Uses string `mapstructure:"uses" yaml:"uses" json:"uses"`

	// A map of inputs that are passed to the called workflow. Any inputs that you pass must match the input specifications defined in the called workflow. Unlike 'jobs.<job_id>.steps[*].with', the inputs you pass with 'jobs.<job_id>.with' are not be available as environment variables in the called workflow. Instead, you can reference the inputs by using the inputs context.
	With interface{} `json:"with,omitempty"`
}
