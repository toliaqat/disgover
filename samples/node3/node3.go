package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/dispatchlabs/disgover"
	"github.com/dispatchlabs/disgover/transport"
)

func main() {
	var contact = disgover.Contact{
		Id: "333333333333333333333333333",
		Endpoint: disgover.Endpoint{
			Port: 9003,
		},
	}

	var disgover *disgover.Disgover = disgover.NewDisgover(
		&contact,
		[]*disgover.Contact{
			&disgover.Contact{
				Id: "111111111111111111111111111",
				Endpoint: disgover.Endpoint{
					Host: "172.17.0.5",
					Port: 9001,
				},
			},
		},
		transport.NewHTTPTransport(contact.Endpoint),
	)
	disgover.Run()

	node2, _ := disgover.Find("222222222222222222222222222", disgover.Contact)

	node2AsBytes, _ := json.Marshal(node2)

	fmt.Println("DISGOVER: Find()")
	fmt.Println("         ", string(node2AsBytes[:]))

	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigs
		fmt.Println()
		fmt.Println(sig)
		done <- true
	}()
	<-done
}
