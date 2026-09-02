// Command genapis generates Kiwoom stock/ETF/US-stock REST API bindings from the official spec.
//
//go:generate go run .
package main

import (
	"encoding/json"
	"fmt"
	"go/format"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"unicode"
)

type ioField struct {
	ItemID string `json:"itemId"`
	ItemNm string `json:"itemNm"`
	Type   string `json:"type"`
	Desc   string `json:"desc"`
}

type apiEntry struct {
	APIID       string    `json:"apiId"`
	APINm       string    `json:"apiNm"`
	Method      string    `json:"method"`
	URL         string    `json:"url"`
	Description string    `json:"description"`
	RequestIO   []ioField `json:"requestIo"`
	ResponseIO  []ioField `json:"responseIo"`
}

var skipManual = map[string]bool{
	"au10001": true, "au10002": true,
	"ka00001": true, "kt00001": true, "kt00018": true,
	"ka10001": true, "ka10008": true, "ka00198": true, "ka10023": true,
	"ka10079": true, "ka10080": true, "ka10081": true, "ka10082": true, "ka10083": true, "ka10094": true,
	"kt10000": true, "kt10001": true, "kt10002": true, "kt10003": true,
	"kt10006": true, "kt10007": true, "kt10008": true, "kt10009": true,
	"ust20000": true, "ust20001": true, "ust20002": true, "ust20003": true,
	"ust21050": true, "ust21070": true, "ust21110": true,
}

var headerFields = map[string]bool{
	"api-id": true, "authorization": true, "cont-yn": true, "next-key": true,
}

func includeAPI(apiID string) bool {
	if skipManual[apiID] {
		return false
	}
	switch {
	case strings.HasPrefix(apiID, "ka300"),
		strings.HasPrefix(apiID, "ka200"),
		strings.HasPrefix(apiID, "ka500"),
		strings.HasPrefix(apiID, "ka501"),
		strings.HasPrefix(apiID, "kt500"),
		strings.HasPrefix(apiID, "ka900"),
		strings.HasPrefix(apiID, "ust313"):
		return false
	}
	if strings.HasPrefix(apiID, "ka100") ||
		strings.HasPrefix(apiID, "ka400") ||
		strings.HasPrefix(apiID, "usa") ||
		strings.HasPrefix(apiID, "ust") ||
		strings.HasPrefix(apiID, "ka001") ||
		strings.HasPrefix(apiID, "ka013") ||
		strings.HasPrefix(apiID, "kt000") {
		return true
	}
	switch apiID {
	case "ka10100", "ka10101", "ka10102", "ka10103":
		return true
	}
	return strings.HasPrefix(apiID, "ka1017")
}

func main() {
	root := findModuleRoot()
	specPath := filepath.Join(root, "_ref", "kiwoom-rest-api-spec.json")
	specBytes, err := os.ReadFile(specPath)
	if err != nil {
		fatal(err)
	}

	var raw map[string]json.RawMessage
	if err := json.Unmarshal(specBytes, &raw); err != nil {
		fatal(err)
	}

	var apis []apiEntry
	for id, body := range raw {
		if !includeAPI(id) {
			continue
		}
		var e apiEntry
		if err := json.Unmarshal(body, &e); err != nil {
			fatal(fmt.Errorf("%s: %w", id, err))
		}
		e.APIID = id
		apis = append(apis, e)
	}
	sort.Slice(apis, func(i, j int) bool { return apis[i].APIID < apis[j].APIID })

	for _, api := range apis {
		if _, ok := apiMethodNames[api.APIID]; !ok {
			fatal(fmt.Errorf("missing method name for %s (%s)", api.APIID, api.APINm))
		}
	}

	if err := writeRegistry(root, apis); err != nil {
		fatal(err)
	}
	if err := writeGeneratedAPIs(root, apis); err != nil {
		fatal(err)
	}
	if err := writeImplementationStatus(root, apis); err != nil {
		fatal(err)
	}

	fmt.Printf("generated %d stock/ETF/US APIs\n", len(apis))
}

func findModuleRoot() string {
	dir, err := os.Getwd()
	if err != nil {
		fatal(err)
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			fatal(fmt.Errorf("go.mod not found"))
		}
		dir = parent
	}
}

