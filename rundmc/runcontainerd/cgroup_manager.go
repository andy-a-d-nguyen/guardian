package runcontainerd

import (
	"encoding/json"
	"os"
	"path/filepath"

	"code.cloudfoundry.org/guardian/rundmc/goci"
)

//go:generate go run github.com/maxbrunsfeld/counterfeiter/v6 -generate
//counterfeiter:generate . CgroupManager
type CgroupManager interface {
	SetUseMemoryHierarchy(handle string) error
	SetUnifiedResources(bundle goci.Bndl) error
}

type cgroupManager struct {
	runcRoot  string
	namespace string
}

type containerState struct {
	CgroupPaths cgroupPaths `json:"cgroup_paths"`
}

type cgroupPaths struct {
	Memory string
}

func NewCgroupManager(runcRoot, namespace string) CgroupManager {
	return cgroupManager{
		runcRoot:  runcRoot,
		namespace: namespace,
	}
}

func (m cgroupManager) SetUseMemoryHierarchy(handle string) error {
	statePath := filepath.Join(m.runcRoot, m.namespace, handle, "state.json")
	stateFile, err := os.Open(statePath)
	if err != nil {
		return err
	}
	defer stateFile.Close()

	var state containerState
	err = json.NewDecoder(stateFile).Decode(&state)
	if err != nil {
		return err
	}

	return os.WriteFile(filepath.Join(state.CgroupPaths.Memory, "memory.use_hierarchy"), []byte("1"), os.ModePerm)
}

func (m cgroupManager) SetUnifiedResources(bundle goci.Bndl) error {
	return m.setUnifiedResources(bundle)
}
