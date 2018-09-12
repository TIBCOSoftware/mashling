package service

import (
	"context"
	"errors"
	"net"
	"testing"

	g "github.com/TIBCOSoftware/mashling/internal/pkg/model/v2/action/service/grpc"
	"github.com/TIBCOSoftware/mashling/internal/pkg/model/v2/types"
	_ "github.com/TIBCOSoftware/mashling/test/gen/grpc/client"
	pb "github.com/TIBCOSoftware/mashling/test/gen/grpc/petstore"
	_ "github.com/TIBCOSoftware/mashling/test/gen/grpc/server"
	"google.golang.org/grpc"
)

// ServerStrct is a stub for your Trigger implementation
type ServerStrct struct {
}

var (
	petMapArr  = make(map[int32]pb.Pet)
	userMapArr = make(map[string]pb.User)
)

// PetById gets a pet by id
func (t *ServerStrct) PetById(ctx context.Context, req *pb.PetByIdRequest) (*pb.PetResponse, error) {
	for _, pet := range petMapArr {
		if pet.Id == req.Id {
			return &pb.PetResponse{Pet: &pet}, nil
		}
	}
	return nil, errors.New("Pet not found")
}

// UserByName gets a user by name
func (t *ServerStrct) UserByName(ctx context.Context, req *pb.UserByNameRequest) (*pb.UserResponse, error) {
	for _, user := range userMapArr {
		if req.Username == user.Username {
			return &pb.UserResponse{User: &user}, nil
		}
	}
	return nil, errors.New("User not found")
}

func TestGRPC(t *testing.T) {
	var petArr = []pb.Pet{
		{
			Id:   2,
			Name: "cat2",
		},
		{
			Id:   3,
			Name: "cat3",
		},
		{
			Id:   4,
			Name: "cat4",
		},
	}
	var userArr = []pb.User{
		{
			Id:       2,
			Username: "user2",
			Email:    "email2",
			Phone:    "phone2",
		},
		{
			Id:       3,
			Username: "user3",
			Email:    "email3",
			Phone:    "phone3",
		},
		{
			Id:       4,
			Username: "user4",
			Email:    "email4",
			Phone:    "phone4",
		},
	}
	for _, pet := range petArr {
		petMapArr[pet.GetId()] = pet
	}
	for _, user := range userArr {
		userMapArr[user.GetUsername()] = user
	}

	addr := ":9000"
	socket, err := net.Listen("tcp", addr)
	if err != nil {
		t.Fatal(err)
	}
	server := grpc.NewServer()
	pb.RegisterPetStoreServiceServer(server, &ServerStrct{})

	done := make(chan bool, 1)
	go func() {
		server.Serve(socket)
		done <- true
	}()
	defer func() {
		server.GracefulStop()
		<-done
	}()

	service := types.Service{
		Type: "grpc",
		Settings: map[string]interface{}{
			"hosturl": "localhost:9000",
		},
	}
	grpcData := map[string]interface{}{
		"methodName":  "PetById",
		"contextdata": context.Background(),
		"reqdata":     &pb.PetByIdRequest{Id: 2},
		"serviceName": "PetStoreService",
		"protoName":   "petstore",
	}
	instance, err := Initialize(service)
	if err != nil {
		t.Fatal(err)
	}
	err = instance.UpdateRequest(map[string]interface{}{
		"grpcMthdParamtrs": grpcData,
	})
	if err != nil {
		t.Fatal(err)
	}
	err = instance.Execute()
	if err != nil {
		t.Fatal(err)
	}
	body := instance.(*g.GRPC).Response.Body
	if pet, ok := body.(*pb.PetResponse); !ok {
		t.Fatal("should be pet response")
	} else if pet.Pet.Name != "cat2" {
		t.Fatal("didn't get correct pet")
	}

	service = types.Service{
		Type: "grpc",
		Settings: map[string]interface{}{
			"hosturl":     "localhost:9000",
			"protoName":   "petstore",
			"serviceName": "PetStoreService",
		},
	}
	instance, err = Initialize(service)
	if err != nil {
		t.Fatal(err)
	}
	err = instance.UpdateRequest(map[string]interface{}{
		"methodName": "PetById",
		"queryParams": map[string]string{
			"id": "2",
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	err = instance.Execute()
	if err != nil {
		t.Fatal(err)
	}
	if a, ok := instance.(*g.GRPC).Response.Body.(map[string]interface{}); !ok {
		t.Fatal("should be a map")
	} else if b := a["pet"]; b == nil {
		t.Fatal("pet should not be nil")
	} else if c, ok := b.(map[string]interface{}); !ok {
		t.Fatal("pet should be a map")
	} else if d := c["name"]; d == nil {
		t.Fatal("name should not be nil")
	} else if e, ok := d.(string); !ok {
		t.Fatal("name should be a string")
	} else if e != "cat2" {
		t.Fatal("name should be equal to cat2")
	}
}
