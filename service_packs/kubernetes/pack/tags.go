package kubernetes_pack

var (
	tags = map[string][]string{
		"@probes/kubernetes":                           []string{"k-cra", "k-gen", "k-iam", "k-iaf"},
		"@probes/kubernetes/general":                   []string{"k-gen"},
		"@probes/kubernetes/iam":                       []string{"k-iam"},
		"@probes/kubernetes/internet_access":           []string{"k-iaf"},
		"@standard/citihub/CHC2-APPDEV135":             []string{"k-cra"},
		"@standard/citihub/CHC2-ITS120":                []string{"k-cra"},
		"@control_type/preventative":                   []string{"k-cra-001", "k-cra-002", "k-cra-003", "k-iam-001", "k-iam-002", "k-iam-003", "k-iaf-001"},
		"@standard/cis":                                []string{"k-gen"},
		"@standard/cis/gke":                            []string{"k-gen"},
		"@standard/cis/gke/5.1.3":                      []string{"k-gen-001"},
		"@standard/cis/gke/5.6.3":                      []string{"k-gen-002"},
		"@standard/cis/gke/6":                          []string{"k-cra"},
		"@standard/cis/gke/6.1":                        []string{"k-cra"},
		"@standard/cis/gke/6.1.3":                      []string{"k-cra-001"},
		"@standard/cis/gke/6.1.4":                      []string{"k-cra-002"},
		"@standard/cis/gke/6.1.5":                      []string{"k-cra-003"},
		"@standard/cis/gke/6.10.1":                     []string{"k-gen-003"},
		"@csp/any":                                     []string{"k-cra", "k-gen", "k-iam"},
		"@csp/azure":                                   []string{"k-iam-001", "k-iam-002", "k-iam-003"},
		"@probes/kubernetes/container_registry_access": []string{"k-cra"},
		"@control_type/inspection":                     []string{"k-gen-001", "k-gen-002", "k-gen-003"},
		"@standard/citihub/CHC2-IAM105":                []string{"k-gen-001", "k-iam"},
		"@standard/citihub/CHC2-ITS115":                []string{"k-gen-003"},
		"@standard/citihub/CHC2-SVD010":                []string{"k-iaf"},
		"@category/iam":                                []string{"k-iam"},
		"@category/internet_access":                    []string{"k-iaf"},
		"@standard/citihub":                            []string{"k-iam", "k-iaf"},
	}
)
