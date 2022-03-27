package fetch_data

import (
	"fmt"
	"github.com/looplab/fsm"
	"github.com/zmisgod/gofun/lunar"
	"testing"
	"time"
)

type Door struct {
	To  string
	FSM *fsm.FSM
}

func NewDoor(to string) *Door {
	d := &Door{
		To: to,
	}

	d.FSM = fsm.NewFSM(
		"closed",
		fsm.Events{
			{Name: "open111", Src: []string{"closed"}, Dst: "open"},
			{Name: "close222", Src: []string{"open"}, Dst: "closed"},
		},
		fsm.Callbacks{
			"enter_state": func(e *fsm.Event) { d.enterState(e) },
		},
	)

	return d
}

func TestNewLunar(t *testing.T) {
	myLunar1, err := lunar.SolarTimeToLunar(time.Date(time.Now().Year(), 10, 18, 0, 0, 0, 0, time.Local))
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Printf("%#v\n", myLunar1)
}

func (d *Door) enterState(e *fsm.Event) {
	fmt.Println(e.Event)
	fmt.Printf("%+v", e)
	fmt.Printf("The door to %s is %s\n", d.To, e.Dst)
}

func TestNewFms(t *testing.T) {
	door := NewDoor("heaven")

	err := door.FSM.Event("open111")
	if err != nil {
		fmt.Println(err)
	}

	err = door.FSM.Event("close222")
	if err != nil {
		fmt.Println(err)
	}
}

func TestNewAPP(t *testing.T) {
	fsm := fsm.NewFSM(
		"closed",
		fsm.Events{
			{Name: "openName", Src: []string{"closed"}, Dst: "open1"},
			{Name: "closeName", Src: []string{"open1"}, Dst: "closed"},
		},
		fsm.Callbacks{},
	)
	fmt.Println(fsm.Current())
	err := fsm.Event("openName")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(fsm.Current())
	err = fsm.Event("closeName")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(fsm.Current())
}
