package test

import (
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/imoowi/go-v-chan/vch"
)

type Student struct {
	Name     string
	NickName string
	Class    string
	Degree   string
}

func TestVch(t *testing.T) {
	//初始化log日志
	vch.InitLog(`../runtime/log`)
	vch.CanLog = true
	//定义通道名
	chanName := `9527-大内密探零零狗`
	//新建通道
	vChanl := vch.NewVChannel[*Student](chanName, 1000)
	// time.Sleep(time.Second)
	i := 0
	go func() {
		for j := 0; j < 200; j++ {
			student := Student{
				Name:   fmt.Sprintf(`名字_%d_%d`, i, j),
				Class:  fmt.Sprintf(`班级-%d-%d`, i, j),
				Degree: fmt.Sprintf(`学位-%d-%d`, i, j),
			}
			vChanl.Add(&student)
			i++
			// time.Sleep(time.Second * 1)
		}
	}()
	time.Sleep(time.Microsecond)
	go func() {
		for i := 0; i < 100; i++ {
			_student := vChanl.Pop()
			log.Println(_student)
			// time.Sleep(time.Second * 1)
		}
	}()
	time.Sleep(time.Second * 120)
	// time.Sleep(time.Second * 110)
}
