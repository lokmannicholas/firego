package remoteConfig

import (
	"compress/gzip"
	"net/http"

	"github.com/lokmannicholas/firego"

	"time"

	"fmt"

	"io"

	"bytes"

	"encoding/json"
)

var httpClient = &http.Client{Transport: &http.Transport{
	MaxIdleConns:       10,
	IdleConnTimeout:    30 * time.Second,
	DisableCompression: true,
}}

func GetRemoteConfig(v interface{}) error {
	httpClient := &http.Client{Transport: &http.Transport{
		MaxIdleConns:       10,
		IdleConnTimeout:    30 * time.Second,
		DisableCompression: true,
	},
		Timeout: time.Second * 10}
	url := fmt.Sprintf("https://firebaseremoteconfig.googleapis.com/v1/projects/%s/remoteConfig", firebase.GetConfig().ProjectID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	token, err := firebase.GetCredentialToken()
	if err != nil {
		return err
	}
	//req.Proto = "HTTP/1.1"
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Add("Accept-Encoding", "gzip")
	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	r, err := gzip.NewReader(resp.Body)
	var out bytes.Buffer
	io.Copy(&out, r)
	r.Close()

	if json.Valid(out.Bytes()) {

		err = json.Unmarshal(out.Bytes(), v)
		if err != nil {
			return err
		}

	}
	return nil
}
