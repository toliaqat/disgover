package disgover

import (
	"sync"
)

// DisGoverService
type DisGoverService struct {
	running bool
}

// NewDisGoverService
func NewDisGoverService() *DisGoverService {
	disGoverService := DisGoverService{false}
	return &disGoverService
}

// Name
func (disGoverService *DisGoverService) Name() string {
	return "DisGoverService"
}

// IsRunning
func (disGoverService *DisGoverService) IsRunning() bool {
	return disGoverService.running
}

// Go
func (disGoverService *DisGoverService) Go(waitGroup *sync.WaitGroup) {
	disGoverService.running = true
}

func (disGoverService *DisGoverService) WithGrpc() *DisGoverService {
	disGoverService.RegisterGrpc()
	return disGoverService
}

//TODO: depending on what tranport layer you are jusing call the approopriate implementation:
//will need some flags set (in the WithProtocol function)


