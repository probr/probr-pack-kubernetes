# Probr Kubernetes Service Pack
## Probes Provenance

The Kubernetes Probr service pack has been built based on the [CIS Kubernetes Benchmark 1.5.1](https://www.cisecurity.org/cis-benchmarks/).  

This pack only implements the benchmarks that can be run _as an end user_ in any managed Kubernetes services. It doesn't look "under the hood" - use a tool like [kube-bench](https://github.com/aquasecurity/kube-bench) for that.

We currently have a complimentary Service Pack for [AKS](https://github.com/probr/probr-pack-aks), with plans for GKE and EKS.

## Controls covered

| CIS ID | CIS Policy Statement | Probr Implementation | Suggested further improvements |
| ------ | ------               | -------------------- | ------------------- |
| 5.2.1	| Minimize the admission of privileged containers	| Attempt to deploy non-compliant pod; run command that should be blocked  | - |
| 5.2.2	| Minimize the admission of containers wishing to share the host process ID namespace	| Attempt to deploy non-compliant pod; run command that should be blocked | - |
| 5.2.3	| Minimize the admission of containers wishing to share the host IPC namespace	| Attempt to deploy non-compliant pod; run command that should be blocked | - |
| 5.2.4	| Minimize the admission of containers wishing to share the host network namespace	| Attempt to deploy non-compliant pod; run command that should be blocked | - |
| 5.2.5	| Minimize the admission of containers with allowPrivilegeEscalation	| Attempt to deploy non-compliant pod; run command that should be blocked | - |
| 5.2.6	| Minimize the admission of root containers	| Attempt to deploy non-compliant pod; run command that should be blocked | - |
| 5.2.7	| Minimize the admission of containers with the NET_RAW capability	| Attempt to deploy non-compliant pod; run command that should be blocked | - |
| 5.2.8	| Minimize the admission of containers with added capabilities	| Attempt to deploy non-compliant pod; run command that should be blocked | - |
| 5.2.9	| Minimize the admission of containers with capabilities assigned	| Attempt to deploy non-compliant pod; run command that should be blocked | - |
| 5.6.2 | Ensure that the seccomp profile is set to docker/default in your pod definitions | Put logic here | - |
| 5.6.4	| The default namespace should not be used |	Attempt to deploy a Pod to the default namespace | - |
| 6.10.1	| Ensure Kubernetes Web UI is Disabled | look for kubernetes dashboard pod in kube-system namespace | - |
| 6.10.3 | Ensure Pod Security Policy is Enabled and set as appropriate | Tests per 5.2.x | - |
