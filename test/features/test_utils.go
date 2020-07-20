package features

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var outputDir *string

//GetProbrRoot ...
func GetProbrRoot() (string, error) {
	//TODO: fix this!! think it's a tad dodgy!
	pwd, _ := os.Getwd()
	log.Printf("[DEBUG] GetProbrRoot pwd is: %v", pwd)

	b := strings.Contains(pwd, "probr")
	if !b {
		return "", fmt.Errorf("could not find 'probr' root directory in %v", pwd)
	}

	s := strings.SplitAfter(pwd, "probr")
	log.Printf("[DEBUG] path(s) after splitting: %v\n", s)

	if len(s) < 1 {
		//expect at least one result
		return "", fmt.Errorf("could not split out 'probr' from directory in %v", pwd)
	}

	if !strings.HasSuffix(s[0], "probr") {
		//the first path should end with "probr"
		return "", fmt.Errorf("first path after split (%v) does not end with 'probr'", s[0])
	}

	return s[0], nil
}

//SetOutputDirectory ...
func SetOutputDirectory(d *string) {
	outputDir = d
}

//GetOutputPath gets the output path for the test based on the output directory
//plus the test name supplied
func GetOutputPath(t *string) (*os.File, error) {
	outPath, err := getOutputDirectory()
	if err != nil {
		return nil, err
	}
	os.Mkdir(*outPath, os.ModeDir)

	//filename is test name (supplied) + .json
	fn := *t + ".json"
	return os.Create(filepath.Join(*outPath, fn))
}

func getOutputDirectory() (*string, error) {
	if outputDir == nil {
		log.Printf("[INFO] output directory not set - attempting to default")
		//default it:
		r, err := GetProbrRoot()
		if err != nil {
			return nil, fmt.Errorf("output directory not set - attempt to default resulted in error: %v", err)
		}

		f := filepath.Join(r, "testoutput")
		outputDir = &f
	}

	log.Printf("[INFO] output directory is: %v", *outputDir)

	return outputDir, nil
}
