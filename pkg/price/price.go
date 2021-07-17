package price

import (
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
)

const (
	coinbaseBuyPriceUrl = "https://api.coinbase.com/v2/prices/ETH-USD/buy"
)

type (
	coinbaseBuyPriceResponseBody struct {
		Data struct {
			Base     string  `json:"base"`
			Currency string  `json:"currency"`
			Amount   float64 `json:"amount,string"`
		} `json:"data"`
	}
)

func UsdPerEth() (float64, error) {
	req, err := http.NewRequest(http.MethodGet, coinbaseBuyPriceUrl, nil)
	if err != nil {
		return 0, errors.Wrap(err, "building request")
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return 0, errors.Wrap(err, "executing request")
	}
	defer func() { _ = res.Body.Close() }()

	var body coinbaseBuyPriceResponseBody
	if err = json.NewDecoder(res.Body).Decode(&body); err != nil {
		return 0, errors.Wrap(err, "decoding response body")
	}

	return body.Data.Amount, nil
}
