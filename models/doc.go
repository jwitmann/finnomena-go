// Package models provides data types for the Finnomena.com Thai mutual fund API.
//
// These types are used by the finnomena-go client for all API responses.
// All structs include proper JSON tags for seamless serialization.
//
// This package was migrated from the standalone finnomena-models module.
//
// Usage:
//
//	import "github.com/jwitmann/finnomena-go/models"
//
//	var fund models.Fund
//	var bars models.BarsResponse
//
// Types are organized by API domain:
//
//   - Fund catalog: Fund, FundsResponse, SymbolInfo
//   - Historical prices: BarsResponse, Bar
//   - Fund information: FundLatest, FundPerformance, FundOverview
//   - Fund details: FundFee, FundVerify, FundPortfolio
//   - Supporting types: Fee, Metric, YearlyReturn, PortfolioSection, PortfolioItem
package models
