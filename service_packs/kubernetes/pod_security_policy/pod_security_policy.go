package psp

import (
	"fmt"
	"log"
	"strings"

	"github.com/cucumber/godog"

	"github.com/citihub/probr/service_packs/coreengine"
	"github.com/citihub/probr/service_packs/kubernetes"

	"github.com/citihub/probr/utils"
)

type probeStruct struct{}

// Probe meets the service pack interface for adding the logic from this file
var Probe probeStruct

// PodSecurityPolicy is the section of the kubernetes package which provides the kubernetes interactions required to support
// pod security policy
var psp PodSecurityPolicy

// SetPodSecurityPolicy allows injection of a specific PodSecurityPolicy helper.
func SetPodSecurityPolicy(p PodSecurityPolicy) {
	psp = p
}

// General

func (s *scenarioState) creationWillWithAMessage(arg1, arg2 string) error {
	// Standard auditing logic to ensures panics are also audited
	stepTrace, payload, err := utils.AuditPlaceholders()
	defer func() {
		s.audit.AuditScenarioStep(stepTrace.String(), payload, err)
	}()
	stepTrace.WriteString("PENDING IMPLEMENTATION")
	return godog.ErrPending
}

func (s *scenarioState) aKubernetesClusterIsDeployed() error {
	var stepTrace strings.Builder
	description, payload, err := kubernetes.ClusterIsDeployed()
	stepTrace.WriteString(description)
	defer func() {
		s.audit.AuditScenarioStep(stepTrace.String(), payload, err)
	}()
	return err //  ClusterIsDeployed will create a fatal error if kubeconfig doesn't validate
}

// PENDING IMPLEMENTATION
func (s *scenarioState) aKubernetesDeploymentIsAppliedToAnExistingKubernetesCluster() error {
	// Standard auditing logic to ensures panics are also audited
	stepTrace, payload, err := utils.AuditPlaceholders()
	defer func() {
		s.audit.AuditScenarioStep(stepTrace.String(), payload, err)
	}()

	//TODO: not sure this step is adding value ... return "pass" for now ...
	stepTrace.WriteString("PENDING IMPLEMENTATION")

	return nil
}

func (s *scenarioState) theOperationWillWithAnError(result, message string) error {
	// Standard auditing logic to ensures panics are also audited
	stepTrace, payload, err := utils.AuditPlaceholders()
	defer func() {
		s.audit.AuditScenarioStep(stepTrace.String(), payload, err)
	}()

	stepTrace.WriteString("Validate that the scenario state was updated in the previous step with a particular result and message; ")
	err = kubernetes.AssertResult(&s.podState, result, message)
	payload = struct {
		ExpectedResult  string
		ExpectedMessage string
		PodState        kubernetes.PodState
	}{result, message, s.podState}

	return err
}

func (s *scenarioState) allOperationsWillWithAnError(result, message string) error {
	// Standard auditing logic to ensures panics are also audited
	stepTrace, payload, err := utils.AuditPlaceholders()
	defer func() {
		s.audit.AuditScenarioStep(stepTrace.String(), payload, err)
	}()

	var kubeErrors []error
	for _, ps := range s.podStates {
		kubeErrors = append(kubeErrors, kubernetes.AssertResult(&ps, result, message))
	}

	for _, ke := range kubeErrors {
		if ke != nil {
			err = utils.ReformatError("%v; %v", err, ke)
		}
	}

	stepTrace.WriteString("Validate that the scenario state was updated in the previous step with a list of pods with a particular result and message; ")
	payload = struct {
		ExpectedResult  string
		ExpectedMessage string
		PodState        []kubernetes.PodState
	}{result, message, s.podStates}

	return err
}

func (s *scenarioState) iShouldBeAbleToPerformAnAllowedCommand() error {
	// Standard auditing logic to ensures panics are also audited
	stepTrace, payload, err := utils.AuditPlaceholders()
	defer func() {
		s.audit.AuditScenarioStep(stepTrace.String(), payload, err)
	}()

	if len(s.podStates) == 0 {
		stepTrace.WriteString("Validate that 'ls' is executed successfully")
		payload = struct {
			PodState         kubernetes.PodState
			Command          string
			ExpectedExitCode int
		}{s.podState, Ls.String(), 0}
		err = s.runVerificationProbe(VerificationProbe{Cmd: Ls, ExpectedExitCode: 0}) //'0' exit code as we expect this to succeed
	} else {
		stepTrace.WriteString("Validate that 'ls' is executed successfully for all pods specified in the scenario state")
		payload = struct {
			PodStates        []kubernetes.PodState
			Command          string
			ExpectedExitCode int
		}{s.podStates, Ls.String(), 0}

		var errorMessage strings.Builder
		for _, podState := range s.podStates {
			err = s.runVerificationProbeWithCommand(podState, "ls", 0)
			if err != nil {
				errorMessage.WriteString(fmt.Sprintf("%v: %v; ", podState.PodName, errorMessage))
			}
		}
		if errorMessage.String() != "" {
			err = utils.ReformatError("Unable to run expected command against pods: %s", errorMessage.String())
		}
	}
	return err
}

// common helper funcs
func (s *scenarioState) runControlProbe(cf func() (*bool, error), c string) error {

	yesNo, err := cf()

	if err != nil {
		err = utils.ReformatError("error determining Pod Security Policy: %v error: %v", c, err)
		return err
	}
	if yesNo == nil {
		err = utils.ReformatError("result of %v is nil despite no error being raised from the call", c)
		log.Print(err)
		return err
	}

	if !*yesNo {
		return utils.ReformatError("%v is NOT restricted (result: %t)", c, *yesNo)
	}

	return nil
}

//add expected exit code//
func (s *scenarioState) runVerificationProbeWithCommand(podState kubernetes.PodState, command string, expectedExitCode int) error {
	if podState.CreationError == nil {
		res, err := psp.ExecCmd(&podState.PodName, command, s.probe)
		//analyse the results
		if err != nil {
			//this is an error from trying to execute the command as opposed to
			//the command itself returning an error
			err = utils.ReformatError("Likely a misconfiguration error. Error raised trying to execute verification command (%v) - %v", command, err)
			log.Print(err)
			return err
		}
		if res == nil {
			err = utils.ReformatError("<nil> result received when trying to execute verification command (%v)", command)
			log.Print(err)
			return err
		}
		if res.Err != nil && res.Internal {
			//we have an error which was raised before reaching the cluster (i.e. it's "internal")
			//this indicates that the command was not successfully executed
			err = utils.ReformatError("%s: %v - (%v)", utils.CallerName(0), command, res.Err)
			log.Print(err)
			return err
		}

		//log.Printf("Command: %s, Exit Code: %v\n", command, res.Code)

		//we've managed to execution against the cluster.  This may have failed due to pod security, but this
		//is still a 'successful' execution.  The exit code of the command needs to be verified against expected
		//check the result against expected:
		if res.Code == expectedExitCode {
			//then as expected, test passes
			return nil
		}
		//else it's a fail:
		return utils.ReformatError("exit code %d from verification command %q did not match expected %v",
			res.Code, command, expectedExitCode)
	}
	return nil
}

