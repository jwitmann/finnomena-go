# Finnomena API Documentation

## Overview

The Finnomena API provides access to Thai mutual fund data. The Go client in `github.com/jwitmann/finnomena-go` wraps these API endpoints.

Base URL: `https://www.finnomena.com/fn3/api/fund/v2/public`

## Go Client

### Creating a Client

```go
import "github.com/jwitmann/finnomena-go"

// Default client (30s timeout)
client := finnomena.NewClient()

// Custom timeout
client := finnomena.NewClientWithTimeout(60 * time.Second)
```

### Client Methods

| Method | Description |
| ------ | ----------- |
| `GetFundsList()` | Returns all available funds |
| `SearchFund(query)` | Find fund by ticker or ID |
| `GetSymbolInfo(symbol)` | Get symbol metadata |
| `GetHistoricalPrices(symbol, resolution, from, to)` | Get OHLCV bars |
| `GetFundLatest(fundID)` | Get latest NAV and change |
| `GetFundPerformance(fundID)` | Get returns, std dev, Sharpe |
| `GetFundOverview(fundID)` | Get 3D metrics (PP, RR, DD) |
| `GetFundFee(fundID)` | Get fund fee structure |
| `GetFundVerify(fundID)` | Get available periods |
| `GetServerTime()` | Get server timestamp |

### Example Usage

```go
client := finnomena.NewClient()

// Get all funds
funds, err := client.GetFundsList()
if err != nil {
    log.Fatal(err)
}

// Search for a fund
fund, err := client.SearchFund("TNEXTGEN-A")
if err != nil {
    log.Fatal(err)
}

// Get latest data
latest, err := client.GetFundLatest(fund.FundID)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("NAV: %.4f, Change: %.2f%%\n", latest.Value, latest.DChange)

// Get historical prices (last 30 days)
to := time.Now()
from := to.AddDate(0, 0, -30)
bars, err := client.GetHistoricalPrices(fund.ShortCode, "D", from, to)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Bars: %d\n", len(bars.Time))

// Get performance
perf, err := client.GetFundPerformance(fund.FundID)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("1Y Return: %.2f%%\n", perf.TotalReturn1Y)

// Get 3D overview
overview, err := client.GetFundOverview(fund.FundID)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("PP: %d, RR: %d, DD: %d\n", overview.PP.Fund, overview.RR.Fund, overview.DD.Fund)
```

## API Endpoints (TradingView-compatible)

It's partly based on the TradingView API

### Data feed configuration data

Request: `GET /tv/config`

Response: Library expects to receive a JSON response of the same structure as a result of JS API.

Also there should be 2 additional properties:

* `supports_search`: Set it to `true` if your data feed supports symbol search and individual symbol resolve logic.
* `supports_group_request`: Set it to `true`  if your data feed provides full information on symbol group only and is not able to perform symbol search or individual symbol resolve.

Either `supports_search` or `supports_group_request` should be set to `true`.

**Remark**: If your data feed doesn't implement this call (doesn't respond or sends 404 error) then the default configuration is being used. Here is the default configuration:

```json
{
    "supported_resolutions": ["1", "5", "15", "30", "60", "1D", "1W", "1M"],
    "supports_group_request": true,
    "supports_marks": false,
    "supports_search": false,
    "supports_timescale_marks": false
}
```

### Funds list

Request: `GET /funds`

A response is an object with the following keys.

* `s`: Status code for the request. Expected values are: `ok` or `error`
* `errmsg`: Error message. Should be present only when `s = 'error'`
* `d`: [symbols data] Array

symbol data:

* `fund_id`: string. Fund ID.
* `short_code`: string. Fund short code aka Symbol name/ticker.
* `name_th`: string. Fund name in Thai.
* `aimc_category_id`: string. AIMC category ID.
* `short_desc`: string. Short description of the fund.
* `is_finnomena_pick`: boolean. Indicates whether the fund is a Finnomena Pick.
* `credit_card_allowed`: boolean. Indicates whether the fund allows credit card purchases.
* `is_in_trending`: boolean. Indicates whether the fund is in trending.
* `sec_is_active`: boolean. Indicates whether the fund is active in SEC.

Example:

