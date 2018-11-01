package restful

import (
	"fmt"
	"util/http"
	"util/safemap"
	"util/sign"
)

type Method int

type FuncAsynCall func([]byte, error)

type MsgProcess interface {
	UnMarshal([]byte, interface{}) error
	Marshal(interface{}) ([]byte, error)
}

const (
	_ Method = iota
	Method_Get
	Method_Post
)

type mapValue struct {
	Name    string
	Url     string
	TMethod Method
}

type RestServer struct {
	safeMap    *safemap.Map
	process    MsgProcess
	baseUrl    string
	apiKey     string
	secrKey    string
	passPhrase string
}

func NewRestServer(baseUrl, apiKey, secrKey, passPhrase string, process MsgProcess) *RestServer {
	return &RestServer{baseUrl: baseUrl, apiKey: apiKey, secrKey: secrKey, passPhrase: passPhrase, safeMap: new(safemap.Map), process: process}
}

func (s *RestServer) RegistInterface(name, url string, method Method) (err error) {

	if nil != s.safeMap.Get(name) {
		return fmt.Errorf("duplicate name:%v", name)
	}
	s.safeMap.Set(name, &mapValue{Name: name, Url: url, TMethod: method})
	return nil
}

func (s *RestServer) CancelInterface(name string) {
	s.safeMap.Del(name)
}

func (s *RestServer) genHeards(url, method string) map[string]string {
	heards := make(map[string]string, 0)
	ts := s.timeStamp()
	heards["OK-ACCESS-KEY"] = s.apiKey
	heards["OK-ACCESS-SIGN"] = sign.HMacSha256(ts+method+url, []byte(s.secrKey))
	heards["OK-ACCESS-TIMESTAMP"] = ts
	heards["OK-ACCESS-PASSPHRASE"] = s.passPhrase
	heards["contentType"] = "application/json"
	return heards
}

func (s *RestServer) SynCall(name string, params map[string]interface{}, paths ...interface{}) (body []byte, err error) {

	v := s.safeMap.Get(name)
	if v == nil {
		return nil, fmt.Errorf("%v not exist,regist first", name)
	}
	switch v.(*mapValue).TMethod {
	case Method_Get:
		extUrl := s.PathsProcess(v.(*mapValue).Url, paths...) + s.getParamProcess(params)
		hearders := s.genHeards(extUrl, "GET")
		fmt.Println(hearders)
		return http.Get(s.baseUrl+extUrl, hearders)
	case Method_Post:
		req, err := s.process.Marshal(params)
		if err != nil {
			return nil, err
		}
		extUrl := s.PathsProcess(v.(*mapValue).Url, paths...)
		hearders := s.genHeards(extUrl, "POST")
		return http.Post(s.baseUrl+extUrl, hearders, req)
	default:
		return nil, fmt.Errorf("method not right")
	}
}

func (s *RestServer) AsynCall(name string, hearders map[string]string, params map[string]interface{}, funcAsynCall FuncAsynCall, paths ...interface{}) {

	v := s.safeMap.Get(name)
	if v == nil {
		funcAsynCall(nil, fmt.Errorf("%v not exist,regist first", name))
	}
	switch v.(*mapValue).TMethod {
	case Method_Get:
		rsp, err := http.Get(s.PathsProcess(v.(*mapValue).Url, paths...)+s.getParamProcess(params), hearders)
		funcAsynCall(rsp, err)
		break
	case Method_Post:
		req, err := s.process.Marshal(params)
		if err != nil {
			funcAsynCall(nil, err)
			break
		}
		rsp, err := http.Post(s.PathsProcess(v.(*mapValue).Url, paths...), hearders, req)
		funcAsynCall(rsp, err)
		break
	default:
		funcAsynCall(nil, fmt.Errorf("method not right"))
	}
}

func (s *RestServer) PathsProcess(url string, params ...interface{}) string {
	return fmt.Sprintf(url, params...)
}

func (s *RestServer) getParamProcess(params map[string]interface{}) string {
	if len(params) == 0 {
		return ""
	}
	var result string = "?"
	for name, param := range params {
		s := name + "=" + fmt.Sprintf("%v", param) + "&"
		result += s
	}
	result = result[:len(result)-1]
	return result
}

type TimeStamp struct {
	Iso   string `json:"iso"`
	Epoch string `json:"epoch"`
}

func (s *RestServer) timeStamp() string {
	var ts TimeStamp
	m := make(map[string]string, 0)
	if body, err := http.Get(s.baseUrl+"/api/general/v3/time", m); err != nil {
		if err := s.process.UnMarshal(body, ts); err != nil {
			return ts.Epoch
		}
	}
	return ""
}
