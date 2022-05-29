package consts

import "fmt"

type ErrInvalidPlatform struct {
	Platform Platform
}

func (e *ErrInvalidPlatform) Error() string {
	return fmt.Sprintf("invalid platform: %s", e.Platform)
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
