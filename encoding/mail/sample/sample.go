package sample

import (
	"encoding/json"
	"io/ioutil"
	"path"
	"strings"
	"fmt"
	"errors"
)

/*
 *
 * Mail Sample representation
 *
 * Headers and Body can be provided together or separate through
 * the following fields (this order is preserved, the latter fields can
 * override the former:
 *
 * 	- EmailFile
 * 	- Headers
 * 	- HeadersFile (if not Headers)
 * 	- Body
 * 	- BodyFile (if not Body)
 * 
 */

type MailSample struct {
	Type    		[]string
	Port			int
	To      		string
	From    		string
	EmailFile    	string
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

func (s *MailSample) buildEmail(file_path string) error {
	var err error
	if len(s.EmailFile) != 0 {
		f, err := ioutil.ReadFile(path.Dir(file_path)+"/"+s.EmailFile)
		if err != nil {
			return err
		}
		data_split := strings.SplitN(string(f), "\n\n", 2)
		if len(data_split) != 2 {
			return errors.New(fmt.Sprintf("Invalid email file: %v", file_path))
		}
		// parse headers
		headers_aux := strings.Split(data_split[0], "\n")
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

		//parse body
		s.Body = data_split[1]
	}
	return err
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
	if err := s.buildEmail(file_path); err != nil {
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