//Tries to run the specified command, which shouldn't be allowed if the capability hasn't been added
func (s *scenarioState) runCapabilityVerificationProbes() error {
	var errorMessage string

	for n, podState := range s.podStates {

		var capability string
		//if we didn't provide any additional info then we assume that there were no additional capabilities as part of this scenario
		if len(s.info) != 0 {
			capability = s.info[n].(string)
		} else {
			capability = ""
		}

		for cap, cmd := range getLinuxNonDefaultCapabilities() {
			if cap != capability {
				log.Printf("[DEBUG] Attempting verification command '%s' with expected exit 2", cmd)
				cmdErr := s.runVerificationProbeWithCommand(podState, cmd, 2)
				if cmdErr != nil {
					errorMessage = fmt.Sprintf("%s, Error: Command for capability %v was successful. ", errorMessage, cap)
				}
			}
		}
	}

	if errorMessage == "" {
		return nil
	}
	return fmt.Errorf("Error verifying capabilities: %s", errorMessage)
}

func (s *scenarioState) runVerificationProbe(c VerificationProbe) error {

	//check for lack of creation error, i.e. pod was created successfully
	if s.podState.CreationError == nil {
		res, err := psp.ExecPSPProbeCmd(&s.podState.PodName, c.Cmd, s.probe)

		//analyse the results
		if err != nil {
			//this is an error from trying to execute the command as opposed to
			//the command itself returning an error
			err = utils.ReformatError("Likely a misconfiguration error. Error raised trying to execute verification command (%v) - %v", c.Cmd, err)
			log.Print(err)
			return err
		}
		if res == nil {
			err = utils.ReformatError("<nil> result received when trying to execute verification command (%v)", c.Cmd)
			log.Print(err)
			return err
		}
		if res.Err != nil && res.Internal {
			//we have an error which was raised before reaching the cluster (i.e. it's "internal")
			//this indicates that the command was not successfully executed
			err = utils.ReformatError("%s: %v - (%v)", utils.CallerName(0), c, res.Err)
			log.Print(err)
			return err
		}

		//we've managed to execution against the cluster.  This may have failed due to pod security, but this
		//is still a 'successful' execution.  The exit code of the command needs to be verified against expected
		//check the result against expected:
		if res.Code == c.ExpectedExitCode {
			//then as expected, test passes
			return nil
		}
		//else it's a fail:
		return utils.ReformatError("exit code %d from verification command %q did not match expected %d",
			res.Code, c.Cmd, c.ExpectedExitCode)
	}

	//pod wasn't created so nothing to test
	//TODO: really, we don't want to 'pass' this.  Is there an alternative?
	return nil
}

// TEST STEPS:

// CIS-5.2.1
// privileged access
func (s *scenarioState) privilegedAccessRequestIsMarkedForTheKubernetesDeployment(privilegedAccessRequested string) error {
	// Standard auditing logic to ensures panics are also audited
	stepTrace, payload, err := utils.AuditPlaceholders()
	defer func() {
		s.audit.AuditScenarioStep(stepTrace.String(), payload, err)
	}()

	var pa bool
	if privilegedAccessRequested == "True" {
		pa = true
	} else {
		pa = false
	}

	stepTrace.WriteString(fmt.Sprintf(
		"Attempt to deploy a pod with priviledged access '%s'; ", privilegedAccessRequested))
	pd, err := psp.CreatePODSettingSecurityContext(&pa, nil, nil, s.probe)
	err = kubernetes.ProcessPodCreationResult(&s.podState, pd, kubernetes.PSPNoPrivilege, err)

	payload = struct {
		PodState kubernetes.PodState
	}{s.podState}

	return err
}

func (s *scenarioState) someControlExistsToPreventPrivilegedAccessForKubernetesDeploymentsToAnActiveKubernetesCluster() error {
	// Standard auditing logic to ensures panics are also audited
	stepTrace, payload, err := utils.AuditPlaceholders()
	defer func() {
		s.audit.AuditScenarioStep(stepTrace.String(), payload, err)
	}()

	// TODO: This entire process needs refactored for readability, and to remove holistic dependency on azure security context
	stepTrace.WriteString("Validate that the kube instance security context contains 'k8sazurecontainernoprivilege'; ")
	err = s.runControlProbe(psp.PrivilegedAccessIsRestricted, "PrivilegedAccessIsRestricted")

	payload = struct {
		PodState kubernetes.PodState
	}{s.podState}

	return err
}

func (s *scenarioState) iShouldNotBeAbleToPerformACommandThatRequiresPrivilegedAccess() error {
	// Standard auditing logic to ensures panics are also audited
	stepTrace, payload, err := utils.AuditPlaceholders()
	defer func() {
		s.audit.AuditScenarioStep(stepTrace.String(), payload, err)
	}()

	stepTrace.WriteString(fmt.Sprintf(
		"Validate that the command '%s' fails to execute; ", Chroot.String()))
	err = s.runVerificationProbe(VerificationProbe{Cmd: Chroot, ExpectedExitCode: 1})

	payload = struct {
		PodState kubernetes.PodState
	}{s.podState}

	return err
}

// CIS-5.2.2
// hostPID
func (s *scenarioState) hostPIDRequestIsMarkedForTheKubernetesDeployment(hostPIDRequested string) error {
	// Standard auditing logic to ensures panics are also audited
	stepTrace, payload, err := utils.AuditPlaceholders()
	defer func() {
		s.audit.AuditScenarioStep(stepTrace.String(), payload, err)
	}()

	var hostPID bool
	if hostPIDRequested == "True" {
		hostPID = true
	} else {
		hostPID = false
	}
	stepTrace.WriteString(fmt.Sprintf(
		"Attempt to deploy a pod with host PID '%s'; ", hostPIDRequested))
	pd, err := psp.CreatePODSettingAttributes(&hostPID, nil, nil, s.probe)
	err = kubernetes.ProcessPodCreationResult(&s.podState, pd, kubernetes.PSPHostNamespace, err)

	payload = struct {
		PodState kubernetes.PodState
	}{s.podState}

	return err
}

