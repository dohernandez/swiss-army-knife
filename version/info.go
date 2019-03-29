package version

import (
	"runtime"
)

// Build information. Populated at build-time.
// nolint:gochecknoglobals
var (
	version   = "dev"
	revision  string
	branch    string
	buildUser string
	buildDate string
	goVersion = runtime.Version()
)

// Information holds command version info.
type Information struct {
	Version   string `json:"version"`
	Revision  string `json:"revision"`
	Branch    string `json:"branch"`
	BuildUser string `json:"build_user"`
	BuildDate string `json:"build_date"`
	GoVersion string `json:"go_version"`
}

// Info returns command version info.
func Info() Information {
	return Information{
		Version:   version,
		Revision:  revision,
		Branch:    branch,
		BuildUser: buildUser,
		BuildDate: buildDate,
		GoVersion: goVersion,
	}
}

// String return version information as string.
func (i Information) String() string {
	return "Version: " + i.Version +
		", Revision: " + i.Revision +
		", Branch: " + i.Branch +
		", BuildUser: " + i.BuildUser +
		", BuildDate: " + i.BuildDate +
		", GoVersion: " + i.GoVersion
}

// Values return version information as string map.
func (i Information) Values() map[string]string {
	return map[string]string{
		"version":    i.Version,
		"revision":   i.Revision,
		"branch":     i.Branch,
		"build_user": i.BuildUser,
		"build_date": i.BuildDate,
		"go_version": i.GoVersion,
	}
}
