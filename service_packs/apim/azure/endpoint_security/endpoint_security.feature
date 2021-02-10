@aapim-es
@probes/azure_apim/endpoint_security
Feature: Ensure APIM policies are in place
    As as an API developer
    I want to ensure that APIM policies are being enforced on a particular API
    So that web traffic to the API is properly secured

    @aapim-es-001
    Scenario Outline: Ensure endpoints deployed to APIM
       Given an API that is deployed to APIM
       When all endpoints are retrieved from APIM
       Then each endpoint has mTLS enabled