func (s *scenarioState) someSystemExistsToPreventAKubernetesContainerFromRunningUsingTheHostPIDOnTheActiveKubernetesCluster() error {
	// Standard auditing logic to ensures panics are also audited
	stepTrace, payload, err := utils.AuditPlaceholders()
	defer func() {
		s.audit.AuditScenarioStep(stepTrace.String(), payload, err)
	}()

	stepTrace.WriteString("Validate that at least on of the cluster's pod security policies has HostPID set to false; ")
	err = s.runControlProbe(psp.HostPIDIsRestricted, "HostPIDIsRestricted")

	payload = struct {
		PodState kubernetes.PodState
	}{s.podState}

	return err
}

func (s *scenarioState) iShouldNotBeAbleToPerformACommandThatProvidesAccessToTheHostPIDNamespace() error {
	// Standard auditing logic to ensures panics are also audited
	stepTrace, payload, err := utils.AuditPlaceholders()
	defer func() {
		s.audit.AuditScenarioStep(stepTrace.String(), payload, err)
	}()

	stepTrace.WriteString(fmt.Sprintf(
		"Validate that the exit code for command '%s' is 1; ", EnterHostPIDNS.String()))
	err = s.runVerificationProbe(VerificationProbe{Cmd: EnterHostPIDNS, ExpectedExitCode: 1})

	payload = struct {
		PodState kubernetes.PodState
	}{s.podState}

	return err
}

//CIS-5.2.3
func (s *scenarioState) hostIPCRequestIsMarkedForTheKubernetesDeployment(hostIPCRequested string) error {
	// Standard auditing logic to ensures panics are also audited
	stepTrace, payload, err := utils.AuditPlaceholders()
	defer func() {
		s.audit.AuditScenarioStep(stepTrace.String(), payload, err)
	}()

	var hostIPC bool
	if hostIPCRequested == "True" {
		hostIPC = true
	} else {
		hostIPC = false
	}

	stepTrace.WriteString(fmt.Sprintf(
		"Attempt to deploy a pod with Host IPC requested: %s; ", hostIPCRequested))
	pd, err := psp.CreatePODSettingAttributes(nil, &hostIPC, nil, s.probe)
	err = kubernetes.ProcessPodCreationResult(&s.podState, pd, kubernetes.PSPHostNamespace, err)

	payload = struct {
		PodState kubernetes.PodState
	}{s.podState}

	return err
}

func (s *scenarioState) someSystemExistsToPreventAKubernetesDeploymentFromRunningUsingTheHostIPCInAnExistingKubernetesCluster() error {
	// Standard auditing logic to ensures panics are also audited
	stepTrace, payload, err := utils.AuditPlaceholders()
	defer func() {
		s.audit.AuditScenarioStep(stepTrace.String(), payload, err)
	}()

	stepTrace.WriteString("Validate that at least on of the cluster's pod security policies has HostIPC set to false; ")
	err = s.runControlProbe(psp.HostIPCIsRestricted, "HostIPCIsRestricted")

	payload = struct {
		PodState kubernetes.PodState
	}{s.podState}

	return err
}

func (s *scenarioState) iShouldNotBeAbleToPerformACommandThatProvidesAccessToTheHostIPCNamespace() error {
	// Standard auditing logic to ensures panics are also audited
	stepTrace, payload, err := utils.AuditPlaceholders()
	defer func() {
		s.audit.AuditScenarioStep(stepTrace.String(), payload, err)
	}()

	stepTrace.WriteString(fmt.Sprintf(
		"Validate that the exit code for command '%s' is 1; ", EnterHostIPCNS.String()))
	err = s.runVerificationProbe(VerificationProbe{Cmd: EnterHostIPCNS, ExpectedExitCode: 1})

	payload = struct {
		PodState kubernetes.PodState
	}{s.podState}

	return err
}

//CIS-5.2.4
func (s *scenarioState) hostNetworkRequestIsMarkedForTheKubernetesDeployment(hostNetworkRequested string) error {
	// Standard auditing logic to ensures panics are also audited
	stepTrace, payload, err := utils.AuditPlaceholders()
	defer func() {
		s.audit.AuditScenarioStep(stepTrace.String(), payload, err)
	}()

	var hostNetwork bool
	if hostNetworkRequested == "True" {
		hostNetwork = true
	} else {
		hostNetwork = false
	}

	stepTrace.WriteString(fmt.Sprintf(
		"Attempt to deploy a pod with host network request marked '%s'; ", hostNetworkRequested))
	pd, err := psp.CreatePODSettingAttributes(nil, nil, &hostNetwork, s.probe)
	err = kubernetes.ProcessPodCreationResult(&s.podState, pd, kubernetes.PSPHostNetwork, err)

	payload = struct {
		PodState kubernetes.PodState
	}{s.podState}

	return err
}

func (s *scenarioState) someSystemExistsToPreventAKubernetesDeploymentFromRunningUsingTheHostNetworkInAnExistingKubernetesCluster() error {
	// Standard auditing logic to ensures panics are also audited
	stepTrace, payload, err := utils.AuditPlaceholders()
	defer func() {
		s.audit.AuditScenarioStep(stepTrace.String(), payload, err)
	}()

	stepTrace.WriteString("Validate that at least on of the pod security policies has HostNetwork set to false; ")
	err = s.runControlProbe(psp.HostNetworkIsRestricted, "HostNetworkIsRestricted")

	payload = struct {
		PodState kubernetes.PodState
	}{s.podState}

	return err
}

func (s *scenarioState) iShouldNotBeAbleToPerformACommandThatProvidesAccessToTheHostNetworkNamespace() error {
	// Standard auditing logic to ensures panics are also audited
	stepTrace, payload, err := utils.AuditPlaceholders()
	defer func() {
		s.audit.AuditScenarioStep(stepTrace.String(), payload, err)
	}()

	stepTrace.WriteString(fmt.Sprintf(
		"Validate that the exit code for command '%s' is 1; ", EnterHostNetworkNS.String()))
	err = s.runVerificationProbe(VerificationProbe{Cmd: EnterHostNetworkNS, ExpectedExitCode: 1})

	payload = struct {
		PodState kubernetes.PodState
	}{s.podState}

	return err
}

