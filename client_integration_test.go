//go:build integration

package finnomena

import (
	"os"
	"testing"
	"time"
)

// getTestFund returns the fund to use for testing
func getTestFund() string {
	if fund := os.Getenv("FINNOMENA_TEST_FUND"); fund != "" {
		return fund
	}
	return "TNEXTGEN-A"
}

func TestIntegration_GetServerTime(t *testing.T) {
	if os.Getenv("FINNOMENA_INTEGRATION") != "1" {
		t.Skip("Skipping integration test. Set FINNOMENA_INTEGRATION=1 to run.")
	}

	client := NewClient()

	t.Log("Getting server time...")
	serverTime, err := client.GetServerTime()
	if err != nil {
		t.Fatalf("GetServerTime() error = %v", err)
	}

	if serverTime <= 0 {
		t.Errorf("Expected positive server time, got %d", serverTime)
	}

	t.Logf("Server time: %d (%s)", serverTime, time.Unix(serverTime, 0).Format(time.RFC3339))
}

func TestIntegration_GetFundsList(t *testing.T) {
	if os.Getenv("FINNOMENA_INTEGRATION") != "1" {
		t.Skip("Skipping integration test. Set FINNOMENA_INTEGRATION=1 to run.")
	}

	client := NewClient()

	t.Log("Getting funds list...")
	funds, err := client.GetFundsList()
	if err != nil {
		t.Fatalf("GetFundsList() error = %v", err)
	}

	if len(funds) == 0 {
		t.Error("Expected non-empty funds list")
	}

	t.Logf("Found %d funds", len(funds))
	if len(funds) > 0 {
		t.Logf("First fund: %s (%s)", funds[0].ShortCode, funds[0].NameTH)
	}
}

func TestIntegration_SearchFund(t *testing.T) {
	if os.Getenv("FINNOMENA_INTEGRATION") != "1" {
		t.Skip("Skipping integration test. Set FINNOMENA_INTEGRATION=1 to run.")
	}

	client := NewClient()
	testFund := getTestFund()

	t.Logf("Searching for fund '%s'...", testFund)
	fund, err := client.SearchFund(testFund)
	if err != nil {
		t.Fatalf("SearchFund(%q) error = %v", testFund, err)
	}

	if fund == nil {
		t.Fatal("Expected fund, got nil")
	}

	if fund.ShortCode != testFund {
		t.Errorf("Expected ShortCode %q, got %q", testFund, fund.ShortCode)
	}

	t.Logf("Found: %s (ID: %s)", fund.ShortCode, fund.FundID)
}

func TestIntegration_GetSymbolInfo(t *testing.T) {
	if os.Getenv("FINNOMENA_INTEGRATION") != "1" {
		t.Skip("Skipping integration test. Set FINNOMENA_INTEGRATION=1 to run.")
	}

	client := NewClient()
	testFund := getTestFund()

	// First get fund to get short code
	fund, err := client.SearchFund(testFund)
	if err != nil {
		t.Fatalf("SearchFund(%q) error = %v", testFund, err)
	}

	t.Logf("Getting symbol info for %s...", fund.ShortCode)
	symbolInfo, err := client.GetSymbolInfo(fund.ShortCode)
	if err != nil {
		t.Fatalf("GetSymbolInfo(%q) error = %v", fund.ShortCode, err)
	}

	if symbolInfo == nil {
		t.Fatal("Expected symbol info, got nil")
	}

	t.Logf("Symbol: %s, Type: %s, Currency: %s", symbolInfo.Name, symbolInfo.Type, symbolInfo.CurrencyCode)
}

func TestIntegration_GetFundLatest(t *testing.T) {
	if os.Getenv("FINNOMENA_INTEGRATION") != "1" {
		t.Skip("Skipping integration test. Set FINNOMENA_INTEGRATION=1 to run.")
	}

	client := NewClient()
	testFund := getTestFund()

	// First get fund to get ID
	fund, err := client.SearchFund(testFund)
	if err != nil {
		t.Fatalf("SearchFund(%q) error = %v", testFund, err)
	}

	t.Logf("Getting fund latest for %s...", fund.FundID)
	latest, err := client.GetFundLatest(fund.FundID)
	if err != nil {
		t.Fatalf("GetFundLatest(%q) error = %v", fund.FundID, err)
	}

	if latest == nil {
		t.Fatal("Expected latest data, got nil")
	}

	if latest.ShortCode != testFund {
		t.Errorf("Expected ShortCode %q, got %q", testFund, latest.ShortCode)
	}

	t.Logf("Value: %.4f, Change: %.2f%%", latest.Value, latest.DChange)
}