func writeRegistry(root string, apis []apiEntry) error {
	var b strings.Builder
	b.WriteString(headerComment())
	b.WriteString("\n\n")
	b.WriteString("package kioom\n\n")
	b.WriteString("import (\n\t\"reflect\"\n\t\"sort\"\n)\n\n")
	b.WriteString("// GeneratedAPI describes a Kiwoom REST endpoint generated from the official spec.\n")
	b.WriteString("type GeneratedAPI struct {\n")
	b.WriteString("\tID          string\n")
	b.WriteString("\tName        string\n")
	b.WriteString("\tDescription string\n")
	b.WriteString("\tMethod      string\n")
	b.WriteString("\tPath        string\n")
	b.WriteString("\tMethodName   string\n")
	b.WriteString("\tRequestType string\n")
	b.WriteString("\tResponseType string\n")
	b.WriteString("}\n\n")
	b.WriteString("// GeneratedAPIRegistry maps Kiwoom TR codes to generated API metadata.\n")
	b.WriteString("var GeneratedAPIRegistry = map[string]GeneratedAPI{\n")
	for _, api := range apis {
		method := methodNameForAPI(api.APIID)
		typeBase := typeBaseFromMethod(method)
		fmt.Fprintf(&b, "\t%q: {ID: %q, Name: %q, Description: %q, Method: %q, Path: %q, MethodName: %q, RequestType: %q, ResponseType: %q},\n",
			api.APIID, api.APIID, api.APINm, api.Description, api.Method, api.URL, method,
			typeBase+"Request", typeBase+"Response")
	}
	b.WriteString("}\n\n")
	b.WriteString("var generatedTypeByName = map[string]reflect.Type{\n")
	for _, api := range apis {
		typeBase := typeBaseFromMethod(methodNameForAPI(api.APIID))
		fmt.Fprintf(&b, "\t%q: reflect.TypeFor[%sRequest](),\n", typeBase+"Request", typeBase)
		fmt.Fprintf(&b, "\t%q: reflect.TypeFor[%sResponse](),\n", typeBase+"Response", typeBase)
	}
	b.WriteString("}\n\n")
	b.WriteString("// GeneratedAPIIDs returns sorted TR codes for generated stock/ETF/US APIs.\n")
	b.WriteString("func GeneratedAPIIDs() []string {\n")
	b.WriteString("\tids := make([]string, 0, len(GeneratedAPIRegistry))\n")
	b.WriteString("\tfor id := range GeneratedAPIRegistry {\n")
	b.WriteString("\t\tids = append(ids, id)\n")
	b.WriteString("\t}\n")
	b.WriteString("\tsort.Strings(ids)\n")
	b.WriteString("\treturn ids\n")
	b.WriteString("}\n")

	return writeGoFile(filepath.Join(root, "api_registry_gen.go"), b.String())
}

func writeGeneratedAPIs(root string, apis []apiEntry) error {
	const maxPerFile = 40
	files := splitAPIs(apis, maxPerFile)
	for i, chunk := range files {
		var b strings.Builder
		b.WriteString(headerComment())
		b.WriteString("\n\n")
		b.WriteString("package kioom\n\n")
		b.WriteString("import (\n\t\"context\"\n\t\"net/http\"\n)\n\n")

	for _, api := range chunk {
			writeAPI(&b, api)
		}

		name := fmt.Sprintf("api_generated_%02d.go", i+1)
		if err := writeGoFile(filepath.Join(root, name), b.String()); err != nil {
			return err
		}
	}
	// remove stale generated files beyond current count
	for i := len(files) + 1; i <= 20; i++ {
		_ = os.Remove(filepath.Join(root, fmt.Sprintf("api_generated_%02d.go", i)))
	}
	return nil
}

func splitAPIs(apis []apiEntry, n int) [][]apiEntry {
	var out [][]apiEntry
	for i := 0; i < len(apis); i += n {
		end := i + n
		if end > len(apis) {
			end = len(apis)
		}
		out = append(out, apis[i:end])
	}
	return out
}

type parsedIO struct {
	scalars []ioField
	lists   []listField
}

type listField struct {
	key    string
	name   string
	fields []ioField
}

func parseIO(fields []ioField, typePrefix string) parsedIO {
	var out parsedIO
	for i := 0; i < len(fields); i++ {
		f := fields[i]
		if headerFields[f.ItemID] {
			continue
		}
		if strings.Contains(strings.ToLower(f.Type), "list") {
			lf := listField{key: f.ItemID, name: listStructName(typePrefix, f.ItemID)}
			for j := i + 1; j < len(fields); j++ {
				if !strings.HasPrefix(fields[j].ItemID, "- ") {
					break
				}
				child := fields[j]
				child.ItemID = strings.TrimPrefix(child.ItemID, "- ")
				lf.fields = append(lf.fields, child)
				i = j
			}
			out.lists = append(out.lists, lf)
			continue
		}
		if strings.HasPrefix(f.ItemID, "- ") {
			continue
		}
		out.scalars = append(out.scalars, f)
	}
	return out
}

