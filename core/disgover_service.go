package core

import (
	"sync"
	"context"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"github.com/dispatchlabs/disgover/proto"
	"github.com/dispatchlabs/disgo_commons/services"
)

// DisGoverService
type DisGoverService struct {
	running bool
}

// NewDisGoverService
func NewDisGoverService() *DisGoverService {
	disGoverService := DisGoverService{false}
	proto.RegisterDisGoverGrpcServer(services.GetGrpcServer(), &disGoverService)
	return &disGoverService
}

// Init
func (disGoverService *DisGoverService) Init() {
	log.WithFields(log.Fields{
		"method": "DisGoverService.Init",
	}).Info("initializing...")
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

// PeerPing
func (disGoverService *DisGoverService) PeerPing(context.Context, *proto.Contact) (*proto.Contact, error) {
	return nil, nil
}

// PeerFind
func (disGoverService *DisGoverService) PeerFind(context.Context, *proto.FindRequest) (*proto.Contact, error) {
	return nil, nil
}