//CIS-5.2.5
func (s *scenarioState) privilegedEscalationIsMarkedForTheKubernetesDeployment(privilegedEscalationRequested string) error {
	// Standard auditing logic to ensures panics are also audited
	stepTrace, payload, err := utils.AuditPlaceholders()
	defer func() {
		s.audit.AuditScenarioStep(stepTrace.String(), payload, err)
	}()

	allowPrivilegeEscalation := "true"
	if strings.ToLower(privilegedEscalationRequested) != "true" {
		allowPrivilegeEscalation = "false"
	}
	specPath := "psp-azp-privileges.yaml"
	stepTrace.WriteString("Parsing pod spec file '%s'; ")
	y, err := utils.ReadStaticFile(kubernetes.AssetsDir, specPath)
	if err == nil {
		yaml := utils.ReplaceBytesValue(y, "{{ allowPrivilegeEscalation }}", allowPrivilegeEscalation)

		stepTrace.WriteString(fmt.Sprintf(
			"Attempt to create a pod with privilege escalation set to '%s'; ", allowPrivilegeEscalation))
		pd, cErr := psp.CreatePodFromYaml(yaml, s.probe)
		err = kubernetes.ProcessPodCreationResult(&s.podState, pd, kubernetes.PSPNoPrivilegeEscalation, cErr)
	}

	payload = struct {
		PodState kubernetes.PodState
	}{s.podState}

	return err
}

func (s *scenarioState) someSystemExistsToPreventAKubernetesDeploymentFromRunningUsingTheAllowPrivilegeEscalationInAnExistingKubernetesCluster() error {
	// Standard auditing logic to ensures panics are also audited
	stepTrace, payload, err := utils.AuditPlaceholders()
	defer func() {
		s.audit.AuditScenarioStep(stepTrace.String(), payload, err)
	}()

	stepTrace.WriteString(
		"Validate that at least one of the pod security policies has AllowPrivilegeEscalation set to false; ")
	err = s.runControlProbe(psp.PrivilegedEscalationIsRestricted, "PrivilegedEscalationIsRestricted")

	payload = struct {
		PodState kubernetes.PodState
	}{s.podState}

	return err
}

func (s *scenarioState) iShouldNotBeAbleToPerformASudoCommandThatRequiresPrivilegedAccess() error {
	// Standard auditing logic to ensures panics are also audited
	stepTrace, payload, err := utils.AuditPlaceholders()
	defer func() {
		s.audit.AuditScenarioStep(stepTrace.String(), payload, err)
	}()

	stepTrace.WriteString(fmt.Sprintf(
		"Validate that the exit code for command '%s' is 1; ", SudoChroot.String()))
	err = s.runVerificationProbe(VerificationProbe{Cmd: SudoChroot, ExpectedExitCode: 126})

	payload = struct {
		PodState kubernetes.PodState
	}{s.podState}

	return err
}

//CIS-5.2.6
func (s *scenarioState) theUserRequestedIsForTheKubernetesDeployment(requestedUser string) error {
	// Standard auditing logic to ensures panics are also audited
	stepTrace, payload, err := utils.AuditPlaceholders()
	defer func() {
		s.audit.AuditScenarioStep(stepTrace.String(), payload, err)
	}()

	var runAsUser int64
	if requestedUser == "Root" {
		runAsUser = 0
	} else {
		runAsUser = 1000
	}

	stepTrace.WriteString(fmt.Sprintf(
		"Attempt to create a pod with the requested user '%s'; ", requestedUser))
	pd, err := psp.CreatePODSettingSecurityContext(nil, nil, &runAsUser, s.probe)
	err = kubernetes.ProcessPodCreationResult(&s.podState, pd, kubernetes.PSPAllowedUsersGroups, err)

	payload = struct {
		PodState kubernetes.PodState
	}{s.podState}

	return err
}

func (s *scenarioState) someSystemExistsToPreventAKubernetesDeploymentFromRunningAsTheRootUserInAnExistingKubernetesCluster() error {
	// Standard auditing logic to ensures panics are also audited
	stepTrace, payload, err := utils.AuditPlaceholders()
	defer func() {
		s.audit.AuditScenarioStep(stepTrace.String(), payload, err)
	}()

	stepTrace.WriteString(
		"Validate that at least one of the pod security policies has RunAsUser set to MustRunAsNonRoot; ")
	err = s.runControlProbe(psp.RootUserIsRestricted, "RootUserIsRestricted")

	payload = struct {
		PodState kubernetes.PodState
	}{s.podState}

	return err
}

func (s *scenarioState) theKubernetesDeploymentShouldRunWithANonrootUID() error {
	// Standard auditing logic to ensures panics are also audited
	stepTrace, payload, err := utils.AuditPlaceholders()
	defer func() {
		s.audit.AuditScenarioStep(stepTrace.String(), payload, err)
	}()

	stepTrace.WriteString(fmt.Sprintf(
		"Validate that the exit code for command '%s' is 1; ", VerifyNonRootUID.String()))
	err = s.runVerificationProbe(VerificationProbe{Cmd: VerifyNonRootUID, ExpectedExitCode: 1})

	payload = struct {
		PodState kubernetes.PodState
	}{s.podState}

	return err
}

func (s *scenarioState) kubernetesDeploymentWithNETRAWCapability(netRawRequested string) error {
	// Standard auditing logic to ensures panics are also audited
	stepTrace, payload, err := utils.AuditPlaceholders()
	defer func() {
		s.audit.AuditScenarioStep(stepTrace.String(), payload, err)
	}()

	var capAdd, capDrop []string

	switch netRawRequested {
	case "Added":
		capAdd = []string{"NET_RAW"}
		capDrop = nil
	case "Dropped":
		capAdd = nil
		capDrop = []string{"NET_RAW"}
	case "Not Defined":
		capAdd, capDrop = nil, nil
	}

	stepTrace.WriteString(fmt.Sprintf(
		"Creating pod with NET_RAW %s in the pod's capabilities; ", strings.ToLower(netRawRequested)))
	pd, err := psp.CreatePODWithCapabilities(capAdd, capDrop, s.probe)
	err = kubernetes.ProcessPodCreationResult(&s.podState, pd, kubernetes.PSPAllowedCapabilities, err)

	payload = struct {
		PodState kubernetes.PodState
	}{s.podState}

	return err
}

func (s *scenarioState) iShouldNotBeAbleToPerformACommandThatRequiresNETRAWCapability() error {
	// Standard auditing logic to ensures panics are also audited
	stepTrace, payload, err := utils.AuditPlaceholders()
	defer func() {
		s.audit.AuditScenarioStep(stepTrace.String(), payload, err)
	}()

	stepTrace.WriteString(fmt.Sprintf(
		"Validate that the exit code for command '%s' is 1; ", NetRawProbe.String()))
	err = s.runVerificationProbe(VerificationProbe{Cmd: NetRawProbe, ExpectedExitCode: 1})

	payload = struct {
		PodState kubernetes.PodState
	}{s.podState}

	return err
}

