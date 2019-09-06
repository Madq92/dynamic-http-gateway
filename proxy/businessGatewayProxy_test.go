package proxy

import (
	"io/ioutil"
	"net/http"
	"net/textproto"
	"net/url"
	"strings"
	"testing"
)

func Test_checkPairs(t *testing.T) {
	type args struct {
		pairs []string
	}
	tests := []struct {
		name    string
		args    args
		want    int
		wantErr bool
	}{
		{
			name:    "pairs ok",
			args:    args{pairs: []string{"key", "value", "key2", "value2"}},
			want:    4,
			wantErr: false,
		},
		{
			name:    "pairs not ok",
			args:    args{pairs: []string{"key", "value", "key2"}},
			want:    3,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := checkPairs(tt.args.pairs...)
			if (err != nil) != tt.wantErr {
				t.Errorf("checkPairs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("checkPairs() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_headersEqFun(t *testing.T) {
	type args struct {
		param []string
	}
	tests := []struct {
		name    string
		args    args
		req     *http.Request
		want    bool
		wantErr bool
	}{
		{
			name:    "headersEq eq test",
			args:    args{param: []string{"key", "value"}},
			req:     &http.Request{Header: map[string][]string{"Key": {"value"}}},
			want:    true,
			wantErr: false,
		},
		{
			name:    "headersEq not eq test",
			args:    args{param: []string{"key", "value2"}},
			req:     &http.Request{Header: map[string][]string{"Key": {"value"}}},
			want:    false,
			wantErr: false,
		},
		{
			name:    "headersEq not has test",
			args:    args{param: []string{"key", "value2"}},
			req:     &http.Request{Header: map[string][]string{"Key2": {"value2"}}},
			want:    false,
			wantErr: false,
		},
		{
			name:    "headersEq error test",
			args:    args{param: []string{"key"}},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := headersEqFun(tt.args.param...)
			if (err != nil) != tt.wantErr {
				t.Errorf("headersEqFun() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				re := got(tt.req)
				if re != tt.want {
					t.Errorf("headersEqFun() = %v, want %v", re, tt.want)
				}
			}
		})
	}
}

func Test_headersHasFun(t *testing.T) {
	type args struct {
		headers []string
	}
	tests := []struct {
		name    string
		args    args
		req     *http.Request
		want    bool
		wantErr bool
	}{
		{
			name:    "headersHasFun ok test",
			args:    args{headers: []string{"key"}},
			req:     &http.Request{Header: map[string][]string{"Key": {"value"}}},
			want:    true,
			wantErr: false,
		},
		{
			name:    "headersHasFun not ok test",
			args:    args{headers: []string{"key"}},
			req:     &http.Request{Header: map[string][]string{"Key1": {"value"}}},
			want:    false,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := headersHasFun(tt.args.headers...)
			if (err != nil) != tt.wantErr {
				t.Errorf("headersHasFun() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				re := got(tt.req)
				if re != tt.want {
					t.Errorf("headersHasFun() = %v, want %v", re, tt.want)
				}
			}
		})
	}
}

func Test_headersRegexFun(t *testing.T) {
	type args struct {
		param []string
	}
	tests := []struct {
		name    string
		args    args
		req     *http.Request
		want    bool
		wantErr bool
	}{
		{
			name: "headersRegexFun all match test",
			args: args{[]string{"allMatch", ".*"}},
			req: &http.Request{Header: map[string][]string{textproto.CanonicalMIMEHeaderKey(
				"allMatch"): {"value"}}},
			want:    true,
			wantErr: false,
		},
		//{
		//	name: "headersRegexFun not match test",
		//	args: args{[]string{"notMatch", "\b123"}},
		//	req: &http.Request{Header: map[string][]string{textproto.CanonicalMIMEHeaderKey(
		//		"notMatch"): {"123a"}}},
		//	want:    true,
		//	wantErr: false,
		//},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := headersRegexFun(tt.args.param...)
			if (err != nil) != tt.wantErr {
				t.Errorf("headersRegexFun() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				re := got(tt.req)
				if re != tt.want {
					t.Errorf("headersRegexFun() = %v, want %v", re, tt.want)
				}
			}
		})
	}
}

func Test_pathPrefixFun(t *testing.T) {
	type args struct {
		param []string
	}
	tests := []struct {
		name    string
		args    args
		req     *http.Request
		want    bool
		wantErr bool
	}{
		{
			name:    "pathPrefixFun ok test",
			args:    args{[]string{"/prefix"}},
			req:     &http.Request{URL: &url.URL{Path: "/aa"}},
			want:    false,
			wantErr: false,
		},
		{
			name:    "pathPrefixFun not ok test",
			args:    args{[]string{"/prefix"}},
			req:     &http.Request{URL: &url.URL{Path: "/prefixsss"}},
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := pathPrefixFun(tt.args.param...)
			if (err != nil) != tt.wantErr {
				t.Errorf("pathPrefixFun() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				re := got(tt.req)
				if re != tt.want {
					t.Errorf("pathPrefixFun() = %v, want %v", re, tt.want)
				}
			}
		})
	}
}

func Test_pathFun(t *testing.T) {
	type args struct {
		param []string
	}
	tests := []struct {
		name    string
		args    args
		req     *http.Request
		want    bool
		wantErr bool
	}{
		{
			name:    "pathPrefixFun ok test",
			args:    args{[]string{"/path"}},
			req:     &http.Request{URL: &url.URL{Path: "/path"}},
			want:    true,
			wantErr: false,
		},
		{
			name:    "pathPrefixFun not ok test",
			args:    args{[]string{"/path"}},
			req:     &http.Request{URL: &url.URL{Path: "/path/aa"}},
			want:    false,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := pathFun(tt.args.param...)
			if (err != nil) != tt.wantErr {
				t.Errorf("pathFun() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				re := got(tt.req)
				if re != tt.want {
					t.Errorf("pathFun() = %v, want %v", re, tt.want)
				}
			}
		})
	}
}

func Test_methodFun(t *testing.T) {
	type args struct {
		methods []string
	}
	tests := []struct {
		name    string
		args    args
		req     *http.Request
		want    bool
		wantErr bool
	}{
		{
			name:    "methodFun ok test",
			args:    args{[]string{"get"}},
			req:     &http.Request{Method: "GET"},
			want:    true,
			wantErr: false,
		},
		{
			name:    "methodFun not ok test",
			args:    args{[]string{"get"}},
			req:     &http.Request{Method: "POST"},
			want:    false,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := methodFun(tt.args.methods...)
			if (err != nil) != tt.wantErr {
				t.Errorf("methodFun() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				re := got(tt.req)
				if re != tt.want {
					t.Errorf("methodFun() = %v, want %v", re, tt.want)
				}
			}
		})
	}
}

func Test_queryFun(t *testing.T) {
	type args struct {
		param []string
	}
	url, _ := url.Parse("http://www.56qq.com/path?key=value")

	tests := []struct {
		name    string
		args    args
		req     *http.Request
		want    bool
		wantErr bool
	}{
		{
			name:    "queryFun ok test",
			args:    args{[]string{"key", "value"}},
			req:     &http.Request{URL: url},
			want:    true,
			wantErr: false,
		},
		{
			name:    "queryFun not ok test",
			args:    args{[]string{"key", "value2"}},
			req:     &http.Request{URL: url},
			want:    false,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := queryFun(tt.args.param...)
			if (err != nil) != tt.wantErr {
				t.Errorf("queryFun() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				re := got(tt.req)
				if re != tt.want {
					t.Errorf("queryFun() = %v, want %v", re, tt.want)
				}
			}
		})
	}
}

func Test_bodyFun(t *testing.T) {
	type args struct {
		param []string
	}

	//var body io.ReadCloser
	//bf.WriteTo(body)

	tests := []struct {
		name    string
		args    args
		want    bool
		req     *http.Request
		wantErr bool
	}{
		{
			name:    "bodyFun ok test",
			args:    args{[]string{"key", "value"}},
			want:    true,
			wantErr: false,
			req: &http.Request{
				Body:   ioutil.NopCloser(strings.NewReader(`{"key":"value"}`)),
				Header: http.Header{"Content-Type": {`application/json`}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := bodyFun(tt.args.param...)
			if (err != nil) != tt.wantErr {
				t.Errorf("bodyFun() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				re := got(tt.req)
				if re != tt.want {
					t.Errorf("queryFun() = %v, want %v", re, tt.want)
				}
			}
		})
	}
}
