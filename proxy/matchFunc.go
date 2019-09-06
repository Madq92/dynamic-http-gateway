package proxy

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/labstack/echo"
	"net/http"
	"regexp"
	"strings"
)

func checkPairs(pairs ...string) (int, error) {
	length := len(pairs)
	if length%2 != 0 {
		return length, fmt.Errorf("number of parameters must be multiple of 2, got %v", pairs)
	}
	return length, nil
}

func headersEqFun(param ...string) (func(req *http.Request) bool, error) {
	length, err := checkPairs(param...)
	if err != nil {
		return nil, err
	}

	m := make(map[string]string, length/2)
	for i := 0; i < length; i += 2 {
		m[param[i]] = param[i+1]
	}

	return func(req *http.Request) bool {
		reqHeader := req.Header
		for k, v := range m {
			if v != reqHeader.Get(k) {
				return false
			}
		}
		return true
	}, nil
}

func headersHasFun(headers ...string) (func(req *http.Request) bool, error) {
	return func(req *http.Request) bool {
		reqHeader := req.Header
		for _, header := range headers {
			val := reqHeader.Get(header)
			if val == "" {
				return false
			}
		}
		return true
	}, nil
}

func headersRegexFun(param ...string) (func(req *http.Request) bool, error) {
	length, err := checkPairs(param...)
	if err != nil {
		return nil, err
	}

	m := make(map[string]*regexp.Regexp, length/2)
	for i := 0; i < length; i += 2 {
		reg, err := regexp.Compile(param[i+1])
		if err != nil {
			return nil, errors.New("regexp error : " + param[i+1])
		}
		m[param[i]] = reg
	}

	return func(req *http.Request) bool {
		reqHeader := req.Header
		for header, re := range m {
			he := reqHeader.Get(header)
			if he == "" {
				return false
			}
			result := re.MatchString(he)
			if !result {
				return false
			}
		}
		return true
	}, nil
}

func pathPrefixFun(param ...string) (func(req *http.Request) bool, error) {
	var paths []string
	for _, path := range param {
		if !strings.HasPrefix(path, "/") {
			path = "/" + path
		}
		paths = append(paths, path)
	}

	return func(req *http.Request) bool {
		reqPath := req.URL.Path
		for _, path := range paths {
			if strings.HasPrefix(reqPath, path) {
				return true
			}
		}

		return false
	}, nil
}

func pathFun(param ...string) (func(req *http.Request) bool, error) {
	pathMap := make(map[string]bool)
	for _, path := range param {
		if !strings.HasPrefix(path, "/") {
			path = "/" + path
		}
		pathMap[path] = true
	}

	return func(req *http.Request) bool {
		if pathMap[req.URL.Path] {
			return true
		}

		return false
	}, nil
}

func methodFun(methods ...string) (func(req *http.Request) bool, error) {
	methodsMap := make(map[string]bool)
	for _, method := range methods {
		methodsMap[strings.ToTitle(method)] = true
	}
	return func(req *http.Request) bool {
		if methodsMap[req.Method] {
			return true
		}
		return false
	}, nil
}

func queryFun(param ...string) (func(req *http.Request) bool, error) {
	length, err := checkPairs(param...)
	if err != nil {
		return nil, err
	}

	m := make(map[string]string, length/2)
	for i := 0; i < length; i += 2 {
		m[param[i]] = param[i+1]
	}

	return func(req *http.Request) bool {
		query := req.URL.Query()
		for k, v := range m {
			if query.Get(k) != v {
				return false
			}
		}
		return true
	}, nil
}

func bodyFun(param ...string) (func(req *http.Request) bool, error) {
	length, err := checkPairs(param...)
	if err != nil {
		return nil, err
	}

	bodyParam := make(map[string]string, length/2)
	for i := 0; i < length; i += 2 {
		bodyParam[param[i]] = param[i+1]
	}
	return func(req *http.Request) bool {
		ctype := req.Header.Get(echo.HeaderContentType)
		if strings.HasPrefix(ctype, echo.MIMEApplicationJSON) {
			bodyMap := make(map[string]string)
			if err := json.NewDecoder(req.Body).Decode(&bodyMap); err != nil {
				return false
			}
			for k := range bodyParam {
				if bodyMap[k] != bodyParam[k] {
					return false
				}
			}
			return true

		}
		return false
	}, nil
}
