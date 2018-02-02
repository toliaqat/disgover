package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/dispatchlabs/disgover"
)

func main() {
	var dsg = disgover.NewDisgover(
		disgover.NewContact(),
		[]*disgover.Contact{},
	)
	// dsg.ThisContact.Id = "NODE-1"
	dsg.ThisContact.Endpoint.Port = 9001
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
