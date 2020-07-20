package features

import (
	"os"
	"fmt"
	"log"
	"strings"
)

//GetProbrRoot ...
func GetProbrRoot() (string, error) {
	//TODO: fix this!! think it's a tad dodgy!
	pwd, _ := os.Getwd()
	log.Printf("[DEBUG] GetProbrRoot pwd is: %v", pwd)	
	
	b := strings.Contains(pwd, "probr")
	if !b {
		return "", fmt.Errorf("could not find 'probr' root directory in %v", pwd)
	}
	
	s := strings.SplitAfter(pwd,"probr")
	log.Printf("[DEBUG] path(s) after splitting: %v\n",s)

	if len(s) < 1 { 
		//expect at least one result
		return "", fmt.Errorf("could not split out 'probr' from directory in %v", pwd)
	}
	
	if  !strings.HasSuffix(s[0], "probr"){
		//the first path should end with "probr"
		return "", fmt.Errorf("first path after split (%v) does not end with 'probr'", s[0])
	}

	return s[0],nil
}