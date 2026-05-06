# go_kioom : A Go package for Kiwoom REST API

[![Go Reference](https://pkg.go.dev/badge/github.com/suapapa/go_kioom.svg)](https://pkg.go.dev/github.com/suapapa/go_kioom)

`go_kioom` 은 키움증권 REST API 명세서에 맞춰 작성된 비공식 Golang 래퍼(Wrapper) 패키지입니다. 

## 구현 현황

Kiwoom REST API의 전체 구현 목록 및 현재 진행 상태는 [IMPLEMENTATION_STATUS.md](IMPLEMENTATION_STATUS.md)에서 확인하실 수 있습니다.

## 사용 방법 (예제)

```go
package main

import (
	"fmt"
	"log"

	"github.com/suapapa/go_kioom" // 올바른 임포트 경로로 수정하여 사용하세요
)

func main() {
	appKey := "YOUR_APP_KEY"
	secretKey := "YOUR_SECRET_KEY"
	useMock := true // 모의투자 환경 사용 여부

	// 1. 클라이언트 생성
	client := kioom.NewClient(appKey, secretKey, useMock)

	// 2. 접근토큰 발급 (자동으로 client.Token에 저장됨)
	tokenRes, err := client.IssueToken()
	if err != nil {
		log.Fatalf("토큰 발급 실패: %v", err)
	}
	fmt.Printf("발급된 토큰: %s\n", tokenRes.Token)

	// 3. 계좌번호 조회
	acctRes, err := client.GetAccountNumber()
	if err != nil {
		log.Printf("계좌번호 조회 실패: %v", err)
	} else {
		fmt.Printf("계좌번호: %s\n", acctRes.AcctNo)
	}

	// 4. 삼성전자 주식 기본정보 조회
	stockRes, err := client.GetStockBasicInfo("005930")
	if err != nil {
		log.Printf("주식정보 조회 실패: %v", err)
	} else {
		fmt.Printf("주식명: %s, 현재가: %s\n", stockRes.StkNm, stockRes.CurPrc)
	}

	// 5. 접근토큰 폐기
	_, err = client.RevokeToken()
	if err != nil {
		log.Printf("토큰 폐기 실패: %v", err)
	} else {
		fmt.Println("토큰 성공적으로 폐기됨")
	}
}
```

## 예제 (Examples)

더 많은 예제는 [cmd/examples/](./cmd/examples/) 디렉토리에서 확인할 수 있습니다.
- `01_auth`: 토큰 발급 및 폐기
- `02_basic_info`: 주식 기본 정보 및 실시간 검색 순위 조회
- `03_monitor`: 계좌 현황 모니터링
- `04_trading`: 주식/신용 매수, 매도, 정정 및 취소 주문

## CLI (`kioom`)

자세한 내용은 [cmd/kioom/README.md](cmd/kioom/README.md)를 참고하세요.

## 확장 가이드
`kioom` 패키지는 손쉽게 새로운 API 엔드포인트를 추가할 수 있도록 `Client` 구조체 내부에 공통 HTTP 요청 로직(`c.newRequest`, `c.do`)을 구현해 두었습니다. 공식 문서를 참고하여 추가하고자 하는 API의 Request, Response 구조체를 만들고, 메서드를 작성해주시면 됩니다.
