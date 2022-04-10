package github

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/mitchellh/mapstructure"
	"gopkg.in/yaml.v3"
)

// Concurrency
type Concurrency struct {

	// To cancel any currently running job or workflow in the same concurrency group, specify cancel-in-progress: true.
	CancelInProgress interface{} `json:"cancel-in-progress,omitempty"`

	// When a concurrent job or workflow is queued, if another job or workflow using the same concurrency group in the repository is in progress, the queued job or workflow will be pending. Any previously pending job or workflow in the concurrency group will be canceled.
	Group string `json:"group"`
}

// Container
type Container struct {

	// If the image's container registry requires authentication to pull the image, you can use credentials to set a map of the username and password. The credentials are the same values that you would provide to the `docker login` command.
	Credentials *Credentials `json:"credentials,omitempty"`

	// Sets an array of environment variables in the container.
	Env interface{} `json:"env,omitempty"`

	// The Docker image to use as the container to run the action. The value can be the Docker Hub image name or a registry name.
	Image string `json:"image"`

	// Additional Docker container resource options. For a list of options, see https://docs.docker.com/engine/reference/commandline/create/#options.
	Options string `json:"options,omitempty"`

	// Sets an array of ports to expose on the container.
	Ports []interface{} `json:"ports,omitempty"`

	// Sets an array of volumes for the container to use. You can use volumes to share data between services or other steps in a job. You can specify named Docker volumes, anonymous Docker volumes, or bind mounts on the host.
	// To specify a volume, you specify the source and destination path: <source>:<destinationPath>
	// The <source> is a volume name or an absolute path on the host machine, and <destinationPath> is an absolute path in the container.
	Volumes []string `json:"volumes,omitempty"`
}

// Credentials If the image's container registry requires authentication to pull the image, you can use credentials to set a map of the username and password. The credentials are the same values that you would provide to the `docker login` command.
type Credentials struct {
	Password string `json:"password,omitempty"`
	Username string `json:"username,omitempty"`
}

// Defaults
type Defaults struct {
	Run *Run `json:"run,omitempty"`
}

// Environment The environment that the job references
type Environment struct {

	// The name of the environment configured in the repo.
	Name string `json:"name"`

	// A deployment URL
	Url string `json:"url,omitempty"`
}

// Jobs A workflow run is made up of one or more jobs. Jobs run in parallel by default. To run jobs sequentially, you can define dependencies on other jobs using the jobs.<job_id>.needs keyword.
// Each job runs in a fresh instance of the virtual environment specified by runs-on.
// You can run an unlimited number of jobs as long as you are within the workflow usage limits. For more information, see https://help.github.com/en/github/automating-your-workflow-with-github-actions/workflow-syntax-for-github-actions#usage-limits.
type Jobs struct {
	NormalJobs               map[string]*NormalJob
	ReusableWorkflowCallJobs map[string]*ReusableWorkflowCallJob
}

// type Jobs map[string]*NormalJob

// NormalJob Each job must have an id to associate with the job. The key job_id is a string and its value is a map of the job's configuration data. You must replace <job_id> with a string that is unique to the jobs object. The <job_id> must start with a letter or _ and contain only alphanumeric characters, -, or _.
type NormalJob struct {

	// Concurrency ensures that only a single job or workflow using the same concurrency group will run at a time. A concurrency group can be any string or expression. The expression can use any context except for the secrets context.
	// You can also specify concurrency at the workflow level.
	// When a concurrent job or workflow is queued, if another job or workflow using the same concurrency group in the repository is in progress, the queued job or workflow will be pending. Any previously pending job or workflow in the concurrency group will be canceled. To also cancel any currently running job or workflow in the same concurrency group, specify cancel-in-progress: true.
	Concurrency interface{} `json:"concurrency,omitempty"`

	// A container to run any steps in a job that don't already specify a container. If you have steps that use both script and container actions, the container actions will run as sibling containers on the same network with the same volume mounts.
	// If you do not set a container, all steps will run directly on the host specified by runs-on unless a step refers to an action configured to run in a container.
	Container interface{} `json:"container,omitempty"`

	// Prevents a workflow run from failing when a job fails. Set to true to allow a workflow run to pass when this job fails.
	ContinueOnError interface{} `json:"continue-on-error,omitempty"`

	// A map of default settings that will apply to all steps in the job.
	Defaults *Defaults `json:"defaults,omitempty"`

	// A map of environment variables that are available to all steps in the job.
	Env interface{} `json:"env,omitempty"`

	// The environment that the job references.
	Environment interface{} `json:"environment,omitempty"`

	// You can use the if conditional to prevent a job from running unless a condition is met. You can use any supported context and expression to create a conditional.
	// Expressions in an if conditional do not require the ${{ }} syntax. For more information, see https://help.github.com/en/articles/contexts-and-expression-syntax-for-github-actions.
	If string `json:"if,omitempty"`

	// The name of the job displayed on GitHub.
	Name  string      `json:"name,omitempty"`
	Needs interface{} `json:"needs,omitempty"`

	// A map of outputs for a job. Job outputs are available to all downstream jobs that depend on this job.
	Outputs     map[string]string `json:"outputs,omitempty"`
	Permissions *PermissionsEvent `json:"permissions,omitempty"`

	// The type of machine to run the job on. The machine can be either a GitHub-hosted runner, or a self-hosted runner.
	RunsOn interface{} `json:"runs-on"`

	// Additional containers to host services for a job in a workflow. These are useful for creating databases or cache services like redis. The runner on the virtual machine will automatically create a network and manage the life cycle of the service containers.
	// When you use a service container for a job or your step uses container actions, you don't need to set port information to access the service. Docker automatically exposes all ports between containers on the same network.
	// When both the job and the action run in a container, you can directly reference the container by its hostname. The hostname is automatically mapped to the service name.
	// When a step does not use a container action, you must access the service using localhost and bind the ports.
	Services map[string]*Container `json:"services,omitempty"`

	// A job contains a sequence of tasks called steps. Steps can run commands, run setup tasks, or run an action in your repository, a public repository, or an action published in a Docker registry. Not all steps run actions, but all actions run as a step. Each step runs in its own process in the virtual environment and has access to the workspace and filesystem. Because steps run in their own process, changes to environment variables are not preserved between steps. GitHub provides built-in steps to set up and complete a job.
	//
	Steps []*StepsItems `json:"steps,omitempty"`

	// A strategy creates a build matrix for your jobs. You can define different variations of an environment to run each job in.
	Strategy *Strategy `json:"strategy,omitempty"`

	// The maximum number of minutes to let a workflow run before GitHub automatically cancels it. Default: 360
	TimeoutMinutes float64 `json:"timeout-minutes,omitempty"`
}

