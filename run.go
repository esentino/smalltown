package main

import (
	"encoding/gob"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var app *tview.Application
var box *tview.Box
var list *tview.List
var grid *tview.Grid

type workType int

const (
	gather_wood workType = iota
	gather_stone
	build_house
	idle
)

type worker struct {
	CurrentWork workType
	Progress    int
}

type Resources struct {
	Wood  int
	Stone int
}

type Building struct {
	Home int
}

var workQueue []workType = []workType{}
var workers []worker = []worker{{idle, 0}}
var building = Building{4}
var resources = Resources{0, 0}

type SaveData struct {
	WorkQueue []workType
	Workers   []worker
	Building  Building
	Resources Resources
}

func updateScreen(screen tcell.Screen, x int, y int, width int, height int) (int, int, int, int) {
	tview.PrintSimple(screen, "Hello World", x+1, y+1)
	for index, worker := range workers {
		tview.PrintSimple(screen, "Worker: "+work_type_to_name(worker.CurrentWork), x+1, y+index+1)
		tview.PrintSimple(screen, "Progress: "+fmt.Sprint(worker.Progress), x+1+width/2, y+index+1)
	}

	for index, work := range workQueue {
		tview.PrintSimple(screen, "Work Queue: "+work_type_to_name(work), x+1, y+len(workers)+index+1)
	}
	tview.PrintSimple(screen, "-- Resources --", x+width/2, y+len(workers)+1)
	tview.PrintSimple(screen, "Wood: "+fmt.Sprint(resources.Wood), x+1+width/2, y+len(workers)+2)
	tview.PrintSimple(screen, "Stone: "+fmt.Sprint(resources.Stone), x+1+width/2, y+len(workers)+3)
	tview.PrintSimple(screen, "-- Buildings --", x+width/2, y+len(workers)+4)
	tview.PrintSimple(screen, "Homes: "+fmt.Sprint(building.Home), x+1+width/2, y+len(workers)+5)
	return 0, 0, 0, 0
}

func work_type_to_name(work workType) string {
	switch work {
	case gather_wood:
		return "wood gathering"
	case gather_stone:
		return "stone gathering"
	case build_house:
		return "house building"
	case idle:
		return "idle"
	}
	return "nill"

}

func queue_wood() {
	workQueue = append(workQueue, gather_wood)
}

func queue_stone() {
	workQueue = append(workQueue, gather_stone)
}

func queue_build_house() {
	if resources.Stone >= 10 && resources.Wood >= 10 {
		resources.Stone -= 10
		resources.Wood -= 10
		if resources.Stone >= 0 && resources.Wood >= 0 {
			workQueue = append(workQueue, build_house)
		} else {
			resources.Stone += 10
			resources.Wood += 10
		}
	}
}

func game_progress() {
	tick := time.NewTicker(1000 * time.Millisecond)
	for range tick.C {
		for i, current_worker := range workers {
			if current_worker.CurrentWork == idle {
				if len(workQueue) > 0 {
					workers[i].CurrentWork = workQueue[0]
					workQueue = workQueue[1:]
					workers[i].Progress = 0
				}
			} else {
				workers[i].Progress += rand.Int() % 10
				if current_worker.Progress >= 100 {
					doneWork(i)
				}
			}
		}
		app.Draw()
	}
}

func doneWork(i int) {
	switch workers[i].CurrentWork {
	case gather_wood:
		resources.Wood += 1
	case gather_stone:
		resources.Stone += 1
	case build_house:
		building.Home += 1
		workers = append(workers, worker{idle, 0})

	}
	workers[i].CurrentWork = idle
	workers[i].Progress = 0
}

func save_and_close() {

	fmt.Println(workers[0].CurrentWork)

	save_data := SaveData{
		WorkQueue: workQueue,
		Workers:   workers,
		Building:  building,
		Resources: resources,
	}
	fmt.Println(save_data.Workers)
	json_data, err := json.Marshal(save_data)
	fmt.Printf(string(json_data))
	if err == nil {
		fmt.Println(string(json_data))
		f, _ := os.OpenFile("save.on", os.O_WRONLY|os.O_CREATE, 0600)
		f.Write(json_data)
		f.Close()
	}
	os.Exit(0)
}

func try_load_save() {
	f, err := os.OpenFile("save.on", os.O_RDONLY, 0600)
	if err == nil {
		encoder := gob.NewDecoder(f)
		var save_data SaveData
		err = encoder.Decode(&save_data)
		if err != nil {

			fmt.Printf("unexpected division error: %s\n", err)
			time.Sleep(8 * time.Second)
			os.Exit(10)
		}
		workQueue = save_data.WorkQueue
		workers = save_data.Workers
		building = save_data.Building
		resources = save_data.Resources
	}
	f.Close()
}

func main() {
	//try_load_save()

	app = tview.NewApplication()

	box = tview.NewBox()
	list = tview.NewList()
	list.AddItem("get more wood", "queue job for wood gathering", '1', queue_wood)
	list.AddItem("get mode stone", "queue job for stone gathering", '2', queue_stone)
	list.AddItem("build house (10 wood, 10 stone) add extra worker", "queue job for house building", '3', queue_build_house)
	list.AddItem("Save and close", "save game to save.on", 'q', save_and_close)

	workers = append(workers, worker{idle, 0})
	workers = append(workers, worker{idle, 0})
	workers = append(workers, worker{idle, 0})
	box.SetBorder(true).SetTitle("Welcome to Smalltown").SetDrawFunc(updateScreen)
	grid = tview.NewGrid()
	grid.AddItem(list, 0, 0, 1, 1, 0, 0, true)
	grid.AddItem(box, 0, 1, 1, 3, 0, 0, false)
	go game_progress()

	if err := app.SetRoot(grid, true).Run(); err != nil {
		panic(err)
	}
}
