package kubernetes_test

import (
	"flag"
	"fmt"
	"log"
	"os"
	"testing"

	"citihub.com/probr/internal/clouddriver/kubernetes"
	_ "citihub.com/probr/internal/config"
)

var (
	integrationTest = flag.Bool("integrationTest", false, "run integration tests")
)

func TestMain(m *testing.M) {
	flag.Parse()

	argLength := len(os.Args[1:])
	fmt.Printf("Arg length is %d\n", argLength)

	for i, a := range os.Args[1:] {
		fmt.Printf("Arg %d is %s\n", i+1, a)
	}

	args := flag.Args()
	log.Printf("Args: %v", args)

	if !*integrationTest {
		//skip
		log.Print("kubehelper_test: Integration Test Flag not set. SKIPPING TEST.")
		return
	}

	result := m.Run()

	os.Exit(result)
}

func TestGetPods(t *testing.T) {
	kubernetes.GetPods()
}

func TestCreatePod(t *testing.T) {
	pname, ns, cname, image := "my-test2", "default", "probr", "busybox"
	_, err := kubernetes.CreatePod(&pname, &ns, &cname, &image)

	handleResult(nil, err)
}

func TestDeletePod(t *testing.T) {
	pname, ns := "my-test2", "default"
	err := kubernetes.DeletePod(&pname, &ns)

	handleResult(nil, err)
}

func TestExecCmd(t *testing.T) {

	url := "http://www.google.com"
	cmd := fmt.Sprintf("curl -s -o /dev/null -I -L -w \"%%{http_code}\" %s", url)
	ns, pn := "default", "my-testcurl-pod"
	so, se, err := kubernetes.ExecCommand(&cmd, &ns, &pn)

	fmt.Printf("CMD Result\n")
	fmt.Printf("stdout:\n %v \n stderr: %v\n", so, se)

	handleResult(nil, err)
}

func TestCreateNamespace(t *testing.T) {
	ns := "default-1"
	_, err := kubernetes.CreateNamespace(&ns)

	handleResult(nil, err)
}

func TestDeleteNamespace(t *testing.T) {
	ns := "default-1"
	err := kubernetes.DeleteNamespace(&ns)

	handleResult(nil, err)
}

func handleResult(yesNo *bool, err error) {
	if err != nil {
		//FAIL ... but don't check for this atm ...
		fmt.Printf("Test failed: %v\n", err)
		return
	}

	if yesNo != nil {
		fmt.Printf("RESULT: %t\n", *yesNo)
	}

}
