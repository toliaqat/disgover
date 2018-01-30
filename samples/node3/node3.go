package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/dispatchlabs/disgover"
	"github.com/dispatchlabs/disgover/transport"
)

func main() {
	var thisID = "333333333333333333333333333"
	var port int64 = 9003

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

	node2, _ := disgover.Find("222222222222222222222222222", disgover.Contact)
	fmt.Println(node2)

	bufio.NewReader(os.Stdin).ReadString('\n')
}
