# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.0] - 2024-03-06

### Added
- Initial release
- Complete Finnomena API client
- HTTP retry logic with exponential backoff
- Support for all API endpoints:
  - GetFundsList - List all available funds
  - GetHistoricalPrices - Historical OHLCV data
  - GetFundLatest - Latest NAV and change
  - GetFundPerformance - Performance metrics
  - GetFundOverview - 3D metrics (PP, RR, DD)
  - GetFundFee - Fee structure
  - GetFundPortfolio - Portfolio holdings
  - GetFundVerify - Available periods
  - GetServerTime - Server timestamp
- Thai-to-English fee translation
- Configurable retry settings
- Comprehensive test suite
- Full API documentation (1137 lines)

### Features
- Automatic retry on network failures
- Exponential backoff: 1s, 2s, 4s
- Configurable timeout (default 30s)
- Zero external dependencies (except models)

[1.0.0]: https://github.com/jwitmann/finnomena-go/releases/tag/v1.0.0
