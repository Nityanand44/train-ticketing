package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"

	"google.golang.org/grpc"
	pb "train-ticketing"
)

type trainServer struct {
	mu    sync.Mutex
	users map[string]*pb.Receipt
	pb.UnimplementedTrainServiceServer
}

func (s *trainServer) PurchaseTicket(ctx context.Context, req *pb.TicketRequest) (*pb.Receipt, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	price := float32(20.0)
	seat := fmt.Sprintf("%02d", len(s.users)+1)

	receipt := &pb.Receipt{
		From:      req.From,
		To:        req.To,
		User:      req.User,
		PricePaid: price,
		Seat:      seat,
	}

	s.users[seat] = receipt

	return receipt, nil
}

func (s *trainServer) GetReceiptDetails(ctx context.Context, req *pb.User) (*pb.Receipt, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	seat := req.Seat
	receipt, exists := s.users[seat]
	if !exists {
		return nil, fmt.Errorf("user with seat %s not found", seat)
	}

	return receipt, nil
}

func (s *trainServer) GetSectionDetails(ctx context.Context, req *pb.SectionRequest) (*pb.SectionDetails, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	section := req.Section
	usersInSection := make([]*pb.User, 0)

	for seat, receipt := range s.users {
		if getSection(seat) == section {
			usersInSection = append(usersInSection, receipt.User)
		}
	}

	return &pb.SectionDetails{Users: usersInSection}, nil
}

func (s *trainServer) ModifyUserSeat(ctx context.Context, req *pb.ModifySeatRequest) (*pb.Receipt, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	seat := req.User.Seat
	receipt, exists := s.users[seat]
	if !exists {
		return nil, fmt.Errorf("user with seat %s not found", seat)
	}

	newSeat := req.NewSeat
	receipt.Seat = newSeat
	s.users[newSeat] = receipt
	delete(s.users, seat)

	return receipt, nil
}

func (s *trainServer) RemoveUser(ctx context.Context, req *pb.User) (*pb.RemoveUserResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	seat := req.Seat
	_, exists := s.users[seat]
	if !exists {
		return nil, fmt.Errorf("user with seat %s not found", seat)
	}

	delete(s.users, seat)

	return &pb.RemoveUserResponse{Success: true}, nil
}

func getSection(seat string) string {
	if seat[0] < 'M' {
		return "A"
	}
	return "B"
}

func main() {
	server := &trainServer{users: make(map[string]*pb.Receipt)}
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	srv := grpc.NewServer()
	pb.RegisterTrainServiceServer(srv, server)

	log.Println("Server is listening on port 50051...")
	if err := srv.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
