package vch

import (
	"fmt"
	"testing"
	"time"
)

type Student struct {
	Name     string
	NickName string
	Class    string
	Degree   string
}

func TestVch(t *testing.T) {

	chName := `可视化通道01`
	// var logFile *os.File
	// var err error
	// logFile, err = os.OpenFile("channel_log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// log.SetOutput(logFile)
	InitLog(`runtime/log`)
	newCh := NewVChannel[*Student](chName, 1000)
	newCh.Log()
	for j := 0; j < 20; j++ {
		student := &Student{
			Name:   fmt.Sprintf(`q名字_%d_%d`, 1, j),
			Class:  fmt.Sprintf(`q班级-%d-%d`, 1, j),
			Degree: fmt.Sprintf(`q学位-%d-%d`, 1, j),
		}
		newCh.Push(chName, student)
	}
	newCh.SetCanLog(true)
	time.Sleep(time.Second)
	i := 0
	go func() {
		for j := 0; j < 200; j++ {
			student := Student{
				Name:   fmt.Sprintf(`名字_%d_%d`, i, j),
				Class:  fmt.Sprintf(`班级-%d-%d`, i, j),
				Degree: fmt.Sprintf(`学位-%d-%d`, i, j),
			}
			newCh.Push(chName, &student)
			i++
			// time.Sleep(time.Second * 1)
		}
	}()
	time.Sleep(time.Microsecond)
	go func() {
		for i := 0; i < 100; i++ {
			_ = newCh.Pull(chName)
			// time.Sleep(time.Second * 1)
		}
	}()
	time.Sleep(time.Second * 120)
	// time.Sleep(time.Second * 110)
}
