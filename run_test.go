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

func TestQueueBuildHouse(t *testing.T) {
	resources.stone = 10
	resources.wood = 10
	queue_build_house()
	if len(workQueue) != 1 || workQueue[0] != build_house {
		panic("failed")
	}
}
