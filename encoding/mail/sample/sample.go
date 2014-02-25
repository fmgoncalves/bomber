package sample

import (
	"encoding/json"
	"io/ioutil"
	"path"
	"strings"
)

type MailSample struct {
	Type    		[]string
	Port			int
	To      		string
	From    		string
	Headers 		[]string
	HeadersFile		string
	Body    		string
	BodyFile    	string
}

func (s MailSample) DefinedHeaders() map[string]struct{} {
	keys_map := make(map[string] struct{})
	for _, val := range s.Headers {
		key := strings.SplitN(val,":",2)[0]
		keys_map[strings.ToLower(key)] = struct{}{};
	}
	return keys_map;
}



func (s *MailSample) buildHeaders(file_path string) error {
	var err error
	if len(s.Headers) == 0 && len(s.HeadersFile) != 0 {
		f, err := ioutil.ReadFile(path.Dir(file_path)+"/"+s.HeadersFile)
		if err != nil {
			return err
		}
		headers_aux := strings.Split(string(f), "\n")
		headers := make([]string, len(headers_aux))
		idx := -1
		for _, val := range headers_aux {
			if ! strings.HasPrefix(val," ") && ! strings.HasPrefix(val,"\t") {
				idx++
			} else {
				headers[idx] += "\r\n"
			}
			headers[idx] += strings.TrimRight(val,"\r")
		}
		headers = headers[0:idx+1]
		s.Headers = headers
	}
	return err
}

func (s *MailSample) buildBody(file_path string) error {
	var err error
	if len(s.Body) == 0 && len(s.BodyFile) != 0 {
		f, err := ioutil.ReadFile(path.Dir(file_path)+"/"+s.BodyFile)
		if err != nil {
			return err
		}
		s.Body = string(f)
	}
	return err
}

func Unmarshal(file_path string) (MailSample, error){
	var s MailSample
	f, err := ioutil.ReadFile(file_path)
	if err != nil {
		return s, err
	}
	err = json.Unmarshal(f, &s)
	if err != nil {
		return s, err
	}
	if err := s.buildHeaders(file_path); err != nil {
		return s, err
	}
	if err := s.buildBody(file_path) ; err != nil {
		return s, err
	}
	return s, nil
}