```json
{
  "status":true,
  "service_code":"51",
  "data":[
    {
      "fund_id":"F000000QSX",
      "short_code":"ES-SET50-A",
      "name_th":"???????????????????? SET50 (??????????????)",
      "aimc_category_id":"LC00002660",
      "short_desc":"?????????? Index Fund ????????????????????????????????? SET50 Index ???????????????????????????????????????????? SET50 ?????????",
      "is_finnomena_pick":false,
      "credit_card_allowed":false,
      "is_in_trending":false,
      "sec_is_active":true,
      "promotions":[]
    },
    {
      "fund_id":"F000000QU8",
      "short_code":"ABG",
      "name_th":"?????????? ????????? ????",
      "aimc_category_id":"LC00002470",
      "short_desc":null,
      "is_finnomena_pick":false,
      "credit_card_allowed":false,
      "is_in_trending":false,
      "sec_is_active":true,
      "promotions":[]
    },
    {
      "fund_id":"F000000QU9",
      "short_code":"ABTOPP",
      "name_th":"?????????? ????????? ???????? ????????????????",
      "aimc_category_id":"LC00002470",
      "short_desc":null,
      "is_finnomena_pick":false,
      "credit_card_allowed":false,
      "is_in_trending":false,
      "sec_is_active":true,
      "promotions":[]
    }
  ]
}
```

### Symbol resolve

Request: `GET /tv/symbols?symbol=<symbol>`

1. `symbol`: string. Symbol name or ticker.

Example: `GET /tv/symbols?symbol=TNEXTGEN-A`

```json
{
  "name":"TNEXTGEN-A",
  "timezone":"Asia/Bangkok",
  "minmov":1,
  "minmove2":0,
  "pricescale":100,
  "pointvalue":0,
  "ticker":"TNEXTGEN-A",
  "description":"",
  "type":"fund",
  "data_status":"endofday",
  "supported-resolutions":["1D","1W","1M","3M","6M","1Y","3Y","5Y","MAX","YTD"],
  "session":"1000-1630",
  "currency_code":"THB"
}
```

**Remark**: This call will be requested if your data feed sent `supports_group_request: false` and `supports_search: true` in the configuration data.

### Bars

Request: `GET /tv/history?symbol=<ticker_name>&from=<unix_timestamp>&to=<unix_timestamp>&resolution=<resolution>`

* `symbol`: symbol name or ticker.
* `from`: unix timestamp (UTC) of leftmost required bar
* `to`: unix timestamp (UTC) of rightmost required bar
* `resolution`: string

Example: `GET /tv/history?symbol=BEAM~0&resolution=D&from=1386493512&to=1395133512`

A response is expected to be an object with some properties listed below. Each property is treated as a table column, as described above.

* `s`: status code. Expected values: `ok` | `error` | `no_data`
* `errmsg`: Error message. Should be present only when `s = 'error'`
* `t`: Bar time. Unix timestamp (UTC)
* `c`: Closing price
* `o`: Opening price (optional)
* `h`: High price (optional)
* `l`: Low price (optional)
* `v`: Volume (optional)
* `nextTime`: Time of the next bar if there is no data (status code is `no_data`) in the requested period (optional)

**Remark**: Bar time for daily bars should be 00:00 UTC and is expected to be a trading day (not a day when the session starts).
Charting Library aligns the time according to the [Session] from SymbolInfo.

**Remark**: Bar time for monthly bars should be 00:00 UTC and is the first trading day of the month.

**Remark**: Prices should be passed as numbers and not as strings in quotation marks.

Example:

```json
{
   "s" : "ok",
   "t" : [1386493512, 1386493572, 1386493632, 1386493692],
   "c" : [42.1, 43.4, 44.3, 42.8]
}
```

```json
{
   "s": "no_data",
   "nextTime": 1386493512
}
```

```json
{
   "s": "ok",
   "t": [1386493512, 1386493572, 1386493632, 1386493692],
   "c": [42.1, 43.4, 44.3, 42.8],
   "o": [41.0, 42.9, 43.7, 44.5],
   "h": [43.0, 44.1, 44.8, 44.5],
   "l": [40.4, 42.1, 42.8, 42.3],
   "v": [12000, 18500, 24000, 45000]
}
```

#### How `nextTime` works

Let's assume that a user opened the chart where `resolution = 1` and the Library requests the following range of data from the data feed `[3 Apr 2015 16:00 UTC+0, 3 Apr 2015 19:00 UTC+0]` for a stock that is traded on the NYSE.
April 3rd was Good Friday which means that the markets were closed.
Library expects the following response from the data feed.

```json
{
  "s": "no_data",
  "nextTime": 1428001140000
}
```

`nextTime` is the time of the closest available bar in the past.

### Server time

Request: `GET /tv/time`

Response: Numeric unix time without milliseconds.

Example: `1445324591`

### Funds latest data

REQUEST: `GET /funds/<fund id>/latest`

Example:

```json
{
  "status":true,
  "service_code":"51",
  "data":
  {
    "fund_id":"F0000161I6",
    "short_code":"TNEXTGEN-A",
    "is_in_trending":false,
    "date":"2026-02-12T00:00:00Z",
    "value":8.6713,
    "amount":1543950235.81,
    "d_change":-3.16,
    "is_finnomena_pick":false,
    "credit_card_allowed":false,
    "performance_ready":true,
    "operation_ready":true,
    "sec_fund_status":"RG",
    "sec_is_active":true
  }
}
```

### Fund verify

Request: `GET /funds/<fund id>/nav/verify`

Example: `GET /funds/F0000161I6/nav/verify`

```json
{
  "status":true,
  "service_code":"51",
  "data":
  {
    "fund_id":"F0000161I6",
    "short_code":"TNEXTGEN-A",
    "range":["1D","1W","1M","3M","6M","YTD","1Y","3Y","5Y","MAX"],
    "default":"1Y"
  }
}
```

### Fund performance

Request: `GET /funds/<fund id>/performance`

```json
{
  "status": true,
  "service_code": "51",
  "data": {
    "fund_id": "F0000161I6",
    "short_code": "TNEXTGEN-A",
    "day_end_date": "2026-02-12T00:00:00Z",
    "total_return_1w": 0.83024,
    "total_return_1m": -15.94157,
    "total_return_3m": -18.95451,
    "total_return_6m": -14.66348,
    "total_return_1y": -6.00217,
    "total_return_3y": 27.68066,
    "total_return_5y": -11.16339,
    "total_return_10y": null,
    "std_3m": null,
    "std_6m": null,
    "std_1y": 28.646,
    "std_3y": 33.322,
    "std_5y": 37.818,
    "std_10y": null,
    "total_return_p_3m": 95,
    "total_return_p_6m": 95,
    "total_return_p_1y": 95,
    "total_return_p_3y": 25,
    "total_return_p_5y": 95,
    "total_return_p_10y": null,
    "std_p_3m": null,
    "std_p_6m": null,
    "std_p_1y": 95,
    "std_p_3y": 95,
    "std_p_5y": 100,
    "std_p_10y": null,
    "sharpe_ratio_p_1y": 75,
    "sharpe_ratio_p_3y": 50,
    "sharpe_ratio_p_5y": 75,
    "sharpe_ratio_p_10y": null,
    "max_drawdown_p_1y": 95,
    "max_drawdown_p_3y": 75,
    "max_drawdown_p_5y": 100,
    "max_drawdown_p_10y": null,
    "total_return_avg_3m": -5.731619,
    "total_return_avg_6m": 5.272869,
    "total_return_avg_1y": 12.252457,
    "total_return_avg_3y": 19.267344,
    "total_return_avg_5y": 0.193716,
    "total_return_avg_10y": null,
    "std_avg_3m": null,
    "std_avg_6m": null,
    "std_avg_1y": 25.421422,
    "std_avg_3y": 24.194519,
    "std_avg_5y": 22.127897,
    "std_avg_10y": null,
    "sharpe_ratio_avg_1y": 0.659889,
    "sharpe_ratio_avg_3y": 0.855037,
    "sharpe_ratio_avg_5y": 0.171385,
    "sharpe_ratio_avg_10y": null,
    "max_drawdown_avg_1y": -16.687556,
    "max_drawdown_avg_3y": -20.665269,
    "max_drawdown_avg_5y": -46.322385,
    "max_drawdown_avg_10y": null,
    "unit_change_1d": -3.185339,
    "unit_change_1w": -0.838743,
    "unit_change_1m": -21.407493,
    "unit_change_3m": -25.075898,
    "unit_change_6m": -17.227929,
    "unit_change_1y": -24.655595,
    "unit_change_3y": 40.448316,
    "unit_change_5y": -67.343156,
    "unit_change_10y": null,
    "sharpe_ratio_1y": 0.217,
    "sharpe_ratio_3y": 0.924,
    "sharpe_ratio_5y": -0.033,
    "sharpe_ratio_10y": null,
    "sharpe_ratio_15y": null,
    "sharpe_ratio_20y": null,
    "max_drawdown_1y": -19.556,
    "max_drawdown_3y": -22.793,
    "max_drawdown_5y": -75.386,
    "max_drawdown_10y": null,
    "max_drawdown_15y": null,
    "max_drawdown_20y": null,
    "net_assets": 1543950235.81,
    "data_date": "2026-02-12T00:00:00Z",
    "returns_year": [
      {
        "year": 2016,
        "value": null
      },
      {
        "year": 2017,
        "value": null
      },
      {
        "year": 2018,
        "value": null
      },
      {
        "year": 2019,
        "value": null
      },
      {
        "year": 2020,
        "value": null
      },
      {
        "year": 2021,
        "value": -16.52016
      },
      {
        "year": 2022,
        "value": -68.00008
      },
      {
        "year": 2023,
        "value": 86.50075
      },
      {
        "year": 2024,
        "value": 35.54682
      },
      {
        "year": 2025,
        "value": 20.03231
      }
    ]
  }
}
```

