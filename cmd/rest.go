package cmd

import (
	"context"
	"encoding/xml"
	"fmt"
	"github.com/spf13/viper"
	"golang.org/x/net/context/ctxhttp"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"time"
)

type ArinApiError struct {
	XMLName        xml.Name `xml:"error"`
	AdditionalInfo string   `xml:"additionalInfo"`
	Code           string   `xml:"code"`
	Components     string   `xml:"components"`
	Message        string   `xml:"message"`
}

func restDelete(ctx context.Context, url string, result interface{}) error {
	var apierr ArinApiError

	client := http.Client{Timeout: time.Second * 10}
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}
	response, err := ctxhttp.Do(ctx, &client, req)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}
	if flagShowHTTPResult {
		fmt.Fprintln(os.Stderr, string(body))
	}
	if err = xml.Unmarshal(body, &apierr); err == nil {
		return fmt.Errorf("%s", apierr.Message)
	}
	if result != nil {
		return xml.Unmarshal(body, result)
	}
	return nil
}

func restGet(ctx context.Context, url string, result interface{}) error {
	var apierr ArinApiError

	client := http.Client{Timeout: time.Second * 10}
	response, err := ctxhttp.Get(ctx, &client, url)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return err
	}
	if flagShowHTTPResult {
		fmt.Fprintln(os.Stderr, string(body))
	}
	if err = xml.Unmarshal(body, &apierr); err == nil {
		return fmt.Errorf("%s", apierr.Message)
	}
	return xml.Unmarshal(body, result)
}

func makeUrl(api string, extra ...string) string {
	s := "https://reg.arin.net/" + api
	if extra != nil {
		for _, p := range extra {
			s += "/" + url.PathEscape(p)
		}
	}
	u, _ := url.Parse(s)
	v := make(url.Values)
	v.Add("apikey", viper.GetString("apikey"))
	u.RawQuery = v.Encode()
	return u.String()
}