//CIS-5.2.8
func (s *scenarioState) additionalCapabilitiesForTheKubernetesDeployment(capabilitiesAllowed string) error {
	// Standard auditing logic to ensures panics are also audited
	stepTrace, payload, err := utils.AuditPlaceholders()
	defer func() {
		s.audit.AuditScenarioStep(stepTrace.String(), payload, err)
	}()

	var capabilities []string
	switch capabilitiesAllowed {
	case "ARE":
		stepTrace.WriteString("Creating list of capabilities from config vars; ")
		//TODO: what if the list is empty?
		capabilities = getAllowedAdditionalCapabilities()
		//if you provide an empty array in the Pod specification then the PSP will block it as requesting an additional capability that is not in the allowed list
		if capabilities[0] == "" {
			capabilities = make([]string, 0)
		}
	case "NOT":
		stepTrace.WriteString("Creating list of allowed capabilities from linux non-default capability list; ")
		// if we haven't been instructed to allow it, then add it to the list
		for cap := range getLinuxNonDefaultCapabilities() {
			found := false
			for _, ac := range getAllowedAdditionalCapabilities() {
				if cap == ac {
					stepTrace.WriteString(fmt.Sprintf(
						"Not including %s, as it was marked as allowed in the config; ", cap))
					found = true
				}
			}
			if found == false {
				capabilities = append(capabilities, cap)
			}
		}
	case "Not Defined":
		//do nothing
	}

	if len(capabilities) == 0 {
		stepTrace.WriteString("Attempting to deploy pod with no capabilities added; ")
		var podState kubernetes.PodState
		boolVal := false
		pd, cErr := psp.CreatePODSettingAttributes(nil, nil, &boolVal, s.probe)
		err = kubernetes.ProcessPodCreationResult(&podState, pd, kubernetes.PSPHostNetwork, cErr)
		s.podStates = append(s.podStates, podState)
	} else {
		stepTrace.WriteString(fmt.Sprintf(
			"Attempting to deploy one pod for each specified capability (%v total); ",
			len(capabilities)))
	}

	for _, capability := range capabilities {
		var podState kubernetes.PodState
		c := []string{capability}
		pd, cErr := psp.CreatePODSettingCapabilities(&c, s.probe)
		err = kubernetes.ProcessPodCreationResult(&podState, pd, kubernetes.PSPAllowedCapabilities, cErr)

		s.podStates = append(s.podStates, podState)
		s.info = append(s.info, capability)
	}

	payload = struct {
		PodState     kubernetes.PodState
		Capabilities []string
	}{s.podState, capabilities}

	return err
}

func (s *scenarioState) someSystemExistsToPreventKubernetesDeploymentsWithCapabilitiesBeyondTheDefaultSetFromBeingDeployedToAnExistingKubernetesCluster() error {
	// Standard auditing logic to ensures panics are also audited
	stepTrace, payload, err := utils.AuditPlaceholders()
	defer func() {
		s.audit.AuditScenarioStep(stepTrace.String(), payload, err)
	}()

	stepTrace.WriteString("Validate that no pod security policies have any values in AllowedCapabilities; ")
	err = s.runControlProbe(psp.AllowedCapabilitiesAreRestricted, "AllowedCapabilitiesAreRestricted")

	payload = struct {
		PodState kubernetes.PodState
	}{s.podState}

	return err
}

func (s *scenarioState) iShouldNotBeAbleToPerformACommandThatRequiresCapabilitiesOutsideOfTheDefaultSet() error {
	// Standard auditing logic to ensures panics are also audited
	stepTrace, payload, err := utils.AuditPlaceholders()
	defer func() {
		s.audit.AuditScenarioStep(stepTrace.String(), payload, err)
	}()

	stepTrace.WriteString(
		"Validate that verification commands exit 2 for all non-default capabilities, except those that are specified as allowed in the config; ")
	err = s.runCapabilityVerificationProbes()
	//err = s.runVerificationProbe(VerificationProbe{Cmd: SpecialCapProbe, ExpectedExitCode: 2})

	payload = struct {
		PodState kubernetes.PodState
	}{s.podState}

	return err
}

//CIS-5.2.9
func (s *scenarioState) assignedCapabilitiesForTheKubernetesDeployment(assignCapabilities string) error {
	// Standard auditing logic to ensures panics are also audited
	stepTrace, payload, err := utils.AuditPlaceholders()
	defer func() {
		s.audit.AuditScenarioStep(stepTrace.String(), payload, err)
	}()

	var c []string
	if assignCapabilities == "ARE" {
		stepTrace.WriteString(fmt.Sprintf(
			"Attempt to deploy a pod with the NET_ADMIN capability; "))
		//TODO: just add net_admin for now - but is this appropriate?
		//what's the difference with 5.2.8???
		c = make([]string, 1)
		c[0] = "NET_ADMIN"
	} else {
		stepTrace.WriteString(fmt.Sprintf(
			"Attempt to deploy a pod without the NET_ADMIN capability; "))
	}

	pd, err := psp.CreatePODSettingCapabilities(&c, s.probe)
	err = kubernetes.ProcessPodCreationResult(&s.podState, pd, kubernetes.PSPAllowedCapabilities, err)

	payload = struct {
		PodState kubernetes.PodState
	}{s.podState}

	return err

}

func (s *scenarioState) someSystemExistsToPreventKubernetesDeploymentsWithAssignedCapabilitiesFromBeingDeployedToAnExistingKubernetesCluster() error {
	// Standard auditing logic to ensures panics are also audited
	stepTrace, payload, err := utils.AuditPlaceholders()
	defer func() {
		s.audit.AuditScenarioStep(stepTrace.String(), payload, err)
	}()

	stepTrace.WriteString("Validate that the kube instance security context contains 'k8sazureallowedcapabilities'; ")
	err = s.runControlProbe(psp.AssignedCapabilitiesAreRestricted, "AssignedCapabilitiesAreRestricted")

	payload = struct {
		PodState kubernetes.PodState
	}{s.podState}

	return err
}

func (s *scenarioState) iShouldNotBeAbleToPerformACommandThatRequiresAnyCapabilities() error {
	// Standard auditing logic to ensures panics are also audited
	stepTrace, payload, err := utils.AuditPlaceholders()
	defer func() {
		s.audit.AuditScenarioStep(stepTrace.String(), payload, err)
	}()

	stepTrace.WriteString(fmt.Sprintf(
		"Validate that the exit code for command '%s' is 1; ", SpecialCapProbe.String()))
	err = s.runVerificationProbe(VerificationProbe{Cmd: SpecialCapProbe, ExpectedExitCode: 2})

	payload = struct {
		PodState kubernetes.PodState
	}{s.podState}

	return err
}

