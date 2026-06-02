package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	kioom "github.com/suapapa/go_kioom"
	"github.com/suapapa/go_kioom/internal/kioomenv"
	"gopkg.in/yaml.v3"
)

var yamlOut bool

func statusPrintf(format string, a ...any) {
	if yamlOut {
		fmt.Fprintf(os.Stderr, format, a...)
	} else {
		fmt.Printf(format, a...)
	}
}

func main() {
	// 1. 환경 변수 로드 (KIOOM_APP_KEY, KIOOM_SECRET_KEY)
	cfg := kioomenv.Load(os.Getenv)
	if err := cfg.RequireAppKeys(); err != nil {
		log.Fatal(err)
	}

	// 2. 플래그 선언 및 분석
	codeFlag := flag.String("code", "", "종목 코드 (6자리)")
	codeShort := flag.String("c", "", "종목 코드 (6자리) (약칭)")
	typeFlag := flag.String("type", "", "차트 종류 (tick, min, day, week, month, year)")
	typeShort := flag.String("t", "", "차트 종류 (tick, min, day, week, month, year) (약칭)")
	yamlFlag := flag.Bool("yaml", false, "YAML 형식으로 출력")
	yamlShort := flag.Bool("y", false, "YAML 형식으로 출력 (약칭)")
	mockFlag := flag.Bool("mock", true, "모의투자 도메인 사용 여부")
	mockShort := flag.Bool("m", false, "모의투자 도메인 사용 여부 (약칭)")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "사용법:\n")
		fmt.Fprintf(os.Stderr, "  go run cmd/examples/05_chart/main.go [options]\n\n")
		fmt.Fprintf(os.Stderr, "옵션:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\n차트 종류 (type/t):\n")
		fmt.Fprintf(os.Stderr, "  tick  - 틱 차트 (1틱)\n")
		fmt.Fprintf(os.Stderr, "  min   - 분봉 차트 (1분봉)\n")
		fmt.Fprintf(os.Stderr, "  day   - 일봉 차트 (기본값)\n")
		fmt.Fprintf(os.Stderr, "  week  - 주봉 차트\n")
		fmt.Fprintf(os.Stderr, "  month - 월봉 차트\n")
		fmt.Fprintf(os.Stderr, "  year  - 년봉 차트\n")
	}

	flag.Parse()

	// 값 분석
	stkCd := "005930" // 기본값 삼성전자
	if *codeFlag != "" {
		stkCd = *codeFlag
	} else if *codeShort != "" {
		stkCd = *codeShort
	}

	chartType := "day" // 기본값
	if *typeFlag != "" {
		chartType = *typeFlag
	} else if *typeShort != "" {
		chartType = *typeShort
	}

	yamlOut = *yamlFlag || *yamlShort
	useMock := *mockFlag && *mockShort

	chartType = strings.ToLower(chartType)
	switch chartType {
	case "tick", "틱", "틱차트":
		chartType = "tick"
	case "min", "minute", "분봉":
		chartType = "min"
	case "day", "daily", "일봉":
		chartType = "day"
	case "week", "weekly", "주봉":
		chartType = "week"
	case "month", "monthly", "월봉":
		chartType = "month"
	case "year", "yearly", "년봉", "연봉":
		chartType = "year"
	default:
		fmt.Fprintf(os.Stderr, "에러: 알 수 없는 차트 종류: %s\n\n", chartType)
		flag.Usage()
		os.Exit(1)
	}

	ctx := context.Background()

	// 3. Kiwoom REST API 클라이언트 초기화
	statusPrintf("클라이언트를 초기화합니다. (차트: %s, 종목코드: %s, 모의투자: %t)\n", chartType, stkCd, useMock)
	var opts []kioom.Option
	if useMock {
		opts = append(opts, kioom.WithMockDomain())
	}
	client := kioom.NewClient(cfg.AppKey, cfg.SecretKey, opts...)

	// 4. 접근토큰 발급
	if yamlOut {
		log.SetOutput(os.Stderr)
	}
	log.Println("접근토큰 발급을 요청합니다...")
	_, err := client.IssueToken(ctx)
	if err != nil {
		log.Fatalf("토큰 발급 실패: %v", err)
	}
	log.Println("접근토큰 발급 성공!")

	// 5. 차트 종류에 따른 쿼리 및 출력
	switch chartType {
	case "tick":
		queryTickChart(ctx, client, stkCd)
	case "min":
		queryMinuteChart(ctx, client, stkCd)
	case "day":
		queryChart(ctx, client, stkCd, "일봉", client.GetStockDailyChart)
	case "week":
		queryChart(ctx, client, stkCd, "주봉", client.GetStockWeeklyChart)
	case "month":
		queryChart(ctx, client, stkCd, "월봉", client.GetStockMonthlyChart)
	case "year":
		queryYearlyChart(ctx, client, stkCd)
	}
}

