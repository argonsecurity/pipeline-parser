package consts

type OutputTarget string

const (
	Stdout OutputTarget = "stdout"
	File   OutputTarget = "file"
)

var OutputTargets = []OutputTarget{
	Stdout,
	File,
}
