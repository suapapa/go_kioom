package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/suapapa/go_kioom"
	"github.com/suapapa/go_kioom/internal/kioomenv"
)

func main() {
	cfg := kioomenv.Load(os.Getenv)
	if err := cfg.RequireAppKeys(); err != nil {
		log.Fatal(err)
	}

	client := kioom.NewClient(cfg.AppKey, cfg.SecretKey, true) // 모의투자

	// 토큰 발급
	_, err := client.IssueToken()
	if err != nil {
		log.Fatalf("토큰 발급 실패: %v", err)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	fmt.Println("종목 실시간 순위 변동 모니터링 (30초 간격)...")
	fmt.Println("종료하려면 Ctrl+C를 누르세요.")

	for {
		select {
		case <-ctx.Done():
			fmt.Println("\n종료합니다.")
			return
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
