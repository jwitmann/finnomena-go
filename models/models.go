package models

import (
	"time"
)

// FundsResponse represents the response from the funds list endpoint
type FundsResponse struct {
	Status      bool   `json:"status"`
	ServiceCode string `json:"service_code"`
	Data        []Fund `json:"data"`
}

// Fund represents a single fund in the list
type Fund struct {
	FundID            string `json:"fund_id"`
	ShortCode         string `json:"short_code"`
	NameTH            string `json:"name_th"`
	AIMCCategoryID    string `json:"aimc_category_id"`
	ShortDesc         string `json:"short_desc"`
	IsFinnomenaPick   bool   `json:"is_finnomena_pick"`
	CreditCardAllowed bool   `json:"credit_card_allowed"`
	IsInTrending      bool   `json:"is_in_trending"`
	SECIsActive       bool   `json:"sec_is_active"`
	Promotions        []any  `json:"promotions"`
	LegalName         string `json:"legal_name,omitempty"`
	FirmName          string `json:"firm_name,omitempty"`
	AIMCCategory      string `json:"aimc_category,omitempty"`
}

// SymbolInfo represents symbol information from the TradingView API
type SymbolInfo struct {
	Name                 string   `json:"name"`
	Timezone             string   `json:"timezone"`
	MinMov               int      `json:"minmov"`
	MinMove2             int      `json:"minmove2"`
	PriceScale           int      `json:"pricescale"`
	PointValue           int      `json:"pointvalue"`
	Ticker               string   `json:"ticker"`
	Description          string   `json:"description"`
	Type                 string   `json:"type"`
	DataStatus           string   `json:"data_status"`
	SupportedResolutions []string `json:"supported-resolutions"`
	Session              string   `json:"session"`
	CurrencyCode         string   `json:"currency_code"`
}

// BarsResponse represents the response from the historical bars endpoint
type BarsResponse struct {
	Status   string    `json:"s"`
	ErrMsg   string    `json:"errmsg,omitempty"`
	Time     []int64   `json:"t"`
	Close    []float64 `json:"c"`
	Open     []float64 `json:"o,omitempty"`
	High     []float64 `json:"h,omitempty"`
	Low      []float64 `json:"l,omitempty"`
	Volume   []float64 `json:"v,omitempty"`
	NextTime int64     `json:"nextTime,omitempty"`
}

// Bar represents a single OHLCV bar
type Bar struct {
	Time   int64
	Open   float64
	High   float64
	Low    float64
	Close  float64
	Volume float64
}

// FundLatestResponse represents the response from the fund latest endpoint
type FundLatestResponse struct {
	Status      bool       `json:"status"`
	ServiceCode string     `json:"service_code"`
	Data        FundLatest `json:"data"`
}

// FundLatest represents the latest fund data
type FundLatest struct {
	FundID            string    `json:"fund_id"`
	ShortCode         string    `json:"short_code"`
	IsInTrending      bool      `json:"is_in_trending"`
	Date              time.Time `json:"date"`
	Value             float64   `json:"value"`
	Amount            float64   `json:"amount"`
	DChange           float64   `json:"d_change"`
	IsFinnomenaPick   bool      `json:"is_finnomena_pick"`
	CreditCardAllowed bool      `json:"credit_card_allowed"`
	PerformanceReady  bool      `json:"performance_ready"`
	OperationReady    bool      `json:"operation_ready"`
	SECFundStatus     string    `json:"sec_fund_status"`
	SECIsActive       bool      `json:"sec_is_active"`
}

// FundPerformanceResponse represents the response from the fund performance endpoint
type FundPerformanceResponse struct {
	Status      bool            `json:"status"`
	ServiceCode string          `json:"service_code"`
	Data        FundPerformance `json:"data"`
}