func TestIntegration_GetHistoricalPrices(t *testing.T) {
	if os.Getenv("FINNOMENA_INTEGRATION") != "1" {
		t.Skip("Skipping integration test. Set FINNOMENA_INTEGRATION=1 to run.")
	}

	client := NewClient()
	testFund := getTestFund()

	// First get fund to get short code
	fund, err := client.SearchFund(testFund)
	if err != nil {
		t.Fatalf("SearchFund(%q) error = %v", testFund, err)
	}

	to := time.Now()
	from := to.AddDate(0, 0, -30)

	t.Logf("Getting historical prices for %s (last 30 days)...", fund.ShortCode)
	bars, err := client.GetHistoricalPrices(fund.ShortCode, "D", from, to)
	if err != nil {
		t.Fatalf("GetHistoricalPrices(%q) error = %v", fund.ShortCode, err)
	}

	t.Logf("Status: %s, Bars: %d", bars.Status, len(bars.Time))
	if len(bars.Time) > 0 {
		t.Logf("First bar: Close=%.4f", bars.Close[0])
	}
}

func TestIntegration_GetFundPerformance(t *testing.T) {
	if os.Getenv("FINNOMENA_INTEGRATION") != "1" {
		t.Skip("Skipping integration test. Set FINNOMENA_INTEGRATION=1 to run.")
	}

	client := NewClient()
	testFund := getTestFund()

	// First get fund to get ID
	fund, err := client.SearchFund(testFund)
	if err != nil {
		t.Fatalf("SearchFund(%q) error = %v", testFund, err)
	}

	t.Logf("Getting fund performance for %s...", fund.FundID)
	perf, err := client.GetFundPerformance(fund.FundID)
	if err != nil {
		t.Fatalf("GetFundPerformance(%q) error = %v", fund.FundID, err)
	}

	if perf == nil {
		t.Fatal("Expected performance data, got nil")
	}

	t.Logf("1Y Return: %.2f%%, 3Y Return: %.2f%%", perf.TotalReturn1Y, perf.TotalReturn3Y)
}

func TestIntegration_GetFundFee_Thai(t *testing.T) {
	if os.Getenv("FINNOMENA_INTEGRATION") != "1" {
		t.Skip("Skipping integration test. Set FINNOMENA_INTEGRATION=1 to run.")
	}

	client := NewClient()
	testFund := getTestFund()

	// First get fund to get ID
	fund, err := client.SearchFund(testFund)
	if err != nil {
		t.Fatalf("SearchFund(%q) error = %v", testFund, err)
	}

	t.Logf("Getting fund fee (Thai) for %s...", fund.FundID)
	fee, err := client.GetFundFee(fund.FundID)
	if err != nil {
		t.Fatalf("GetFundFee(%q) error = %v", fund.FundID, err)
	}

	if fee == nil {
		t.Fatal("Expected fee data, got nil")
	}

	t.Logf("Fund: %s, Fees count: %d", fee.ShortCode, len(fee.Fees))
	for i, f := range fee.Fees {
		if i < 3 {
			t.Logf("  - %s: %s %s", f.Description, f.Rate, f.Unit)
		}
	}
	if len(fee.Fees) > 3 {
		t.Logf("  ... and %d more", len(fee.Fees)-3)
	}
}

func TestIntegration_GetFundFee_English(t *testing.T) {
	if os.Getenv("FINNOMENA_INTEGRATION") != "1" {
		t.Skip("Skipping integration test. Set FINNOMENA_INTEGRATION=1 to run.")
	}

	client := NewClient()
	testFund := getTestFund()

	// First get fund to get ID
	fund, err := client.SearchFund(testFund)
	if err != nil {
		t.Fatalf("SearchFund(%q) error = %v", testFund, err)
	}

	t.Logf("Getting fund fee (English) for %s...", fund.FundID)
	fee, err := client.GetFundFee(fund.FundID)
	if err != nil {
		t.Fatalf("GetFundFee(%q) error = %v", fund.FundID, err)
	}

	if fee == nil {
		t.Fatal("Expected fee data, got nil")
	}

	// Translate to English
	for i := range fee.Fees {
		TranslateFee(&fee.Fees[i], true)
	}

	t.Logf("Fund: %s, Fees count: %d", fee.ShortCode, len(fee.Fees))
	for i, f := range fee.Fees {
		if i < 3 {
			t.Logf("  - %s: %s %s", f.Description, f.Rate, f.Unit)
		}
	}
	if len(fee.Fees) > 3 {
		t.Logf("  ... and %d more", len(fee.Fees)-3)
	}
}

func TestIntegration_GetFundOverview(t *testing.T) {
	if os.Getenv("FINNOMENA_INTEGRATION") != "1" {
		t.Skip("Skipping integration test. Set FINNOMENA_INTEGRATION=1 to run.")
	}

	client := NewClient()
	testFund := getTestFund()

	// First get fund to get ID
	fund, err := client.SearchFund(testFund)
	if err != nil {
		t.Fatalf("SearchFund(%q) error = %v", testFund, err)
	}

	t.Logf("Getting fund overview for %s...", fund.FundID)
	overview, err := client.GetFundOverview(fund.FundID)
	if err != nil {
		t.Fatalf("GetFundOverview(%q) error = %v", fund.FundID, err)
	}

	if overview == nil {
		t.Fatal("Expected overview data, got nil")
	}

	t.Logf("Category: %s, Finno Score: %d", overview.AIMCCategory, overview.FinnoScore)
}

