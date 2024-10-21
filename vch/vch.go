package vch

import (
	"fmt"
	"reflect"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

type VChannel[T any] struct {
	lock     *sync.RWMutex
	ch       map[string]chan T
	counter  map[string]*ChCounter
	canlog   bool
	chaction chan *ChAction[T]
}

// new a visualization channel
func NewVChannel[T any](chName string, size int) *VChannel[T] {
	newCh := &VChannel[T]{
		lock:     &sync.RWMutex{},
		ch:       make(map[string]chan T),
		counter:  make(map[string]*ChCounter),
		chaction: make(chan *ChAction[T], 10000),
	}
	newCh.ch[chName] = make(chan T, size)
	newCh.counter[chName] = NewChCounter()
	return newCh
}

// set log flag
func (vc *VChannel[T]) SetCanLog(canlog bool) {
	vc.canlog = canlog
}

// judge if can log
func (vc *VChannel[T]) IsCanLog() bool {
	return vc.canlog
}
func (vc *VChannel[T]) DoAction(action *ChAction[T]) {
	action.ActionTime = time.Now()
	vc.chaction <- action
}
func (vc *VChannel[T]) Log() {
	go func() {
		for {
			if vc.canlog {
				cell := <-vc.chaction
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
				time.Sleep(time.Second * 1)
			}
		}
	}()
}

// push a cell into channel
func (vc *VChannel[T]) Push(chName string, cell T) {
	vc.lock.Lock()
	defer vc.lock.Unlock()
	if _ch, ok := vc.ch[chName]; ok {
		_ch <- cell
		// counter++
		vc.counter[chName].Incr()
		action := &ChAction[T]{
			ActionType: ActionTypeIn,
			ChName:     chName,
			ChCounter:  vc.counter[chName].ChCounter(),
			Cell:       cell,
		}
		vc.DoAction(action)
	}
}

// pull a cell from channel
func (vc *VChannel[T]) Pull(chName string) T {
	vc.lock.RLock()
	defer vc.lock.RUnlock()
	var cell T
	ch := vc.ch[chName]
	if cell, ok := <-ch; ok {
		// counter--
		vc.counter[chName].Decr()
		action := &ChAction[T]{
			ActionType: ActionTypeOut,
			ChName:     chName,
			ChCounter:  vc.counter[chName].ChCounter(),
			Cell:       cell,
		}
		vc.DoAction(action)
		return cell
	}
	return cell
}
