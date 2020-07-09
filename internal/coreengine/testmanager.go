package coreengine

import (
	"sync"
	"github.com/google/uuid"

	api "citihub.com/probr/api/probrapi"

	
)

//TestStore ...
type TestStore struct {
	Tests map[uuid.UUID]*api.ProbrTest
	Lock  sync.Mutex
}

//GetAvailableTests - return the universe of available tests
func GetAvailableTests() []api.ProbrTest {

	//not sure if this is needed
	//TODO: get this from the defined BDD (how?), hold on to it and rtn a pointer
	//TODO: for now, just knock something up and return
	s := "kube"
	uuid1 := uuid.New().String()
	uuid2 := uuid.New().String()
	p := []api.ProbrTest { 
			{ Category: &s, Uuid: &uuid1 },
			{ Category: &s, Uuid: &uuid2 },
		} 

	return p
}