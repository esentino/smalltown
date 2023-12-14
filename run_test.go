package main

import (
	"fmt"
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
	resources.stone = 10
	resources.wood = 10
	queue_build_house()
	if len(workQueue) != 1 || workQueue[0] != build_house {
		panic("failed")
	}
}

func TestQueueBuildHouseNotEnoughtResource(t *testing.T) {
	workQueue = nil
	resources.stone = 0
	resources.wood = 0
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
	go closeWindow()
	if app != nil {
		panic("app is not nil")
	}
	main()

}

func TestDoneWork(t *testing.T) {
	workers = []worker{{gather_wood, 100}, {gather_stone, 100}, {build_house, 100}}
	doneWork(0)
	if workers[0].currentWork != idle && workers[0].progress != 0 {
		panic("worker[0] not idle")
	}

	doneWork(1)
	if workers[1].currentWork != idle && workers[1].progress != 0 {
		panic("worker[1] not idle")
	}

	doneWork(2)
	if workers[2].currentWork != idle && workers[2].progress != 0 {
		panic("worker[2] not idle")
	}

}