// PermissionsEvent
type PermissionsEvent struct {
	Actions            string `json:"actions,omitempty"`
	Checks             string `json:"checks,omitempty"`
	Contents           string `json:"contents,omitempty"`
	Deployments        string `json:"deployments,omitempty"`
	Discussions        string `json:"discussions,omitempty"`
	IdToken            string `json:"id-token,omitempty"`
	Issues             string `json:"issues,omitempty"`
	Packages           string `json:"packages,omitempty"`
	Pages              string `json:"pages,omitempty"`
	PullRequests       string `json:"pull-requests,omitempty"`
	RepositoryProjects string `json:"repository-projects,omitempty"`
	SecurityEvents     string `json:"security-events,omitempty"`
	Statuses           string `json:"statuses,omitempty"`
}

// Ref
type Ref struct {
	Branches       []string `json:"branches,omitempty"`
	BranchesIgnore []string `json:"branches-ignore,omitempty"`
	Paths          []string `json:"paths,omitempty"`
	PathsIgnore    []string `json:"paths-ignore,omitempty"`
	Tags           []string `json:"tags,omitempty"`
	TagsIgnore     []string `json:"tags-ignore,omitempty"`
}

// ReusableWorkflowCallJob Each job must have an id to associate with the job. The key job_id is a string and its value is a map of the job's configuration data. You must replace <job_id> with a string that is unique to the jobs object. The <job_id> must start with a letter or _ and contain only alphanumeric characters, -, or _.
type ReusableWorkflowCallJob struct {

	// You can use the if conditional to prevent a job from running unless a condition is met. You can use any supported context and expression to create a conditional.
	// Expressions in an if conditional do not require the ${{ }} syntax. For more information, see https://help.github.com/en/articles/contexts-and-expression-syntax-for-github-actions.
	If string `json:"if,omitempty"`

	// The name of the job displayed on GitHub.
	Name        string            `json:"name,omitempty"`
	Needs       interface{}       `json:"needs,omitempty"`
	Permissions *PermissionsEvent `json:"permissions,omitempty"`

	// When a job is used to call a reusable workflow, you can use 'secrets' to provide a map of secrets that are passed to the called workflow. Any secrets that you pass must match the names defined in the called workflow.
	Secrets interface{} `json:"secrets,omitempty"`

	// The location and version of a reusable workflow file to run as a job, of the form './{path/to}/{localfile}.yml' or '{owner}/{repo}/{path}/{filename}@{ref}'. {ref} can be a SHA, a release tag, or a branch name. Using the commit SHA is the safest for stability and security.
	Uses string `json:"uses"`

	// A map of inputs that are passed to the called workflow. Any inputs that you pass must match the input specifications defined in the called workflow. Unlike 'jobs.<job_id>.steps[*].with', the inputs you pass with 'jobs.<job_id>.with' are not be available as environment variables in the called workflow. Instead, you can reference the inputs by using the inputs context.
	With interface{} `json:"with,omitempty"`
}

// Workflow
type Workflow struct {

	// Concurrency ensures that only a single job or workflow using the same concurrency group will run at a time. A concurrency group can be any string or expression. The expression can use any context except for the secrets context.
	// You can also specify concurrency at the workflow level.
	// When a concurrent job or workflow is queued, if another job or workflow using the same concurrency group in the repository is in progress, the queued job or workflow will be pending. Any previously pending job or workflow in the concurrency group will be canceled. To also cancel any currently running job or workflow in the same concurrency group, specify cancel-in-progress: true.
	Concurrency interface{} `json:"concurrency,omitempty"`

	// A map of default settings that will apply to all jobs in the workflow.
	Defaults *Defaults `json:"defaults,omitempty"`

	// A map of environment variables that are available to all jobs and steps in the workflow.
	Env interface{} `json:"env,omitempty"`

	// A workflow run is made up of one or more jobs. Jobs run in parallel by default. To run jobs sequentially, you can define dependencies on other jobs using the jobs.<job_id>.needs keyword.
	// Each job runs in a fresh instance of the virtual environment specified by runs-on.
	// You can run an unlimited number of jobs as long as you are within the workflow usage limits. For more information, see https://help.github.com/en/github/automating-your-workflow-with-github-actions/workflow-syntax-for-github-actions#usage-limits.
	Jobs *Jobs `yaml:"jobs"`

	// The name of your workflow. GitHub displays the names of your workflows on your repository's actions page. If you omit this field, GitHub sets the name to the workflow's filename.
	Name string `json:"name,omitempty"`

	// The name of the GitHub event that triggers the workflow. You can provide a single event string, array of events, array of event types, or an event configuration map that schedules a workflow or restricts the execution of a workflow to specific files, tags, or branch changes. For a list of available events, see https://help.github.com/en/github/automating-your-workflow-with-github-actions/events-that-trigger-workflows.
	On          interface{} `json:"on"`
	Permissions interface{} `json:"permissions,omitempty"`
}

// Run
type Run struct {
	Shell            interface{} `json:"shell,omitempty"`
	WorkingDirectory string      `json:"working-directory,omitempty"`
}

