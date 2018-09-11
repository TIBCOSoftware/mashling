// This file registers with grpc service. This file was auto-generated by mashling at
// 2018-09-06 13:05:41.831820537 -0600 MDT m=+0.488161437
package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"

	servInfo "github.com/TIBCOSoftware/mashling/ext/flogo/trigger/grpc"
	pb "github.com/TIBCOSoftware/mashling/test/gen/grpc/petstore"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type serviceImplpetstorePetStoreService struct {
	trigger     *servInfo.GRPCTrigger
	serviceInfo *servInfo.ServiceInfo
}

var serviceInfopetstorePetStoreService = &servInfo.ServiceInfo{
	ProtoName:   "petstore",
	ServiceName: "PetStoreService",
}

func init() {
	servInfo.ServiceRegistery.RegisterServerService(&serviceImplpetstorePetStoreService{serviceInfo: serviceInfopetstorePetStoreService})
}

//RunRegisterServerService registers server method implimentaion with grpc
func (s *serviceImplpetstorePetStoreService) RunRegisterServerService(serv *grpc.Server, trigger *servInfo.GRPCTrigger) {
	service := &serviceImplpetstorePetStoreService{
		trigger:     trigger,
		serviceInfo: serviceInfopetstorePetStoreService,
	}
	pb.RegisterPetStoreServiceServer(serv, service)
}

func (s *serviceImplpetstorePetStoreService) PetById(ctx context.Context, req *pb.PetByIdRequest) (res *pb.PetResponse, err error) {

	methodName := "PetById"

	grpcData := make(map[string]interface{})
	grpcData["methodName"] = methodName
	grpcData["contextdata"] = ctx
	grpcData["reqdata"] = req

	_, replyData, err := s.trigger.CallHandler(grpcData)

	if err != nil {
		log.Println("error: ", err)
	}

	typeHandRes := fmt.Sprintf("%T", replyData)
	if strings.Compare(typeHandRes, "*status.statusError") == 0 {
		return res, replyData.(error)
	}
	typeMethodRes := fmt.Sprintf("%T", res)
	if strings.Compare(typeHandRes, typeMethodRes) == 0 {
		res = replyData.(*pb.PetResponse)
	} else {
		var errValue = replyData.(map[string]interface{})["error"]
		if errValue != nil && len(errValue.(string)) != 0 {
			return res, errors.New(errValue.(string))
		} else {
			rDBytes, err := json.Marshal(replyData)
			if err != nil {
				log.Println("error: ", err)
			}

			err = json.Unmarshal(rDBytes, &res)
			if err != nil {
				log.Println("error: ", err)
			}
		}
	}
	log.Println("response: ", res)

	//User implimentation area

	return res, err
}

func (s *serviceImplpetstorePetStoreService) UserByName(ctx context.Context, req *pb.UserByNameRequest) (res *pb.UserResponse, err error) {

	methodName := "UserByName"

	grpcData := make(map[string]interface{})
	grpcData["methodName"] = methodName
	grpcData["contextdata"] = ctx
	grpcData["reqdata"] = req

	_, replyData, err := s.trigger.CallHandler(grpcData)

	if err != nil {
		log.Println("error: ", err)
	}

	typeHandRes := fmt.Sprintf("%T", replyData)
	if strings.Compare(typeHandRes, "*status.statusError") == 0 {
		return res, replyData.(error)
	}
	typeMethodRes := fmt.Sprintf("%T", res)
	if strings.Compare(typeHandRes, typeMethodRes) == 0 {
		res = replyData.(*pb.UserResponse)
	} else {
		var errValue = replyData.(map[string]interface{})["error"]
		if errValue != nil && len(errValue.(string)) != 0 {
			return res, errors.New(errValue.(string))
		} else {
			rDBytes, err := json.Marshal(replyData)
			if err != nil {
				log.Println("error: ", err)
			}

			err = json.Unmarshal(rDBytes, &res)
			if err != nil {
				log.Println("error: ", err)
			}
		}
	}
	log.Println("response: ", res)

	//User implimentation area

	return res, err
}

func (s *serviceImplpetstorePetStoreService) ServiceInfo() *servInfo.ServiceInfo {
	return s.serviceInfo
}
