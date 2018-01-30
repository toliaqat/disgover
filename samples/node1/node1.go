package main

import (
	"bufio"
	"os"

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

	bufio.NewReader(os.Stdin).ReadString('\n')
}
