# 📊 주식 차트 시각화 스크립트 (Kiwoom Chart Plotter)

이 디렉토리는 Go 클라이언트(`go_kioom`)가 조회한 주식 차트 데이터를 아름답고 전문적인 캔들스틱(봉차트) 이미지로 시각화해 주는 Python 스크립트를 포함하고 있습니다.

---

## ✨ 주요 특징

1. **모든 차트 타입 지원**: 틱(Tick), 분봉(Minute), 일봉(Daily), 주봉(Weekly), 월봉(Monthly), 년봉(Yearly) 등 모든 주식 차트 인터벌을 완벽히 지원합니다.
2. **자동 차트 탐색**: YAML 파일 내의 루트 키 구조(`stkdtpolechartqry`, `stkminpolechartqry` 등)를 분석하여 차트 종류를 자동으로 인식하고 알맞은 라벨 및 시간 형식을 적용합니다.
3. **프리미엄 다크 테마 디자인**: 금융 트레이딩 소프트웨어 수준의 고품질 어두운 테마 환경(`#0b0c10`, `#12131c` 배경)과 부드러운 격자선(Grid)을 제공합니다.
4. **이동평균선(MA) 및 거래량 시각화**: 
   - 가격 차트 위에 5일, 20일, 60일, 120일 이동평균선(MA)을 깔끔하게 표현합니다.
   - 하단 서브플롯에 거래량 막대그래프와 거래량 5일 이동평균선(Vol MA5)을 함께 시각화합니다.
5. **KRX 표준 색상 대응**: 한국 거래소(KRX) 표준에 맞춰 상승(양봉)은 빨간색(`#f03e3e`), 하락(음봉)은 파란색(`#1971c2`)으로 표현합니다.
6. **최신 종목 정보 및 가격 변동 표시**: 차트 우측에 최근 거래일의 종가, 전일 대비 변동폭 및 등락률을 깔끔한 뱃지(Badge) 형태로 렌더링합니다.

---

## 🛠️ 사전 준비 및 설치

이 스크립트는 `YAML` 데이터 파싱을 위해 `PyYAML`, 데이터 정제를 위해 `pandas`, 시각화를 위해 `matplotlib`를 사용합니다.

### 1. 일반적인 설치 방법 (`pip` 사용)

필요한 라이브러리를 아래 명령어로 설치합니다.
```bash
pip install pyyaml pandas matplotlib
```

### 2. 모던 Python 도구 (`uv` 사용)

본 스크립트는 PEP 723 인라인 스크립트 메타데이터를 포함하고 있어, 최신 Python 패키지 관리자인 `uv`가 설치되어 있다면 의존성을 수동으로 설치하지 않고도 즉시 실행할 수 있습니다.
```bash
# uv를 사용해 필요한 패키지를 자동으로 준비하고 실행
uv run script/plot_chart.py script/samsung_chart_day.yaml
```

---

## 📥 1 단계: 차트 데이터 YAML 생성 방법

시각화 스크립트를 실행하려면 `go_kioom` 클라이언트를 사용하여 차트 데이터를 YAML 형태로 저장해야 합니다. `cmd/examples/05_chart` 예제 프로그램은 `-yaml` 플래그를 지원하여 터미널 표준 출력으로 YAML 데이터를 내보낼 수 있습니다.

### 차트 종류별 YAML 데이터 생성 예시 (삼성전자: `005930`)

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

---

## 📈 2 단계: 차트 시각화 스크립트 사용법

데이터 준비가 완료되면 다양한 방식으로 스크립트를 실행해 차트를 그릴 수 있습니다.

### 1. 기본 실행 (대화형 창 띄우기)

YAML 데이터를 전달하면 대화형 차트 윈도우가 표시됩니다. 차트 줌인/아웃 및 영역 이동이 가능합니다.
```bash
python3 script/plot_chart.py script/samsung_chart_day.yaml
```

