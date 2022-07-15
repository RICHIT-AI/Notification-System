package main

import (
	"flag"
	"fmt"
	"golang.org/x/net/context"
	"log"
	"net"

	pb "github.com/richit-ai/notification-system/pb/v1"
	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

// server is used to implement helloworld.GreeterServer.
type server struct {
	pb.UnsafeNotificationServiceServer
}

func (s *server) GetUserNotifications(ctx context.Context, request *pb.GetUserNotificationsRequest) (*pb.GetUserNotificationsResponse, error) {
	log.Printf("Received Get User Notifications: %v", request.GetUserId())
	return nil, nil
}

func (s *server) SendNotification(ctx context.Context, request *pb.SendNotificationRequest) (*pb.SendNotificationResponse, error) {
	log.Printf("Received Send Notification: %v", request.GetNotification())

	return &pb.SendNotificationResponse{}, nil
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterNotificationServiceServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