### fund overview

request: `GET funds/<fund id>/3d`

Example:

```json
{
  "status":true,
  "service_code":"51",
  "data":
  {
    "fund_id":"F00001KB6Z",
    "short_code":"TISCOAI",
    "aimc_category":"Technology Equity",
    "aimc_category_id":"LC00002858",
    "aimc_category_name_th":"หุ้นกลุ่มเทคโนโลยี",
    "finno_score":0,
    "data_date":"2026-02-15T04:02:51Z",
    "pp":
    {
      "fund":50,
      "avg":57,
      "text":"fair"
    },
    "rr":
    {
      "fund":48,
      "avg":55,
      "text":"fair"
    },
    "dd":
    {
      "fund":41,
      "avg":39,
      "text":"fair"
    }
  }
}
```

### fund fee

request: `GET funds/<fund id>/fee`

Fund fee is a list of fees that are charged by the fund.

Fee description translation from thai to english:

```json
"fees_dict":{
        "ค่าใช้จ่ายอื่นๆ": "other fee",
        "ค่าธรรมเนียมการขายหน่วยลงทุน (Front-end Fee)": "purchase fee",
        "ค่าธรรมเนียมการจัดการ": "management fee",
        "ค่าธรรมเนียมการรับซื้อคืนหน่วยลงทุน (Back-end Fee)": "redemption fee",
        "ค่าธรรมเนียมการสับเปลี่ยนหน่วยลงทุนเข้า (SWITCHING IN)": "switch in fee",
        "ค่าธรรมเนียมการสับเปลี่ยนหน่วยลงทุนออก (SWITCHING OUT)": "switch out fee",
        "ค่าธรรมเนียมและค่าใช้จ่ายรวมทั้งหมด": "total expense ratio",
        "ค่าธรรมเนียมการโอนหน่วยลงทุน": "unit transfer fee",
        "ค่าธรรมเนียมนายทะเบียนหน่วย": "registrar fee",
        "ค่าธรรมเนียมผู้ดูแลผลประโยชน์": "trustee fee"
}
```

Example:

```json
{
  "status": true,
  "service_code": "51",
  "data": {
    "fund_id": "F00001KB6Z",
    "short_code": "TISCOAI",
    "fees": [
      {
        "last_update": null,
        "class_abbr_name": "main",
        "description": "ค่าใช้จ่ายอื่นๆ",
        "rate": "1.30",
        "rate_unit": "ต่อปี ของมูลค่าทรัพย์สินสุทธิของกองทุน",
        "actual_value": "",
        "actual_value_unit": "",
        "other_description": "",
        "unit": "%"
      },
      {
        "last_update": null,
        "class_abbr_name": "main",
        "description": "ค่าธรรมเนียมการขายหน่วยลงทุน (Front-end Fee)",
        "rate": "2.5",
        "rate_unit": "",
        "actual_value": "1",
        "actual_value_unit": "",
        "other_description": "บริษัทจัดการอาจเรียกเก็บค่าธรรมเนียมการขายและค่าธรรมเนียมการรับซื้อคืนกับผู้ลงทุนแต่ละกลุ่มไม่เท่ากัน ทั้งนี้ ผู้ลงทุนสามารถดูรายละเอียดเพิ่มเติมได้ที่หนังสือชี้ชวนส่วนข้อมูลกองทุนรวม",
        "unit": "%"
      },
      {
        "last_update": null,
        "class_abbr_name": "main",
        "description": "ค่าธรรมเนียมการจัดการ",
        "rate": "3",
        "rate_unit": "",
        "actual_value": "1.07",
        "actual_value_unit": "",
        "other_description": "บริษัทจัดการอาจพิจารณาเปลี่ยนแปลงค่าธรรมเนียมที่เรียกเก็บจริงเพื่อให้สอดคล้องกับกลยุทธ์หรือค่าใช้จ่ายในการบริหารจัดการ และรวมค่าใช้จ่ายเป็นข้อมูลของรอบปีบัญชีล่าสุดหรือประมาณการเบื้องต้น (กรณียังไม่ครบรอบปีบัญชี)",
        "unit": "% ต่อปี"
      },
      {
        "last_update": null,
        "class_abbr_name": "main",
        "description": "ค่าธรรมเนียมการรับซื้อคืนหน่วยลงทุน (Back-end Fee)",
        "rate": "2.5",
        "rate_unit": "",
        "actual_value": "0",
        "actual_value_unit": "",
        "other_description": "บริษัทจัดการอาจเรียกเก็บค่าธรรมเนียมการขายและค่าธรรมเนียมการรับซื้อคืนกับผู้ลงทุนแต่ละกลุ่มไม่เท่ากัน ทั้งนี้ ผู้ลงทุนสามารถดูรายละเอียดเพิ่มเติมได้ที่หนังสือชี้ชวนส่วนข้อมูลกองทุนรวม",
        "unit": "%"
      },
      {
        "last_update": null,
        "class_abbr_name": "main",
        "description": "ค่าธรรมเนียมการสับเปลี่ยนหน่วยลงทุนเข้า (SWITCHING IN)",
        "rate": "2.5",
        "rate_unit": "",
        "actual_value": "0",
        "actual_value_unit": "",
        "other_description": "บริษัทจัดการอาจเรียกเก็บค่าธรรมเนียมการขายและค่าธรรมเนียมการรับซื้อคืนกับผู้ลงทุนแต่ละกลุ่มไม่เท่ากัน ทั้งนี้ ผู้ลงทุนสามารถดูรายละเอียดเพิ่มเติมได้ที่หนังสือชี้ชวนส่วนข้อมูลกองทุนรวม",
        "unit": "%"
      },
      {
        "last_update": null,
        "class_abbr_name": "main",
        "description": "ค่าธรรมเนียมการสับเปลี่ยนหน่วยลงทุนออก (SWITCHING OUT)",
        "rate": "2.5",
        "rate_unit": "",
        "actual_value": "0",
        "actual_value_unit": "",
        "other_description": "บริษัทจัดการอาจเรียกเก็บค่าธรรมเนียมการขายและค่าธรรมเนียมการรับซื้อคืนกับผู้ลงทุนแต่ละกลุ่มไม่เท่ากัน ทั้งนี้ ผู้ลงทุนสามารถดูรายละเอียดเพิ่มเติมได้ที่หนังสือชี้ชวนส่วนข้อมูลกองทุนรวม",
        "unit": "%"
      },
      {
        "last_update": null,
        "class_abbr_name": "main",
        "description": "ค่าธรรมเนียมการโอนหน่วยลงทุน",
        "rate": "999999",
        "rate_unit": "",
        "actual_value": "0",
        "actual_value_unit": "",
        "other_description": "บริษัทจัดการอาจเรียกเก็บค่าธรรมเนียมการขายและค่าธรรมเนียมการรับซื้อคืนกับผู้ลงทุนแต่ละกลุ่มไม่เท่ากัน ทั้งนี้ ผู้ลงทุนสามารถดูรายละเอียดเพิ่มเติมได้ที่หนังสือชี้ชวนส่วนข้อมูลกองทุนรวม",
        "unit": "บาท"
      },
      {
        "last_update": null,
        "class_abbr_name": "main",
        "description": "ค่าธรรมเนียมนายทะเบียนหน่วย",
        "rate": "0.50",
        "rate_unit": "ต่อปี ของมูลค่าทรัพย์สินสุทธิของกองทุน",
        "actual_value": "",
        "actual_value_unit": "",
        "other_description": "",
        "unit": "%"
      },
      {
        "last_update": null,
        "class_abbr_name": "main",
        "description": "ค่าธรรมเนียมผู้ดูแลผลประโยชน์",
        "rate": "0.20",
        "rate_unit": "ต่อปี ของมูลค่าทรัพย์สินสุทธิของกองทุน",
        "actual_value": "",
        "actual_value_unit": "",
        "other_description": "",
        "unit": "%"
      },
      {
        "last_update": null,
        "class_abbr_name": "main",
        "description": "ค่าธรรมเนียมและค่าใช้จ่ายรวมทั้งหมด",
        "rate": "5",
        "rate_unit": "",
        "actual_value": "1.27",
        "actual_value_unit": "",
        "other_description": "บริษัทจัดการอาจพิจารณาเปลี่ยนแปลงค่าธรรมเนียมที่เรียกเก็บจริงเพื่อให้สอดคล้องกับกลยุทธ์หรือค่าใช้จ่ายในการบริหารจัดการ และรวมค่าใช้จ่ายเป็นข้อมูลของรอบปีบัญชีล่าสุดหรือประมาณการเบื้องต้น (กรณียังไม่ครบรอบปีบัญชี)",
        "unit": "% ต่อปี"
      }
    ]
  }
}
```

