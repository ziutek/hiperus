package hiperus

import (
	"github.com/ziutek/soap"
	"net/http"
)

const realm = "PLATFORM_MNGM"

//var HiperusURL = "http://10.128.102.5:8080/hiperusapi.php"

var httpClient = &http.Client{Transport: new(http.Transport)}

type Session struct {
	url string
	id  string
}

func (s *Session) do(action string, arg interface{}) (*soap.Element, error) {
	sr, err := newSoapRequest("PLATFORM_MNGM", action, arg, s.id)
	if err != nil {
		return nil, err
	}

	hr, err := http.NewRequest("POST", s.url, sr)
	if err != nil {
		return nil, err
	}
	hr.Header.Set("Content-Type", "text/xml; charset=utf-8")
	hr.Header.Set("SOAPAction", s.url+"#request")

	resp, err := http.DefaultClient.Do(hr)
	if err != nil {
		return nil, err
	}
	/*if action == "GetBilling" {
		io.Copy(os.Stdout, resp.Body)
		return nil, nil
	}*/
	return parseResponse(resp.Body)
}

func (s *Session) cmd(action string, arg interface{}) (*soap.Element, error) {
	ret, err := s.do(action, arg)
	if err != nil {
		return nil, err
	}
	return ret.Get("result_set")
}

type login struct {
	Username string `soap:"username"`
	Password string `soap:"password"`
	Domain   string `soap:"domain"`
}

func NewSession(url, username, password, domain string) (*Session, error) {
	s := &Session{url: url}
	ret, err := s.do("Login", login{username, password, domain})
	if err != nil {
		return nil, err
	}
	sessid, err := ret.Get("sessid")
	if err != nil {
		return nil, err
	}
	s.id, err = sessid.Str()
	if err != nil {
		return nil, err
	}
	return s, nil
}
