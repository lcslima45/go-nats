package main

import (
	"awesome/episode2/models"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/nats-io/nats.go"
)

func main() {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatalf("can't connect to NATS: %v", err)
	}
	defer nc.Close()

	go func() {
		count := 0
		for {
			payload := models.Payload{
				Data:  "Hello World!",
				Count: count,
			}
			data, err := json.Marshal(payload)
			if err != nil {
				log.Println("Error marshaling:", err)
				continue
			}
			msg, err := nc.Request("intros", data, 2*time.Second)
			if err != nil {
				log.Printf("Request error: %v", err)
			} else {
				log.Printf("Got reply: %s", string(msg.Data))
			}
			count++
			time.Sleep(2 * time.Second)
		}
	}()

	_, err = nc.Subscribe("intros", func(m *nats.Msg) {
		pl := &models.Payload{}
		if err := json.Unmarshal(m.Data, pl); err != nil {
			fmt.Println("Error unmarshalling payload:", err)
			return
		}
		replyData := fmt.Sprintf("ack message # %v", pl.Count)
		m.Respond([]byte(replyData))
		fmt.Printf("I got message %s count %v\n", pl.Data, pl.Count)
	})
	if err != nil {
		log.Fatalf("Subscription error: %v", err)
	}

	select {}
}
