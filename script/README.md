# 주식 차트 시각화 스크립트 (Kiwoom Chart Plotter)

Go 클라이언트(`go_kioom`)로 받은 주식 차트 데이터를 캔들스틱(봉차트) 이미지로 그려 주는 Python 스크립트 모음입니다.

## 주요 특징

1. **모든 차트 타입 지원**: 틱(Tick), 분봉(Minute), 일봉(Daily), 주봉(Weekly), 월봉(Monthly), 년봉(Yearly) 차트를 모두 다룹니다.
2. **자동 차트 탐색**: YAML 루트 키(`stkdtpolechartqry`, `stkminpolechartqry` 등)를 읽어 차트 종류를 알아내고, 맞는 라벨·시간 형식을 적용합니다.
3. **다크 테마**: 금융 트레이딩 UI에 가까운 어두운 배경(`#0b0c10`, `#12131c`)과 격자선(Grid)을 씁니다.
4. **이동평균선(MA) 및 거래량**:
   - 가격 차트에 5·20·60·120일 이동평균선(MA)을 겹쳐 그립니다.
   - 하단 서브플롯에 거래량 막대와 거래량 5일 이동평균(Vol MA5)을 함께 표시합니다.
5. **KRX 표준 색상**: 상승(양봉)은 빨간색(`#f03e3e`), 하락(음봉)은 파란색(`#1971c2`)입니다.
6. **최신 종목 정보**: 차트 우측에 최근 거래일 종가, 전일 대비 변동폭·등락률을 뱃지 형태로 표시합니다.

## 사전 준비 및 설치

`PyYAML`, `pandas`, `matplotlib`가 필요합니다.

### 1. pip로 설치

```bash
pip install pyyaml pandas matplotlib
```

### 2. uv 사용

스크립트에 PEP 723 인라인 메타데이터가 있어 `uv`가 있으면 의존성을 따로 설치하지 않고 바로 실행할 수 있습니다.

```bash
uv run script/plot_chart.py script/samsung_chart_day.yaml
```

## 1단계: 차트 데이터 YAML 만들기

시각화 전에 `go_kioom`으로 차트 데이터를 YAML로 저장해야 합니다. `cmd/examples/05_chart` 예제는 `-yaml` 플래그로 YAML을 stdout에 냅니다.

### 차트 종류별 예시 (삼성전자: `005930`)

* **일봉 (Daily)**
  ```bash
  go run cmd/examples/05_chart/main.go -type day -code 005930 -yaml > script/samsung_chart_day.yaml
  ```
* **분봉 (1분봉, Minute)**
  ```bash
  go run cmd/examples/05_chart/main.go -type min -code 005930 -yaml > script/samsung_chart_min.yaml
  ```
* **틱차트 (1틱, Tick)**
  ```bash
  go run cmd/examples/05_chart/main.go -type tick -code 005930 -yaml > script/samsung_chart_tick.yaml
  ```
* **주봉/월봉/년봉 (Week / Month / Year)**
  ```bash
  go run cmd/examples/05_chart/main.go -type week -code 005930 -yaml > script/samsung_chart_week.yaml
  go run cmd/examples/05_chart/main.go -type month -code 005930 -yaml > script/samsung_chart_month.yaml
  go run cmd/examples/05_chart/main.go -type year -code 005930 -yaml > script/samsung_chart_year.yaml
  ```

## 2단계: 차트 그리기

### 1. 기본 실행 (대화형 창)

YAML 파일을 넘기면 대화형 차트 창이 뜹니다. 줌·패닝이 가능합니다.

```bash
python3 script/plot_chart.py script/samsung_chart_day.yaml
```

### 2. 파일로 저장 (`-o` / `--output`)

창을 띄우지 않고 300 DPI 이미지로 저장합니다.

```bash
python3 script/plot_chart.py script/samsung_chart_day.yaml -o script/samsung_chart_day.png
```

### 3. 봉 개수 지정 (`--days`)

표시할 최신 봉 개수를 바꿉니다 (기본값: `120`).

```bash
python3 script/plot_chart.py script/samsung_chart_day.yaml --days 60 -o script/samsung_chart_day_60.png
```

### 4. 파이프로 바로 실행 (추천)

YAML 파일을 거치지 않고 Go 출력을 Python으로 바로 넘길 수 있습니다. 인자로 `-`를 주면 stdin에서 읽습니다.

* **일봉 즉시 시각화**:
  ```bash
  go run cmd/examples/05_chart/main.go -type day -code 005930 -yaml | python3 script/plot_chart.py -
  ```
* **분봉 실시간 저장**:
  ```bash
  go run cmd/examples/05_chart/main.go -type min -code 005930 -yaml | python3 script/plot_chart.py - -o script/samsung_live_min.png
  ```

## CLI 옵션

### `plot_chart.py`

```
사용법:
  python3 script/plot_chart.py [yaml_file] [--days DAYS] [-o OUTPUT]

위치 인자:
  yaml_file             파싱할 YAML 차트 데이터 파일 (기본값: script/samsung_chart_day.yaml).
                        '-'이면 stdin에서 읽습니다.

선택 인자:
  -h, --help            도움말 출력 후 종료.
  --days DAYS           표시할 최신 trading period/봉 개수 (기본값: 120).
  -o OUTPUT, --output OUTPUT
                        대화형 창 대신 지정 경로에 이미지 저장 (예: -o output.png).
```

### `search_ticker.py`

```
사용법:
  python3 script/search_ticker.py query [-q]

위치 인자:
  query                 검색할 종목명 또는 일부 (대소문자 무시).

선택 인자:
  -h, --help            도움말 출력 후 종료.
  -q, --quiet           가장 일치하는 종목의 티커(종목코드)만 출력 (셸 연동용).
```

## 종목 코드 검색 (`search_ticker.py`)

차트 조회에는 종목코드(티커)가 필요합니다. `search_ticker.py`로 종목명에서 코드를 찾을 수 있습니다.

* **인터랙티브 검색**:
  ```bash
  python3 script/search_ticker.py 삼성
  ```
* **티커만 출력 (스크립트 연동)**:
  ```bash
  python3 script/search_ticker.py 삼성전자 -q
  # 출력: 005930
  ```

## 파일 구조

* **`plot_chart.py`**: 차트를 그리는 메인 스크립트
* **`search_ticker.py`**: 종목명으로 티커 코드 검색
* **`*.yaml`**: Go 예제로 만든 차트 원본 데이터(YAML)
* **`*.png`**: 시각화 결과 예시 이미지
