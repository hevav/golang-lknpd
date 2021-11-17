package lknpd

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"
)

type NalogClient struct {
	API           string
	INN           string
	Token         string
	TokenExpireIn time.Time
	RefreshToken  string
	DeviceInfo    DeviceInfo
}

func CreateClient(deviceId string) *NalogClient {
	return &NalogClient{
		API:           "https://lknpd.nalog.ru/api/v1/",
		INN:           "",
		Token:         "",
		TokenExpireIn: time.Now(),
		RefreshToken:  "",
		DeviceInfo: DeviceInfo{
			AppVersion:     "1.0.0",
			SourceDeviceId: deviceId,
			SourceType:     "WEB",
			MetaDetails: MetaDetails{
				UserAgent: "Mozilla/5.0 (Macintosh; Intel Mac OS X 11_2_2) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/88.0.4324.192 Safari/537.36",
			},
		},
	}
}

func (cl *NalogClient) Auth(login, password string) error {
	resp, err := cl.CallAPI("auth/token", AuthPayload{
		Username:   login,
		Password:   password,
		DeviceInfo: cl.DeviceInfo,
	}, false)

	if err != nil {
		return err
	}

	return cl.DecodeAuth(resp)
}

func (cl *NalogClient) AddIncome(income Income) (*Receipt, error) {
	resp, err := cl.CallAPI("income", income, true)
	if err != nil {
		return nil, err
	}

	var decoded IncomeResponse
	err = json.NewDecoder(resp).Decode(&decoded)

	if err != nil {
		return nil, err
	}

	return &Receipt{
		UUID:     decoded.ApprovedReceiptUuid,
		InfoURL:  cl.API + "receipt/" + cl.INN + "/" + decoded.ApprovedReceiptUuid + "/json",
		PrintURL: cl.API + "receipt/" + cl.INN + "/" + decoded.ApprovedReceiptUuid + "/print",
	}, nil
}

func (cl *NalogClient) GetToken() (string, error) {
	if cl.Token == "" {
		return "", errors.New("unauthorized")
	}

	if cl.TokenExpireIn.After(time.Now()) {
		return cl.Token, nil
	}

	resp, err := cl.CallAPI("/auth/token", RefreshTokenPayload{
		DeviceInfo:   cl.DeviceInfo,
		RefreshToken: cl.RefreshToken,
	}, false)

	if err != nil {
		return "", err
	}

	err = cl.DecodeAuth(resp)
	return cl.Token, err
}

func (cl *NalogClient) DecodeAuth(resp io.ReadCloser) error {
	var decoded AuthResponse
	err := json.NewDecoder(resp).Decode(&decoded)
	if err != nil {
		return err
	}

	if decoded.RefreshToken != "" {
		cl.RefreshToken = decoded.RefreshToken
	}

	if decoded.Token == "" || decoded.TokenExpireIn == "" {
		return errors.New("server responded with no token")
	}

	cl.Token = decoded.Token
	cl.TokenExpireIn, err = time.Parse(TimeFormat, decoded.TokenExpireIn)

	if err != nil {
		return err
	}

	if decoded.Profile.INN != "" {
		cl.INN = decoded.Profile.INN
	}

	return nil
}

func (cl *NalogClient) CallAPI(method string, payload interface{}, sendToken bool) (io.ReadCloser, error) {
	url := cl.API + method
	var httpMethod string
	var body io.Reader

	if payload == nil {
		httpMethod = "GET"
		body = nil
	} else {
		obj, err := json.Marshal(payload)
		if err != nil {
			return nil, err
		}
		body = bytes.NewBuffer(obj)
	}

	request, err := http.NewRequest(httpMethod, url, body)
	if err != nil {
		return nil, err
	}

	if sendToken {
		request.Header.Set("Authorization", "Bearer "+cl.Token)
		request.Header.Set("Referer", "https://lknpd.nalog.ru/sales")
	} else {
		request.Header.Set("Referer", "https://lknpd.nalog.ru/")
	}

	request.Header.Set("Accept", "application/json, text/plain, */*")
	request.Header.Set("Accept-Language", "ru-RU,ru;q=0.9,en-US;q=0.8,en;q=0.7")
	request.Header.Set("Cache-Control", "no-cache")
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Pragma", "no-cache")
	request.Header.Set("Referer-Policy", "strict-origin-when-cross-origin")

	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}

	return resp.Body, nil
}
