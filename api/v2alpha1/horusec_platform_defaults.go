package v2alpha1

import (
	"encoding/json"
	"io/ioutil"
)

func MergeWithDefaultValues(horus *HorusecPlatform) (*HorusecPlatform, error) {
	data, err := ioutil.ReadFile("defaults.json")
	if err != nil {
		return nil, err
	}

	merged := new(HorusecPlatform)
	if err = json.Unmarshal(data, &merged.Spec); err != nil {
		return nil, err
	}

	jb, err := json.Marshal(horus)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(jb, &merged)
	if err != nil {
		return nil, err
	}

	return merged, nil
}
