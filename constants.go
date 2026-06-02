package kioom

const (
	// LiveDomain is the base URL for the live Kiwoom investment environment.
	LiveDomain = "https://api.kiwoom.com"
	// MockDomain is the base URL for the mock/simulated Kiwoom investment environment.
	MockDomain = "https://mockapi.kiwoom.com"
)

// API IDs for Kiwoom REST API
const (
	// Auth
	API_IssueToken  = "au10001"
	API_RevokeToken = "au10002"

	// Account
	API_GetAccountNumber  = "ka00001"
	API_GetDeposit        = "kt00001"
	API_GetAccountBalance = "kt00018"

	// Stock Info
	API_GetStockBasicInfo = "ka10001"
	API_GetRealtimeRank   = "ka00198"
	API_GetTickChart      = "ka10079"
	API_GetMinuteChart    = "ka10080"
	API_GetDailyChart     = "ka10081"
	API_GetWeeklyChart    = "ka10082"
	API_GetMonthlyChart   = "ka10083"
	API_GetYearlyChart    = "ka10094"

	// Orders
	API_OrderBuy        = "kt10000"
	API_OrderSell       = "kt10001"
	API_OrderModify     = "kt10002"
	API_OrderCancel     = "kt10003"
	API_CreditOrderBuy  = "kt10006"
	API_CreditOrderSell = "kt10007"
	API_CreditModify    = "kt10008"
	API_CreditCancel    = "kt10009"
)
