package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"

	"github.com/dispatchlabs/disgover"
	"github.com/dispatchlabs/disgover/transport"
)

func main() {
	var contact = disgover.Contact{
		Id: "333333333333333333333333333",
		Endpoint: disgover.Endpoint{
			Port: 9003,
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

	node2AsBytes, _ := json.Marshal(node2)

	fmt.Println("DISGOVER: Find()")
	fmt.Println("         ", string(node2AsBytes[:]))

	bufio.NewReader(os.Stdin).ReadString('\n')
}