### Fund dividend

request: `GET funds/<fund id>/dividend`

Fund dividend information.

Example:

* No dividend

```json
{
  "status":true,
  "service_code":"51",
  "data":
  {
    "fund":"F0000161I6",
    "short_code":"TNEXTGEN-A",
    "dividends":[]
  }
```

* With dividend

```json
{
  "status": true,
  "service_code": "51",
  "data": {
    "fund": "F000000ROP",
    "short_code": "KFSDIV",
    "dividends": [
      {
        "xd_date": "2025-11-07T00:00:00Z",
        "value": "0.2",
        "pay_date": "2025-11-18T00:00:00Z"
      },
      {
        "xd_date": "2024-11-08T00:00:00Z",
        "value": "0.2",
        "pay_date": "2024-11-19T00:00:00Z"
      },
      {
        "xd_date": "2023-02-10T00:00:00Z",
        "value": "0.25",
        "pay_date": "2023-02-21T00:00:00Z"
      },
      {
        "xd_date": "2022-05-13T00:00:00Z",
        "value": "0.1",
        "pay_date": "2022-05-25T00:00:00Z"
      },
      {
        "xd_date": "2022-02-10T00:00:00Z",
        "value": "0.1",
        "pay_date": "2022-02-22T00:00:00Z"
      },
      {
        "xd_date": "2021-11-08T00:00:00Z",
        "value": "0.1",
        "pay_date": "2021-11-17T00:00:00Z"
      },
      {
        "xd_date": "2021-05-14T00:00:00Z",
        "value": "0.25",
        "pay_date": "2021-05-25T00:00:00Z"
      },
      {
        "xd_date": "2021-02-10T00:00:00Z",
        "value": "0.2",
        "pay_date": "2021-02-22T00:00:00Z"
      },
      {
        "xd_date": "2020-08-13T00:00:00Z",
        "value": "0.5",
        "pay_date": "2020-08-24T00:00:00Z"
      },
      {
        "xd_date": "2019-08-06T00:00:00Z",
        "value": "0.2",
        "pay_date": "2019-08-16T00:00:00Z"
      },
      {
        "xd_date": "2019-05-14T00:00:00Z",
        "value": "0.25",
        "pay_date": "2019-05-24T00:00:00Z"
      },
      {
        "xd_date": "2018-11-08T00:00:00Z",
        "value": "0.5",
        "pay_date": "2018-11-19T00:00:00Z"
      },
      {
        "xd_date": "2018-05-14T00:00:00Z",
        "value": "0.1",
        "pay_date": "2018-05-23T00:00:00Z"
      },
      {
        "xd_date": "2018-02-09T00:00:00Z",
        "value": "0.4",
        "pay_date": "2018-02-20T00:00:00Z"
      },
      {
        "xd_date": "2017-11-06T00:00:00Z",
        "value": "0.5",
        "pay_date": "2017-11-15T00:00:00Z"
      },
      {
        "xd_date": "2017-08-11T00:00:00Z",
        "value": "0.2",
        "pay_date": "2017-08-23T00:00:00Z"
      },
      {
        "xd_date": "2017-05-15T00:00:00Z",
        "value": "0.25",
        "pay_date": "2017-05-24T00:00:00Z"
      },
      {
        "xd_date": "2017-02-10T00:00:00Z",
        "value": "0.1",
        "pay_date": "2017-02-22T00:00:00Z"
      },
      {
        "xd_date": "2016-11-08T00:00:00Z",
        "value": "0.25",
        "pay_date": "2016-11-17T00:00:00Z"
      },
      {
        "xd_date": "2016-05-16T00:00:00Z",
        "value": "0.5",
        "pay_date": "2016-05-26T00:00:00Z"
      },
      {
        "xd_date": "2016-02-10T00:00:00Z",
        "value": "0.3",
        "pay_date": "2016-02-19T00:00:00Z"
      },
      {
        "xd_date": "2015-11-09T00:00:00Z",
        "value": "0.3",
        "pay_date": "2015-11-18T00:00:00Z"
      },
      {
        "xd_date": "2015-08-13T00:00:00Z",
        "value": "0.3",
        "pay_date": "2015-08-24T00:00:00Z"
      },
      {
        "xd_date": "2015-05-14T00:00:00Z",
        "value": "0.3",
        "pay_date": "2015-05-25T00:00:00Z"
      },
      {
        "xd_date": "2015-02-10T00:00:00Z",
        "value": "0.3",
        "pay_date": "2015-02-19T00:00:00Z"
      },
      {
        "xd_date": "2014-11-07T00:00:00Z",
        "value": "0.3",
        "pay_date": "2014-11-18T00:00:00Z"
      },
      {
        "xd_date": "2014-08-13T00:00:00Z",
        "value": "0.4",
        "pay_date": "2014-08-22T00:00:00Z"
      },
      {
        "xd_date": "2014-05-14T00:00:00Z",
        "value": "0.3",
        "pay_date": "2014-05-23T00:00:00Z"
      },
      {
        "xd_date": "2014-02-10T00:00:00Z",
        "value": "0.4",
        "pay_date": "2014-02-20T00:00:00Z"
      },
      {
        "xd_date": "2013-11-08T00:00:00Z",
        "value": "0.4",
        "pay_date": "2013-11-19T00:00:00Z"
      },
      {
        "xd_date": "2013-08-13T00:00:00Z",
        "value": "0.5",
        "pay_date": "2013-08-22T00:00:00Z"
      },
      {
        "xd_date": "2013-04-19T00:00:00Z",
        "value": "1",
        "pay_date": "2013-04-30T00:00:00Z"
      },
      {
        "xd_date": "2013-02-08T00:00:00Z",
        "value": "0.75",
        "pay_date": "2013-02-22T00:00:00Z"
      },
      {
        "xd_date": "2012-11-08T00:00:00Z",
        "value": "1",
        "pay_date": "2012-11-22T00:00:00Z"
      },
      {
        "xd_date": "2012-08-10T00:00:00Z",
        "value": "0.5",
        "pay_date": "2012-08-27T00:00:00Z"
      },
      {
        "xd_date": "2012-05-14T00:00:00Z",
        "value": "0.5",
        "pay_date": "2012-05-28T00:00:00Z"
      },
      {
        "xd_date": "2012-02-14T00:00:00Z",
        "value": "0.5",
        "pay_date": "2012-02-28T00:00:00Z"
      },
      {
        "xd_date": "2011-11-10T00:00:00Z",
        "value": "0.5",
        "pay_date": "2011-11-24T00:00:00Z"
      },
      {
        "xd_date": "2011-08-11T00:00:00Z",
        "value": "0.5",
        "pay_date": "2011-08-26T00:00:00Z"
      },
      {
        "xd_date": "2011-05-13T00:00:00Z",
        "value": "0.5",
        "pay_date": "2011-05-31T00:00:00Z"
      },
      {
        "xd_date": "2011-02-14T00:00:00Z",
        "value": "0.5",
        "pay_date": "2011-03-01T00:00:00Z"
      },
      {
        "xd_date": "2010-11-08T00:00:00Z",
        "value": "1",
        "pay_date": "2010-11-22T00:00:00Z"
      },
      {
        "xd_date": "2010-08-11T00:00:00Z",
        "value": "0.5",
        "pay_date": "2010-08-27T00:00:00Z"
      },
      {
        "xd_date": "2010-05-17T00:00:00Z",
        "value": "0.5",
        "pay_date": "2010-06-01T00:00:00Z"
      },
      {
        "xd_date": "2010-02-12T00:00:00Z",
        "value": "0.3",
        "pay_date": "2010-02-26T00:00:00Z"
      },
      {
        "xd_date": "2009-08-11T00:00:00Z",
        "value": "0.5",
        "pay_date": "2009-08-19T00:00:00Z"
      },
      {
        "xd_date": "2008-02-11T00:00:00Z",
        "value": "0.4",
        "pay_date": "2008-02-20T00:00:00Z"
      }
    ]
  }
}
```