func writeAPI(b *strings.Builder, api apiEntry) {
	method := methodNameForAPI(api.APIID)
	typeBase := typeBaseFromMethod(method)
	req := parseIO(api.RequestIO, typeBase+"Request")
	res := parseIO(api.ResponseIO, typeBase+"Response")

	fmt.Fprintf(b, "// %sRequest is the request payload for %s (%s).\n", typeBase, api.APIID, api.APINm)
	fmt.Fprintf(b, "type %sRequest struct {\n", typeBase)
	for _, f := range req.scalars {
		writeStructField(b, f)
	}
	b.WriteString("}\n\n")

	for _, lf := range req.lists {
		fmt.Fprintf(b, "// %s is a list row in %sRequest.\n", lf.name, typeBase)
		fmt.Fprintf(b, "type %s struct {\n", lf.name)
		for _, f := range lf.fields {
			writeStructField(b, f)
		}
		b.WriteString("}\n\n")
	}

	fmt.Fprintf(b, "// %sResponse is the response payload for %s (%s).\n", typeBase, api.APIID, api.APINm)
	fmt.Fprintf(b, "type %sResponse struct {\n", typeBase)
	hasReturnCode, hasReturnMsg := false, false
	for _, f := range res.scalars {
		if f.ItemID == "return_code" {
			hasReturnCode = true
		}
		if f.ItemID == "return_msg" {
			hasReturnMsg = true
		}
		writeStructField(b, f)
	}
	for _, lf := range res.lists {
		fmt.Fprintf(b, "\t%s []%s `json:\"%s\"`\n", exportIdent(lf.key), lf.name, lf.key)
	}
	if !hasReturnCode {
		writeStructField(b, ioField{ItemID: "return_code", Type: "Number"})
	}
	if !hasReturnMsg {
		writeStructField(b, ioField{ItemID: "return_msg", Type: "String"})
	}
	b.WriteString("}\n\n")

	for _, lf := range res.lists {
		fmt.Fprintf(b, "// %s is a list row in %sResponse.\n", lf.name, typeBase)
		fmt.Fprintf(b, "type %s struct {\n", lf.name)
		for _, f := range lf.fields {
			writeStructField(b, f)
		}
		b.WriteString("}\n\n")
	}

	desc := strings.ReplaceAll(api.Description, "\n", " ")
	fmt.Fprintf(b, "// %s calls Kiwoom %s — %s.\n", method, api.APIID, api.APINm)
	fmt.Fprintf(b, "// %s\n", desc)
	fmt.Fprintf(b, "func (c *Client) %s(ctx context.Context, req *%sRequest) (*%sResponse, error) {\n", method, typeBase, typeBase)
	fmt.Fprintf(b, "\thttpReq, err := c.newRequest(ctx, http.MethodPost, %q, %q, req)\n", api.URL, api.APIID)
	b.WriteString("\tif err != nil {\n\t\treturn nil, err\n\t}\n")
	fmt.Fprintf(b, "\tvar res %sResponse\n", typeBase)
	b.WriteString("\tif err := c.do(httpReq, &res); err != nil {\n\t\treturn nil, err\n\t}\n")
	b.WriteString("\treturn &res, nil\n")
	b.WriteString("}\n\n")
}

func writeStructField(b *strings.Builder, f ioField) {
	goType := "string"
	if strings.EqualFold(f.Type, "Number") || strings.EqualFold(f.Type, "Int") {
		goType = "int"
	}
	name := exportIdent(f.ItemID)
	tag := f.ItemID
	fmt.Fprintf(b, "\t%s %s `json:\"%s\"`\n", name, goType, tag)
}

func listStructName(typePrefix, key string) string {
	return typePrefix + exportIdent(key) + "Item"
}

func typePrefix(apiID string) string {
	return strings.ToUpper(apiID[:1]) + apiID[1:]
}

