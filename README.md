# finnomena-go

Go client for the [Finnomena.com](https://www.finnomena.com) API.

## Installation

```bash
go get github.com/jwitmann/finnomena-go
```

## Usage

```go
import finnomena "github.com/jwitmann/finnomena-go"

// Create client
client := finnomena.NewClient()

// Get all funds
funds, err := client.GetFundsList()

// Get historical prices
from := time.Now().AddDate(-1, 0, 0)
to := time.Now()
bars, err := client.GetHistoricalPrices("FUND-A", "D", from, to)

// Get fund info
latest, perf, overview, err := client.GetFundInfo("F000001")
```

## Features

- All Finnomena API endpoints
- Automatic retry with exponential backoff
- Thai-to-English fee translation
- Zero external dependencies (except models package)

## Retry Configuration

```go
client := finnomena.NewClient()
client.SetRetryConfig(5, 2*time.Second)
```

Default: 3 retries with 1s, 2s, 4s exponential backoff.

## API Coverage

- GetFundsList - All available funds
- GetHistoricalPrices - OHLCV bars
- GetFundLatest - Current NAV
- GetFundPerformance - Returns, Sharpe, drawdown
- GetFundOverview - 3D metrics
- GetFundFee - Fee structure
- GetFundPortfolio - Holdings and allocation

## Fee Translation Example

Thai fund fees are returned in Thai language. Use `TranslateFee` to convert them to English:

```go
// Get fund fee information
fee, err := client.GetFundFee("F000001")
if err != nil {
    log.Fatal(err)
}

// Translate Thai fee descriptions to English
for i := range fee.Fees {
    finnomena.TranslateFee(&fee.Fees[i], true) // true = use English names
    fmt.Printf("%s: %s %s\n", 
        fee.Fees[i].Description,  // Now in English
        fee.Fees[i].Rate,
        fee.Fees[i].Unit)
}

// Output:
// Management Fee: 1.50 % per year
// Purchase Fee: 2.00 %
// Redemption Fee: 0.00 %
```

Available translations:
- `ค่าธรรมเนียมการจัดการ` → `management fee`
- `ค่าธรรมเนียมการขายหน่วยลงทุน (Front-end Fee)` → `purchase fee`
- `ค่าธรรมเนียมการรับซื้อคืนหน่วยลงทุน (Back-end Fee)` → `redemption fee`
- And more...

## Related

- finnomena-models - Data types

## License

MIT
