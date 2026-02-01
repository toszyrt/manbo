package xiaoai

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func NewMiClient(config *XiaoaiConfig) *MiClient {
	c := &MiClient{
		HTTP: &http.Client{
			Timeout: 30 * time.Second,
		},
		Config:        config,
		UserAgentMiot: "iOS-14.4-6.0.103-iPhone12,3",
		UserAgentMina: "MiHome/6.0.103 (com.xiaomi.mihome; iOS 14.4.0)",
	}
	if c.Config == nil {
		exePath, err := os.Executable()
		if err != nil {
			log.Fatal(err)
		}
		dir := filepath.Dir(exePath)
		tokenPath := filepath.Join(dir, "xiaoai.json")
		data, err := os.ReadFile(tokenPath)
		if err != nil {
			log.Fatal(err)
		}
		var cfg XiaoaiConfig
		if err := json.Unmarshal(data, &cfg); err != nil {
			log.Fatal(err)
		}
		c.Config = &cfg
	}
	if !allStringsNonEmpty(c.Config) {
		log.Fatal("XiaoaiConfig has empty fields")
	}
	return c
}

func (c *MiClient) miotPost(api, uri string, payload map[string]any) ([]byte, error) {
	url := api + uri
	form, err := signData(uri, payload, c.Config.Miot.Ssecurity)
	req, _ := http.NewRequest(http.MethodPost, url, strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", c.UserAgentMiot)
	req.Header.Set("x-xiaomi-protocal-flag-cli", "PROTOCAL-HTTP2")
	req.AddCookie(&http.Cookie{Name: "userId", Value: c.Config.UserID})
	req.AddCookie(&http.Cookie{Name: "serviceToken", Value: c.Config.Miot.ServiceToken})

	log.Printf("[DEBUG] MiClient miot http-post ⬆ send: %s, %s\n", url, form.Encode())
	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("http %d", resp.StatusCode)
	}
	log.Printf("[DEBUG] MiClient miot http-post ⬇ recv: %s\n", body)
	return body, nil
}

func (c *MiClient) MiotDeviceList() (*MiotDeviceListResponse, error) {
	api := "https://api.io.mi.com/app"
	uri := "/home/device_list"
	payload := map[string]any{
		"getVirtualModel": false,
		"getHuamiDevices": 0,
	}
	body, err := c.miotPost(api, uri, payload)
	if err != nil {
		return nil, err
	}
	var resp MiotDeviceListResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *MiClient) MiotAction(did string, siid, aiid int, in []any) error {
	api := "https://api.io.mi.com/app"
	uri := "/miotspec/action"
	payload := map[string]any{
		"params": map[string]any{
			"did":  did,
			"siid": siid,
			"aiid": aiid,
			"in":   in,
		},
	}
	_, err := c.miotPost(api, uri, payload)
	if err != nil {
		return err
	}
	return nil
}

func (c *MiClient) minaGet(api, uri string, cookies ...*http.Cookie) ([]byte, error) {
	url := api + uri
	req, _ := http.NewRequest(http.MethodGet, url, nil)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.UserAgentMina)
	req.AddCookie(&http.Cookie{Name: "userId", Value: c.Config.UserID})
	req.AddCookie(&http.Cookie{Name: "micoapi_slh", Value: c.Config.Mina.MicoapiSlh})
	req.AddCookie(&http.Cookie{Name: "micoapi_ph", Value: c.Config.Mina.MicoapiPh})
	req.AddCookie(&http.Cookie{Name: "serviceToken", Value: c.Config.Mina.ServiceToken})
	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}

	log.Printf("[DEBUG] MiClient mina http-get ⬆ send: %s\n", url)
	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("http %d", resp.StatusCode)
	}
	log.Printf("[DEBUG] MiClient mina http-get ⬇ recv: %s\n", body)
	return body, nil
}

func (c *MiClient) MinaDeviceList() (*MinaDeviceListResponse, error) {
	api := "https://api.mina.mi.com"
	uri := "/admin/v2/device_list?master=0"
	body, err := c.minaGet(api, uri)
	if err != nil {
		return nil, err
	}
	var resp MinaDeviceListResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (c *MiClient) MinaGetLatestConversation(hardware, deviceID string) (*MinaConversationData, error) {
	q := url.Values{}
	q.Set("source", "dialogu")
	q.Set("hardware", hardware)
	q.Set("timestamp", fmt.Sprint(time.Now().UnixMilli()))
	q.Set("limit", "1")
	api := "https://userprofile.mina.mi.com"
	uri := "/device_profile/v2/conversation?" + q.Encode()
	cookies := []*http.Cookie{
		{Name: "deviceId", Value: deviceID},
	}
	body, err := c.minaGet(api, uri, cookies...)
	if err != nil {
		return nil, err
	}
	var resp MinaGetLatestConversationResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		return nil, err
	}
	var data MinaConversationData
	if err := json.Unmarshal([]byte(resp.Data), &data); err != nil {
		return nil, err
	}
	return &data, nil
}