func queryTickChart(ctx context.Context, client *kioom.Client, stkCd string) {
	req := &kioom.StockTickChartRequest{
		StkCd:      stkCd,
		TicScope:   "1", // 1틱
		UpdStkpcTp: "1", // 수정주가 적용
	}

	statusPrintf("\n--- [%s] 틱(1틱) 차트 데이터 조회 ---\n", stkCd)
	res, err := client.GetStockTickChart(ctx, req)
	if err != nil {
		log.Fatalf("틱 차트 조회 실패: %v", err)
	}
	if res.ReturnCode != 0 {
		log.Fatalf("API 오류 [%d]: %s", res.ReturnCode, res.ReturnMsg)
	}

	if yamlOut {
		b, err := yaml.Marshal(res)
		if err != nil {
			log.Fatalf("YAML 변환 실패: %v", err)
		}
		os.Stdout.Write(b)
		return
	}

	displayCount := 10
	poles := res.StkTicChartQry
	if len(poles) < displayCount {
		displayCount = len(poles)
	}

	fmt.Printf("총 %d개의 틱 데이터를 수신했습니다. 최근 %d개 데이터를 표시합니다:\n\n", len(poles), displayCount)
	fmt.Printf("%-19s | %-10s | %-10s | %-10s | %-10s | %-12s | %-12s\n",
		"체결 시간 (Time)", "시가 (Open)", "고가 (High)", "저가 (Low)", "종가 (Close)", "전일대비 (Diff)", "거래량 (Volume)")
	fmt.Println("------------------------------------------------------------------------------------------------------------------")

	for i := 0; i < displayCount; i++ {
		p := poles[i]
		formattedTime := p.CntrTm
		if len(p.CntrTm) == 14 {
			t, err := time.Parse("20060102150405", p.CntrTm)
			if err == nil {
				formattedTime = t.Format("2006-01-02 15:04:05")
			}
		}

		signSym := " "
		switch p.PredPreSig {
		case "1", "2":
			signSym = "▲"
		case "4", "5":
			signSym = "▼"
		case "3":
			signSym = "―"
		}

		fmt.Printf("%-19s | %-10s | %-10s | %-10s | %-10s | %s%-9s | %-12s\n",
			formattedTime, p.OpenPric, p.HighPric, p.LowPric, p.CurPrc, signSym, p.PredPre, p.TrdeQty)
	}
	fmt.Println("------------------------------------------------------------------------------------------------------------------")
}

func queryMinuteChart(ctx context.Context, client *kioom.Client, stkCd string) {
	req := &kioom.StockMinuteChartRequest{
		StkCd:      stkCd,
		TicScope:   "1", // 1분봉
		UpdStkpcTp: "1", // 수정주가 적용
	}

	statusPrintf("\n--- [%s] 분봉(1분봉) 차트 데이터 조회 ---\n", stkCd)
	res, err := client.GetStockMinuteChart(ctx, req)
	if err != nil {
		log.Fatalf("분봉 차트 조회 실패: %v", err)
	}
	if res.ReturnCode != 0 {
		log.Fatalf("API 오류 [%d]: %s", res.ReturnCode, res.ReturnMsg)
	}

	if yamlOut {
		b, err := yaml.Marshal(res)
		if err != nil {
			log.Fatalf("YAML 변환 실패: %v", err)
		}
		os.Stdout.Write(b)
		return
	}

	displayCount := 10
	poles := res.StkMinPoleChartQry
	if len(poles) < displayCount {
		displayCount = len(poles)
	}

	fmt.Printf("총 %d개의 분봉 데이터를 수신했습니다. 최근 %d개 데이터를 표시합니다:\n\n", len(poles), displayCount)
	fmt.Printf("%-19s | %-10s | %-10s | %-10s | %-10s | %-12s | %-12s\n",
		"체결 시간 (Time)", "시가 (Open)", "고가 (High)", "저가 (Low)", "종가 (Close)", "전일대비 (Diff)", "거래량 (Volume)")
	fmt.Println("------------------------------------------------------------------------------------------------------------------")

	for i := 0; i < displayCount; i++ {
		p := poles[i]
		formattedTime := p.CntrTm
		if len(p.CntrTm) == 14 {
			t, err := time.Parse("20060102150405", p.CntrTm)
			if err == nil {
				formattedTime = t.Format("2006-01-02 15:04:05")
			}
		}

		signSym := " "
		switch p.PredPreSig {
		case "1", "2":
			signSym = "▲"
		case "4", "5":
			signSym = "▼"
		case "3":
			signSym = "―"
		}

		fmt.Printf("%-19s | %-10s | %-10s | %-10s | %-10s | %s%-9s | %-12s\n",
			formattedTime, p.OpenPric, p.HighPric, p.LowPric, p.CurPrc, signSym, p.PredPre, p.TrdeQty)
	}
	fmt.Println("------------------------------------------------------------------------------------------------------------------")
}

