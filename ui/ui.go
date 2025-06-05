package ui

import (
	"fmt"

	"tidyTodo/internal/common"
	"tidyTodo/internal/db"
	"tidyTodo/internal/task"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type item struct {
	flex *tview.Flex
	id   int64
}

type App struct {
	app *tview.Application
}

var (
	rows []item
	DB   *db.DB
	flex *tview.Flex
)

func AppInit() (App, error) {
	var err error
	flex = tview.NewFlex().SetDirection(tview.FlexRow)
	app := tview.NewApplication().EnableMouse(true)

	// tview.Styles.PrimaryTextColor = tcell.ColorBlack
	tview.Styles.ContrastBackgroundColor = tcell.ColorBlack

	DB, err = db.Open()
	common.DB = DB
	if err != nil {
		return App{}, fmt.Errorf("app init: %w", err)
	}

	tasks, err := DB.GetTasks()
	if err != nil {
		return App{}, fmt.Errorf("app init: %w", err)
	}

	for _, tsk := range tasks {
		flex.AddItem(MakeRow(tsk), 3, 0, false)
	}

	flex.SetBorder(true)

	return App{
		app: app,
	}, nil
}

func (a *App) Run() error {
	defer DB.Close()
	return a.app.SetRoot(flex, true).Run()
}

func MakeRow(tsk task.Task) tview.Primitive {
	row := tview.NewFlex().SetDirection(tview.FlexColumn)
	done := tview.NewButton("")
	title := tview.NewTextView().SetTextAlign(tview.AlignLeft)
	delelteBtn := tview.NewButton("[white:red]delete[-:-]")
	editBtn := tview.NewButton("[white:blue]edit[-:-]")

	row.SetBorder(true)

	checkDoneAndSetTitleAndDoneBtn := func(title *tview.TextView, tsk task.Task) {
		if tsk.Done {
			title.SetText(strike(tsk.Title))
			done.SetLabel(`|X|`)
		} else {
			title.SetText(unstrike(tsk.Title))
			done.SetLabel("| |")
		}
	}

	done.SetSelectedFunc(func() {
		common.UpdateTaskStatus(&tsk)
		checkDoneAndSetTitleAndDoneBtn(title, tsk)
	})

	checkDoneAndSetTitleAndDoneBtn(title, tsk)

	row.AddItem(done, 3, 0, false)
	row.AddItem(title, 0, 1, false)
	row.AddItem(editBtn, 5, 0, false)
	row.AddItem(delelteBtn, 7, 0, false)

	rows = append(rows, item{
		flex: row,
		id:   tsk.ID,
	})

	return row
}

/*
strike done label
*/
func strike(label string) string {
	var striked string

	for _, r := range label {
		striked += string(r) + "\u0336"
	}

	return striked
}

/*
unstrike undone label
*/
func unstrike(label string) string {
	var unstriked string

	for _, r := range label {
		if r != '\u0336' {
			unstriked += string(r)
		}
	}

	return unstriked
}
