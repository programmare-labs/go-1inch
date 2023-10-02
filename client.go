package go1inch

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
)

const (
	inchURL = "https://api.1inch.dev/v5.2/"
)

type Network string

const (
	Eth         Network = "eth"
	Bsc         Network = "bsc"
	Matic       Network = "matic"
	Optimism    Network = "optimism"
	Arbitrum    Network = "arbitrum"
	GnosisChain Network = "gnosis"
	Avalanche   Network = "avalanche"
	Fantom      Network = "fantom"
	Klaytn      Network = "klaytn"
	Aurora      Network = "aurora"
	ZkSyncEra   Network = "zksync"
	Base        Network = "base"
)

var (
	networks = map[Network]string{
		Eth:         "1",
		Bsc:         "56",
		Matic:       "137",
		Optimism:    "10",
		Arbitrum:    "42161",
		GnosisChain: "100",
		Avalanche:   "43114",
		Fantom:      "250",
		Klaytn:      "8217",
		Aurora:      "1313161554",
		ZkSyncEra:   "324",
		Base:        "8453",
	}
)

func NewClient() *Client {
	return &Client{
		Http: &http.Client{},
	}
}

func setQueryParam(endpoint *string, params []map[string]interface{}) {
	var first = true
	for _, param := range params {
		for i := range param {
			if first {
				*endpoint = fmt.Sprintf("%s?%s=%v", *endpoint, i, param[i])
				first = false
			} else {
				*endpoint = fmt.Sprintf("%s&%s=%v", *endpoint, i, param[i])
			}
		}
	}
}

func (c *Client) doRequest(ctx context.Context, net Network, endpoint, method string, expRes interface{}, reqData interface{}, opts ...map[string]interface{}) (int, error) {
	n, ok := networks[net]
	if !ok {
		return 0, errors.New("invalid network")
	}
	callURL := fmt.Sprintf("%s%s%s", inchURL, n, endpoint)

	var dataReq []byte
	var err error

	if reqData != nil {
		dataReq, err = json.Marshal(reqData)
		if err != nil {
			return 0, err
		}
	}

	if len(opts) > 0 && len(opts[0]) > 0 {
		setQueryParam(&callURL, opts)
	}
	req, err := http.NewRequestWithContext(ctx, method, callURL, bytes.NewBuffer(dataReq))
	if err != nil {
		return 0, err
	}

	req.Header.Add("Content-type", "application/json")

	apiKEY := fmt.Sprintf("Bearer %s", os.Getenv("ONEINCH_API_KEY"))
	req.Header.Add("Authorization", apiKEY)

	resp, err := c.Http.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	switch resp.StatusCode {
	case 200:
		if expRes != nil {
			err = json.Unmarshal(body, expRes)
			if err != nil {
				return 0, err
			}
		}
		return resp.StatusCode, nil

	default:
		return resp.StatusCode, fmt.Errorf("%s", body)
	}
}
