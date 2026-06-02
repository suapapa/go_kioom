package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	kioom "github.com/suapapa/go_kioom"
	"github.com/suapapa/go_kioom/internal/kioomenv"
)

func main() {
	mockFlag := flag.Bool("mock", false, "Use mock domain")
	flag.Parse()

	cfg := kioomenv.Load(os.Getenv)
	if err := cfg.RequireAppKeys(); err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()

	// 1. 클라이언트 생성
	var opts []kioom.Option
	if *mockFlag {
		opts = append(opts, kioom.WithMockDomain())
	}
	client := kioom.NewClient(cfg.AppKey, cfg.SecretKey, opts...)

	// 2. 접근토큰 발급 (자동으로 client.Token에 저장됨)
	log.Println("접근토큰 발급을 요청합니다...")
	tokenRes, err := client.IssueToken(ctx)
	if err != nil {
		log.Fatalf("토큰 발급 실패: %v", err)
	}
	fmt.Printf("✅ 발급된 토큰: %s\n", tokenRes.Token)
	fmt.Printf("✅ 토큰 유효기간: %s\n", tokenRes.ExpiresDt)

	// 3. 접근토큰 폐기
	log.Println("발급받은 토큰을 폐기합니다...")
	_, err = client.RevokeToken(ctx)
	if err != nil {
		log.Printf("⚠️ 토큰 폐기 실패: %v", err)
	} else {
		fmt.Println("✅ 토큰 성공적으로 폐기됨")
	}
}
