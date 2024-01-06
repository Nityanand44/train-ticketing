package main

import (
	"context"
	"google.golang.org/grpc"
	"log"
	pb "train-ticketing"
)

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to server: %v", err)
	}
	defer conn.Close()

	client := pb.NewTrainServiceClient(conn)

	purchaseTicket(client)
	getReceiptDetails(client)
	getSectionDetails(client)
	modifyUserSeat(client)
	removeUser(client)

}

func purchaseTicket(client pb.TrainServiceClient) {
	req := &pb.TicketRequest{
		From: "London",
		To:   "France",
		User: &pb.User{
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john.doe@example.com",
		},
	}

	receipt, err := client.PurchaseTicket(context.Background(), req)
	if err != nil {
		log.Fatalf("Failed to purchase ticket: %v", err)
	}

	log.Printf("Purchase successful. Receipt: from:\"%s\" to:\"%s\" user:%+v price_paid:$%v seat:\"%s\"",
		req.From, req.To, req.User, receipt.PricePaid, receipt.Seat)
}

func getReceiptDetails(client pb.TrainServiceClient) {
	user := &pb.User{
		Seat: "01",
	}

	receipt, err := client.GetReceiptDetails(context.Background(), user)
	if err != nil {
		log.Fatalf("Failed to get receipt details: %v", err)
	}

	if receipt != nil {
		log.Printf("Receipt details: from:\"%s\" to:\"%s\" user:%+v price_paid:$%v seat:\"%s\"",
			receipt.From, receipt.To, receipt.User, receipt.PricePaid, receipt.Seat)
	} else {
		log.Println("Receipt not found.")
	}
}

func getSectionDetails(client pb.TrainServiceClient) {
	section := &pb.SectionRequest{
		Section: "A",
	}

	sectionDetails, err := client.GetSectionDetails(context.Background(), section)
	if err != nil {
		log.Fatalf("Failed to get section details: %v", err)
	}

	log.Printf("Section A details: %+v", sectionDetails)
}

func modifyUserSeat(client pb.TrainServiceClient) {
	user := &pb.User{
		Seat: "01",
	}

	modifyRequest := &pb.ModifySeatRequest{
		User:    user,
		NewSeat: "05",
	}

	receipt, err := client.ModifyUserSeat(context.Background(), modifyRequest)
	if err != nil {
		log.Fatalf("Failed to modify user seat: %v", err)
	}

	log.Printf("User seat modified. New receipt details: from:\"%s\" to:\"%s\" user:%+v price_paid:$%v seat:\"%s\"",
		receipt.From, receipt.To, receipt.User, receipt.PricePaid, receipt.Seat)
}

func removeUser(client pb.TrainServiceClient) {
	user := &pb.User{
		Seat: "05",
	}

	removeResponse, err := client.RemoveUser(context.Background(), user)
	if err != nil {
		log.Fatalf("Failed to remove user: %v", err)
	}

	if removeResponse.Success {
		log.Println("User removed successfully.")
	} else {
		log.Println("User removal failed.")
	}
}
