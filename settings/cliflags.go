package settings

// VersionFlags ...
type VersionFlags struct {
	Verbose *bool
}

// RunFlags ...
type RunFlags struct {
	VarsFile       *string
	WriteDirectory *string
	LogLevel       *string
	ResultsFormat  *string
	Tags           *string
	KubeConfig     *string
}

//VersionCliFlags ...
var VersionCliFlags VersionFlags = VersionFlags{}

// RunCliFlags ...
var RunCliFlags RunFlags = RunFlags{}
