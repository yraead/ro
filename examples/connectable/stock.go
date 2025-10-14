// Copyright 2025 samber.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// https://github.com/samber/ro/blob/main/licenses/LICENSE.apache.md
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.


package main

import (
	"encoding/json"
	"io"
	"net/http"
)

// {
// 	"Meta Data": {
// 		"1. Information": "Daily Time Series with Splits and Dividend Events",
// 		"2. Symbol": "MSFT",
// 		"3. Last Refreshed": "2024-12-26",
// 		"4. Output Size": "Compact",
// 		"5. Time Zone": "US/Eastern"
// 	},
// 	"Time Series (Daily)": {
// 		"2024-12-26": {
// 			"1. open": "438.94",
// 			"2. high": "440.92",
// 			"3. low": "436.63",
// 			"4. close": "438.11",
// 			"5. adjusted close": "438.11",
// 			"6. volume": "8180780",
// 			"7. dividend amount": "0.0000",
// 			"8. split coefficient": "1.0"
// 		},
// 		"2024-12-24": {
// 			"1. open": "434.65",
// 			"2. high": "439.6",
// 			"3. low": "434.19",
// 			"4. close": "439.33",
// 			"5. adjusted close": "439.33",
// 			"6. volume": "7164543",
// 			"7. dividend amount": "0.0000",
// 			"8. split coefficient": "1.0"
// 		},
// 		"2024-12-23": {
// 			"1. open": "436.74",
// 			"2. high": "437.65",
// 			"3. low": "432.83",
// 			"4. close": "435.25",
// 			"5. adjusted close": "435.25",
// 			"6. volume": "19152519",
// 			"7. dividend amount": "0.0000",
// 			"8. split coefficient": "1.0"
// 		}
// 	}
// }

type stock struct {
	TS map[string]map[string]string `json:"Time Series (Daily)"`
}

func getMSFTStock() (map[string]map[string]string, error) {
	// send http request

	url := "https://www.alphavantage.co/query?function=TIME_SERIES_DAILY_ADJUSTED&symbol=MSFT&apikey=demo"
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	// parse json response
	var result stock

	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return result.TS, nil
}
