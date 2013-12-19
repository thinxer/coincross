package btcchina

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

func (bc *BTCChina) request(method string, params []interface{}, reply interface{}) (err error) {
	tonce := time.Now().UnixNano() / 1000
	data := map[string]interface{}{
		"id":            fmt.Sprintf("%d", tonce),
		"tonce":         tonce,
		"accesskey":     bc.apikey,
		"requestmethod": "post",
		"method":        method,
		"params":        params,
	}

	var message bytes.Buffer
	fields := strings.Split("tonce accesskey requestmethod id method", " ")
	for _, field := range fields {
		message.WriteString(fmt.Sprintf("%s=%v&", field, data[field]))
	}
	message.WriteString(fmt.Sprintf("params=%s", php_implode(params)))
	h := hmac.New(sha1.New, bc.secret)
	h.Write(message.Bytes())
	digest := hex.EncodeToString(h.Sum(nil))

	data_json, _ := json.Marshal(data)
	req, _ := http.NewRequest("POST", ENDPOINT, bytes.NewReader(data_json))
	req.SetBasicAuth(bc.apikey, digest)
	req.Header.Set("Json-Rpc-Tonce", fmt.Sprintf("%d", tonce))
	r, err := bc.client.Do(req)
	if err == nil {
		decoder := json.NewDecoder(r.Body)
		var response struct {
			Result interface{}
			Id     string
		}
		response.Result = reply
		err = decoder.Decode(&response)
		r.Body.Close()
	}
	return
}

func getjson(client *http.Client, url string, v interface{}) (err error) {
	res, err := client.Get(url)
	if err != nil {
		return
	}
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(v)
	res.Body.Close()
	return
}

func php_float(v interface{}) string {
	s := fmt.Sprintf("%f", v)
	s = strings.TrimRight(s, "0")
	s = strings.TrimRight(s, ".")
	return s
}

func php_implode(values []interface{}) string {
	parts := make([]string, 0)
	for _, v := range values {
		switch v := v.(type) {
		case bool:
			if v {
				parts = append(parts, "1")
			} else {
				parts = append(parts, "")
			}
		case float32, float64:
			parts = append(parts, php_float(v))
		case string:
			parts = append(parts, v)
		default:
			parts = append(parts, fmt.Sprintf("%v", v))
		}
	}
	return strings.Join(parts, ",")
}
