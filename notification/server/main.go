package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/go-redis/redis/v9"
	redisDB "github.com/richit-ai/notification-system/internal/db/redis"
	"golang.org/x/net/context"
	"log"
	"net"
	"time"

	pb "github.com/richit-ai/notification-system/pb/v1"
	"google.golang.org/grpc"
)

var (
	port = flag.Int("port", 50051, "The server port")
)

// server is used to implement helloworld.GreeterServer.
type server struct {
	pb.UnsafeNotificationServiceServer
	RedisClient *redis.Client
}

func (s *server) GetUserNotifications(ctx context.Context, request *pb.GetUserNotificationsRequest) (*pb.GetUserNotificationsResponse, error) {
	Notifications := make([]*pb.Notification, 0)

	log.Printf("Received Get User Notifications: %v", request.GetUserId())
	smts := s.RedisClient.Keys(ctx, "/users/"+request.UserId+"/*")
	keys, err := smts.Result()
	if err != nil {
		return nil, err
	}

	for _, key := range keys {
		var notification pb.Notification

		fmt.Println(key)

		val := s.RedisClient.Get(ctx, key)
		result, err := val.Result()
		if err != nil {
			fmt.Println(err)
			continue
			//return nil, err
		}

		err = json.Unmarshal([]byte(result), &notification)
		if err != nil {
			continue
		}

		Notifications = append(Notifications, &notification)
	}

	return &pb.GetUserNotificationsResponse{
		Notifications: Notifications,
	}, nil
}

func (s *server) SendNotification(ctx context.Context, request *pb.SendNotificationRequest) (*pb.SendNotificationResponse, error) {
	log.Printf("Received Send Notification: %v", request.GetNotification())

	key := fmt.Sprintf("/users/%s/%s", request.Notification.UserId, request.Notification.Id)
	data, err := json.Marshal(request.Notification)
	if err != nil {
		return nil, err
	}
	smts := s.RedisClient.Set(ctx, key, data, time.Hour*24*90)
	res := smts.String()
	fmt.Println(res)

	return &pb.SendNotificationResponse{}, nil
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()

	pb.RegisterNotificationServiceServer(s, &server{
		RedisClient: redisDB.InitRedis(),
	})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
