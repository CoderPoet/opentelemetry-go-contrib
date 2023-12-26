package request

// Version is the current release version of the request metrics instrumentation.
func Version() string {
	return "0.46.0"
	// This string is updated by the pre_release.sh script during release
}

// SemVersion is the semantic version to be supplied to meter creation.
func SemVersion() string {
	return "semver:" + Version()
}
