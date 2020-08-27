package features

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const rootDirName = "probr"

var outputDir *string

//GetRootDir ...
func GetRootDir() (string, error) {
	//TODO: fix this!! think it's a tad dodgy!
	pwd, _ := os.Getwd()
	log.Printf("[DEBUG] GetRootDir pwd is: %v", pwd)

	b := strings.Contains(pwd, rootDirName)
	if !b {
		return "", fmt.Errorf("could not find '%v' root directory in %v", rootDirName, pwd)
	}

	s := strings.SplitAfter(pwd, rootDirName)
	log.Printf("[DEBUG] path(s) after splitting: %v\n", s)

	if len(s) < 1 {
		//expect at least one result
		return "", fmt.Errorf("could not split out '%v' from directory in %v", rootDirName, pwd)
	}

	if !strings.HasSuffix(s[0], rootDirName) {
		//the first path should end with "probr"
		return "", fmt.Errorf("first path after split (%v) does not end with '%v'", s[0], rootDirName)
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
	if outputDir == nil || len(*outputDir) < 1 {
		log.Printf("[INFO] output directory not set - attempting to default")
		//default it:
		r, err := GetRootDir()
		if err != nil {
			return nil, fmt.Errorf("output directory not set - attempt to default resulted in error: %v", err)
		}

		f := filepath.Join(r, "testoutput")
		outputDir = &f
	}

	log.Printf("[INFO] output directory is: %v", *outputDir)

	return outputDir, nil
}

// LogAndReturnError logs the given string and raise an error with the same string.  This is useful in Godog steps
// where an error is displayed in the test report but not logged.
func LogAndReturnError(e string, v ...interface{}) error {
	var b  strings.Builder
	b.WriteString("[ERROR] ")
	b.WriteString(e)

	s := fmt.Sprintf(b.String(), v...)
	log.Print(s)

	return fmt.Errorf(s)
}
