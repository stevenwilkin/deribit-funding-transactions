package deribit

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"time"
)

type Deribit struct {
	ApiId     string
	ApiSecret string
	Test      bool
}

func (d *Deribit) hostname() string {
	if d.Test {
		return "test.deribit.com"
	} else {
		return "www.deribit.com"
	}
}

func (d *Deribit) get(path string, params url.Values, result interface{}) error {
	u := url.URL{
		Scheme:   "https",
		Host:     d.hostname(),
		Path:     path,
		RawQuery: params.Encode()}

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		slog.Warn(err.Error())
		return err
	}

	token := base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%s:%s", d.ApiId, d.ApiSecret)))

	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", token))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		slog.Warn(err.Error())
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		slog.Warn(err.Error())
		return err
	}

	json.Unmarshal(body, result)

	return nil
}

func (d *Deribit) GetTransactions() []Transaction {
	var response transactionsResponse

	startTime := time.Now().AddDate(0, 0, -30).UnixNano() / int64(time.Millisecond)
	endTime := time.Now().UnixNano() / int64(time.Millisecond)

	v := url.Values{
		"currency":        {"BTC"},
		"start_timestamp": {fmt.Sprintf("%d", startTime)},
		"end_timestamp":   {fmt.Sprintf("%d", endTime)},
		"query":           {"settlement"}}

	err := d.get("/api/v2/private/get_transaction_log", v, &response)

	if err != nil {
		return []Transaction{}
	}

	return response.Result.Logs
}

func NewDeribitFromEnv() *Deribit {
	return &Deribit{
		ApiId:     os.Getenv("DERIBIT_API_ID"),
		ApiSecret: os.Getenv("DERIBIT_API_SECRET")}
}
