package storagepack

var (
	tags = map[string][]string{
		"@probes/storage":                      []string{"@s-az-aw", "@s-az-ear", "@s-az-eif"},
		"@probes/storage/access_whitelisting":  []string{"@s-az-aw"},
		"@standard/citihub/CHC2-CHC2-SVD030":   []string{"@s-az-aw"},
		"@csp.aws":                             []string{"@s-az-aw"},
		"@csp.azure":                           []string{"@s-az-aw", "@s-az-ear-001", "@s-az-ear-002", "@s-az-eif-001", "@s-az-eif-002"},
		"@control_type/detective":              []string{"@s-az-aw-001", "@s-az-ear-002", "@s-az-eif-002"},
		"@control_type/preventative":           []string{"@s-az-aw-002", "@s-az-ear-001", "@s-az-eif-001"},
		"@probes/storage/encryption_at_rest":   []string{"@s-az-ear"},
		"@standard/citihub/CHC2-SVD001":        []string{"@s-az-ear", "@s-az-eif"},
		"@standard/citihub/CHC2-AGP140":        []string{"@s-az-ear", "@s-az-eif"},
		"@standard/citihub/CHC2-EUC001":        []string{"@s-az-ear"},
		"@probes/storage/encryption_in_flight": []string{"@s-az-eif"},
	}
)
