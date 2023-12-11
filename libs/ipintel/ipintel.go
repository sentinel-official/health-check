package ipintel

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/brianvoe/gofakeit/v6"
)

func Score(transport *http.Transport, ipAddr string) (float64, error) {
	client := &http.Client{
		Transport: transport,
		Timeout:   15 * time.Second,
	}

	fakeit := gofakeit.New(time.Now().UnixNano())
	urlPath := fmt.Sprintf("https://check.getipintel.net/check.php?ip=%s&contact=%s&format=json&flags=m", ipAddr, fakeit.Email())

	resp, err := client.Get(urlPath)
	if err != nil {
		return 0, err
	}

	defer func() {
		if err = resp.Body.Close(); err != nil {
			panic(err)
		}
	}()

	var m map[string]interface{}
	if err = json.NewDecoder(resp.Body).Decode(&m); err != nil {
		return 0, err
	}

	result, err := strconv.ParseFloat(m["result"].(string), 64)
	if err != nil {
		return 0, err
	}

	return result, nil
}
