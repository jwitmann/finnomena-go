# finnomena-go

[![Go Reference](https://pkg.go.dev/badge/github.com/jwitmann/finnomena-go.svg)](https://pkg.go.dev/github.com/jwitmann/finnomena-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/jwitmann/finnomena-go?t=1)](https://goreportcard.com/report/github.com/jwitmann/finnomena-go)
[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://opensource.org/licenses/MIT)

Go client for the [Finnomena.com](https://www.finnomena.com) Thai mutual fund API.

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

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## Related

- [finnomena-models](https://github.com/jwitmann/finnomena-models) - Data types for this client
- [thai-market-data](https://github.com/jwitmann/thai-market-data) - Thai market data (AIMC, SET)

## Disclaimer

This is an unofficial client for the Finnomena.com API. It is not affiliated with or endorsed by Finnomena.

## License

MIT License - see [LICENSE](LICENSE) file for details