// FundPerformance represents fund performance data
type FundPerformance struct {
	FundID         string         `json:"fund_id"`
	ShortCode      string         `json:"short_code"`
	DayEndDate     time.Time      `json:"day_end_date"`
	TotalReturn1W  float64        `json:"total_return_1w"`
	TotalReturn1M  float64        `json:"total_return_1m"`
	TotalReturn3M  float64        `json:"total_return_3m"`
	TotalReturn6M  float64        `json:"total_return_6m"`
	TotalReturn1Y  float64        `json:"total_return_1y"`
	TotalReturn3Y  float64        `json:"total_return_3y"`
	TotalReturn5Y  float64        `json:"total_return_5y"`
	TotalReturn10Y *float64       `json:"total_return_10y"`
	Std1Y          float64        `json:"std_1y"`
	Std3Y          float64        `json:"std_3y"`
	Std5Y          float64        `json:"std_5y"`
	NetAssets      float64        `json:"net_assets"`
	DataDate       time.Time      `json:"data_date"`
	ReturnsYear    []YearlyReturn `json:"returns_year"`
	SharpeRatio1Y  float64        `json:"sharpe_ratio_1y"`
	SharpeRatio3Y  float64        `json:"sharpe_ratio_3y"`
	SharpeRatio5Y  float64        `json:"sharpe_ratio_5y"`
	MaxDrawdown1Y  float64        `json:"max_drawdown_1y"`
	MaxDrawdown3Y  float64        `json:"max_drawdown_3y"`
	MaxDrawdown5Y  float64        `json:"max_drawdown_5y"`
}

// YearlyReturn represents a single year's return
type YearlyReturn struct {
	Year  int      `json:"year"`
	Value *float64 `json:"value"`
}

// FundOverviewResponse represents the response from the fund overview endpoint
type FundOverviewResponse struct {
	Status      bool         `json:"status"`
	ServiceCode string       `json:"service_code"`
	Data        FundOverview `json:"data"`
}

// FundOverview represents fund overview data
type FundOverview struct {
	FundID             string    `json:"fund_id"`
	ShortCode          string    `json:"short_code"`
	AIMCCategory       string    `json:"aimc_category"`
	AIMCCategoryID     string    `json:"aimc_category_id"`
	AIMCCategoryNameTH string    `json:"aimc_category_name_th"`
	FinnoScore         int       `json:"finno_score"`
	DataDate           time.Time `json:"data_date"`
	PP                 Metric    `json:"pp"`
	RR                 Metric    `json:"rr"`
	DD                 Metric    `json:"dd"`
}

// Metric represents a metric with fund, avg, and text values
type Metric struct {
	Fund int    `json:"fund"`
	Avg  int    `json:"avg"`
	Text string `json:"text"`
}

// FundFeeResponse represents the response from the fund fee endpoint
type FundFeeResponse struct {
	Status      bool    `json:"status"`
	ServiceCode string  `json:"service_code"`
	Data        FundFee `json:"data"`
}

// FundFee represents fund fee data
type FundFee struct {
	FundID    string `json:"fund_id"`
	ShortCode string `json:"short_code"`
	Fees      []Fee  `json:"fees"`
}

// Fee represents a single fee entry
type Fee struct {
	LastUpdate       *string `json:"last_update"`
	ClassAbbrName    string  `json:"class_abbr_name"`
	Description      string  `json:"description"`
	Rate             string  `json:"rate"`
	RateUnit         string  `json:"rate_unit"`
	ActualValue      string  `json:"actual_value"`
	ActualValueUnit  string  `json:"actual_value_unit"`
	OtherDescription string  `json:"other_description"`
	Unit             string  `json:"unit"`
}

// FundVerifyResponse represents the response from the fund verify endpoint
type FundVerifyResponse struct {
	Status      bool       `json:"status"`
	ServiceCode string     `json:"service_code"`
	Data        FundVerify `json:"data"`
}

// FundVerify represents fund verification data with available periods
type FundVerify struct {
	FundID    string   `json:"fund_id"`
	ShortCode string   `json:"short_code"`
	Range     []string `json:"range"`
	Default   string   `json:"default"`
}

// FundPortfolioResponse represents the response from the fund portfolio endpoint
type FundPortfolioResponse struct {
	Status      bool          `json:"status"`
	ServiceCode string        `json:"service_code"`
	Data        FundPortfolio `json:"data"`
}

// FundPortfolio represents fund portfolio data
type FundPortfolio struct {
	FundID            string           `json:"fund_id"`
	ShortCode         string           `json:"short_code"`
	TopHoldings       PortfolioSection `json:"top_holdings"`
	GlobalStockSector PortfolioSection `json:"global_stock_sector"`
	AssetAllocation   PortfolioSection `json:"asset_allocation"`
	RegionalExposure  PortfolioSection `json:"regional_exposure"`
}

// PortfolioSection represents a section of portfolio data
type PortfolioSection struct {
	DataDate time.Time       `json:"data_date"`
	Elements []PortfolioItem `json:"elements"`
}

// PortfolioItem represents a single portfolio item
type PortfolioItem struct {
	Name      string  `json:"name"`
	Percent   float64 `json:"percent"`
	ShortCode string  `json:"short_code"`
	LinkURL   string  `json:"link_url"`
	Color     string  `json:"color"`
}
