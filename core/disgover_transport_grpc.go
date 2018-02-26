package core

import (
	"fmt"
	log "github.com/sirupsen/logrus"

	proto "github.com/dispatchlabs/disgover/proto"
	"github.com/dispatchlabs/disgo_commons/services"
	"github.com/dispatchlabs/disgo_commons/types"
	"github.com/dispatchlabs/disgo_commons/crypto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

func (disGoverService *DisGoverService) RegisterGrpc() *DisGoverService {
	proto.RegisterDisgoverRPCServer(services.GetGrpcServer(), disGoverService)
	return disGoverService
}

// PeerPing
func (disGoverService *DisGoverService) PeerPingGrpc(context.Context, *proto.Contact) (*proto.Contact, error) {
	//convert proto to domain object
	//do call to disgover.go PeerPing
	//take result value and convert back to proto type and return

	return nil, nil
}

// PeerFind
func (disGoverService *DisGoverService) PeerFindGrpc(ctx context.Context, findRequest *proto.FindRequest) (*proto.Contact, error) {
	fmt.Println(fmt.Sprintf("Disgover-TRACE: PeerFind(): %s", findRequest.ContactId))
	idToFind := crypto.ToWalletAddressString(crypto.ToWalletAddress(findRequest.ContactId))

	foundContact, err := disGoverService.PeerFind(idToFind, convertToDomain(findRequest.Sender))
	if err != nil {
		return nil, err
	}
	return convertToProto(foundContact), nil
}

func FindPeerWithGrpcClient(idTofind string, contact *types.Contact) *types.Contact {
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", contact.Endpoint.Host, contact.Endpoint.Port), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("cannot dial server: %v", err)
	}

	client := proto.NewDisgoverRPCClient(conn)
	convertedId, err := crypto.AddressStringToBytes(idTofind)
	if err != nil {
		panic(err)
	}
	response, _ := client.PeerFindGrpc(context.Background(), NewFindRequest(convertedId, contact))
 	if response == nil {
 		log.Error("Could not find desired contact")
 		return nil
	}
	fmt.Println("Disgover-TRACE: findViaPeers() RESULT")
	return convertToDomain(response)
}

func PeerPingWithGrpcClient(contactToPing *types.Contact, thisContact *types.Contact) *types.Contact {
	conn, err := grpc.Dial(fmt.Sprintf("%s:%d", contactToPing.Endpoint.Host, contactToPing.Endpoint.Port), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("cannot dial server: %v", err)
	}
	client := proto.NewDisgoverRPCClient(conn)
	response, _ := client.PeerPingGrpc(context.Background(), convertToProto(thisContact))
	if response == nil {
		log.Error("Could not ping desired contact")
		return nil
	}
	return convertToDomain(response)
}


func NewFindRequest(idTofind []byte, contact *types.Contact) *proto.FindRequest {
	return &proto.FindRequest{
		ContactId: idTofind,
		Sender:    convertToProto(contact),
	}
}

/*
 *  Simple conversion functions from / to proto generated objects and domain level objects
 */
func convertToDomain(sender *proto.Contact) *types.Contact {
	endpoint := &types.Endpoint{
		Host: sender.Endpoint.Host,
		Port: sender.Endpoint.Port,
	}
	contact := types.Contact{
		Address:	crypto.ToWalletAddressString(crypto.ToWalletAddress(sender.Id)),
		Endpoint: 	endpoint,
	}
	return &contact
}

func convertToProto(contact *types.Contact) *proto.Contact  {
	result := &proto.Contact {
		Id:		[]byte(contact.Address),
		Endpoint: &proto.Endpoint{
			Host: contact.Endpoint.Host,
			Port: contact.Endpoint.Port,
		},
	}
	return result
}
