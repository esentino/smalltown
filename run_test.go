package main

import (
	"fmt"
	"testing"
)

func TestWorkTypeToName(t *testing.T) {
	workTypes := map[workType]string{
		gather_wood:  "wood gathering",
		gather_stone: "stone gathering",
		build_house:  "house building",
		idle:         "idle",
	}

	for key, _ := range workTypes {
		name := work_type_to_name(key)
		if name != workTypes[key] {
			panic(fmt.Sprintf("name %s - expected %s", name, workTypes[key]))
		}
	}

}
