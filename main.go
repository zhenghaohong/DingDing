package main

import (
	"fmt"
	"github.com/robfig/cron"
	"golang.org/x/net/websocket"
	"net/http"
	"time"
)
func main() {
	startWebSocket()

}

func WsHandle(ws *websocket.Conn) {
	for {
		var reply string
		if err := websocket.Message.Receive(ws, &reply); err != nil {
			fmt.Println("receive failed:", err)
			break
		}
		fmt.Println("received from client: " + reply)
		go tickWriters(ws)
	}

}
func startWebSocket() {
	// 接受websocket的路由地址
	http.Handle("/dingDing", websocket.Handler(WsHandle))
	if err := http.ListenAndServe(":9100", nil); err != nil {
		fmt.Printf("webSocket 失败:%+v",err)
		return
	}

}


func tickWriters(ws *websocket.Conn) {
	//for {
	//	second :=time.Now().Second()      //秒
	//	fmt.Printf("当前秒: %+v\n",second)

	//	msgArr := "{'isClock':'true'}"
	//	if second == 20 {
	//		if err := websocket.Message.Send(ws, msgArr); err != nil {
	//			fmt.Println(time.Now().Format("2006-01-02 15:04:05")," send failed:", err)
	//			break
	//		}
	//		fmt.Printf("间隔一秒发送成功:%+v\n",msgArr)
	//	}
	//	//休息一秒
	//	time.Sleep(time.Second)
	//}

	for  {
		msgArr := "{'isClock':'true'}"
		i := 0
		c := cron.New()
		spec := "0 */1 * * * *"   // 每一分钟，
		//spec := "*/5 * * * * ?" 	// 每5秒
		//spec := "0 20 8 * * ?"  // 每天8点20分
		c.AddFunc(spec, func() {
			i++
			fmt.Println("cron running:",i)
			if err := websocket.Message.Send(ws, msgArr); err != nil {
				fmt.Println(time.Now().Format("2006-01-02 15:04:05")," send failed:", err)
				return
			}
		})
		c.Start()
		select{}
	}
}




