package config

// SetOutputType
func (e *Config) SetOutputType(s string) {
	e.OutputType = s
}

// SetKubeConfigPath ...
func (e *Config) SetKubeConfigPath(p string) {
	e.KubeConfigPath = p
}
