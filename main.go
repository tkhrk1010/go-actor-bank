// main.go
package main

import (
	"log"
	"sync"
	"time"
	"github.com/asynkron/protoactor-go/actor"
	// accountパッケージをインポート
	"github.com/tkhrk1010/go-actor-bank/account"
	// approvalパッケージをインポート
	"github.com/tkhrk1010/go-actor-bank/approval"
)

// in-memory DB for user balances
var userBalances = map[string]float64{
	"user1": 1000,
	"user2": 2000,
	"user3": 3000,
}

func main() {
	system := actor.NewActorSystem()
	rootContext := system.Root

	// Approvalアクターを生成
	approvalProps := actor.PropsFromProducer(func() actor.Actor {
		return approval.NewApprovalActor(userBalances)
	})
	approvalPID := system.Root.Spawn(approvalProps)

	// ユーザーアカウントを生成
	users := map[string]*actor.PID{}
	users["user1"] = system.Root.Spawn(actor.PropsFromProducer(func() actor.Actor {
    return account.NewAccountActor(1000, approvalPID)
	}))
	log.Printf("user1: %v balance = %.2f", users["user1"], userBalances["user1"])

	users["user2"] = system.Root.Spawn(actor.PropsFromProducer(func() actor.Actor {
			return account.NewAccountActor(2000, approvalPID)
	}))
	log.Printf("user2: %v balance = %.2f", users["user2"], userBalances["user2"])

	users["user3"] = system.Root.Spawn(actor.PropsFromProducer(func() actor.Actor {
			return account.NewAccountActor(3000, approvalPID)
	}))
	log.Printf("user3: %v balance = %.2f", users["user3"], userBalances["user3"])

	// 同時に複数の引き出しリクエストを処理
	var wg sync.WaitGroup
	wg.Add(3)

	go func() {
		defer wg.Done()
		result := withdraw(rootContext, users["user1"], 500.0, "user1")
		log.Printf("withdraw user1: result = %v, balance = %.2f\n", result, userBalances["user1"])
	}()

	go func() {
		defer wg.Done()
		result := withdraw(rootContext,users["user2"], 2100.0, "user2")
		log.Printf("withdraw user2: result = %v, balance = %.2f\n", result, userBalances["user2"])
	}()

	go func() {
		defer wg.Done()
		result := withdraw(rootContext, users["user3"], 2000.0, "user3")
		log.Printf("withdraw user3: result = %v, balance = %.2f\n", result, userBalances["user3"])
	}()

	wg.Wait()
}

// 引き出しリクエストを送信して結果を取得
func withdraw(rootContext *actor.RootContext, accountPID *actor.PID, amount float64, userID string) bool {

	withdrawRequest := &account.WithdrawRequest{Amount: amount, UserID: userID}
	log.Printf("Sending Withdraw Request to Actor: %v, Amount: %.2f\n", accountPID, amount)
	response, err := rootContext.RequestFuture(accountPID, withdrawRequest, 5*time.Second).Result()
	if err != nil {
			log.Println("Withdraw request timeout or error:", err)
			return false
	}

	if resp, ok := response.(*account.WithdrawResponse); ok {
		return resp.Approved
	}

	return false
}