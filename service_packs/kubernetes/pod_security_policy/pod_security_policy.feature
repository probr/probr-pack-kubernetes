@k-psp
@probes/kubernetes/psp
Feature: Maximise security through Pod Security Policies

    As a Cloud Security Administrator
    I want to ensure that a stringent set of Pod Security Policies are present
    So that a policy of least privilege can be enforced in order to prevent malicious attacks on my organization

    Background:
        Given a Kubernetes cluster exists which we can deploy into

    @k-psp-001
    Scenario: Prevent a deployment from running with privileged access

        Pods that request Privileged mode (using the security context of the container spec)
        will get operating system administrative capabilities - almost the same privileges that are
        accessible outside of a container.
        
        Security Standard References:
            - https://kubernetes.io/docs/concepts/policy/pod-security-policy/#privileged
            - CIS Kubernetes Benchmark v1.6.0 - 5.2.5

        Then pod creation "succeeds" with "allowPrivilegeEscalation" set to "false" in the pod spec
        And pod creation "fails" with "allowPrivilegeEscalation" set to "true" in the pod spec

    @k-psp-002
    Scenario Outline: Prevent execution of commands that require privileged access

        By default Pods that don't specify whether Privileged mode is set within the security context
        of the container spec should not have the ability to perform privileged commands.

        Security Standard References:
            - https://kubernetes.io/docs/concepts/policy/pod-security-policy/#privileged
            - CIS Kubernetes Benchmark v1.6.0 - 5.2.5

        When pod creation "succeeds" with "allowPrivilegeEscalation" set to "<VALUE>" in the pod spec
        Then the execution of a "non-privileged" command inside the Pod is "successful"
        But the execution of a "privileged" command inside the Pod is "rejected"

        Examples:
            | VALUE                     |
            | not have a value provided |
            | false                     |

    @k-psp-003
    Scenario Outline: Prevent a deployment from running in the host's process tree namespace

        HostPID controls whether a Pod's containers can share the host process ID namespace.
        If paired with ptrace, this can be used to escalate privileges outside of the container.

        Security Standard References:
            - https://kubernetes.io/docs/concepts/policy/pod-security-policy/#host-namespaces
            - CIS Kubernetes Benchmark v1.6.0 - 5.2.2

        When pod creation "succeeds" with "hostPID" set to "false" in the pod spec
        Then pod creation "fails" with "hostPID" set to "true" in the pod spec

    @k-psp-004
    Scenario: Prevent execution of commands that allow privileged access by default

        By default Pods that don't specify a value for hostPID should not have the ability to
        gain access to processes outside of the Pod's process tree.

        Security Standard References:
            - https://kubernetes.io/docs/concepts/policy/pod-security-policy/#host-namespaces
            - CIS Kubernetes Benchmark v1.6.0 - 5.2.2

        When pod creation "succeeds" with "hostPID" set to "<VALUE>" in the pod spec
        And the execution of a "non-privileged" command inside the Pod is "successful"
        Then a "process" inspection should only show the container processes

        Examples:
            | VALUE                     |
            | not have a value provided |
            | false                     |

    @k-psp-005
    Scenario: Prevent a deployment from running with access to the shared host IPC namespace

        HostIPC controls whether a Pod's containers can share the host IPC namespace, 
        allowing container processes to communicate with other processes on the host.

        Security Standard References:
            - https://kubernetes.io/docs/concepts/policy/pod-security-policy/#host-namespaces
            - CIS Kubernetes Benchmark v1.6.0 - 5.2.3

        When pod creation "succeeds" with "hostIPC" set to "false" in the pod spec
        Then pod creation "fails" with "hostIPC" set to "true" in the pod spec

    @k-psp-006
    Scenario: Prevent a deployment from running with access to the shared host IPC namespace

        By default Pods that don't specify whether Host IPC namespace mode is set should not
        be able to access the shared host IPC namespace.

        Security Standard References:
            - https://kubernetes.io/docs/concepts/policy/pod-security-policy/#host-namespaces
            - CIS Kubernetes Benchmark v1.6.0 - 5.2.3

        When pod creation "succeeds" with "hostIPC" set to "<VALUE>" in the pod spec
        And the execution of a "non-privileged" command inside the Pod is "successful"
        Then a "namespace" inspection should only show the container processes

        Examples:
            | VALUE                     |
            | not have a value provided |
            | false                     |
