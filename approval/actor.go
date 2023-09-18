package approval

import (
	// "log"
	"github.com/asynkron/protoactor-go/actor"
	// account import
	"github.com/tkhrk1010/go-actor-bank/account"
)

type ApprovalActor struct{
	balances map[string]float64
}

func NewApprovalActor(balances map[string]float64) *ApprovalActor {
	return &ApprovalActor{balances: balances}
}

func (a *ApprovalActor) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *account.WithdrawRequest:
			// log.Printf("ApprovalActor received withdrawal request from AccountActor %v: Amount = %.2f", msg.Sender, msg.Amount)

			// Assuming msg.Sender.String() gives us a unique ID for the user (like "user1", "user2", etc.)
			// You may need to adjust this based on your actual implementation.
			userID := msg.UserID

			// Check if the user exists in the in-memory DB and if the withdrawal amount is permissible
			if balance, exists := a.balances[userID]; exists && balance-msg.Amount >= 0 {
				context.Respond(&account.WithdrawResponse{Approved: true})
			} else {
				context.Respond(&account.WithdrawResponse{Approved: false})
			}
	}
}
