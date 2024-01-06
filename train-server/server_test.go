package main

import (
	"testing"
	pb "train-ticketing"

	"github.com/stretchr/testify/assert"
)

func TestPurchaseTicket(t *testing.T) {
	server := &trainServer{users: make(map[string]*pb.Receipt)}

	purchaseRequest := &pb.TicketRequest{
		From: "London",
		To:   "France",
		User: &pb.User{
			FirstName: "John",
			LastName:  "Doe",
			Email:     "john.doe@example.com",
		},
	}

	receipt, err := server.PurchaseTicket(nil, purchaseRequest)

	assert.NoError(t, err)
	assert.NotNil(t, receipt)
	assert.Equal(t, purchaseRequest.From, receipt.From)
	assert.Equal(t, purchaseRequest.To, receipt.To)
	assert.Equal(t, purchaseRequest.User, receipt.User)
	assert.Equal(t, float32(20.0), receipt.PricePaid)
	assert.NotEmpty(t, receipt.Seat)
	assert.Contains(t, server.users, receipt.Seat)
}

func TestGetReceiptDetailsNotFound(t *testing.T) {
	server := &trainServer{users: make(map[string]*pb.Receipt)}

	user := &pb.User{Seat: "NonExistentSeat"}

	receipt, err := server.GetReceiptDetails(nil, user)

	assert.Error(t, err)
	assert.Nil(t, receipt)
	assert.Contains(t, err.Error(), "user with seat NonExistentSeat not found")
}

func TestModifyUserSeat(t *testing.T) {
	server := &trainServer{users: make(map[string]*pb.Receipt)}

	initialUser := &pb.User{
		FirstName: "John",
		LastName:  "Doe",
		Seat:      "01",
	}

	initialReceipt := &pb.Receipt{
		From:      "London",
		To:        "France",
		User:      initialUser,
		PricePaid: 20.0,
		Seat:      "01",
	}

	server.users[initialUser.Seat] = initialReceipt

	modifyRequest := &pb.ModifySeatRequest{
		User:    &pb.User{Seat: "01"},
		NewSeat: "B1",
	}

	modifiedReceipt, err := server.ModifyUserSeat(nil, modifyRequest)

	assert.NoError(t, err)
	assert.NotNil(t, modifiedReceipt)
	assert.Equal(t, modifyRequest.NewSeat, modifiedReceipt.Seat)
	assert.NotContains(t, server.users, initialUser.Seat)
	assert.Contains(t, server.users, modifiedReceipt.Seat)
}

func TestRemoveUser(t *testing.T) {
	server := &trainServer{users: make(map[string]*pb.Receipt)}

	initialUser := &pb.User{
		FirstName: "John",
		LastName:  "Doe",
		Seat:      "01",
	}

	initialReceipt := &pb.Receipt{
		From:      "London",
		To:        "France",
		User:      initialUser,
		PricePaid: 20.0,
		Seat:      "01",
	}

	server.users[initialUser.Seat] = initialReceipt

	removeRequest := &pb.User{Seat: "01"}

	removeResponse, err := server.RemoveUser(nil, removeRequest)

	assert.NoError(t, err)
	assert.NotNil(t, removeResponse)
	assert.True(t, removeResponse.Success)
	assert.NotContains(t, server.users, initialUser.Seat)
}
