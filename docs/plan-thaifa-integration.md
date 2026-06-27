# Plan: FinnomenaAdapter for ThaiFA Unified Backend

## Status (created 2026-06-27)

Pending — implementation depends on ThaiFA's `internal/provider/interface.go` being created first.

## Goal

ThaiFA needs a `FinnomenaAdapter` that wraps `*finnomena.Client` and implements both:
1. **`FundAllProvider`** — fund data (6 concurrent API calls merged into `FundAll`)
2. **`FallbackProvider`** — OHLCV prices (direct delegation to `GetHistoricalPrices`)

The adapter lives in ThaiFA (`internal/provider/adapters.go`), not in finnomena-go. This document describes what the adapter needs from finnomena-go's API.

## What the Adapter Needs from finnomena-go

### FundAllProvider Methods

| ThaiFA method | finnomena-go call | Notes |
|---|---|---|
| `GetFundsList(ctx)` | `client.GetFundsList()` | Returns `[]models.Fund` |
| `GetFundAll(ctx, fundID)` | 6 concurrent calls merged (see below) | Core method |
| `GetFundNav(ctx, fundID, shortCode, from, to)` | `client.GetHistoricalPrices(shortCode, "D", from, to)` | Returns `[]models.Bar` |
| `SearchFund(ctx, query)` | `client.SearchFund(query)` | Returns `[]models.Fund` |
| `GetSymbolInfo(ctx, symbol)` | `client.GetSymbolInfo(symbol)` | Returns `*models.SymbolInfo` |
| `ResolveFundCode(ctx, shortCode)` | No-op | Finnomena shortCode = FundID |

### FallbackProvider Methods

| ThaiFA method | finnomena-go call | Notes |
|---|---|---|
| `CanHandle(ctx, shortCode)` | Always `true` | Finnomena supports all funds |
| `FetchBars(ctx, shortCode, from, to)` | `client.GetHistoricalPrices(shortCode, "D", from, to)` | Returns `*models.BarsResponse` |

### GetFundAll — 6 Concurrent Calls

```go
func (a *FinnomenaAdapter) GetFundAll(ctx context.Context, fundID string) (*provider.FundAll, error) {
    // All 6 calls run concurrently for performance
    var wg sync.WaitGroup
    wg.Add(6)

    // Channel for errors (collected but tolerated)
    errs := make([]error, 6)

    // 1. GetFundLatest — NAV, dates, SEC flags
    latestCh := make(chan *models.FundLatest, 1)
    go func() {
        defer wg.Done()
        l, err := a.client.GetFundLatest(fundID)
        if err != nil { errs[0] = err; return }
        latestCh <- l
    }()

    // 2. GetFundPerformance — returns, risk, MaxDrawdown
    perfCh := make(chan *models.FundPerformance, 1)
    go func() {
        defer wg.Done()
        p, err := a.client.GetFundPerformance(fundID)
        if err != nil { errs[1] = err; return }
        perfCh <- p
    }()

    // 3. GetFundOverview — FinnoScore, fund info
    overviewCh := make(chan *models.FundOverview, 1)
    go func() {
        defer wg.Done()
        o, err := a.client.GetFundOverview(fundID)
        if err != nil { errs[2] = err; return }
        overviewCh <- o
    }()

    // 4. GetFundFee — fees (Thai strings, parsed by adapter)
    feeCh := make(chan *models.FundFee, 1)
    go func() {
        defer wg.Done()
        f, err := a.client.GetFundFee(fundID, false)
        if err != nil { errs[3] = err; return }
        feeCh <- f
    }()

    // 5. GetFundPortfolio — holdings, sectors, allocation
    portfolioCh := make(chan *models.FundPortfolio, 1)
    go func() {
        defer wg.Done()
        pt, err := a.client.GetFundPortfolio(fundID)
        if err != nil { errs[4] = err; return }
        portfolioCh <- pt
    }()

    // 6. GetFundDividend — dividend history
    dividendCh := make(chan *models.FundDividend, 1)
    go func() {
        defer wg.Done()
        d, err := a.client.GetFundDividend(fundID)
        if err != nil { errs[5] = err; return }
        dividendCh <- d
    }()

    // Wait for all, close channels
    go func() { wg.Wait(); close(latestCh); close(perfCh); ... }()

    // Collect successful results
    var latest *models.FundLatest
    var perf *models.FundPerformance
    // ... etc

    // Tolerate partial failures — return what we have
    if latest == nil && perf == nil {
        return nil, fmt.Errorf("all API calls failed for %s: %w", fundID, errors.Join(errs...))
    }

    return mergeToFundAll(fundID, latest, perf, overview, fee, portfolio, dividend), nil
}
```

### Field Mapping: Finnomena → FundAll

