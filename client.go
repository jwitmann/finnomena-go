package finnomena

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/jwitmann/finnomena-models"
)

const (
	BaseURL = "https://www.finnomena.com/fn3/api/fund/v2/public"

	// Default retry configuration
	DefaultMaxRetries      = 3
	DefaultRetryDelay      = 1 * time.Second
	RetryBackoffMultiplier = 2
)

// Client represents the Finnomena API client
type Client struct {
	httpClient *http.Client
	baseURL    string
	maxRetries int
	retryDelay time.Duration
}

// NewClient creates a new Finnomena API client
func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		baseURL:    BaseURL,
		maxRetries: DefaultMaxRetries,
		retryDelay: DefaultRetryDelay,
	}
}

// NewClientWithTimeout creates a new Finnomena API client with custom timeout
func NewClientWithTimeout(timeout time.Duration) *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: timeout,
		},
		baseURL:    BaseURL,
		maxRetries: DefaultMaxRetries,
		retryDelay: DefaultRetryDelay,
	}
}

// NewClientWithRetry creates a new Finnomena API client with custom retry configuration
func NewClientWithRetry(maxRetries int, retryDelay time.Duration) *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		baseURL:    BaseURL,
		maxRetries: maxRetries,
		retryDelay: retryDelay,
	}
}

// SetRetryConfig configures the retry behavior for the client
func (c *Client) SetRetryConfig(maxRetries int, retryDelay time.Duration) {
	c.maxRetries = maxRetries
	c.retryDelay = retryDelay
}

// isRetryableError determines if an error should trigger a retry
// Retries on: 5xx server errors, network errors, timeouts
// Does not retry on: 4xx client errors
func isRetryableError(err error, statusCode int) bool {
	if err != nil {
		// Retry on network errors and timeouts
		return true
	}
	// Don't retry on 4xx client errors
	if statusCode >= 400 && statusCode < 500 {
		return false
	}
	// Retry on 5xx server errors
	return statusCode >= 500
}