// StepsItems
type StepsItems struct {

	// Prevents a job from failing when a step fails. Set to true to allow a job to pass when this step fails.
	ContinueOnError interface{} `json:"continue-on-error,omitempty"`

	// Sets environment variables for steps to use in the virtual environment. You can also set environment variables for the entire workflow or a job.
	Env interface{} `json:"env,omitempty"`

	// A unique identifier for the step. You can use the id to reference the step in contexts. For more information, see https://help.github.com/en/articles/contexts-and-expression-syntax-for-github-actions.
	Id string `json:"id,omitempty"`

	// You can use the if conditional to prevent a step from running unless a condition is met. You can use any supported context and expression to create a conditional.
	// Expressions in an if conditional do not require the ${{ }} syntax. For more information, see https://help.github.com/en/articles/contexts-and-expression-syntax-for-github-actions.
	If string `json:"if,omitempty"`

	// A name for your step to display on GitHub.
	Name string `json:"name,omitempty"`

	// Runs command-line programs using the operating system's shell. If you do not provide a name, the step name will default to the text specified in the run command.
	// Commands run using non-login shells by default. You can choose a different shell and customize the shell used to run commands. For more information, see https://help.github.com/en/actions/automating-your-workflow-with-github-actions/workflow-syntax-for-github-actions#using-a-specific-shell.
	// Each run keyword represents a new process and shell in the virtual environment. When you provide multi-line commands, each line runs in the same shell.
	Run   string      `json:"run,omitempty"`
	Shell interface{} `json:"shell,omitempty"`

	// The maximum number of minutes to run the step before killing the process.
	TimeoutMinutes float64 `json:"timeout-minutes,omitempty"`

	// Selects an action to run as part of a step in your job. An action is a reusable unit of code. You can use an action defined in the same repository as the workflow, a public repository, or in a published Docker container image (https://hub.docker.com/).
	// We strongly recommend that you include the version of the action you are using by specifying a Git ref, SHA, or Docker tag number. If you don't specify a version, it could break your workflows or cause unexpected behavior when the action owner publishes an update.
	// - Using the commit SHA of a released action version is the safest for stability and security.
	// - Using the specific major action version allows you to receive critical fixes and security patches while still maintaining compatibility. It also assures that your workflow should still work.
	// - Using the master branch of an action may be convenient, but if someone releases a new major version with a breaking change, your workflow could break.
	// Some actions require inputs that you must set using the with keyword. Review the action's README file to determine the inputs required.
	// Actions are either JavaScript files or Docker containers. If the action you're using is a Docker container you must run the job in a Linux virtual environment. For more details, see https://help.github.com/en/articles/virtual-environments-for-github-actions.
	Uses string `json:"uses,omitempty"`

	// A map of the input parameters defined by the action. Each input parameter is a key/value pair. Input parameters are set as environment variables. The variable is prefixed with INPUT_ and converted to upper case.
	With             interface{} `json:"with,omitempty"`
	WorkingDirectory string      `json:"working-directory,omitempty"`
}

// Strategy A strategy creates a build matrix for your jobs. You can define different variations of an environment to run each job in.
type Strategy struct {

	// When set to true, GitHub cancels all in-progress jobs if any matrix job fails. Default: true
	FailFast bool `json:"fail-fast,omitempty"`

	// A build matrix is a set of different configurations of the virtual environment. For example you might run a job against more than one supported version of a language, operating system, or tool. Each configuration is a copy of the job that runs and reports a status.
	// You can specify a matrix by supplying an array for the configuration options. For example, if the GitHub virtual environment supports Node.js versions 6, 8, and 10 you could specify an array of those versions in the matrix.
	// When you define a matrix of operating systems, you must set the required runs-on keyword to the operating system of the current job, rather than hard-coding the operating system name. To access the operating system name, you can use the matrix.os context parameter to set runs-on. For more information, see https://help.github.com/en/articles/contexts-and-expression-syntax-for-github-actions.
	Matrix interface{} `json:"matrix"`

	// The maximum number of jobs that can run simultaneously when using a matrix job strategy. By default, GitHub will maximize the number of jobs run in parallel depending on the available runners on GitHub-hosted virtual machines.
	MaxParallel float64 `json:"max-parallel,omitempty"`
}

func (strct *Concurrency) MarshalJSON() ([]byte, error) {
	buf := bytes.NewBuffer(make([]byte, 0))
	buf.WriteString("{")
	comma := false
	// Marshal the "cancel-in-progress" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"cancel-in-progress\": ")
	if tmp, err := json.Marshal(strct.CancelInProgress); err != nil {
		return nil, err
	} else {
		buf.Write(tmp)
	}
	comma = true
	// "Group" field is required
	// only required object types supported for marshal checking (for now)
	// Marshal the "group" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"group\": ")
	if tmp, err := json.Marshal(strct.Group); err != nil {
		return nil, err
	} else {
		buf.Write(tmp)
	}
	comma = true

	buf.WriteString("}")
	rv := buf.Bytes()
	return rv, nil
}

func (j *Jobs) UnmarshalYAML(node *yaml.Node) error {
	var v map[string]any
	if err := node.Decode(&v); err != nil {
		return err
	}

	normalJobs := make(map[string]*NormalJob, 0)
	reusableWorkflowCallJobs := make(map[string]*ReusableWorkflowCallJob, 0)

	for k, v := range v {
		var normalJob *NormalJob
		var reusableWorkflowCallJob *ReusableWorkflowCallJob
		if err := mapstructure.Decode(v, &normalJob); err == nil {
			normalJobs[k] = normalJob
			continue
		} else if err := mapstructure.Decode(v, &reusableWorkflowCallJob); err == nil {
			reusableWorkflowCallJobs[k] = reusableWorkflowCallJob
			continue
		} else {
			return errors.New("unable to unmarshal jobs")
		}

	}
	*j = Jobs{
		NormalJobs:               normalJobs,
		ReusableWorkflowCallJobs: reusableWorkflowCallJobs,
	}
	return nil
}

func (strct *Concurrency) UnmarshalJSON(b []byte) error {
	groupReceived := false
	var jsonMap map[string]json.RawMessage
	if err := json.Unmarshal(b, &jsonMap); err != nil {
		return err
	}
	// parse all the defined properties
	for k, v := range jsonMap {
		switch k {
		case "cancel-in-progress":
			if err := json.Unmarshal([]byte(v), &strct.CancelInProgress); err != nil {
				return err
			}
		case "group":
			if err := json.Unmarshal([]byte(v), &strct.Group); err != nil {
				return err
			}
			groupReceived = true
		default:
			return fmt.Errorf("additional property not allowed: \"" + k + "\"")
		}
	}
	// check if group (a required property) was received
	if !groupReceived {
		return errors.New("\"group\" is required but was not present")
	}
	return nil
}

