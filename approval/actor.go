package approval

import (
	"log"
	"github.com/asynkron/protoactor-go/actor"
	// account import
	"github.com/tkhrk1010/go-actor-bank/account"
)

type ApprovalActor struct{}

func NewApprovalActor() *ApprovalActor {
	return &ApprovalActor{}
}

func (a *ApprovalActor) Receive(context actor.Context) {
	switch msg := context.Message().(type) {
	case *account.WithdrawRequest:
			log.Printf("ApprovalActor received withdrawal request from AccountActor %v: Amount = %.2f", msg.Sender, msg.Amount)

			// 引き出しリクエストを処理し、承認または拒否の応答を返す
			// 審査条件
			if msg.Amount <= 1000 { 
					context.Respond(&account.WithdrawResponse{Approved: true}) 
			} else {
					context.Respond(&account.WithdrawResponse{Approved: false}) 
			}
	}
}
