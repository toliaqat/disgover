package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/dispatchlabs/disgover"
)

func main() {
	var seedNodeIP = os.Getenv("SEED_NODE_IP")
	fmt.Printf("SEED_NODE_IP: %s\n", seedNodeIP)

	var thisNode = disgover.NewContact()
	thisNode.Id = "333333333333333333333333333"
	// thisNode.Endpoint.Host = ""
	// thisNode.Endpoint.Port = 9003

	var dsg *disgover.Disgover = disgover.NewDisgover(
		thisNode,
		[]*disgover.Contact{
			&disgover.Contact{
				Id: "111111111111111111111111111",
				Endpoint: &disgover.Endpoint{
					Host: seedNodeIP,
					Port: 9001,
				},
			},
		},
	)
	dsg.Run()

	node2, _ := dsg.Find("222222222222222222222222222", thisNode)

	fmt.Println("DISGOVER: Find()")
	if node2 == nil {
		fmt.Println("          NOT FOUND")
	} else {
		fmt.Println(fmt.Sprintf("          %s, on [%s : %d]", node2.Id, node2.Endpoint.Host, node2.Endpoint.Port))
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
