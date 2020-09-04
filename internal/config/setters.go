package config

// SetOutputType ...
func (e *ConfigVars) SetOutputType(s string) {
	e.OutputType = s
}

// SetKubeConfigPath ...
func (e *ConfigVars) SetKubeConfigPath(p string) {
	e.KubeConfigPath = p
}

// SetTags ...
func (e *ConfigVars) SetTags(t string) {
	e.Tests.Tags = t
}