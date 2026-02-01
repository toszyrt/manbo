package xiaoai

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"math/rand"
	"net/url"
	"reflect"
	"time"
)

func allStringsNonEmpty(v any) bool {
	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Pointer {
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		return false
	}

	for i := 0; i < val.NumField(); i++ {
		f := val.Field(i)
		switch f.Kind() {
		case reflect.String:
			if f.String() == "" {
				return false
			}
		case reflect.Struct:
			if !allStringsNonEmpty(f.Interface()) {
				return false
			}
		}
	}
	return true
}

func signNonce(ssecurityB64, nonceB64 string) (string, error) {
	ssec, err := base64.StdEncoding.DecodeString(ssecurityB64)
	if err != nil {
		return "", err
	}
	nonce, err := base64.StdEncoding.DecodeString(nonceB64)
	if err != nil {
		return "", err
	}
	h := sha256.New()
	h.Write(ssec)
	h.Write(nonce)
	return base64.StdEncoding.EncodeToString(h.Sum(nil)), nil
}

func signData(uri string, data any, ssecurityB64 string) (url.Values, error) {
	b, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	dataStr := string(b)

	// nonce = base64(8 random bytes + (time/60) 4 bytes)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	rand8 := make([]byte, 8)
	_, _ = r.Read(rand8)
	tm := uint32(time.Now().Unix() / 60)
	tb := []byte{byte(tm >> 24), byte(tm >> 16), byte(tm >> 8), byte(tm)}
	nonceRaw := append(rand8, tb...)
	nonceB64 := base64.StdEncoding.EncodeToString(nonceRaw)

	snonceB64, err := signNonce(ssecurityB64, nonceB64)
	if err != nil {
		return nil, err
	}
	snonceRaw, _ := base64.StdEncoding.DecodeString(snonceB64)

	msg := uri + "&" + snonceB64 + "&" + nonceB64 + "&data=" + dataStr
	mac := hmac.New(sha256.New, snonceRaw)
	mac.Write([]byte(msg))
	sign := base64.StdEncoding.EncodeToString(mac.Sum(nil))

	form := url.Values{}
	form.Set("_nonce", nonceB64)
	form.Set("data", dataStr)
	form.Set("signature", sign)
	return form, nil
}
