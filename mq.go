package main

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"os"
	"time"

	"github.com/robfig/cron"
)

var (
	IP 	  = "159.75.36.133"
	Port  = 1883
	TopicDingDing = "clock/dingding"
)
// d
var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("Received messagels: %s from topic: %s\n", msg.Payload(), msg.Topic())

	fmt.Println("reload Data:", string(msg.Payload()))
}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Println("Connected")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Printf("Connect lost: %v\n", err.Error())
}

func main() {
	var broker = IP
	var port = Port
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", broker, port))
	opts.SetClientID("mqttx_5f070c0522")
	opts.SetDefaultPublishHandler(messagePubHandler)
	opts.OnConnect = connectHandler

	opts.OnConnectionLost = connectLostHandler

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		fmt.Printf("tokenError:%+v",token.Error())
		return
		//panic(token.Error())
	}
	Sub(client)
	publish(client) // 定时器

	//client.Disconnect(250)
	select{}
}

func Sub(client mqtt.Client) {
	topic := TopicDingDing
	if token := client.Subscribe(topic, 0, nil); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}
	client.Subscribe(topic, 0x00, messagePubHandler)
	fmt.Printf("Subscribe topic " + topic + " success\n")
}


func publish(client mqtt.Client) {
	for  {
		msgArr := "{'isClock':true}"
		i := 0
		c := cron.New()
		//spec := "0 */1 * * * *"   // 每一分钟，
		spec := "*/5 * * * * ?" 	// 每5秒
		//spec := "0 20 8 * * ?"  // 每天8点20分
		c.AddFunc(spec, func() {
			i++
			fmt.Println("cron running:",i)
			token := client.Publish(TopicDingDing, 0, false, msgArr)
			token.Wait()
			time.Sleep(time.Second)
		})
		c.Start()
		select{}
	}

}


