package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/suapapa/go_kioom"
)

func main() {
	appKey := os.Getenv("KIOOM_APP_KEY")
	secretKey := os.Getenv("KIOOM_SECRET_KEY")

	if appKey == "" || secretKey == "" {
		log.Fatal("KIOOM_APP_KEY and KIOOM_SECRET_KEY environment variables are required.")
	}

	client := kioom.NewClient(appKey, secretKey, true) // 모의투자

	// 토큰 발급
	_, err := client.IssueToken()
	if err != nil {
		log.Fatalf("토큰 발급 실패: %v", err)
	}

	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	fmt.Println("종목 실시간 순위 변동 모니터링 (30초 간격)...")
	fmt.Println("종료하려면 Ctrl+C를 누르세요.")

	for {
		select {
		case <-ticker.C:
			// "5": 30초 실시간 조회 순위
			rankRes, err := client.GetRealtimeStockRank("5")
			if err != nil {
				log.Printf("⚠️ 조회 실패: %v", err)
				continue
			}

			fmt.Printf("\n[%s] 30초 실시간 조회 순위 TOP 5\n", time.Now().Format("15:04:05"))
			for i, item := range rankRes.ItemInqRank {
				if i >= 5 {
					break
				}
				change := item.RankChg
				if item.RankChgSign == "2" {
					change = "▲" + change
				} else if item.RankChgSign == "5" {
					change = "▼" + change
				} else {
					change = "-" + change
				}
				fmt.Printf("%d. %-12s (%s) [%s]\n", i+1, item.StkNm, item.StkCd, change)
			}
		}
	}
}
