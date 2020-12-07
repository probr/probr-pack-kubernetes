@probes/kubernetes
@probes/kubernetes/internet_access
@category/internet_access
@standard/citihub
@standard/citihub/CHC2-SVD010
@csp/any
Feature: Egress control of a kubernetes cluster
    As a Security Auditor
    I want to ensure that containers running inside Kubernetes clusters cannot directly access the Internet
    So that Internet traffic can be inspected and controlled

    @probes/kubernetes/internet_access/1.0 @control_type/preventative 
    Scenario Outline: Test outgoing connectivity of a deployed pod
        Given a Kubernetes cluster is deployed
        And a pod is deployed in the cluster
        When a process inside the pod establishes a direct http(s) connection to "<url>"
        Then access is "<result>"

        Examples:
            | url                             | result    |
            | www.google.com        | blocked |
            | www.microsoft.com | blocked |
            | www.ubuntu.com        | blocked |