### Fund portfolio

request: `GET funds/<fund id>/portfolio`

Fund portfolio information, i.e. what it is composed of (allocations).

if the name is in thai, it's prefixed with "หุ้นสามัญของ" : "ordinary shares of", the rest is the SET company name.
We can convert the thai name into english using the SET data API. See [SET data API](SET.md) for more information.

Example:

```json
{
  "status": true,
  "service_code": "51",
  "data": {
    "fund_id": "F00001KO3X",
    "short_code": "THDRMF-P",
    "top_holdings": {
      "data_date": "2026-01-01T00:00:00Z",
      "elements": [
        {
          "name": "หุ้นสามัญของบริษัท ปตท. จำกัด (มหาชน)",
          "percent": 10.09,
          "short_code": "",
          "link_url": "https://www.google.com/search?q=",
          "color": "#005125"
        },
        {
          "name": "หุ้นสามัญของธนาคารกสิกรไทย จำกัด (มหาชน)",
          "percent": 9.09,
          "short_code": "",
          "link_url": "https://www.google.com/search?q=",
          "color": "#007f3b"
        },
        {
          "name": "หุ้นสามัญของธนาคารกรุงไทย จำกัด (มหาชน)",
          "percent": 8.85,
          "short_code": "",
          "link_url": "https://www.google.com/search?q=",
          "color": "#00ad50"
        },
        {
          "name": "หุ้นสามัญของบริษัท เอสซีบี เอกซ์ จำกัด (มหาชน)",
          "percent": 8.73,
          "short_code": "",
          "link_url": "https://www.google.com/search?q=",
          "color": "#00e76b"
        },
        {
          "name": "หุ้นสามัญของธนาคารกรุงเทพ จำกัด (มหาชน)",
          "percent": 7.04,
          "short_code": "",
          "link_url": "https://www.google.com/search?q=",
          "color": "#80f3b5"
        }
      ]
    },
    "global_stock_sector": {
      "data_date": "2025-12-31T00:00:00Z",
      "elements": [
        {
          "name": "บริการด้านการเงิน",
          "percent": 48.87655,
          "color": "#bf74d2"
        },
        {
          "name": "พลังงาน",
          "percent": 18.47445,
          "color": "#803592"
        },
        {
          "name": "สินค้าฟุ่มเฟือย/ตามวัฏจักร",
          "percent": 12.38514,
          "color": "#552362"
        },
        {
          "name": "อสังหาริมทรัพย์",
          "percent": 10.71974,
          "color": "#286880"
        },
        {
          "name": "วัสดุทั่วไป",
          "percent": 5.47055,
          "color": "#7cdbff"
        },
        {
          "name": "สาธารณูปโภคพื้นฐาน",
          "percent": 2.79443,
          "color": "#40ed90"
        },
        {
          "name": "การแพทย์",
          "percent": 1.27914,
          "color": "#00ad50"
        }
      ]
    },
    "asset_allocation": {
      "data_date": "2026-01-01T00:00:00Z",
      "elements": [
        {
          "name": "หุ้น",
          "percent": 99.65927,
          "color": "#e6ed39"
        },
        {
          "name": "เงินฝากธนาคาร P/N และ B/E",
          "percent": 1.04099,
          "color": "#007436"
        },
        {
          "name": "สินทรัพย์อื่นๆ/หนี้สินอื่นๆ",
          "percent": -0.70026,
          "color": "#4cc5f2"
        }
      ]
    },
    "regional_exposure": {
      "data_date": "2025-12-31T00:00:00Z",
      "elements": [
        {
          "name": "Asia - Emerging",
          "percent": 100
        },
        {
          "name": "Emerging Market",
          "percent": 100
        }
      ]
    }
  }
}
```

