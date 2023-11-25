package main

import (
	"fmt"
	"math/rand"
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
	currentWork workType
	progress    int
}

type Resources struct {
	wood  int
	stone int
}

type Building struct {
	home int
}

var workQueue []workType = []workType{}
var workers []worker = []worker{{idle, 0}}
var building = Building{4}
var resources = Resources{0, 0}

func updateScreen(screen tcell.Screen, x int, y int, width int, height int) (int, int, int, int) {
	tview.PrintSimple(screen, "Hello World", x+1, y+1)
	for index, worker := range workers {
		tview.PrintSimple(screen, "Worker: "+work_type_to_name(worker.currentWork), x+1, y+index+1)
		tview.PrintSimple(screen, "Progress: "+fmt.Sprint(worker.progress), x+1+width/2, y+index+1)
	}

	for index, work := range workQueue {
		tview.PrintSimple(screen, "Work Queue: "+work_type_to_name(work), x+1, y+len(workers)+index+1)
	}
	tview.PrintSimple(screen, "-- Resources --", x+width/2, y+len(workers)+1)
	tview.PrintSimple(screen, "Wood: "+fmt.Sprint(resources.wood), x+1+width/2, y+len(workers)+2)
	tview.PrintSimple(screen, "Stone: "+fmt.Sprint(resources.stone), x+1+width/2, y+len(workers)+3)
	tview.PrintSimple(screen, "-- Buildings --", x+width/2, y+len(workers)+4)
	tview.PrintSimple(screen, "Homes: "+fmt.Sprint(building.home), x+1+width/2, y+len(workers)+5)
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
	workQueue = append(workQueue, build_house)
}

func game_progress() {
	tick := time.NewTicker(1000 * time.Millisecond)
	for range tick.C {
		for i, current_worker := range workers {
			if current_worker.currentWork == idle {
				if len(workQueue) > 0 {
					workers[i].currentWork = workQueue[0]
					workQueue = workQueue[1:]
					workers[i].progress = 0
				}
			} else {
				workers[i].progress += rand.Int() % 10
				if current_worker.progress >= 100 {
					doneWork(i)
				}
			}
		}
		app.Draw()
	}
}

func doneWork(i int) {
	switch workers[i].currentWork {
	case gather_wood:
		resources.wood += 1
	case gather_stone:
		resources.stone += 1
	case build_house:
		building.home++
		workers = append(workers, worker{idle, 0})

	}
	workers[i].currentWork = idle
	workers[i].progress = 0
}

func main() {
	app = tview.NewApplication()

	box = tview.NewBox()
	list = tview.NewList()
	list.AddItem("get more wood", "queue job for wood gathering", '1', queue_wood)
	list.AddItem("get mode stone", "queue job for stone gathering", '2', queue_stone)
	list.AddItem("build house (10 wood, 10 stone) add extra worker", "queue job for house building", '3', queue_build_house)
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
