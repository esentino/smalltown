package main

import (
	"fmt"
	"os"
	"testing"
	"time"
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
	resources.Stone = 10
	resources.Wood = 10
	queue_build_house()
	if len(workQueue) != 1 || workQueue[0] != build_house {
		panic("failed")
	}
}

func TestQueueBuildHouseNotEnoughtResource(t *testing.T) {
	workQueue = nil
	resources.Stone = 0
	resources.Wood = 0
	queue_build_house()
	if len(workQueue) == 1 {
		panic("failed")
	}
}

func closeWindow() {
	time.Sleep(1_000_000)
	app.Stop()
}

func TestMainFunction(t *testing.T) {
	skipCI(t)
	// go closeWindow()
	// if app != nil {
	// 	panic("app is not nil")
	// }
	// main()
}

func skipCI(t *testing.T) {
	if os.Getenv("CI") != "" {
		t.Skip("Skipping testing in CI environment")
	}
}

func TestDoneWork(t *testing.T) {
	workers = []worker{{gather_wood, 100}, {gather_stone, 100}, {build_house, 100}}
	doneWork(0)
	if workers[0].CurrentWork != idle && workers[0].Progress != 0 {
		panic("worker[0] not idle")
	}

	doneWork(1)
	if workers[1].CurrentWork != idle && workers[1].Progress != 0 {
		panic("worker[1] not idle")
	}

	doneWork(2)
	if workers[2].CurrentWork != idle && workers[2].Progress != 0 {
		panic("worker[2] not idle")
	}

}
