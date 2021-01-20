// This tool prints all environment variables.
// It can be used along with setenvvars_windows.ps1 to confirm local environment setup.

package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {

	fmt.Println("Fetching all env variables")
	for _, element := range os.Environ() {
		variable := strings.Split(element, "=")
		fmt.Println(variable[0], "=>", variable[1])
	}

	fmt.Println("")

	fmt.Println("Fetching specific env variables")
	fmt.Println("AZURE_TENANT_ID=>", os.Getenv("AZURE_TENANT_ID"))
	fmt.Println("AZURE_SUBSCRIPTION_ID=>", os.Getenv("AZURE_SUBSCRIPTION_ID"))
	fmt.Println("AZURE_CLIENT_ID=>", os.Getenv("AZURE_CLIENT_ID"))
	fmt.Println("AZURE_CLIENT_SECRET=>", os.Getenv("AZURE_CLIENT_SECRET"))
}
