package main

import (
	"fmt"
	"github.com/google/uuid"
	pb "github.com/richit-ai/notification-system/pb/v1"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"time"
)

func main() {
	var (
		opts []grpc.CallOption
	)
	{
		insecure.NewCredentials()
	}

	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		fmt.Printf("[ERROR][GRPC] connection error...\n")
		return
	}

	client := pb.NewNotificationServiceClient(conn)
	ctx := context.Background()

	uid, _ := uuid.NewUUID()
	dt := time.Now()
	dtStr := dt.Format("01-02-2006T15:04:05Z")

	notification := &pb.Notification{
		Id:        uid.String(),
		UserId:    "1",
		Service:   "R360",
		Type:      "Error|Success|Warning",
		CreatedAt: dtStr,
		Message:   "Message",
		Unread:    false,
	}
	sendNotification, err := client.SendNotification(ctx, &pb.SendNotificationRequest{
		Notification: notification,
	}, opts...)
	if err != nil {
		fmt.Printf("[ERROR][GRPC] response error...\n")
		return
	}

	resp := sendNotification.String()
	fmt.Printf("[GRPC] response... %s\n", resp)
}
