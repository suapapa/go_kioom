# CLI (`kioom`)

에이전트/자동화 환경에서 사용하기 쉬운 CLI를 제공합니다.

### 빌드

```bash
go build -o bin/kioom ./cmd/kioom
```

### 공통 옵션

- `--app-key`: 키움 앱 키 (`KIOOM_APP_KEY` 환경 변수로 대체 가능)
- `--secret-key`: 키움 시크릿 키 (`KIOOM_SECRET_KEY` 환경 변수로 대체 가능)
- `--token`: Bearer 토큰 (`KIOOM_TOKEN` 환경 변수로 대체 가능)
- `--mock`: 모의투자 도메인 사용
- `--output json|pretty`: 출력 형식 (기본 `json`)

### 에이전트 친화 규약

- 요청 입력은 `--json` 플래그로 전달 (API request struct와 1:1 매핑)
- 응답은 항상 JSON envelope 형식
  - 성공: `{"ok":true,"data":...}`
  - 실패: `{"ok":false,"error":{"code":"...","message":"..."}}`
- `schema` 명령으로 런타임 스키마 조회 가능

### 사용 예시

```bash
# 토큰 발급 (모의투자)
kioom --mock auth issue

# 계좌번호 조회
kioom --mock --token "$KIOOM_TOKEN" account number

# 주식 기본정보 조회
kioom --mock stock basic --json '{"stk_cd":"005930"}'

# 주문 매수
kioom --mock order buy --json '{"dmst_stex_tp":"KRX","stk_cd":"005930","ord_qty":"1","ord_uv":"0","trde_tp":"3","cond_uv":""}'

# 특정 명령의 request/response 스키마 확인
kioom schema order.buy
```
