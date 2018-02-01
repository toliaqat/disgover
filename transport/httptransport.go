package transport

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"github.com/nic0lae/JerryMouse/Servers"
	"github.com/dispatchlabs/disgover"
)

// HTTPTransport -
type HTTPTransport struct {
	Endpoint          disgover.Endpoint
	OnPeerRPCDelegate *disgover.DisgoverRpcDelegate
	server            *Servers.ApiServer
}

// NewHTTPTransport -
func NewHTTPTransport(endpoint disgover.Endpoint) *HTTPTransport {
	return &HTTPTransport{
		Endpoint:          endpoint,
		OnPeerRPCDelegate: nil,
		server:            Servers.Api(),
	}
}

// ExecRPC -
func (transport *HTTPTransport) ExecRPC(destination *disgover.Contact, rpc disgover.DisgoverRpc) []byte {
	rpcAsBytes, err := json.Marshal(rpc)

	url := fmt.Sprintf("http://%s:%d/disgover", destination.Endpoint.Host, destination.Endpoint.Port)

	fmt.Println("TRACE: ExecRPC()")
	fmt.Println("       URL-> ", url)
	fmt.Println("   Payload-> ", string(rpcAsBytes[:]))

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(rpcAsBytes))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("ERROR: ExecRPC(): " + err.Error())
		return nil
	}

	body, _ := ioutil.ReadAll(resp.Body)

	return body
}

// OnPeerRPC -
func (transport *HTTPTransport) OnPeerRPC(delegate disgover.DisgoverRpcDelegate) {
	transport.OnPeerRPCDelegate = &delegate
}

// Listen -
func (transport *HTTPTransport) Listen() {
	if transport.server != nil {

		transport.server.SetLowLevelHandlers([]Servers.LowLevelHandler{
			Servers.LowLevelHandler{
				Route:   "/disgover",
				Handler: func(rw http.ResponseWriter, r *http.Request) { disgoverRequestHandler(transport, rw, r) },
				Verb:    "POST",
			},
		})

		go transport.server.Run(fmt.Sprintf("%s:%d", transport.Endpoint.Host, transport.Endpoint.Port))
	}
}

func disgoverRequestHandler(transport *HTTPTransport, rw http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	result, _ := (*transport.OnPeerRPCDelegate)(body)

	rw.Write(result)
}
