package main

import (
	"os"

	"fmt"
	"os/signal"
	"syscall"

	"github.com/dispatchlabs/disgover"
	"github.com/dispatchlabs/disgover/transport"
)

func main() {
	var node1 = disgover.Contact{
		Id: "111111111111111111111111111",
		Endpoint: disgover.Endpoint{
			Port: 9001,
		},
	}

	var disgover *disgover.Disgover = disgover.NewDisgover(
		&node1,
		[]*disgover.Contact{},
		transport.NewHTTPTransport(node1.Endpoint),
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