func (strct *Container) MarshalJSON() ([]byte, error) {
	buf := bytes.NewBuffer(make([]byte, 0))
	buf.WriteString("{")
	comma := false
	// Marshal the "credentials" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"credentials\": ")
	if tmp, err := json.Marshal(strct.Credentials); err != nil {
		return nil, err
	} else {
		buf.Write(tmp)
	}
	comma = true
	// Marshal the "env" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"env\": ")
	if tmp, err := json.Marshal(strct.Env); err != nil {
		return nil, err
	} else {
		buf.Write(tmp)
	}
	comma = true
	// "Image" field is required
	// only required object types supported for marshal checking (for now)
	// Marshal the "image" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"image\": ")
	if tmp, err := json.Marshal(strct.Image); err != nil {
		return nil, err
	} else {
		buf.Write(tmp)
	}
	comma = true
	// Marshal the "options" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"options\": ")
	if tmp, err := json.Marshal(strct.Options); err != nil {
		return nil, err
	} else {
		buf.Write(tmp)
	}
	comma = true
	// Marshal the "ports" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"ports\": ")
	if tmp, err := json.Marshal(strct.Ports); err != nil {
		return nil, err
	} else {
		buf.Write(tmp)
	}
	comma = true
	// Marshal the "volumes" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"volumes\": ")
	if tmp, err := json.Marshal(strct.Volumes); err != nil {
		return nil, err
	} else {
		buf.Write(tmp)
	}
	comma = true

	buf.WriteString("}")
	rv := buf.Bytes()
	return rv, nil
}

func (strct *Container) UnmarshalJSON(b []byte) error {
	imageReceived := false
	var jsonMap map[string]json.RawMessage
	if err := json.Unmarshal(b, &jsonMap); err != nil {
		return err
	}
	// parse all the defined properties
	for k, v := range jsonMap {
		switch k {
		case "credentials":
			if err := json.Unmarshal([]byte(v), &strct.Credentials); err != nil {
				return err
			}
		case "env":
			if err := json.Unmarshal([]byte(v), &strct.Env); err != nil {
				return err
			}
		case "image":
			if err := json.Unmarshal([]byte(v), &strct.Image); err != nil {
				return err
			}
			imageReceived = true
		case "options":
			if err := json.Unmarshal([]byte(v), &strct.Options); err != nil {
				return err
			}
		case "ports":
			if err := json.Unmarshal([]byte(v), &strct.Ports); err != nil {
				return err
			}
		case "volumes":
			if err := json.Unmarshal([]byte(v), &strct.Volumes); err != nil {
				return err
			}
		default:
			return fmt.Errorf("additional property not allowed: \"" + k + "\"")
		}
	}
	// check if image (a required property) was received
	if !imageReceived {
		return errors.New("\"image\" is required but was not present")
	}
	return nil
}

func (strct *Defaults) MarshalJSON() ([]byte, error) {
	buf := bytes.NewBuffer(make([]byte, 0))
	buf.WriteString("{")
	comma := false
	// Marshal the "run" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"run\": ")
	if tmp, err := json.Marshal(strct.Run); err != nil {
		return nil, err
	} else {
		buf.Write(tmp)
	}
	comma = true

	buf.WriteString("}")
	rv := buf.Bytes()
	return rv, nil
}

func (strct *Defaults) UnmarshalJSON(b []byte) error {
	var jsonMap map[string]json.RawMessage
	if err := json.Unmarshal(b, &jsonMap); err != nil {
		return err
	}
	// parse all the defined properties
	for k, v := range jsonMap {
		switch k {
		case "run":
			if err := json.Unmarshal([]byte(v), &strct.Run); err != nil {
				return err
			}
		default:
			return fmt.Errorf("additional property not allowed: \"" + k + "\"")
		}
	}
	return nil
}

func (strct *Environment) MarshalJSON() ([]byte, error) {
	buf := bytes.NewBuffer(make([]byte, 0))
	buf.WriteString("{")
	comma := false
	// "Name" field is required
	// only required object types supported for marshal checking (for now)
	// Marshal the "name" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"name\": ")
	if tmp, err := json.Marshal(strct.Name); err != nil {
		return nil, err
	} else {
		buf.Write(tmp)
	}
	comma = true
	// Marshal the "url" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"url\": ")
	if tmp, err := json.Marshal(strct.Url); err != nil {
		return nil, err
	} else {
		buf.Write(tmp)
	}
	comma = true

	buf.WriteString("}")
	rv := buf.Bytes()
	return rv, nil
}

func (strct *Environment) UnmarshalJSON(b []byte) error {
	nameReceived := false
	var jsonMap map[string]json.RawMessage
	if err := json.Unmarshal(b, &jsonMap); err != nil {
		return err
	}
	// parse all the defined properties
	for k, v := range jsonMap {
		switch k {
		case "name":
			if err := json.Unmarshal([]byte(v), &strct.Name); err != nil {
				return err
			}
			nameReceived = true
		case "url":
			if err := json.Unmarshal([]byte(v), &strct.Url); err != nil {
				return err
			}
		default:
			return fmt.Errorf("additional property not allowed: \"" + k + "\"")
		}
	}
	// check if name (a required property) was received
	if !nameReceived {
		return errors.New("\"name\" is required but was not present")
	}
	return nil
}

func (strct *Jobs) MarshalJSON() ([]byte, error) {
	buf := bytes.NewBuffer(make([]byte, 0))
	buf.WriteString("{")

	buf.WriteString("}")
	rv := buf.Bytes()
	return rv, nil
}

