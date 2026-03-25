package main

import (
	"fmt"
	"log"
	"os"

	"github.com/suapapa/go_kioom"
)

func main() {
	appKey := os.Getenv("KIOOM_APP_KEY")
	secretKey := os.Getenv("KIOOM_SECRET_KEY")

	if appKey == "" || secretKey == "" {
		log.Fatal("KIOOM_APP_KEY and KIOOM_SECRET_KEY environment variables are required.")
	}

	useMock := true // 모의투자 환경 사용 여부

	// 1. 클라이언트 생성
	client := kioom.NewClient(appKey, secretKey, useMock)

	// 2. 접근토큰 발급 (자동으로 client.Token에 저장됨)
	log.Println("접근토큰 발급을 요청합니다...")
	tokenRes, err := client.IssueToken()
	if err != nil {
		log.Fatalf("토큰 발급 실패: %v", err)
	}
	fmt.Printf("✅ 발급된 토큰: %s\n", tokenRes.Token)
	fmt.Printf("✅ 토큰 유효기간: %s\n", tokenRes.ExpiresDt)

	// 3. 접근토큰 폐기
	log.Println("발급받은 토큰을 폐기합니다...")
	_, err = client.RevokeToken()
	if err != nil {
		log.Printf("⚠️ 토큰 폐기 실패: %v", err)
	} else {
		fmt.Println("✅ 토큰 성공적으로 폐기됨")
	}
}