func exportIdent(itemID string) string {
	itemID = strings.TrimSpace(itemID)
	if itemID == "return_code" {
		return "ReturnCode"
	}
	if itemID == "return_msg" {
		return "ReturnMsg"
	}
	parts := strings.Split(itemID, "_")
	var b strings.Builder
	for _, p := range parts {
		if p == "" {
			continue
		}
		if p[0] >= '0' && p[0] <= '9' {
			b.WriteString("N")
		}
		b.WriteString(strings.ToUpper(p[:1]))
		if len(p) > 1 {
			b.WriteString(p[1:])
		}
	}
	if b.Len() == 0 {
		return "Field"
	}
	r := []rune(b.String())
	if !unicode.IsLetter(r[0]) {
		return "F" + b.String()
	}
	return b.String()
}

func headerComment() string {
	return "// Code generated by cmd/genapis; DO NOT EDIT."
}

func writeGoFile(path, src string) error {
	formatted, err := format.Source([]byte(src))
	if err != nil {
		_ = os.WriteFile(path+".raw", []byte(src), 0o644)
		return fmt.Errorf("format %s: %w", path, err)
	}
	return os.WriteFile(path, formatted, 0o644)
}

func writeImplementationStatus(root string, apis []apiEntry) error {
	manualDone := []struct{ id, name, link, fn string }{
		{"au10001", "접근토큰 발급", "./auth.go", "IssueToken"},
		{"au10002", "접근토큰폐기", "./auth.go", "RevokeToken"},
		{"ka00001", "계좌번호조회", "./account.go", "GetAccountNumber"},
		{"kt00001", "예수금상세현황요청", "./account.go", "GetDeposit"},
		{"kt00018", "계좌평가잔고내역요청", "./account.go", "GetAccountBalance"},
		{"ka00198", "실시간종목조회순위", "./stock.go", "GetRealtimeStockRank"},
		{"ka10001", "주식기본정보요청", "./stock.go", "GetStockBasicInfo"},
		{"ka10008", "주식외국인종목별매매동향", "./fundamental.go", "GetStockForeignInvestor"},
		{"ka10023", "거래량급증요청", "./stock.go", "GetVolumeSurge"},
		{"ka10079", "주식틱차트조회요청", "./chart.go", "GetStockTickChart"},
		{"ka10080", "주식분봉차트조회요청", "./chart.go", "GetStockMinuteChart"},
		{"ka10081", "주식일봉차트조회요청", "./chart.go", "GetStockDailyChart"},
		{"ka10082", "주식주봉차트조회요청", "./chart.go", "GetStockWeeklyChart"},
		{"ka10083", "주식월봉차트조회요청", "./chart.go", "GetStockMonthlyChart"},
		{"ka10094", "주식년봉차트조회요청", "./chart.go", "GetStockYearlyChart"},
		{"kt10000", "주식 매수주문", "./order.go", "OrderBuy"},
		{"kt10001", "주식 매도주문", "./order.go", "OrderSell"},
		{"kt10002", "주식 정정주문", "./order.go", "OrderModify"},
		{"kt10003", "주식 취소주문", "./order.go", "OrderCancel"},
		{"kt10006", "신용 매수주문", "./order.go", "CreditOrderBuy"},
		{"kt10007", "신용 매도주문", "./order.go", "CreditOrderSell"},
		{"kt10008", "신용 정정주문", "./order.go", "CreditOrderModify"},
		{"kt10009", "신용 취소주문", "./order.go", "CreditOrderCancel"},
		{"ust20000", "미국주식 매수 주문", "./us_order.go", "USOrderBuy"},
		{"ust20001", "미국주식 매도 주문", "./us_order.go", "USOrderSell"},
		{"ust20002", "미국주식 정정 주문", "./us_order.go", "USOrderModify"},
		{"ust20003", "미국주식 취소 주문", "./us_order.go", "USOrderCancel"},
		{"ust21050", "미국주식 원장 미체결", "./us_account.go", "GetUSOpenOrders"},
		{"ust21070", "미국주식 원장잔고확인", "./us_account.go", "GetUSAccountBalance"},
		{"ust21110", "해외주식 예수금", "./us_account.go", "GetUSDeposit"},
	}

	genSet := map[string]apiEntry{}
	for _, a := range apis {
		genSet[a.APIID] = a
	}

	var b strings.Builder
	b.WriteString("# Kiwoom REST API Implementation Status\n\n")
	b.WriteString("Implementation progress for the Kiwoom REST API Go wrapper.\n\n")
	b.WriteString("Official API reference: [`_ref/kiwoom-rest-api-spec.json`](./_ref/kiwoom-rest-api-spec.json) (347 endpoints).\n\n")

	total := len(manualDone) + len(apis)
	b.WriteString(fmt.Sprintf("**Summary:** %d / 347 endpoints implemented (%d%%).\n\n", total, total*100/347))

	sections := map[string][]string{
		"🔑 OAuth 인증":        {"au10001", "au10002"},
		"👤 계좌":             {},
		"📈 종목정보/시세":        {},
		"🏆 순위정보":           {},
		"🌐 기관/외국인":         {"ka10008", "ka10131"},
		"📉 차트":             {},
		"📑 조건검색":           {},
		"📥 주문":             {"kt10000", "kt10001", "kt10002", "kt10003"},
		"💳 신용주문":           {"kt10006", "kt10007", "kt10008", "kt10009"},
		"🏷 ETF":            {},
		"🇺🇸 미국주식/ETF":      {},
	}

	// classify generated APIs into sections by prefix/name
	for _, api := range apis {
		id := api.APIID
		switch {
		case strings.HasPrefix(id, "kt000") || strings.HasPrefix(id, "ka1007") || strings.HasPrefix(id, "ka1008") || strings.HasPrefix(id, "usa216") || strings.HasPrefix(id, "usa217") || strings.HasPrefix(id, "ust21"):
			sections["👤 계좌"] = append(sections["👤 계좌"], id)
		case strings.HasPrefix(id, "ka400"):
			sections["🏷 ETF"] = append(sections["🏷 ETF"], id)
		case strings.HasPrefix(id, "usa") || strings.HasPrefix(id, "ust"):
			sections["🇺🇸 미국주식/ETF"] = append(sections["🇺🇸 미국주식/ETF"], id)
		case strings.HasPrefix(id, "ka1017") || strings.HasPrefix(id, "usa202"):
			sections["📑 조건검색"] = append(sections["📑 조건검색"], id)
		case strings.Contains(id, "chart") || strings.HasPrefix(id, "usa060") || id == "ka10060" || id == "ka10064":
			sections["📉 차트"] = append(sections["📉 차트"], id)
		case strings.HasPrefix(id, "ka001") || (strings.HasPrefix(id, "ka100") && !strings.HasPrefix(id, "ka1007") && !strings.HasPrefix(id, "ka1008")) || strings.HasPrefix(id, "ka013") || strings.HasPrefix(id, "ka1010"):
			if strings.HasPrefix(id, "ka1002") || strings.HasPrefix(id, "ka1003") || id == "ka10023" || strings.HasPrefix(id, "ka1004") || strings.HasPrefix(id, "ka1005") {
				sections["🏆 순위정보"] = append(sections["🏆 순위정보"], id)
			} else {
				sections["📈 종목정보/시세"] = append(sections["📈 종목정보/시세"], id)
			}
		default:
			sections["📈 종목정보/시세"] = append(sections["📈 종목정보/시세"], id)
		}
	}

	manualMap := map[string]struct{ name, link, fn string }{}
	for _, m := range manualDone {
		manualMap[m.id] = struct{ name, link, fn string }{m.name, m.link, m.fn}
	}

	writeSection := func(title string, ids []string) {
		if len(ids) == 0 {
			return
		}
		sort.Strings(ids)
		b.WriteString("## " + title + "\n")
		for _, id := range ids {
			if m, ok := manualMap[id]; ok {
				fmt.Fprintf(&b, "- [x] `%s` - [%s](%s) (`%s`)\n", id, m.name, m.link, m.fn)
				continue
			}
			if api, ok := genSet[id]; ok {
				method := methodNameForAPI(id)
				fmt.Fprintf(&b, "- [x] `%s` - [%s](./api_registry_gen.go) (`%s`)\n", id, api.APINm, method)
			}
		}
		b.WriteString("\n")
	}

	// add manual-only account entries not in gen
	sections["👤 계좌"] = append(sections["👤 계좌"], "ka00001", "kt00001", "kt00018")

	order := []string{"🔑 OAuth 인증", "👤 계좌", "📈 종목정보/시세", "🏆 순위정보", "🌐 기관/외국인", "📉 차트", "📑 조건검색", "📥 주문", "💳 신용주문", "🏷 ETF", "🇺🇸 미국주식/ETF"}
	for _, title := range order {
		writeSection(title, sections[title])
	}

	return os.WriteFile(filepath.Join(root, "IMPLEMENTATION_STATUS.md"), []byte(b.String()), 0o644)
}

func fatal(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
