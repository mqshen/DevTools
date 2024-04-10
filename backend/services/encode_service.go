package services

import (
	"crypto/rsa"
	"devtools/backend/types"
	"encoding/base64"
	"encoding/json"
	"github.com/golang-jwt/jwt/v5"
	"html"
	"net/url"
	"sync"
)

var (
	jwtTestDefaultKey *rsa.PublicKey
	defaultKeyFunc    jwt.Keyfunc = func(t *jwt.Token) (interface{}, error) { return jwtTestDefaultKey, nil }
)

type encodeService struct {
}

var encoder *encodeService
var onceEncoder sync.Once

func EncodeService() *encodeService {
	if encoder == nil {
		onceEncoder.Do(func() {
			encoder = &encodeService{}
		})
	}
	return encoder
}

func (e *encodeService) EncodeFromHTML(content string) (resp types.JSResp) {
	encoded := html.EscapeString(content)
	resp.Data = encoded
	resp.Success = true
	return
}

func (e *encodeService) DecodeFromHTML(content string) (resp types.JSResp) {
	encoded := html.UnescapeString(content)
	resp.Data = encoded
	resp.Success = true
	return
}

func (e *encodeService) EncodeFromURL(content string) (resp types.JSResp) {
	encoded := url.QueryEscape(content)
	resp.Data = encoded
	resp.Success = true
	return
}

func (e *encodeService) DecodeFromURL(content string) (resp types.JSResp) {
	encoded, err := url.QueryUnescape(content)
	if err != nil {
		resp.Msg = err.Error()
		return
	}
	resp.Data = encoded
	resp.Success = true
	return
}

func (e *encodeService) EncodeFromBase64(content string) (resp types.JSResp) {
	encoded := base64.StdEncoding.EncodeToString([]byte(content))
	resp.Data = encoded
	resp.Success = true
	return
}

func (e *encodeService) DecodeFromBase64(content string) (resp types.JSResp) {
	encoded, err := base64.StdEncoding.DecodeString(content)
	if err != nil {
		resp.Msg = err.Error()
		return
	}
	resp.Data = string(encoded)
	resp.Success = true
	return
}

func (e *encodeService) EncodeFromJWT(token string) (resp types.JSResp) {
	parser := new(jwt.Parser)
	unverified, _, err := parser.ParseUnverified(token, jwt.MapClaims{})
	if err != nil {
		unverified, _, err = parser.ParseUnverified(token, &jwt.RegisteredClaims{})
	}
	if err != nil {
		resp.Msg = err.Error()
		return
	}
	header, err := convertToJson(unverified.Header)
	if err != nil {
		resp.Msg = err.Error()
		return
	}
	payload, err := convertToJson(unverified.Claims)
	if err != nil {
		resp.Msg = err.Error()
		return
	}
	resp.Data = map[string]any{
		"header":  header,
		"payload": payload,
	}
	resp.Success = true
	return
}

func (e *encodeService) DecodeFromJWT(content string) (resp types.JSResp) {
	encoded, err := base64.StdEncoding.DecodeString(content)
	if err != nil {
		resp.Msg = err.Error()
		return
	}
	resp.Data = string(encoded)
	resp.Success = true
	return
}

func convertToJson(data interface{}) (string, error) {
	b, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return "", err
	}
	return string(b), nil
}
