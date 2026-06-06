# finnomena-go

[![Go Reference](https://pkg.go.dev/badge/github.com/jwitmann/finnomena-go.svg)](https://pkg.go.dev/github.com/jwitmann/finnomena-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/jwitmann/finnomena-go?t=1)](https://goreportcard.com/report/github.com/jwitmann/finnomena-go)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)

Go client for the [Finnomena.com](https://www.finnomena.com) Thai mutual fund API.

## Installation

```bash
go get github.com/jwitmann/finnomena-go
```

## Quick Start

```go
import (
    "context"
    "time"
    
    finnomena "github.com/jwitmann/finnomena-go"
)

func main() {
    client := finnomena.NewClient()
    
    // Get all funds
    funds, err := client.GetFundsList(context.Background())
    
    // Get historical prices
    from := time.Now().AddDate(-1, 0, 0)
    to := time.Now()
    bars, err := client.GetHistoricalPrices(context.Background(), "FUND-A", "D", from, to)
    
    // Get fund latest NAV
    latest, err := client.GetFundLatest(context.Background(), "F000001")
}
```

## Data Types

All API response types are available in the `models` subpackage:

```go
import "github.com/jwitmann/finnomena-go/models"

// Direct access to types
var fund models.Fund
var bars models.BarsResponse
```

> **Note:** The standalone `finnomena-models` module has been merged into this package. Update your import from `github.com/jwitmann/finnomena-models` to `github.com/jwitmann/finnomena-go/models`.

## Features

- All Finnomena API endpoints
- Context support for cancellation/timeouts
- Automatic retry with exponential backoff
- Thai-to-English fee translation
- Zero external dependencies

## Retry Configuration

```go
client := finnomena.NewClient()
client.SetRetryConfig(5, 2*time.Second)
```

Default: 3 retries with 1s, 2s, 4s exponential backoff.

## API Coverage

### Fund Data
- `GetFundsList(ctx)` - All available funds
- `SearchFund(query)` - Find fund by short code or ID
- `GetSymbolInfo(ctx, symbol)` - Trading symbol metadata

### Historical Prices
- `GetHistoricalPrices(ctx, symbol, resolution, from, to)` - OHLCV bars

### Fund Information
- `GetFundLatest(ctx, fundID)` - Current NAV and change
- `GetFundPerformance(ctx, fundID)` - Returns, Sharpe, drawdown
- `GetFundOverview(ctx, fundID)` - 3D metrics (PP, RR, DD scores)
- `GetFundFee(ctx, fundID)` - Fee structure
- `GetFundPortfolio(ctx, fundID)` - Holdings and allocation
- `GetFundVerify(ctx, fundID)` - Available data periods

### Utility
- `GetServerTime(ctx)` - Server timestamp

## Fee Translation Example

Thai fund fees are returned in Thai language. Use `TranslateFee` to convert them to English:

```go
fee, err := client.GetFundFee(context.Background(), "F000001")
if err != nil {
    log.Fatal(err)
}

for i := range fee.Fees {
    finnomena.TranslateFee(&fee.Fees[i], true)
    fmt.Printf("%s: %s %s\n", 
        fee.Fees[i].Description,
        fee.Fees[i].Rate,
        fee.Fees[i].Unit)
}
```

## Related

- [finnomena-models](https://github.com/jwitmann/finnomena-models) - **Deprecated**. Types merged into this repo.
- [thai-market-data](https://github.com/jwitmann/thai-market-data) - Thai market data (AIMC, SET)

## Disclaimer

This is an unofficial client for the Finnomena.com API. It is not affiliated with or endorsed by Finnomena.

## License

MIT License - see [LICENSE](LICENSE) file for details
