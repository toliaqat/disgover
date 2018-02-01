package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/dispatchlabs/disgover"
)

func main() {
	var thisNode = disgover.NewContact()
	thisNode.Id = "111111111111111111111111111"
	// thisNode.Endpoint.Host = ""
	// thisNode.Endpoint.Port = 9001

	var dsg *disgover.Disgover = disgover.NewDisgover(
		thisNode,
		[]*disgover.Contact{},
	)
	dsg.Run()

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
