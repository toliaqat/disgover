package main

import (
	"os"
	"net"
	"fmt"
	"os/signal"
	"syscall"

	"github.com/dispatchlabs/disgover"
	"github.com/dispatchlabs/disgover/transport"
)

func main() {
	name, err := os.Hostname()
	if err != nil {
		fmt.Printf("Oops: %v\n", err)
		return
	}

	addrs, err := net.LookupHost(name)
	if err != nil {
		fmt.Printf("Oops: %v\n", err)
		return
	}
	fmt.Printf("Local IP: %s\n", addrs[0])

	var node1 = disgover.Contact{
		Id: "111111111111111111111111111",
		Endpoint: disgover.Endpoint{
			Port: 9001,
			Host: addrs[0],
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
