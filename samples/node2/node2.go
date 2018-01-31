package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/dispatchlabs/disgover"
	"github.com/dispatchlabs/disgover/transport"
)

func main() {
	var contact = disgover.Contact{
		Id: "222222222222222222222222222",
		Endpoint: disgover.Endpoint{
			Port: 9002,
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
