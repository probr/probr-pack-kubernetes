    @k-cra-001
     Scenario: Ensure the cluster service account has read only access to the authorized container registry 
       When I attempt to push to the container registry using the cluster identity 
       Then the push request is rejected due to authorization 