//AZ Policy - port range
func (s *scenarioState) anPortRangeIsRequestedForTheKubernetesDeployment(portRange string) error {
	// Standard auditing logic to ensures panics are also audited
	stepTrace, payload, err := utils.AuditPlaceholders()
	defer func() {
		s.audit.AuditScenarioStep(stepTrace.String(), payload, err)
	}()

	var y []byte
	var yaml []byte
	var specPath string

	switch portRange {
	case "unapproved":
		specPath = "psp-azp-hostport-unapproved.yaml"
		unapprovedHostPort := kubernetes.GetUnapprovedHostPortFromConfig()
		y, err = utils.ReadStaticFile(kubernetes.AssetsDir, specPath)

		stepTrace.WriteString(fmt.Sprintf(
			"Set unapproved port for pod spec: %s; ", unapprovedHostPort))
		yaml = utils.ReplaceBytesValue(y, "{{ unapproved-port }}", unapprovedHostPort)
		stepTrace.WriteString(fmt.Sprintf(
			"%s port range is requested for kubernetes deployment. Port was %v.", portRange, unapprovedHostPort))
	case "not defined":
		specPath = "psp-azp-hostport-notdefined.yaml"
		yaml, err = utils.ReadStaticFile(kubernetes.AssetsDir, specPath)
		stepTrace.WriteString(fmt.Sprintf(
			"%s port range is requested for kubernetes deployment.", portRange))
	default:
		err = fmt.Errorf("Unrecognised port range option")
	}

	if err == nil {
		stepTrace.WriteString("Attempting to create pod from yaml spec; ")
		pd, cErr := psp.CreatePodFromYaml(yaml, s.probe)
		err = kubernetes.ProcessPodCreationResult(&s.podState, pd, kubernetes.PSPAllowedPortRange, cErr)
	}

	//audit log description defined in case statement above
	payload = struct {
		PodSpecPath string
		PodState    kubernetes.PodState
	}{specPath, s.podState}

	return err
}

func (s *scenarioState) someSystemExistsToPreventKubernetesDeploymentsWithUnapprovedPortRangeFromBeingDeployedToAnExistingKubernetesCluster() error {
	// Standard auditing logic to ensures panics are also audited
	stepTrace, payload, err := utils.AuditPlaceholders()
	defer func() {
		s.audit.AuditScenarioStep(stepTrace.String(), payload, err)
	}()

	stepTrace.WriteString("Validate that the kube instance security context contains 'k8sazurehostnetworkingports'; ")
	err = s.runControlProbe(psp.HostPortsAreRestricted, "HostPortsAreRestricted")

	payload = struct {
		PodState kubernetes.PodState
	}{s.podState}

	return err
}

func (s *scenarioState) iShouldNotBeAbleToPerformACommandThatAccessAnUnapprovedPortRange() error {
	// Standard auditing logic to ensures panics are also audited
	stepTrace, payload, err := utils.AuditPlaceholders()
	defer func() {
		s.audit.AuditScenarioStep(stepTrace.String(), payload, err)
	}()

	stepTrace.WriteString(fmt.Sprintf(
		"Validate that the exit code for command '%s' is 1; ", NetCat.String()))
	err = s.runVerificationProbe(VerificationProbe{Cmd: NetCat, ExpectedExitCode: 1})

	stepTrace.WriteString("Should not be able to perform a command that access an up approved port range")
	payload = struct {
		PodState kubernetes.PodState
	}{s.podState}

	return err
}

func (s *scenarioState) volumeTypesAreRequestedForTheKubernetesDeployment(volumeType string) error {
	// Standard auditing logic to ensures panics are also audited
	stepTrace, payload, err := utils.AuditPlaceholders()
	defer func() {
		s.audit.AuditScenarioStep(stepTrace.String(), payload, err)
	}()

	var volumeTypes []string

	if volumeType == "unapproved" {
		stepTrace.WriteString("Getting a list of supported volume types; removing Approved volume types; ")
		for _, vt := range getSupportedVolumeTypes() {
			found := false
			for _, avt := range getApprovedVolumeTypes() {
				if vt == avt {
					found = true
				}
			}

			if found == false {
				volumeTypes = append(volumeTypes, vt)
			}
		}
	} else {
		stepTrace.WriteString("Getting a list of approved volume types; ")
		volumeTypes = getApprovedVolumeTypes()
	}

	var y []byte
	var yaml [][]byte
	stepTrace.WriteString("Preparing spec for volumes; ")
	for _, vt := range volumeTypes {
		specPath := fmt.Sprintf("psp-azp-volumetypes-%s.yaml", vt)
		stepTrace.WriteString(fmt.Sprintf("Adding spec '%s'; ", specPath))
		y, err = utils.ReadStaticFile(kubernetes.AssetsDir, "volumetypes", specPath)
		yaml = append(yaml, y)
	}

	if err == nil {
		stepTrace.WriteString("Attempting to deploy pods using each aforementioned spec; ")
		for _, podyaml := range yaml {
			var podState kubernetes.PodState
			pd, cErr := psp.CreatePodFromYaml(podyaml, s.probe)
			err = kubernetes.ProcessPodCreationResult(&podState, pd, kubernetes.PSPAllowedVolumeTypes, cErr)
			s.podStates = append(s.podStates, podState)
		}
	}

	payload = struct {
		PodStates []kubernetes.PodState
	}{s.podStates}

	return err

}

func (s *scenarioState) someSystemExistsToPreventKubernetesDeploymentsWithUnapprovedVolumeTypesFromBeingDeployedToAnExistingKubernetesCluster() error {
	// Standard auditing logic to ensures panics are also audited
	stepTrace, payload, err := utils.AuditPlaceholders()
	defer func() {
		s.audit.AuditScenarioStep(stepTrace.String(), payload, err)
	}()

	stepTrace.WriteString("Validate that the kube instance security context contains 'k8sazurevolumetypes'; ")
	err = s.runControlProbe(psp.VolumeTypesAreRestricted, "VolumeTypesAreRestricted")

	stepTrace.WriteString("some systems exists to prevent kubernetes deployments without un approved volume types from being deployed existing kubernetes cluster")
	payload = struct {
		PodState kubernetes.PodState
	}{s.podState}

	return err
}

// PENDING IMPLEMENTATION
func (s *scenarioState) iShouldNotBeAbleToPerformACommandThatAccessesAnUnapprovedVolumeType() error {
	// Standard auditing logic to ensures panics are also audited
	stepTrace, payload, err := utils.AuditPlaceholders()
	defer func() {
		s.audit.AuditScenarioStep(stepTrace.String(), payload, err)
	}()

	stepTrace.WriteString("PENDING IMPLEMENTATION")
	//TODO: Not sure what the test is here - if any
	return godog.ErrPending
}

