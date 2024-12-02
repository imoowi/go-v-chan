package vch

import (
	"fmt"
	"reflect"
	"sync"
	"time"

	"github.com/imoowi/comer/utils/maker"
	"github.com/sirupsen/logrus"
)

type VChannel[T any] struct {
	name     string
	lock     *sync.RWMutex
	ch       map[string]chan T
	counter  map[string]*ChCounter
	chAction chan *ChAction[T]
}

// new a visualization channel
func NewVChannel[T any](chName string, size int) *VChannel[T] {
	if chName == `` {
		chName = maker.MakeSn(`vch_`)
	}
	newCh := &VChannel[T]{
		name:     chName,
		lock:     &sync.RWMutex{},
		ch:       make(map[string]chan T),
		counter:  make(map[string]*ChCounter),
		chAction: make(chan *ChAction[T], 10000),
	}
	newCh.ch[chName] = make(chan T, size)
	newCh.counter[chName] = NewChCounter()
	newCh.Log()
	return newCh
}

func (vc *VChannel[T]) DoAction(action *ChAction[T]) {
	action.ActionTime = time.Now()
	vc.chAction <- action
}
func (vc *VChannel[T]) Log() {
	go func() {
		for {
			if CanLog {
				cell := <-vc.chAction
				sType := reflect.TypeOf(cell.Cell)
				cell.CellType = fmt.Sprintf(`%v`, sType)
				/*
					cellJson, _ := json.Marshal(cell)
					if cell.ActionType == ActionTypeIn {
						log.Println(`<-`, string(cellJson))
					} else {
						log.Println(`->`, string(cellJson))
					}
					//*/
				Log := MyLog.WithFields(logrus.Fields{
					`chName`:  cell.ChName,
					`chCount`: cell.ChCounter,
					`obj`:     cell,
				})
				go Log.Info()
				// time.Sleep(time.Second * 1)
			}
		}
	}()
}

// add a cell into channel
func (vc *VChannel[T]) Add(cell T) {
	vc.lock.Lock()
	defer vc.lock.Unlock()
	if _ch, ok := vc.ch[vc.name]; ok {
		_ch <- cell
		// counter++
		vc.counter[vc.name].Incr()
		action := &ChAction[T]{
			ActionType: ActionTypeIn,
			ChName:     vc.name,
			ChCounter:  vc.counter[vc.name].ChCounter(),
			Cell:       cell,
		}
		vc.DoAction(action)
	}
}

// pop a cell from channel
func (vc *VChannel[T]) Pop() T {
	vc.lock.RLock()
	defer vc.lock.RUnlock()
	var cell T
	ch := vc.ch[vc.name]
	if cell, ok := <-ch; ok {
		// counter--
		vc.counter[vc.name].Decr()
		action := &ChAction[T]{
			ActionType: ActionTypeOut,
			ChName:     vc.name,
			ChCounter:  vc.counter[vc.name].ChCounter(),
			Cell:       cell,
		}
		vc.DoAction(action)
		return cell
	}
	return cell
}
func (vc *VChannel[T]) Len() int {
	return vc.counter[vc.name].ChCounter()
}
