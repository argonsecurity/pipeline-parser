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

type ErrInvalidYamlTag struct {
	Tag string
}

func (e *ErrInvalidYamlTag) Error() string {
	return fmt.Sprintf("invalid yaml tag: %s", e.Tag)
}

func NewErrInvalidYamlTag(tag string) error {
	return &ErrInvalidYamlTag{Tag: tag}
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
