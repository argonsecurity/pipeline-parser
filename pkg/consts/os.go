package consts

import "github.com/argonsecurity/pipeline-parser/pkg/models"

const (
	Macos1015     = "macos-10.15"
	Macos11       = "macos-11"
	MacosLatest   = "macos-latest"
	SelfHosted    = "self-hosted"
	Ubuntu1804    = "ubuntu-18.04"
	Ubuntu2004    = "ubuntu-20.04"
	UbuntuLatest  = "ubuntu-latest"
	Windows2016   = "windows-2016"
	Windows2019   = "windows-2019"
	Windows2022   = "windows-2022"
	WindowsLatest = "windows-latest"

	ArmArch = "arm32"
	X64Arch = "x64"
	X32Arch = "x32"
)

var (
	WindowsKeywords = []string{"windows", Windows2016, Windows2019, Windows2022, WindowsLatest}
	LinuxKeywords   = []string{"linux", "ubuntu", "debian", Ubuntu1804, Ubuntu2004, UbuntuLatest}
	MacKeywords     = []string{"macos", "darwin", "osx", Macos1015, Macos11, MacosLatest}

	ArchKeywords = []string{ArmArch, X64Arch, X32Arch}

	OsToKeywords = map[models.OS][]string{
		models.WindowsOS: WindowsKeywords,
		models.LinuxOS:   LinuxKeywords,
		models.MacOS:     MacKeywords,
	}
)