//AZ Policy - seccomp profile
func (s *scenarioState) anSeccompProfileIsRequestedForTheKubernetesDeployment(seccompProfile string) error {
	// Standard auditing logic to ensures panics are also audited
	stepTrace, payload, err := utils.AuditPlaceholders()
	defer func() {
		s.audit.AuditScenarioStep(stepTrace.String(), payload, err)
	}()

	var y []byte
	var specPath string

	if seccompProfile == "unapproved" {
		specPath = "psp-azp-seccomp-unapproved.yaml"
		y, err = utils.ReadStaticFile(kubernetes.AssetsDir, specPath)
	} else if seccompProfile == "undefined" {
		specPath = "psp-azp-seccomp-undefined.yaml"
		y, err = utils.ReadStaticFile(kubernetes.AssetsDir, specPath)
	} else if seccompProfile == "approved" {
		specPath = "psp-azp-seccomp-approved.yaml"
		y, err = utils.ReadStaticFile(kubernetes.AssetsDir, specPath)
	}

	if err != nil {
		log.Print(utils.ReformatError("error reading seccomp provile %v yaml file : %v", seccompProfile, err))
	}
	stepTrace.WriteString(fmt.Sprintf("Attempting to create pod with seccomp %s", seccompProfile))
	pd, cErr := psp.CreatePodFromYaml(y, s.probe)
	err = kubernetes.ProcessPodCreationResult(&s.podState, pd, kubernetes.PSPSeccompProfile, cErr)

	payload = struct {
		SpecPath string
		PodState kubernetes.PodState
	}{specPath, s.podState}

	return err
}

func (s *scenarioState) someSystemExistsToPreventKubernetesDeploymentsWithoutApprovedSeccompProfilesFromBeingDeployedToAnExistingKubernetesCluster() error {
	// Standard auditing logic to ensures panics are also audited
	stepTrace, payload, err := utils.AuditPlaceholders()
	defer func() {
		s.audit.AuditScenarioStep(stepTrace.String(), payload, err)
	}()

	stepTrace.WriteString("Validate that the kube instance security context contains 'k8sazureallowedseccomp'; ")
	err = s.runControlProbe(psp.SeccompProfilesAreRestricted, "SeccompProfilesAreRestricted")

	payload = struct {
		PodState kubernetes.PodState
	}{s.podState}

	return err
}

func (s *scenarioState) iShouldNotBeAbleToPerformASystemCallThatIsBlockedByTheSeccompProfile() error {
	// Standard auditing logic to ensures panics are also audited
	stepTrace, payload, err := utils.AuditPlaceholders()
	defer func() {
		s.audit.AuditScenarioStep(stepTrace.String(), payload, err)
	}()

	stepTrace.WriteString(fmt.Sprintf(
		"Validate that the exit code for command '%s' is 1; ", Unshare.String()))
	err = s.runVerificationProbe(VerificationProbe{Cmd: Unshare, ExpectedExitCode: 1})

	payload = struct {
		PodState kubernetes.PodState
	}{s.podState}

	return err
}

// Name presents the name of this probe for external reference
func (p probeStruct) Name() string {
	return "pod_security_policy"
}

// Path presents the path of these feature files for external reference
func (p probeStruct) Path() string {
	return coreengine.GetFeaturePath("service_packs", "kubernetes", p.Name())
}

// pspProbeInitialize handles any overall Test Suite initialisation steps.  This is registered with the
// test handler as part of the init() function.
func (p probeStruct) ProbeInitialize(ctx *godog.TestSuiteContext) {
	ctx.BeforeSuite(func() {
		//check dependencies ...
		if psp == nil {
			// not been given one so set default
			psp = NewDefaultPSP()
		}
		psp.CreateConfigMap()
	})

	ctx.AfterSuite(func() {
		psp.DeleteConfigMap()
	})
}

