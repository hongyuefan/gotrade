package restful

import (
	"fmt"

	"util/http"
	"util/safemap"
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
	Name        string
	Url         string
	TMethod     Method
	ContentType string
}

type RestServer struct {
	safeMap *safemap.Map
	process MsgProcess
	baseUrl string
}

func NewRestServer(baseUrl string, process MsgProcess) *RestServer {
	return &RestServer{baseUrl: baseUrl, safeMap: new(safemap.Map), process: process}
}

func (s *RestServer) RegistInterface(name, url, contentType string, method Method) (err error) {

	if nil != s.safeMap.Get(name) {
		return fmt.Errorf("duplicate name:%v", name)
	}
	s.safeMap.Set(name, &mapValue{Name: name, Url: s.baseUrl + url, TMethod: method, ContentType: contentType})

	return nil
}

func (s *RestServer) CancelInterface(name string) {
	s.safeMap.Del(name)
}

func (s *RestServer) SynCall(name string, params map[string]interface{}, paths ...interface{}) (body []byte, err error) {

	v := s.safeMap.Get(name)
	if v == nil {
		return nil, fmt.Errorf("%v not exist,regist first", name)
	}
	switch v.(*mapValue).TMethod {
	case Method_Get:
		return http.Get(s.PathsProcess(v.(*mapValue).Url, paths...)+s.getParamProcess(params), v.(*mapValue).ContentType)
	case Method_Post:
		req, err := s.process.Marshal(params)
		if err != nil {
			return nil, err
		}
		return http.Post(s.PathsProcess(v.(*mapValue).Url, paths...), v.(*mapValue).ContentType, req)
	default:
		return nil, fmt.Errorf("method not right")
	}
}

func (s *RestServer) AsynCall(name string, params map[string]interface{}, funcAsynCall FuncAsynCall, paths ...interface{}) {

	v := s.safeMap.Get(name)
	if v == nil {
		funcAsynCall(nil, fmt.Errorf("%v not exist,regist first", name))
	}
	switch v.(*mapValue).TMethod {
	case Method_Get:
		rsp, err := http.Get(s.PathsProcess(v.(*mapValue).Url, paths...)+s.getParamProcess(params), v.(*mapValue).ContentType)
		funcAsynCall(rsp, err)
		break
	case Method_Post:
		req, err := s.process.Marshal(params)
		if err != nil {
			funcAsynCall(nil, err)
			break
		}
		rsp, err := http.Post(s.PathsProcess(v.(*mapValue).Url, paths...), v.(*mapValue).ContentType, req)
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
