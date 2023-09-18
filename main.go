// main.go
package main

import (
	"fmt"
	"time"
	"github.com/asynkron/protoactor-go/actor"
	// accountパッケージをインポート
	"github.com/tkhrk1010/go-actor-bank/account"
	// approvalパッケージをインポート
	"github.com/tkhrk1010/go-actor-bank/approval"
)

func main() {
	// Actorシステムを作成
	system := actor.NewActorSystem()

	// ApprovalActorを起動
	approvalProps := actor.PropsFromProducer(func() actor.Actor {
		return approval.NewApprovalActor()
	})
	approvalActor, _ := system.Root.SpawnNamed(approvalProps, "approval")

	// AccountActorを起動
	accountProps := actor.PropsFromProducer(func() actor.Actor {
		return account.NewAccountActor(approvalActor)
	})
	accountActor, _ := system.Root.SpawnNamed(accountProps, "account")

	// 引き出しリクエストを送信
	withdrawRequest := &account.WithdrawRequest{Amount: 300}
	response, err := system.Root.RequestFuture(accountActor, withdrawRequest, 5*time.Second).Result()
	if err != nil {
		fmt.Println("引き出しリクエストエラー:", err)
		return
	}

	// 引き出しリクエストの結果を処理
	switch resp := response.(type) {
	case *account.WithdrawResponse:
		if resp.Approved {
			fmt.Println("引き出しが承認されました")
		} else {
			fmt.Println("引き出しが拒否されました")
		}
	}

	// Actorシステムを終了
	system.Shutdown()
}