func (strct *NormalJob) MarshalJSON() ([]byte, error) {
	buf := bytes.NewBuffer(make([]byte, 0))
	buf.WriteString("{")
	comma := false
	// Marshal the "concurrency" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"concurrency\": ")
	if tmp, err := json.Marshal(strct.Concurrency); err != nil {
		return nil, err
	} else {
		buf.Write(tmp)
	}
	comma = true
	// Marshal the "container" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"container\": ")
	if tmp, err := json.Marshal(strct.Container); err != nil {
		return nil, err
	} else {
		buf.Write(tmp)
	}
	comma = true
	// Marshal the "continue-on-error" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"continue-on-error\": ")
	if tmp, err := json.Marshal(strct.ContinueOnError); err != nil {
		return nil, err
	} else {
		buf.Write(tmp)
	}
	comma = true
	// Marshal the "defaults" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"defaults\": ")
	if tmp, err := json.Marshal(strct.Defaults); err != nil {
		return nil, err
	} else {
		buf.Write(tmp)
	}
	comma = true
	// Marshal the "env" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"env\": ")
	if tmp, err := json.Marshal(strct.Env); err != nil {
		return nil, err
	} else {
		buf.Write(tmp)
	}
	comma = true
	// Marshal the "environment" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"environment\": ")
	if tmp, err := json.Marshal(strct.Environment); err != nil {
		return nil, err
	} else {
		buf.Write(tmp)
	}
	comma = true
	// Marshal the "if" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"if\": ")
	if tmp, err := json.Marshal(strct.If); err != nil {
		return nil, err
	} else {
		buf.Write(tmp)
	}
	comma = true
	// Marshal the "name" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"name\": ")
	if tmp, err := json.Marshal(strct.Name); err != nil {
		return nil, err
	} else {
		buf.Write(tmp)
	}
	comma = true
	// Marshal the "needs" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"needs\": ")
	if tmp, err := json.Marshal(strct.Needs); err != nil {
		return nil, err
	} else {
		buf.Write(tmp)
	}
	comma = true
	// Marshal the "outputs" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"outputs\": ")
	if tmp, err := json.Marshal(strct.Outputs); err != nil {
		return nil, err
	} else {
		buf.Write(tmp)
	}
	comma = true
	// Marshal the "permissions" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"permissions\": ")
	if tmp, err := json.Marshal(strct.Permissions); err != nil {
		return nil, err
	} else {
		buf.Write(tmp)
	}
	comma = true
	// "RunsOn" field is required
	// only required object types supported for marshal checking (for now)
	// Marshal the "runs-on" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"runs-on\": ")
	if tmp, err := json.Marshal(strct.RunsOn); err != nil {
		return nil, err
	} else {
		buf.Write(tmp)
	}
	comma = true
	// Marshal the "services" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"services\": ")
	if tmp, err := json.Marshal(strct.Services); err != nil {
		return nil, err
	} else {
		buf.Write(tmp)
	}
	comma = true
	// Marshal the "steps" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"steps\": ")
	if tmp, err := json.Marshal(strct.Steps); err != nil {
		return nil, err
	} else {
		buf.Write(tmp)
	}
	comma = true
	// Marshal the "strategy" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"strategy\": ")
	if tmp, err := json.Marshal(strct.Strategy); err != nil {
		return nil, err
	} else {
		buf.Write(tmp)
	}
	comma = true
	// Marshal the "timeout-minutes" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"timeout-minutes\": ")
	if tmp, err := json.Marshal(strct.TimeoutMinutes); err != nil {
		return nil, err
	} else {
		buf.Write(tmp)
	}
	comma = true

	buf.WriteString("}")
	rv := buf.Bytes()
	return rv, nil
}

func (strct *NormalJob) UnmarshalYAML(node *yaml.Node) error {
	var jsonMap map[string]json.RawMessage
	if err := node.Decode(&jsonMap); err != nil {
		return err
	}
	return strct.unmarshal(jsonMap)
}

func (strct *NormalJob) UnmarshalJSON(b []byte) error {
	var jsonMap map[string]json.RawMessage
	if err := json.Unmarshal(b, &jsonMap); err != nil {
		return err
	}

	return strct.unmarshal(jsonMap)
}

func (strct *NormalJob) unmarshal(jsonMap map[string]json.RawMessage) error {
	var runsOnReceived bool
	for k, v := range jsonMap {
		switch k {
		case "concurrency":
			if err := json.Unmarshal([]byte(v), &strct.Concurrency); err != nil {
				return err
			}
		case "container":
			if err := json.Unmarshal([]byte(v), &strct.Container); err != nil {
				return err
			}
		case "continue-on-error":
			if err := json.Unmarshal([]byte(v), &strct.ContinueOnError); err != nil {
				return err
			}
		case "defaults":
			if err := json.Unmarshal([]byte(v), &strct.Defaults); err != nil {
				return err
			}
		case "env":
			if err := json.Unmarshal([]byte(v), &strct.Env); err != nil {
				return err
			}
		case "environment":
			if err := json.Unmarshal([]byte(v), &strct.Environment); err != nil {
				return err
			}
		case "if":
			if err := json.Unmarshal([]byte(v), &strct.If); err != nil {
				return err
			}
		case "name":
			if err := json.Unmarshal([]byte(v), &strct.Name); err != nil {
				return err
			}
		case "needs":
			if err := json.Unmarshal([]byte(v), &strct.Needs); err != nil {
				return err
			}
		case "outputs":
			if err := json.Unmarshal([]byte(v), &strct.Outputs); err != nil {
				return err
			}
		case "permissions":
			if err := json.Unmarshal([]byte(v), &strct.Permissions); err != nil {
				return err
			}
		case "runs-on":
			if err := json.Unmarshal([]byte(v), &strct.RunsOn); err != nil {
				return err
			}
			runsOnReceived = true
		case "services":
			if err := json.Unmarshal([]byte(v), &strct.Services); err != nil {
				return err
			}
		case "steps":
			if err := json.Unmarshal([]byte(v), &strct.Steps); err != nil {
				return err
			}
		case "strategy":
			if err := json.Unmarshal([]byte(v), &strct.Strategy); err != nil {
				return err
			}
		case "timeout-minutes":
			if err := json.Unmarshal([]byte(v), &strct.TimeoutMinutes); err != nil {
				return err
			}
		default:
			return fmt.Errorf("additional property not allowed: \"" + k + "\"")
		}
	}
	// check if runs-on (a required property) was received
	if !runsOnReceived {
		return errors.New("\"runs-on\" is required but was not present")
	}
	return nil
}