type chartQueryFunc[T any] func(context.Context, *kioom.StockChartRequest) (T, error)

func queryChart[T any](ctx context.Context, client *kioom.Client, stkCd, title string, queryFn chartQueryFunc[T]) {
	req := &kioom.StockChartRequest{
		StkCd:      stkCd,
		BaseDt:     time.Now().Format("20060102"),
		UpdStkpcTp: "1", // 수정주가 적용
	}

	statusPrintf("\n--- [%s] %s 차트 데이터 조회 ---\n", stkCd, title)
	res, err := queryFn(ctx, req)
	if err != nil {
		log.Fatalf("%s 차트 조회 실패: %v", title, err)
	}

	var poles []kioom.StockChartPole
	var returnCode int
	var returnMsg string

	switch r := any(res).(type) {
	case *kioom.StockDailyChartResponse:
		poles = r.StkDtPoleChartQry
		returnCode = r.ReturnCode
		returnMsg = r.ReturnMsg
	case *kioom.StockWeeklyChartResponse:
		poles = r.StkStkPoleChartQry
		returnCode = r.ReturnCode
		returnMsg = r.ReturnMsg
	case *kioom.StockMonthlyChartResponse:
		poles = r.StkMthPoleChartQry
		returnCode = r.ReturnCode
		returnMsg = r.ReturnMsg
	}

	if returnCode != 0 {
		log.Fatalf("API 오류 [%d]: %s", returnCode, returnMsg)
	}

	if yamlOut {
		b, err := yaml.Marshal(res)
		if err != nil {
			log.Fatalf("YAML 변환 실패: %v", err)
		}
		os.Stdout.Write(b)
		return
	}

	printPoles(title, poles)
}

func queryYearlyChart(ctx context.Context, client *kioom.Client, stkCd string) {
	req := &kioom.StockChartRequest{
		StkCd:      stkCd,
		BaseDt:     time.Now().Format("20060102"),
		UpdStkpcTp: "1",
	}

	statusPrintf("\n--- [%s] 년봉 차트 데이터 조회 ---\n", stkCd)
	res, err := client.GetStockYearlyChart(ctx, req)
	if err != nil {
		log.Fatalf("년봉 차트 조회 실패: %v", err)
	}
	if res.ReturnCode != 0 {
		log.Fatalf("API 오류 [%d]: %s", res.ReturnCode, res.ReturnMsg)
	}

	if yamlOut {
		b, err := yaml.Marshal(res)
		if err != nil {
			log.Fatalf("YAML 변환 실패: %v", err)
		}
		os.Stdout.Write(b)
		return
	}

	printPoles("년봉", res.StkYrPoleChartQry)
}

func printPoles(title string, poles []kioom.StockChartPole) {
	displayCount := 10
	if len(poles) < displayCount {
		displayCount = len(poles)
	}

	fmt.Printf("총 %d개의 %s 데이터를 수신했습니다. 최근 %d개 데이터를 표시합니다:\n\n", len(poles), title, displayCount)
	fmt.Printf("%-10s | %-10s | %-10s | %-10s | %-10s | %-12s | %-12s\n",
		"일자 (Date)", "시가 (Open)", "고가 (High)", "저가 (Low)", "종가 (Close)", "전일대비 (Diff)", "거래량 (Volume)")
	fmt.Println("------------------------------------------------------------------------------------------------------------------")

	for i := 0; i < displayCount; i++ {
		p := poles[i]
		formattedDate := p.Dt
		if len(p.Dt) == 8 {
			t, err := time.Parse("20060102", p.Dt)
			if err == nil {
				formattedDate = t.Format("2006-01-02")
			}
		}

		signSym := " "
		switch p.PredPreSig {
		case "1", "2":
			signSym = "▲"
		case "4", "5":
			signSym = "▼"
		case "3":
			signSym = "―"
		}

		diffStr := signSym + p.PredPre
		if p.PredPre == "" {
			diffStr = "N/A"
		}

		fmt.Printf("%-10s | %-10s | %-10s | %-10s | %-10s | %-12s | %-12s\n",
			formattedDate, p.OpenPric, p.HighPric, p.LowPric, p.CurPrc, diffStr, p.TrdeQty)
	}
	fmt.Println("------------------------------------------------------------------------------------------------------------------")
}
