package stopper

import (
	"encoding/json"
	"os"
	"path/filepath"

	"code.cloudfoundry.org/guardian/rundmc/cgroups"
)

type state struct {
	CgroupPaths map[string]string `json:"cgroup_paths"`
}

type resolver struct {
	stateStore string
}

func NewRuncStateCgroupPathResolver(stateStorePath string) *resolver {
	return &resolver{
		stateStore: stateStorePath,
	}
}

func (r resolver) Resolve(name, subsystem string) (string, error) {
	stateJson, err := os.Open(filepath.Join(r.stateStore, name, "state.json"))
	if err != nil {
		return "", err
	}
	defer stateJson.Close()

	var s state
	if err := json.NewDecoder(stateJson).Decode(&s); err != nil {
		return "", err
	}

	if cgroups.IsCgroup2UnifiedMode() {
		return s.CgroupPaths[""], nil
	}
	return s.CgroupPaths["devices"], nil
}