func (strct *PermissionsEvent) MarshalJSON() ([]byte, error) {
	buf := bytes.NewBuffer(make([]byte, 0))
	buf.WriteString("{")
	comma := false
	// Marshal the "actions" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"actions\": ")
	if tmp, err := json.Marshal(strct.Actions); err != nil {
		return nil, err
	} else {
		buf.Write(tmp)
	}
	comma = true
	// Marshal the "checks" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"checks\": ")
	if tmp, err := json.Marshal(strct.Checks); err != nil {
		return nil, err
	} else {
		buf.Write(tmp)
	}
	comma = true
	// Marshal the "contents" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"contents\": ")
	if tmp, err := json.Marshal(strct.Contents); err != nil {
		return nil, err
	} else {
		buf.Write(tmp)
	}
	comma = true
	// Marshal the "deployments" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"deployments\": ")
	if tmp, err := json.Marshal(strct.Deployments); err != nil {
		return nil, err
	} else {
		buf.Write(tmp)
	}
	comma = true
	// Marshal the "discussions" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"discussions\": ")
	if tmp, err := json.Marshal(strct.Discussions); err != nil {
		return nil, err
	} else {
		buf.Write(tmp)
	}
	comma = true
	// Marshal the "id-token" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"id-token\": ")
	if tmp, err := json.Marshal(strct.IdToken); err != nil {
		return nil, err
	} else {
		buf.Write(tmp)
	}
	comma = true
	// Marshal the "issues" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"issues\": ")
	if tmp, err := json.Marshal(strct.Issues); err != nil {
		return nil, err
	} else {
		buf.Write(tmp)
	}
	comma = true
	// Marshal the "packages" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"packages\": ")
	if tmp, err := json.Marshal(strct.Packages); err != nil {
		return nil, err
	} else {
		buf.Write(tmp)
	}
	comma = true
	// Marshal the "pages" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"pages\": ")
	if tmp, err := json.Marshal(strct.Pages); err != nil {
		return nil, err
	} else {
		buf.Write(tmp)
	}
	comma = true
	// Marshal the "pull-requests" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"pull-requests\": ")
	if tmp, err := json.Marshal(strct.PullRequests); err != nil {
		return nil, err
	} else {
		buf.Write(tmp)
	}
	comma = true
	// Marshal the "repository-projects" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"repository-projects\": ")
	if tmp, err := json.Marshal(strct.RepositoryProjects); err != nil {
		return nil, err
	} else {
		buf.Write(tmp)
	}
	comma = true
	// Marshal the "security-events" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"security-events\": ")
	if tmp, err := json.Marshal(strct.SecurityEvents); err != nil {
		return nil, err
	} else {
		buf.Write(tmp)
	}
	comma = true
	// Marshal the "statuses" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"statuses\": ")
	if tmp, err := json.Marshal(strct.Statuses); err != nil {
		return nil, err
	} else {
		buf.Write(tmp)
	}
	comma = true

	buf.WriteString("}")
	rv := buf.Bytes()
	return rv, nil
}

func (strct *PermissionsEvent) UnmarshalJSON(b []byte) error {
	var jsonMap map[string]json.RawMessage
	if err := json.Unmarshal(b, &jsonMap); err != nil {
		return err
	}
	// parse all the defined properties
	for k, v := range jsonMap {
		switch k {
		case "actions":
			if err := json.Unmarshal([]byte(v), &strct.Actions); err != nil {
				return err
			}
		case "checks":
			if err := json.Unmarshal([]byte(v), &strct.Checks); err != nil {
				return err
			}
		case "contents":
			if err := json.Unmarshal([]byte(v), &strct.Contents); err != nil {
				return err
			}
		case "deployments":
			if err := json.Unmarshal([]byte(v), &strct.Deployments); err != nil {
				return err
			}
		case "discussions":
			if err := json.Unmarshal([]byte(v), &strct.Discussions); err != nil {
				return err
			}
		case "id-token":
			if err := json.Unmarshal([]byte(v), &strct.IdToken); err != nil {
				return err
			}
		case "issues":
			if err := json.Unmarshal([]byte(v), &strct.Issues); err != nil {
				return err
			}
		case "packages":
			if err := json.Unmarshal([]byte(v), &strct.Packages); err != nil {
				return err
			}
		case "pages":
			if err := json.Unmarshal([]byte(v), &strct.Pages); err != nil {
				return err
			}
		case "pull-requests":
			if err := json.Unmarshal([]byte(v), &strct.PullRequests); err != nil {
				return err
			}
		case "repository-projects":
			if err := json.Unmarshal([]byte(v), &strct.RepositoryProjects); err != nil {
				return err
			}
		case "security-events":
			if err := json.Unmarshal([]byte(v), &strct.SecurityEvents); err != nil {
				return err
			}
		case "statuses":
			if err := json.Unmarshal([]byte(v), &strct.Statuses); err != nil {
				return err
			}
		default:
			return fmt.Errorf("additional property not allowed: \"" + k + "\"")
		}
	}
	return nil
}

func (strct *ReusableWorkflowCallJob) MarshalJSON() ([]byte, error) {
	buf := bytes.NewBuffer(make([]byte, 0))
	buf.WriteString("{")
	comma := false
	// Marshal the "if" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"if\": ")
	if tmp, err := json.Marshal(strct.If); err != nil {
		return nil, err
	} else {
		buf.Write(tmp)
	}
	comma = true
	// Marshal the "name" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"name\": ")
	if tmp, err := json.Marshal(strct.Name); err != nil {
		return nil, err
	} else {
		buf.Write(tmp)
	}
	comma = true
	// Marshal the "needs" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"needs\": ")
	if tmp, err := json.Marshal(strct.Needs); err != nil {
		return nil, err
	} else {
		buf.Write(tmp)
	}
	comma = true
	// Marshal the "permissions" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"permissions\": ")
	if tmp, err := json.Marshal(strct.Permissions); err != nil {
		return nil, err
	} else {
		buf.Write(tmp)
	}
	comma = true
	// Marshal the "secrets" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"secrets\": ")
	if tmp, err := json.Marshal(strct.Secrets); err != nil {
		return nil, err
	} else {
		buf.Write(tmp)
	}
	comma = true
	// "Uses" field is required
	// only required object types supported for marshal checking (for now)
	// Marshal the "uses" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"uses\": ")
	if tmp, err := json.Marshal(strct.Uses); err != nil {
		return nil, err
	} else {
		buf.Write(tmp)
	}
	comma = true
	// Marshal the "with" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"with\": ")
	if tmp, err := json.Marshal(strct.With); err != nil {
		return nil, err
	} else {
		buf.Write(tmp)
	}
	comma = true

	buf.WriteString("}")
	rv := buf.Bytes()
	return rv, nil
}

