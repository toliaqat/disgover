package main

import (
	"bufio"
	"os"

	"github.com/dispatchlabs/disgover"
	"github.com/dispatchlabs/disgover/transport"
)

func main() {
	var thisID = "222222222222222222222222222"
	var port int64 = 9002

	var contact = disgover.Contact{
		Id: thisID,
		Endpoint: disgover.Endpoint{
			Port: port,
		},
	}

	var disgover *disgover.Disgover = disgover.NewDisgover(
		&contact,
		[]*disgover.Contact{
			&disgover.Contact{
				Id: "111111111111111111111111111",
				Endpoint: disgover.Endpoint{
					Host: "localhost",
					Port: 9001,
				},
			},
		},
		transport.NewHTTPTransport(contact.Endpoint),
	)
	disgover.Run()

	bufio.NewReader(os.Stdin).ReadString('\n')
}