### Fund portfolio (with translation)

Request: `GET /funds/<fund id>/portfolio`

This endpoint returns the fund's portfolio composition including top holdings, sector allocation, asset allocation, and regional exposure.

When `UseEnglishNames` setting is enabled (default), Thai company names in `top_holdings` and Thai sector names in `global_stock_sector` are automatically translated to English using the SET data API.

**Response Structure:**

```json
{
  "fund_id": "F00001KO3X",
  "short_code": "THDRMF-P",
  "top_holdings": {
    "data_date": "2026-01-01T00:00:00Z",
    "elements": [
      {
        "name": "PTT PUBLIC COMPANY LIMITED",
        "percent": 10.09,
        "short_code": "",
        "link_url": "https://www.google.com/search?q=",
        "color": "#005125"
      }
    ]
  },
  "global_stock_sector": {
    "data_date": "2025-12-31T00:00:00Z",
    "elements": [
      {
        "name": "Financial Services",
        "percent": 48.87655,
        "color": "#bf74d2"
      }
    ]
  },
  "asset_allocation": {
    "data_date": "2026-01-01T00:00:00Z",
    "elements": [
      {
        "name": "หุ้น",
        "percent": 99.65927,
        "color": "#e6ed39"
      }
    ]
  },
  "regional_exposure": {
    "data_date": "2025-12-31T00:00:00Z",
    "elements": [
      {
        "name": "Asia - Emerging",
        "percent": 100
      }
    ]
  }
}
```

**Translation:**

- `top_holdings.elements[].name`: Thai company names are translated to English (e.g., "หุ้นสามัญของบริษัท ปตท. จำกัด (มหาชน)" → "PTT PUBLIC COMPANY LIMITED")
- `global_stock_sector.elements[].name`: Thai sector names are translated to English (e.g., "บริการด้านการเงิน" → "Financial Services")
- `asset_allocation.elements[].name`: Not translated (generic Thai terms)
- `regional_exposure.elements[].name`: Not translated (already in English)
