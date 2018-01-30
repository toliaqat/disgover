package main

import (
	"bufio"
	"os"

	"github.com/dispatchlabs/disgover"
	"github.com/dispatchlabs/disgover/transport"
)

func main() {
	var contact = disgover.Contact{
		Id: "222222222222222222222222222",
		Endpoint: disgover.Endpoint{
			Port: 9002,
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
