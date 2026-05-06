package main

import (
	"fmt"
	"log"
	"os"

	"github.com/suapapa/go_kioom"
	"github.com/suapapa/go_kioom/internal/kioomenv"
)

func main() {
	cfg := kioomenv.Load(os.Getenv)
	if err := cfg.RequireAppKeys(); err != nil {
		log.Fatal(err)
	}

	client := kioom.NewClient(cfg.AppKey, cfg.SecretKey, true) // 모의투자 사용

	// 토큰 발급
	_, err := client.IssueToken()
	if err != nil {
		log.Fatalf("토큰 발급 실패: %v", err)
	}

	// 1. 계좌번호 조회
	fmt.Println("\n--- 계좌 정보 조회 ---")
	acctRes, err := client.GetAccountNumber()
	if err != nil {
		log.Printf("계좌번호 조회 실패: %v", err)
	} else {
		fmt.Printf("계좌번호: %s\n", acctRes.AcctNo)
		fmt.Printf("응답 메시지: %s (%d)\n", acctRes.ReturnMsg, acctRes.ReturnCode)
	}

	// 2. 주식 기본정보 조회 (삼성전자: 005930)
	fmt.Println("\n--- 주식 기본 정보 조회 (삼성전자) ---")
	stockRes, err := client.GetStockBasicInfo("005930")
	if err != nil {
		log.Printf("주식정보 조회 실패: %v", err)
	} else {
		fmt.Printf("종목명: %s (%s)\n", stockRes.StkNm, stockRes.StkCd)
		fmt.Printf("현재가: %s\n", stockRes.CurPrc)
		fmt.Printf("등락율: %s%%\n", stockRes.FluRt)
		fmt.Printf("시가총액: %s\n", stockRes.Mac)
	}

	// 3. 실시간 종목 조회 순위 (당일 누적: 4)
	fmt.Println("\n--- 실시간 종목 조회 순위 (TOP 10) ---")
	rankRes, err := client.GetRealtimeStockRank("4")
	if err != nil {
		log.Printf("실시간 순위 조회 실패: %v", err)
	} else {
		for i, item := range rankRes.ItemInqRank {
			if i >= 10 {
				break
			}
			fmt.Printf("%2d위. %s (%s) - 현재가: %s\n", i+1, item.StkNm, item.StkCd, item.PastCurrPrc)
		}
	}
}
