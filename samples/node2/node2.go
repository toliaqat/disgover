package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	//"github.com/dispatchlabs/disgover"
	//"github.com/dispatchlabs/disgover/proto"
)

func main() {
	var seedNodeIP = os.Getenv("SEED_NODE_IP") // Needed when run from Kubernetes
	if len(seedNodeIP) == 0 {
		seedNodeIP = "127.0.0.1"
	}
	fmt.Printf("SEED_NODE_IP: %s\n", seedNodeIP)

	//var dsg = disgover.NewDisgover(
	//	disgover.NewContact(),
	//	[]*proto.Contact{
	//		&proto.Contact{
	//			// Id: "NODE-1",
	//			Endpoint: &proto.Endpoint{
	//				Host: seedNodeIP,
	//				Port: 9001,
	//			},
	//		},
	//	},
	//)
	//dsg.ThisContact.Id = "NODE-2"
	//dsg.ThisContact.Endpoint.Port = 9002
	//dsg.Run()

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