func TestIntegration_GetFundVerify(t *testing.T) {
	if os.Getenv("FINNOMENA_INTEGRATION") != "1" {
		t.Skip("Skipping integration test. Set FINNOMENA_INTEGRATION=1 to run.")
	}

	client := NewClient()
	testFund := getTestFund()

	// First get fund to get ID
	fund, err := client.SearchFund(testFund)
	if err != nil {
		t.Fatalf("SearchFund(%q) error = %v", testFund, err)
	}

	t.Logf("Getting fund verify for %s...", fund.FundID)
	verify, err := client.GetFundVerify(fund.FundID)
	if err != nil {
		t.Fatalf("GetFundVerify(%q) error = %v", fund.FundID, err)
	}

	if verify == nil {
		t.Fatal("Expected verify data, got nil")
	}

	t.Logf("Verify data retrieved successfully")
}

func TestIntegration_GetFundPortfolio(t *testing.T) {
	if os.Getenv("FINNOMENA_INTEGRATION") != "1" {
		t.Skip("Skipping integration test. Set FINNOMENA_INTEGRATION=1 to run.")
	}

	client := NewClient()
	testFund := getTestFund()

	// First get fund to get ID
	fund, err := client.SearchFund(testFund)
	if err != nil {
		t.Fatalf("SearchFund(%q) error = %v", testFund, err)
	}

	t.Logf("Getting fund portfolio for %s...", fund.FundID)
	portfolio, err := client.GetFundPortfolio(fund.FundID)
	if err != nil {
		t.Fatalf("GetFundPortfolio(%q) error = %v", fund.FundID, err)
	}

	if portfolio == nil {
		t.Fatal("Expected portfolio data, got nil")
	}

	t.Logf("Top Holdings: %d, Global Stock Sector: %d, Asset Allocation: %d, Regional Exposure: %d",
		len(portfolio.TopHoldings.Elements),
		len(portfolio.GlobalStockSector.Elements),
		len(portfolio.AssetAllocation.Elements),
		len(portfolio.RegionalExposure.Elements))
}

func TestIntegration_EndToEndWorkflow(t *testing.T) {
	if os.Getenv("FINNOMENA_INTEGRATION") != "1" {
		t.Skip("Skipping integration test. Set FINNOMENA_INTEGRATION=1 to run.")
	}

	client := NewClient()
	testFund := getTestFund()

	t.Logf("Running end-to-end workflow for %s...", testFund)

	// Step 1: Search for fund
	t.Log("1. Searching for fund...")
	fund, err := client.SearchFund(testFund)
	if err != nil {
		t.Fatalf("SearchFund failed: %v", err)
	}
	t.Logf("   Found: %s (ID: %s)", fund.ShortCode, fund.FundID)

	// Step 2: Get symbol info
	t.Log("2. Getting symbol info...")
	symbolInfo, err := client.GetSymbolInfo(fund.ShortCode)
	if err != nil {
		t.Fatalf("GetSymbolInfo failed: %v", err)
	}
	t.Logf("   Symbol: %s, Type: %s", symbolInfo.Name, symbolInfo.Type)

	// Step 3: Get latest data
	t.Log("3. Getting latest data...")
	latest, err := client.GetFundLatest(fund.FundID)
	if err != nil {
		t.Fatalf("GetFundLatest failed: %v", err)
	}
	t.Logf("   Value: %.4f", latest.Value)

	// Step 4: Get historical prices
	t.Log("4. Getting historical prices...")
	to := time.Now()
	from := to.AddDate(0, 0, -7)
	bars, err := client.GetHistoricalPrices(fund.ShortCode, "D", from, to)
	if err != nil {
		t.Fatalf("GetHistoricalPrices failed: %v", err)
	}
	t.Logf("   Retrieved %d bars", len(bars.Time))

	// Step 5: Get performance
	t.Log("5. Getting performance...")
	perf, err := client.GetFundPerformance(fund.FundID)
	if err != nil {
		t.Fatalf("GetFundPerformance failed: %v", err)
	}
	t.Logf("   1Y: %.2f%%", perf.TotalReturn1Y)

	// Step 6: Get fees
	t.Log("6. Getting fees...")
	fee, err := client.GetFundFee(fund.FundID)
	if err != nil {
		t.Fatalf("GetFundFee failed: %v", err)
	}
	t.Logf("   %d fees", len(fee.Fees))

	t.Log("End-to-end workflow completed successfully!")
}