### 2. 파일로 저장하기 (`-o` / `--output` 옵션)

대화형 창을 띄우지 않고 고해상도(300 DPI) 이미지 파일로 직접 내보내고자 할 때 사용합니다.
```bash
python3 script/plot_chart.py script/samsung_chart_day.yaml -o script/samsung_chart_day.png
```

### 3. 봉 표시 개수 지정 (`--days` 옵션)

차트에 표시할 최신 봉(Candlestick)의 개수를 조절할 수 있습니다 (기본값: `120`).
```bash
# 최근 60 거래일/봉 데이터만 시각화
python3 script/plot_chart.py script/samsung_chart_day.yaml --days 60 -o script/samsung_chart_day_60.png
```

### 4. 파이프라인(Pipeline)을 통한 실시간 시각화 (추천)

파일을 저장하는 단계를 거치지 않고, Go 클라이언트 출력을 Python 시각화 도구로 직접 파이프 연결하여 즉시 실행할 수 있습니다. 스크립트의 인자로 `-`를 입력하면 표준 입력(stdin)으로부터 데이터를 읽어옵니다.

* **일봉 즉시 시각화**:
  ```bash
  go run cmd/examples/05_chart/main.go -type day -code 005930 -yaml | python3 script/plot_chart.py -
  ```
* **분봉 차트 실시간 파일 저장**:
  ```bash
  go run cmd/examples/05_chart/main.go -type min -code 005930 -yaml | python3 script/plot_chart.py - -o script/samsung_live_min.png
  ```

---

## ⚙️ CLI 옵션 상세 (CLI Options)

### 1. `plot_chart.py`
```
사용법:
  python3 script/plot_chart.py [yaml_file] [--days DAYS] [-o OUTPUT]

위치 인자 (Positional Arguments):
  yaml_file             파싱할 YAML 차트 데이터 파일 경로 (기본값: script/samsung_chart_day.yaml).
                        '-'를 입력하면 표준 입력(stdin)을 수신합니다.

선택 인자 (Optional Arguments):
  -h, --help            도움말 메시지를 출력하고 종료합니다.
  --days DAYS           차트에 시각화할 최신 trading period/봉의 개수 (기본값: 120).
  -o OUTPUT, --output OUTPUT
                        차트를 대화형 화면에 띄우는 대신, 지정한 이미지 경로로 저장합니다 (예: -o output.png).
```

### 2. `search_ticker.py`
```
사용법:
  python3 script/search_ticker.py query [-q]

위치 인자 (Positional Arguments):
  query                 검색할 종목명 또는 종목명 일부 (대소문자 구분 없음).

선택 인자 (Optional Arguments):
  -h, --help            도움말 메시지를 출력하고 종료합니다.
  -q, --quiet           가장 일치하는 종목의 티커(종목코드)만 출력합니다 (셸 스크립트 연동에 유용).
```

---

## 🔍 종목 코드로 티커 검색 예시 (`search_ticker.py`)

주식 차트를 조회하려면 종목코드(티커)가 필요합니다. `search_ticker.py`를 사용해 종목 이름으로 신속하게 코드를 찾을 수 있습니다.

* **인터랙티브 검색**:
  ```bash
  python3 script/search_ticker.py 삼성
  ```
* **티커 코드만 바로 추출 (스크립트 연동용)**:
  ```bash
  python3 script/search_ticker.py 삼성전자 -q
  # 출력: 005930
  ```

---

## 📂 파일 구조 참고

* **`plot_chart.py`**: 차트를 생성하는 메인 Python 스크립트입니다.
* **`search_ticker.py`**: 종목명으로 티커 코드를 검색하는 스크립트입니다.
* **`*.yaml`**: Go 예제 실행을 통해 생성되는 주식 차트 로우 데이터(YAML) 백업입니다.
* **`*.png`**: 스크립트의 시각화 결과를 내보낸 고해상도 캔들스틱 차트 예시 이미지들입니다.