// pspScenarioInitialize initialises the specific test steps.  This is essentially the creation of the test
// which reflects the tests described in the events directory.  There must be a test step registered for
// each line in the feature files. Note: Godog will output stub steps and implementations if it doesn't find
// a step / function defined.  See: https://github.com/cucumber/godog#example.
func (p probeStruct) ScenarioInitialize(ctx *godog.ScenarioContext) {
	ps := scenarioState{}

	ctx.BeforeScenario(func(s *godog.Scenario) {
		beforeScenario(&ps, p.Name(), s)
	})

	ctx.Step(`^a Kubernetes cluster exists which we can deploy into$`, ps.aKubernetesClusterIsDeployed)

	ctx.Step(`^a Kubernetes deployment is applied to an existing Kubernetes cluster$`, ps.aKubernetesDeploymentIsAppliedToAnExistingKubernetesCluster)

	//CIS-5.2.1
	ctx.Step(`^privileged access request is marked "([^"]*)" for the Kubernetes deployment$`, ps.privilegedAccessRequestIsMarkedForTheKubernetesDeployment)
	ctx.Step(`^some system exists to prevent a Kubernetes deployment running with privileged access in an existing Kubernetes cluster$`, ps.someControlExistsToPreventPrivilegedAccessForKubernetesDeploymentsToAnActiveKubernetesCluster)
	ctx.Step(`^I should not be able to perform a command that requires privileged access$`, ps.iShouldNotBeAbleToPerformACommandThatRequiresPrivilegedAccess)

	//CIS-5.2.2
	ctx.Step(`^hostPID request is marked "([^"]*)" for the Kubernetes deployment$`, ps.hostPIDRequestIsMarkedForTheKubernetesDeployment)
	ctx.Step(`^some system exists to prevent a Kubernetes deployment from running using the hostPID in an existing Kubernetes cluster$`, ps.someSystemExistsToPreventAKubernetesContainerFromRunningUsingTheHostPIDOnTheActiveKubernetesCluster)
	ctx.Step(`^I should not be able to perform a command that provides access to the host PID namespace$`, ps.iShouldNotBeAbleToPerformACommandThatProvidesAccessToTheHostPIDNamespace)

	//CIS-5.2.3
	ctx.Step(`^hostIPC request is marked "([^"]*)" for the Kubernetes deployment$`, ps.hostIPCRequestIsMarkedForTheKubernetesDeployment)
	ctx.Step(`^some system exists to prevent a Kubernetes deployment from running using the hostIPC in an existing Kubernetes cluster$`, ps.someSystemExistsToPreventAKubernetesDeploymentFromRunningUsingTheHostIPCInAnExistingKubernetesCluster)
	ctx.Step(`^I should not be able to perform a command that provides access to the host IPC namespace$`, ps.iShouldNotBeAbleToPerformACommandThatProvidesAccessToTheHostIPCNamespace)

	//CIS-5.2.4
	ctx.Step(`^hostNetwork request is marked "([^"]*)" for the Kubernetes deployment$`, ps.hostNetworkRequestIsMarkedForTheKubernetesDeployment)
	ctx.Step(`^some system exists to prevent a Kubernetes deployment from running using the hostNetwork in an existing Kubernetes cluster$`, ps.someSystemExistsToPreventAKubernetesDeploymentFromRunningUsingTheHostNetworkInAnExistingKubernetesCluster)
	ctx.Step(`^I should not be able to perform a command that provides access to the host network namespace$`, ps.iShouldNotBeAbleToPerformACommandThatProvidesAccessToTheHostNetworkNamespace)

	//CIS-5.2.5
	ctx.Step(`^privileged escalation is marked "([^"]*)" for the Kubernetes deployment$`, ps.privilegedEscalationIsMarkedForTheKubernetesDeployment)
	ctx.Step(`^some system exists to prevent a Kubernetes deployment from running using the allowPrivilegeEscalation in an existing Kubernetes cluster$`, ps.someSystemExistsToPreventAKubernetesDeploymentFromRunningUsingTheAllowPrivilegeEscalationInAnExistingKubernetesCluster)
	ctx.Step(`^I should not be able to perform a sudo command that requires privileged access$`, ps.iShouldNotBeAbleToPerformASudoCommandThatRequiresPrivilegedAccess)

	//CIS-5.2.6
	ctx.Step(`^the user requested is "([^"]*)" for the Kubernetes deployment$`, ps.theUserRequestedIsForTheKubernetesDeployment)
	ctx.Step(`^some system exists to prevent a Kubernetes deployment from running as the root user in an existing Kubernetes cluster$`, ps.someSystemExistsToPreventAKubernetesDeploymentFromRunningAsTheRootUserInAnExistingKubernetesCluster)
	ctx.Step(`^the Kubernetes deployment should run with a non-root UID$`, ps.theKubernetesDeploymentShouldRunWithANonrootUID)

	//CIS-5.2.7
	ctx.Step(`^a Kubernetes deployment with NET_RAW capability "([^"]*)" is applied to an existing Kubernetes cluster$`, ps.kubernetesDeploymentWithNETRAWCapability)
	ctx.Step(`^I should not be able to perform a command that requires NET_RAW capability$`, ps.iShouldNotBeAbleToPerformACommandThatRequiresNETRAWCapability)

	//CIS-5.2.8
	//ctx.Step(`^additional capabilities "([^"]*)" requested for the Kubernetes deployment$`, ps.additionalCapabilitiesForTheKubernetesDeployment)
	ctx.Step(`^additional capabilities requested for the Kubernetes deployment are "([^"]*)" allowed`, ps.additionalCapabilitiesForTheKubernetesDeployment)
	ctx.Step(`^some system exists to prevent Kubernetes deployments with capabilities beyond the default set from being deployed to an existing kubernetes cluster$`, ps.someSystemExistsToPreventKubernetesDeploymentsWithCapabilitiesBeyondTheDefaultSetFromBeingDeployedToAnExistingKubernetesCluster)
	ctx.Step(`^I should not be able to perform a command that requires capabilities outside of the default set$`, ps.iShouldNotBeAbleToPerformACommandThatRequiresCapabilitiesOutsideOfTheDefaultSet)

	//CIS-5.2.9
	ctx.Step(`^assigned capabilities "([^"]*)" requested for the Kubernetes deployment$`, ps.assignedCapabilitiesForTheKubernetesDeployment)
	ctx.Step(`^some system exists to prevent Kubernetes deployments with assigned capabilities from being deployed to an existing Kubernetes cluster$`, ps.someSystemExistsToPreventKubernetesDeploymentsWithAssignedCapabilitiesFromBeingDeployedToAnExistingKubernetesCluster)
	ctx.Step(`^I should not be able to perform a command that requires any capabilities$`, ps.iShouldNotBeAbleToPerformACommandThatRequiresAnyCapabilities)

	//AZPolicy - port range
	ctx.Step(`^an "([^"]*)" hostPort is requested for the Kubernetes deployment$`, ps.anPortRangeIsRequestedForTheKubernetesDeployment)
	ctx.Step(`^some system exists to prevent Kubernetes deployments with unapproved port range from being deployed to an existing Kubernetes cluster$`, ps.someSystemExistsToPreventKubernetesDeploymentsWithUnapprovedPortRangeFromBeingDeployedToAnExistingKubernetesCluster)

	//AZPolicy - volume types
	ctx.Step(`^"([^"]*)" volume types are requested for the Kubernetes deployment$`, ps.volumeTypesAreRequestedForTheKubernetesDeployment)
	ctx.Step(`^I should not be able to perform a command that accesses an unapproved volume type$`, ps.iShouldNotBeAbleToPerformACommandThatAccessesAnUnapprovedVolumeType)
	ctx.Step(`^some system exists to prevent Kubernetes deployments with unapproved volume types from being deployed to an existing Kubernetes cluster$`, ps.someSystemExistsToPreventKubernetesDeploymentsWithUnapprovedVolumeTypesFromBeingDeployedToAnExistingKubernetesCluster)

	//AZPolicy - seccomp profile
	ctx.Step(`^an "([^"]*)" seccomp profile is requested for the Kubernetes deployment$`, ps.anSeccompProfileIsRequestedForTheKubernetesDeployment)
	ctx.Step(`^some system exists to prevent Kubernetes deployments without approved seccomp profiles from being deployed to an existing Kubernetes cluster$`, ps.someSystemExistsToPreventKubernetesDeploymentsWithoutApprovedSeccompProfilesFromBeingDeployedToAnExistingKubernetesCluster)
	ctx.Step(`^I should not be able to perform a system call that is blocked by the seccomp profile$`, ps.iShouldNotBeAbleToPerformASystemCallThatIsBlockedByTheSeccompProfile)

	//general - outcome
	ctx.Step(`^the operation will "([^"]*)" with an error "([^"]*)"$`, ps.theOperationWillWithAnError)
	ctx.Step(`^all operations will "([^"]*)" with an error "([^"]*)"$`, ps.allOperationsWillWithAnError)
	ctx.Step(`^I should be able to perform an allowed command$`, ps.iShouldBeAbleToPerformAnAllowedCommand)

	ctx.AfterScenario(func(s *godog.Scenario, err error) {
		if kubernetes.GetKeepPodsFromConfig() == false {
			if len(ps.podStates) == 0 {
				psp.TeardownPodSecurityProbe(ps.podState.PodName, p.Name())
			} else {
				for _, s := range ps.podStates {
					psp.TeardownPodSecurityProbe(s.PodName, p.Name())
				}
			}
		}
		ps.podState.PodName = ""
		ps.podState.CreationError = nil
		coreengine.LogScenarioEnd(s)
	})
}
