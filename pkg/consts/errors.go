package consts

import "fmt"

type ErrInvalidPlatform struct {
	Platform Platform
}

func (e *ErrInvalidPlatform) Error() string {
	return fmt.Sprintf("invalid platform: %s", e.Platform)
}
