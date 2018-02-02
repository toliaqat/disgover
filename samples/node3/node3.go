package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/dispatchlabs/disgover"
)

func main() {
	var seedNodeIP = os.Getenv("SEED_NODE_IP") // Needed when run from Kubernetes
	if len(seedNodeIP) == 0 {
		seedNodeIP = "127.0.0.1"
	}
	fmt.Printf("SEED_NODE_IP: %s\n", seedNodeIP)

	var dsg = disgover.NewDisgover(
		disgover.NewContact(),
		[]*disgover.Contact{
			&disgover.Contact{
				// Id: "NODE-1",
				Endpoint: &disgover.Endpoint{
					Host: seedNodeIP,
					Port: 9001,
				},
			},
		},
	)
	dsg.ThisContact.Id = "NODE-3"
	dsg.ThisContact.Endpoint.Port = 9003
	dsg.Run()

	node2, _ := dsg.Find("NODE-2", dsg.ThisContact)

	if node2 == nil {
		fmt.Println("DISGOVER: Find() -> NOT FOUND")
	} else {
		fmt.Println(fmt.Sprintf("DISGOVER: Find() -> %s on [%s : %d]", node2.Id, node2.Endpoint.Host, node2.Endpoint.Port))
	}

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