// doRequest performs an HTTP GET request with retry logic and returns the response body
func (c *Client) doRequest(endpoint string) ([]byte, error) {
	var lastErr error
	var lastStatusCode int

	for attempt := 0; attempt < c.maxRetries; attempt++ {
		if attempt > 0 {
			// Exponential backoff: delay * (2^(attempt-1))
			// Attempt 1: delay * 1, Attempt 2: delay * 2, Attempt 3: delay * 4
			delay := c.retryDelay * time.Duration(1<<uint(attempt-1))
			time.Sleep(delay)
		}

		resp, err := c.httpClient.Get(c.baseURL + endpoint)
		if err != nil {
			lastErr = err
			// Check if error is retryable (network error, timeout)
			if !isRetryableError(err, 0) {
				return nil, fmt.Errorf("failed to make request: %w", err)
			}
			continue
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()

		if err != nil {
			lastErr = err
			if !isRetryableError(err, resp.StatusCode) {
				return nil, fmt.Errorf("failed to read response body: %w", err)
			}
			continue
		}

		if resp.StatusCode != http.StatusOK {
			lastStatusCode = resp.StatusCode
			lastErr = fmt.Errorf("unexpected status code: %d", resp.StatusCode)
			// Don't retry on 4xx client errors
			if !isRetryableError(nil, resp.StatusCode) {
				return nil, lastErr
			}
			continue
		}

		return body, nil
	}

	return nil, fmt.Errorf("failed after %d attempts (last status: %d): %w", c.maxRetries, lastStatusCode, lastErr)
}

// GetFundsList retrieves the list of all available funds
func (c *Client) GetFundsList() ([]models.Fund, error) {
	body, err := c.doRequest("/funds")
	if err != nil {
		return nil, err
	}

	var response models.FundsResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if !response.Status {
		return nil, fmt.Errorf("API returned error status")
	}

	return response.Data, nil
}

// GetSymbolInfo retrieves symbol information for a specific fund
func (c *Client) GetSymbolInfo(symbol string) (*models.SymbolInfo, error) {
	params := url.Values{}
	params.Set("symbol", symbol)

	body, err := c.doRequest("/tv/symbols?" + params.Encode())
	if err != nil {
		return nil, err
	}

	var symbolInfo models.SymbolInfo
	if err := json.Unmarshal(body, &symbolInfo); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &symbolInfo, nil
}

// GetHistoricalPrices retrieves historical OHLCV bars for a fund
func (c *Client) GetHistoricalPrices(symbol, resolution string, from, to time.Time) (*models.BarsResponse, error) {
	params := url.Values{}
	params.Set("symbol", symbol)
	params.Set("resolution", resolution)
	params.Set("from", strconv.FormatInt(from.Unix(), 10))
	params.Set("to", strconv.FormatInt(to.Unix(), 10))

	body, err := c.doRequest("/tv/history?" + params.Encode())
	if err != nil {
		return nil, err
	}

	var bars models.BarsResponse
	if err := json.Unmarshal(body, &bars); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &bars, nil
}

// GetFundLatest retrieves the latest data for a specific fund
func (c *Client) GetFundLatest(fundID string) (*models.FundLatest, error) {
	body, err := c.doRequest(fmt.Sprintf("/funds/%s/latest", fundID))
	if err != nil {
		return nil, err
	}

	var response models.FundLatestResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if !response.Status {
		return nil, fmt.Errorf("API returned error status")
	}

	return &response.Data, nil
}

// GetFundPerformance retrieves performance data for a specific fund
func (c *Client) GetFundPerformance(fundID string) (*models.FundPerformance, error) {
	body, err := c.doRequest(fmt.Sprintf("/funds/%s/performance", fundID))
	if err != nil {
		return nil, err
	}

	var response models.FundPerformanceResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if !response.Status {
		return nil, fmt.Errorf("API returned error status")
	}

	return &response.Data, nil
}

// GetFundOverview retrieves overview data for a specific fund
func (c *Client) GetFundOverview(fundID string) (*models.FundOverview, error) {
	body, err := c.doRequest(fmt.Sprintf("/funds/%s/3d", fundID))
	if err != nil {
		return nil, err
	}

	var response models.FundOverviewResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if !response.Status {
		return nil, fmt.Errorf("API returned error status")
	}

	return &response.Data, nil
}

// GetFundFee retrieves fee information for a specific fund
func (c *Client) GetFundFee(fundID string) (*models.FundFee, error) {
	body, err := c.doRequest(fmt.Sprintf("/funds/%s/fee", fundID))
	if err != nil {
		return nil, err
	}

	var response models.FundFeeResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if !response.Status {
		return nil, fmt.Errorf("API returned error status")
	}

	return &response.Data, nil
}

// GetFundVerify retrieves verification data including available periods for a fund
func (c *Client) GetFundVerify(fundID string) (*models.FundVerify, error) {
	body, err := c.doRequest(fmt.Sprintf("/funds/%s/nav/verify", fundID))
	if err != nil {
		return nil, err
	}

	var response models.FundVerifyResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if !response.Status {
		return nil, fmt.Errorf("API returned error status")
	}

	return &response.Data, nil
}

// GetServerTime retrieves the server time from the API
func (c *Client) GetServerTime() (int64, error) {
	body, err := c.doRequest("/tv/time")
	if err != nil {
		return 0, err
	}

	var timestamp int64
	if err := json.Unmarshal(body, &timestamp); err != nil {
		return 0, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return timestamp, nil
}

// SearchFund searches for a fund by short code or name
func (c *Client) SearchFund(query string) (*models.Fund, error) {
	funds, err := c.GetFundsList()
	if err != nil {
		return nil, err
	}

	for _, fund := range funds {
		if fund.ShortCode == query || fund.FundID == query {
			return &fund, nil
		}
	}

	return nil, fmt.Errorf("fund not found: %s", query)
}

func (c *Client) GetFundPortfolio(fundID string) (*models.FundPortfolio, error) {
	body, err := c.doRequest(fmt.Sprintf("/funds/%s/portfolio", fundID))
	if err != nil {
		return nil, err
	}

	var response models.FundPortfolioResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if !response.Status {
		return nil, fmt.Errorf("API returned error status")
	}

	return &response.Data, nil
}