| FundAll field | finnomena source | Conversion |
|---|---|---|
| `FundID` | `fundID` (parameter) | Pass-through |
| `ShortCode` | `fundID` (same in Finnomena) | Pass-through |
| `FundNameTH` | `FundOverview.FundNameTH` | Direct |
| `FundNameEN` | — (not in Finnomena API) | Use AIMC data or leave empty |
| `NAV` | `FundLatest.Value` | Rename `Value` → `NAV` |
| `NAVDate` | `FundLatest.Date` | `time.Time` → `YYYYMMDD` string |
| `Return1Y` | `FundPerformance.TotalReturn1Y` | Rename |
| `Return3Y` | `FundPerformance.TotalReturn3Y` | Rename |
| `Return5Y` | `FundPerformance.TotalReturn5Y` | Rename |
| `MaxDrawdown1Y` | `FundPerformance.MaxDrawdown1Y` | Direct |
| `SharpeRatio1Y` | `FundPerformance.SharpeRatio1Y` | Direct |
| `SD1Y` | `FundPerformance.Std1Y` | Rename `Std` → `SD` |
| `Beta1Y` | — (not in Finnomena API) | Leave zero |
| `FrontEndFee` | Parse `FundFee.Fees[]` matching `"Front-end Fee"` | Extract `rate` string → `float64` |
| `BackEndFee` | Parse `FundFee.Fees[]` matching `"Back-end Fee"` | Same pattern |
| `ManagementFee` | Parse `FundFee.Fees[]` matching `"ค่าธรรมเนียมการจัดการ"` | Same pattern |
| `TopHoldings` | `FundPortfolio.TopHoldings` | Drop `ShortCode`, `LinkURL`, `Color` |
| `SectorAllocation` | `FundPortfolio.SectorAllocation` | Drop `ShortCode`, `LinkURL`, `Color` |
| `AssetAllocation` | `FundPortfolio.AssetAllocation` | Drop `ShortCode`, `LinkURL`, `Color` |
| `DividendHistory` | `FundDividend.Dividends[]` | Convert `models.Dividend` → `provider.DividendEntry` |
| `FinnoScore` | `FundOverview.FinnoScore` | Direct (Finnomena-only feature) |
| `SECIsActive` | `FundLatest.SECIsActive` | Direct |
| `IsInTrending` | `FundLatest.IsInTrending` | Direct |
| `IsFinnomenaPick` | `FundLatest.IsFinnomenaPick` | Direct |

### Fee Parsing Detail

Finnomena returns fees as Thai strings in `FundFee.Fees[]`:

```go
func parseFee(fees []models.FundFeeItem, keyword string) float64 {
    for _, f := range fees {
        if strings.Contains(f.Description, keyword) {
            // rate is like "0.50%" or "0"
            rate := strings.TrimSuffix(f.Rate, "%")
            val, _ := strconv.ParseFloat(rate, 64)
            return val
        }
    }
    return 0
}

// Usage:
frontEndFee := parseFee(fee.Fees, "Front-end")
backEndFee := parseFee(fee.Fees, "Back-end")
mgmtFee := parseFee(fee.Fees, "ค่าธรรมเนียมการจัดการ")
switchInFee := parseFee(fee.Fees, "Switching In")
switchOutFee := parseFee(fee.Fees, "Switching Out")
```

### SymbolInfo Mapping

Finnomena's `GetSymbolInfo` returns a different structure than ThaiFA's `SymbolInfo`. The adapter maps:

```go
func mapSymbolInfo(fi *models.FundInfo, fund *models.Fund) *provider.SymbolInfo {
    return &provider.SymbolInfo{
        Name:                 fund.NameTH,
        Timezone:             "Asia/Bangkok",
        MinMov:               1,
        MinMove2:             1,
        PriceScale:           100,
        PointValue:           1,
        Ticker:               fund.ShortCode,
        Description:          fund.NameTH,
        Type:                 "fund",
        DataStatus:           "realtime",
        SupportedResolutions: []string{"D", "W", "M"},
        Session:              "0930-1630",
        CurrencyCode:         "THB",
    }
}
```

## Key Design Decisions

1. **Adapter lives in ThaiFA, not finnomena-go** — finnomena-go stays a clean API client
2. **6 concurrent calls** — matches WM's single-call semantics (all data fetched eagerly)
3. **Partial failure tolerance** — return what we have, don't fail entirely
4. **FundVerify stays separate** — no WM equivalent; use `client.GetFundVerify()` directly in FundService
5. **ShortCode = FundID** — Finnomena uses the same string for both; no mapping needed

## Testing

The adapter should be tested against real Finnomena API responses:

```go
func TestFinnomenaAdapter_GetFundAll(t *testing.T) {
    if testing.Short() {
        t.Skip("skipping integration test")
    }
    client := finnomena.NewClient()
    adapter := NewFinnomenaAdapter(client)

    fundAll, err := adapter.GetFundAll(context.Background(), "TISCOJP")
    require.NoError(t, err)
    assert.Equal(t, "TISCOJP", fundAll.ShortCode)
    assert.NotZero(t, fundAll.NAV)
    assert.NotZero(t, fundAll.FrontEndFee)
    // MaxDrawdown should be non-zero for TISCOJP (equity fund)
    assert.NotNil(t, fundAll.MaxDrawdown1Y)
}
```

## Cross-reference

- ThaiFA `docs/plans/unified-adapter-plan.md`: main unified adapter plan
- ThaiFA `internal/provider/interface.go`: FundAllProvider interface definition
- ThaiFA `internal/service/fallback_provider.go`: FallbackProvider interface
- finnomena-go `docs/API.md`: full API reference
- finnomena-go `models/`: all data types
