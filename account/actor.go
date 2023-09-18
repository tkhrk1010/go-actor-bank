package account

import (
	"time"
	"log"

	"github.com/asynkron/protoactor-go/actor"
)

// 口座
type AccountActor struct {
	balance float64
	approval *actor.PID
}

func NewAccountActor(balance float64, approval *actor.PID) *AccountActor {
	return &AccountActor{
		balance: balance,
		approval: approval,
	}
}

type WithdrawRequest struct {
	Amount float64
	Sender *actor.PID
	UserID  string
}

type WithdrawResponse struct {
	Approved bool
}

func (a *AccountActor) Receive(context actor.Context) {
    switch msg := context.Message().(type) {
    case *WithdrawRequest:
				// log.Printf("Account balance = %.2f", a.balance)
        // log.Printf("Received withdrawal request: Amount = %.2f", msg.Amount)

				// 送信元のAccountActorのPIDを設定
				msg.Sender = context.Self()
        // 引き出しリクエストを審査アクターに送信
        response, err := context.RequestFuture(a.approval, msg, 5*time.Second).Result()
        if err != nil {
            context.Respond(&WithdrawResponse{Approved: false})

						log.Printf("user: %v denied due to timeout or error: %v", msg.UserID, err)
            return
        }

        // 審査アクターからの応答を処理
        switch resp := response.(type) {
        case *WithdrawResponse:
            if resp.Approved {
                if msg.Amount <= a.balance {
                    a.balance -= msg.Amount
                    context.Respond(&WithdrawResponse{Approved: true})

                    // log.Println("Withdrawal request approved")
                } else {
                    context.Respond(&WithdrawResponse{Approved: false})

										log.Printf("user: %v failed due to insufficient balance", msg.UserID)
                }
            } else {
                context.Respond(&WithdrawResponse{Approved: false})

								log.Printf("user: %v denied approval", msg.UserID)
            }
        }
    }
}