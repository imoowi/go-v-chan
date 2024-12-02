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
	chanName = ``
	//新建通道
	vChanl := vch.NewVChannel[*Student](chanName, 100000)
	// time.Sleep(time.Second)
	for i := 0; i < 100; i++ {
		go func(i int) {
			for j := 0; j < 20; j++ {
				student := Student{
					Name:   fmt.Sprintf(`名字_%d_%d`, i, j),
					Class:  fmt.Sprintf(`班级-%d-%d`, i, j),
					Degree: fmt.Sprintf(`学位-%d-%d`, i, j),
				}
				vChanl.Add(&student)
				time.Sleep(time.Microsecond * 10)
			}
		}(i)
		time.Sleep(time.Microsecond * 100)
		go func() {
			_student := vChanl.Pop()
			log.Println(_student)
		}()
		time.Sleep(time.Microsecond * 100)
		go func() {
			_student := vChanl.Pop()
			log.Println(_student)
		}()
		time.Sleep(time.Microsecond * 100)
		go func() {
			_student := vChanl.Pop()
			log.Println(_student)
		}()
		time.Sleep(time.Microsecond * 100)
		go func() {
			_student := vChanl.Pop()
			log.Println(_student)
		}()
		log.Println(`vChanl.Len()=`, vChanl.Len())
	}
	// time.Sleep(time.Microsecond*100 * 10)
	time.Sleep(time.Second * 30)
}
