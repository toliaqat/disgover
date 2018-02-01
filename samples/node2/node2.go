package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"net"
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

	var contact = disgover.Contact{
		Id: "222222222222222222222222222",
		Endpoint: disgover.Endpoint{
			Port: 9001,
			Host: addrs[0],
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
