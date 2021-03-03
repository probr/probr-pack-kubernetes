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
