package kubernetes_pack

var (
	tags = map[string][]string{
		"@probes/kubernetes":                           []string{"k-cra"},
		"@standard/citihub/CHC2-APPDEV135":             []string{"k-cra"},
		"@standard/citihub/CHC2-ITS120":                []string{"k-cra"},
		"@control_type/preventative":                   []string{"k-cra-001", "k-cra-002", "k-cra-003"},
		"@standard/cis/gke/6":                          []string{"k-cra"},
		"@standard/cis/gke/6.1":                        []string{"k-cra"},
		"@standard/cis/gke/6.1.3":                      []string{"k-cra-001"},
		"@standard/cis/gke/6.1.4":                      []string{"k-cra-002"},
		"@standard/cis/gke/6.1.5":                      []string{"k-cra-003"},
		"@csp/any":                                     []string{"k-cra"},
		"@probes/kubernetes/container_registry_access": []string{"k-cra"},
	}
)