func (strct *ReusableWorkflowCallJob) UnmarshalJSON(b []byte) error {
	usesReceived := false
	var jsonMap map[string]json.RawMessage
	if err := json.Unmarshal(b, &jsonMap); err != nil {
		return err
	}
	// parse all the defined properties
	for k, v := range jsonMap {
		switch k {
		case "if":
			if err := json.Unmarshal([]byte(v), &strct.If); err != nil {
				return err
			}
		case "name":
			if err := json.Unmarshal([]byte(v), &strct.Name); err != nil {
				return err
			}
		case "needs":
			if err := json.Unmarshal([]byte(v), &strct.Needs); err != nil {
				return err
			}
		case "permissions":
			if err := json.Unmarshal([]byte(v), &strct.Permissions); err != nil {
				return err
			}
		case "secrets":
			if err := json.Unmarshal([]byte(v), &strct.Secrets); err != nil {
				return err
			}
		case "uses":
			if err := json.Unmarshal([]byte(v), &strct.Uses); err != nil {
				return err
			}
			usesReceived = true
		case "with":
			if err := json.Unmarshal([]byte(v), &strct.With); err != nil {
				return err
			}
		default:
			return fmt.Errorf("additional property not allowed: \"" + k + "\"")
		}
	}
	// check if uses (a required property) was received
	if !usesReceived {
		return errors.New("\"uses\" is required but was not present")
	}
	return nil
}

func (strct *Workflow) MarshalJSON() ([]byte, error) {
	buf := bytes.NewBuffer(make([]byte, 0))
	buf.WriteString("{")
	comma := false
	// Marshal the "concurrency" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"concurrency\": ")
	if tmp, err := json.Marshal(strct.Concurrency); err != nil {
		return nil, err
	} else {
		buf.Write(tmp)
	}
	comma = true
	// Marshal the "defaults" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"defaults\": ")
	if tmp, err := json.Marshal(strct.Defaults); err != nil {
		return nil, err
	} else {
		buf.Write(tmp)
	}
	comma = true
	// Marshal the "env" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"env\": ")
	if tmp, err := json.Marshal(strct.Env); err != nil {
		return nil, err
	} else {
		buf.Write(tmp)
	}
	comma = true
	// "Jobs" field is required
	if strct.Jobs == nil {
		return nil, errors.New("jobs is a required field")
	}
	// Marshal the "jobs" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"jobs\": ")
	if tmp, err := json.Marshal(strct.Jobs); err != nil {
		return nil, err
	} else {
		buf.Write(tmp)
	}
	comma = true
	// Marshal the "name" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"name\": ")
	if tmp, err := json.Marshal(strct.Name); err != nil {
		return nil, err
	} else {
		buf.Write(tmp)
	}
	comma = true
	// "On" field is required // only required object types supported for marshal checking (for now)

	// Marshal the "on" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"on\": ")
	if tmp, err := json.Marshal(strct.On); err != nil {
		return nil, err
	} else {
		buf.Write(tmp)
	}
	comma = true
	// Marshal the "permissions" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"permissions\": ")
	if tmp, err := json.Marshal(strct.Permissions); err != nil {
		return nil, err
	} else {
		buf.Write(tmp)
	}
	comma = true

	buf.WriteString("}")
	rv := buf.Bytes()
	return rv, nil
}

func (strct *Workflow) UnmarshalJSON(b []byte) error {
	jobsReceived := false
	onReceived := false
	var jsonMap map[string]json.RawMessage
	if err := json.Unmarshal(b, &jsonMap); err != nil {
		return err
	}
	// parse all the defined properties
	for k, v := range jsonMap {
		switch k {
		case "concurrency":
			if err := json.Unmarshal([]byte(v), &strct.Concurrency); err != nil {
				return err
			}
		case "defaults":
			if err := json.Unmarshal([]byte(v), &strct.Defaults); err != nil {
				return err
			}
		case "env":
			if err := json.Unmarshal([]byte(v), &strct.Env); err != nil {
				return err
			}
		case "jobs":
			if err := json.Unmarshal([]byte(v), &strct.Jobs); err != nil {
				return err
			}
			jobsReceived = true
		case "name":
			if err := json.Unmarshal([]byte(v), &strct.Name); err != nil {
				return err
			}
		case "on":
			if err := json.Unmarshal([]byte(v), &strct.On); err != nil {
				return err
			}
			onReceived = true
		case "permissions":
			if err := json.Unmarshal([]byte(v), &strct.Permissions); err != nil {
				return err
			}
		default:
			return fmt.Errorf("additional property not allowed: \"" + k + "\"")
		}
	}
	// check if jobs (a required property) was received
	if !jobsReceived {
		return errors.New("\"jobs\" is required but was not present")
	}
	// check if on (a required property) was received
	if !onReceived {
		return errors.New("\"on\" is required but was not present")
	}
	return nil
}

func (strct *Run) MarshalJSON() ([]byte, error) {
	buf := bytes.NewBuffer(make([]byte, 0))
	buf.WriteString("{")
	comma := false
	// Marshal the "shell" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"shell\": ")
	if tmp, err := json.Marshal(strct.Shell); err != nil {
		return nil, err
	} else {
		buf.Write(tmp)
	}
	comma = true
	// Marshal the "working-directory" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"working-directory\": ")
	if tmp, err := json.Marshal(strct.WorkingDirectory); err != nil {
		return nil, err
	} else {
		buf.Write(tmp)
	}
	comma = true

	buf.WriteString("}")
	rv := buf.Bytes()
	return rv, nil
}

func (strct *Run) UnmarshalJSON(b []byte) error {
	var jsonMap map[string]json.RawMessage
	if err := json.Unmarshal(b, &jsonMap); err != nil {
		return err
	}
	// parse all the defined properties
	for k, v := range jsonMap {
		switch k {
		case "shell":
			if err := json.Unmarshal([]byte(v), &strct.Shell); err != nil {
				return err
			}
		case "working-directory":
			if err := json.Unmarshal([]byte(v), &strct.WorkingDirectory); err != nil {
				return err
			}
		default:
			return fmt.Errorf("additional property not allowed: \"" + k + "\"")
		}
	}
	return nil
}

