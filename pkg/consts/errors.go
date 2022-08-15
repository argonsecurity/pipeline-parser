package consts

import "fmt"

type ErrInvalidPlatform struct {
	Platform Platform
}

func (e *ErrInvalidPlatform) Error() string {
	return fmt.Sprintf("invalid platform: %s. Supported platforms: %v", e.Platform, Platforms)
}

func NewErrInvalidPlatform(platform Platform) error {
	return &ErrInvalidPlatform{Platform: platform}
}

type ErrInvalidOutputTarget struct {
	OutputTarget OutputTarget
}

func (e *ErrInvalidOutputTarget) Error() string {
	return fmt.Sprintf("invalid output target: %s. Supported output targets: %v", e.OutputTarget, OutputTargets)
}

func NewErrInvalidOutputTarget(outputTarget OutputTarget) error {
	return &ErrInvalidOutputTarget{OutputTarget: outputTarget}
}

type ErrInvalidYaml struct {
	Message string
}

func (e *ErrInvalidYaml) Error() string {
	return fmt.Sprintf("invalid yaml: %s", e.Message)
}

func NewErrInvalidYaml(message string) error {
	return &ErrInvalidYaml{Message: message}
}

type ErrInvalidYamlTag struct {
	Tag  string
	Type string
}

func (e *ErrInvalidYamlTag) Error() string {
	return fmt.Sprintf("invalid yaml tag '%s' for type '%s'", e.Tag, e.Type)
}

func NewErrInvalidYamlTag(tag string, structType string) error {
	return &ErrInvalidYamlTag{Tag: tag, Type: structType}
}

type ErrInvalidArgumentsCount struct {
	Count int
}

func (e *ErrInvalidArgumentsCount) Error() string {
	return fmt.Sprintf("invalid number of arguments: %d. Expected minimum 1 argument", e.Count)
}

func NewErrInvalidArgumentsCount(count int) error {
	return &ErrInvalidArgumentsCount{Count: count}
}