func (strct *StepsItems) MarshalJSON() ([]byte, error) {
	buf := bytes.NewBuffer(make([]byte, 0))
	buf.WriteString("{")
	comma := false
	// Marshal the "continue-on-error" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"continue-on-error\": ")
	if tmp, err := json.Marshal(strct.ContinueOnError); err != nil {
		return nil, err
	} else {
		buf.Write(tmp)
	}
	comma = true
	// Marshal the "env" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"env\": ")
	if tmp, err := json.Marshal(strct.Env); err != nil {
		return nil, err
	} else {
		buf.Write(tmp)
	}
	comma = true
	// Marshal the "id" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"id\": ")
	if tmp, err := json.Marshal(strct.Id); err != nil {
		return nil, err
	} else {
		buf.Write(tmp)
	}
	comma = true
	// Marshal the "if" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"if\": ")
	if tmp, err := json.Marshal(strct.If); err != nil {
		return nil, err
	} else {
		buf.Write(tmp)
	}
	comma = true
	// Marshal the "name" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"name\": ")
	if tmp, err := json.Marshal(strct.Name); err != nil {
		return nil, err
	} else {
		buf.Write(tmp)
	}
	comma = true
	// Marshal the "run" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"run\": ")
	if tmp, err := json.Marshal(strct.Run); err != nil {
		return nil, err
	} else {
		buf.Write(tmp)
	}
	comma = true
	// Marshal the "shell" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"shell\": ")
	if tmp, err := json.Marshal(strct.Shell); err != nil {
		return nil, err
	} else {
		buf.Write(tmp)
	}
	comma = true
	// Marshal the "timeout-minutes" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"timeout-minutes\": ")
	if tmp, err := json.Marshal(strct.TimeoutMinutes); err != nil {
		return nil, err
	} else {
		buf.Write(tmp)
	}
	comma = true
	// Marshal the "uses" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"uses\": ")
	if tmp, err := json.Marshal(strct.Uses); err != nil {
		return nil, err
	} else {
		buf.Write(tmp)
	}
	comma = true
	// Marshal the "with" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"with\": ")
	if tmp, err := json.Marshal(strct.With); err != nil {
		return nil, err
	} else {
		buf.Write(tmp)
	}
	comma = true
	// Marshal the "working-directory" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"working-directory\": ")
	if tmp, err := json.Marshal(strct.WorkingDirectory); err != nil {
		return nil, err
	} else {
		buf.Write(tmp)
	}
	comma = true

	buf.WriteString("}")
	rv := buf.Bytes()
	return rv, nil
}

func (strct *StepsItems) UnmarshalJSON(b []byte) error {
	var jsonMap map[string]json.RawMessage
	if err := json.Unmarshal(b, &jsonMap); err != nil {
		return err
	}
	// parse all the defined properties
	for k, v := range jsonMap {
		switch k {
		case "continue-on-error":
			if err := json.Unmarshal([]byte(v), &strct.ContinueOnError); err != nil {
				return err
			}
		case "env":
			if err := json.Unmarshal([]byte(v), &strct.Env); err != nil {
				return err
			}
		case "id":
			if err := json.Unmarshal([]byte(v), &strct.Id); err != nil {
				return err
			}
		case "if":
			if err := json.Unmarshal([]byte(v), &strct.If); err != nil {
				return err
			}
		case "name":
			if err := json.Unmarshal([]byte(v), &strct.Name); err != nil {
				return err
			}
		case "run":
			if err := json.Unmarshal([]byte(v), &strct.Run); err != nil {
				return err
			}
		case "shell":
			if err := json.Unmarshal([]byte(v), &strct.Shell); err != nil {
				return err
			}
		case "timeout-minutes":
			if err := json.Unmarshal([]byte(v), &strct.TimeoutMinutes); err != nil {
				return err
			}
		case "uses":
			if err := json.Unmarshal([]byte(v), &strct.Uses); err != nil {
				return err
			}
		case "with":
			if err := json.Unmarshal([]byte(v), &strct.With); err != nil {
				return err
			}
		case "working-directory":
			if err := json.Unmarshal([]byte(v), &strct.WorkingDirectory); err != nil {
				return err
			}
		default:
			return fmt.Errorf("additional property not allowed: \"" + k + "\"")
		}
	}
	return nil
}

func (strct *Strategy) MarshalJSON() ([]byte, error) {
	buf := bytes.NewBuffer(make([]byte, 0))
	buf.WriteString("{")
	comma := false
	// Marshal the "fail-fast" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"fail-fast\": ")
	if tmp, err := json.Marshal(strct.FailFast); err != nil {
		return nil, err
	} else {
		buf.Write(tmp)
	}
	comma = true
	// "Matrix" field is required
	// only required object types supported for marshal checking (for now)
	// Marshal the "matrix" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"matrix\": ")
	if tmp, err := json.Marshal(strct.Matrix); err != nil {
		return nil, err
	} else {
		buf.Write(tmp)
	}
	comma = true
	// Marshal the "max-parallel" field
	if comma {
		buf.WriteString(",")
	}
	buf.WriteString("\"max-parallel\": ")
	if tmp, err := json.Marshal(strct.MaxParallel); err != nil {
		return nil, err
	} else {
		buf.Write(tmp)
	}
	comma = true

	buf.WriteString("}")
	rv := buf.Bytes()
	return rv, nil
}

func (strct *Strategy) UnmarshalJSON(b []byte) error {
	matrixReceived := false
	var jsonMap map[string]json.RawMessage
	if err := json.Unmarshal(b, &jsonMap); err != nil {
		return err
	}
	// parse all the defined properties
	for k, v := range jsonMap {
		switch k {
		case "fail-fast":
			if err := json.Unmarshal([]byte(v), &strct.FailFast); err != nil {
				return err
			}
		case "matrix":
			if err := json.Unmarshal([]byte(v), &strct.Matrix); err != nil {
				return err
			}
			matrixReceived = true
		case "max-parallel":
			if err := json.Unmarshal([]byte(v), &strct.MaxParallel); err != nil {
				return err
			}
		default:
			return fmt.Errorf("additional property not allowed: \"" + k + "\"")
		}
	}
	// check if matrix (a required property) was received
	if !matrixReceived {
		return errors.New("\"matrix\" is required but was not present")
	}
	return nil
}
