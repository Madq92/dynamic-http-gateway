package proxy

import (
	"dynamic-http-gateway/common"
	"dynamic-http-gateway/log"
	"fmt"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"io"
	"math/rand"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"time"
)

type (
	Service struct {
		Name       string       `json:"name"`
		Targets    []string     `json:"targets"`
		MatchRules []MatchRules `json:"matchRules"`
	}
	MatchRules struct {
		Name  string   `json:"name"`
		Param []string `json:"param"`
	}

	multiTargetProxyBalancer struct {
		matchTargets []*matchTarget
	}
	matchTarget struct {
		targets  []*middleware.ProxyTarget
		matchFns []matchFunc
		random   *rand.Rand
	}
	matchFunc func(req *http.Request) bool
)

var (
	matchFunMap = map[string]func(matchRuleBody ...string) (func(req *http.Request) bool, error){
		"headerseq":    headersEqFun,
		"headershas":   headersHasFun,
		"headersregex": headersRegexFun,
		"path":         pathFun,
		"pathprefix":   pathPrefixFun,
		"method":       methodFun,
		"query":        queryFun,
	}
)

func NewMiddlewareFuncWithConfig(service []Service) (echo.MiddlewareFunc, error) {
	if len(service) == 0 {
		return nil, fmt.Errorf("service not be null")
	}
	balancer, err := newMultiTargetProxyBalancer(service)
	if err != nil {
		return nil, err
	}
	c := middleware.DefaultProxyConfig
	c.Balancer = balancer
	return proxyWithConfig(c), nil
}

func proxyWithConfig(config middleware.ProxyConfig) echo.MiddlewareFunc {
	// Defaults
	if config.Skipper == nil {
		config.Skipper = middleware.DefaultLoggerConfig.Skipper
	}
	if config.Balancer == nil {
		panic("echo: proxy middleware requires balancer")
	}
	//config.rewriteRegex = map[*regexp.Regexp]string{}
	//
	//// Initialize
	//for k, v := range config.Rewrite {
	//	k = strings.Replace(k, "*", "(\\S*)", -1)
	//	config.rewriteRegex[regexp.MustCompile(k)] = v
	//}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			if config.Skipper(c) {
				return next(c)
			}
			tgt := config.Balancer.Next(c)
			if nil == tgt {
				return next(c)
			}

			req := c.Request()
			res := c.Response()

			c.Set(config.ContextKey, tgt)

			// Fix host
			req.Host = tgt.URL.Host

			//// Rewrite
			//for k, v := range config.rewriteRegex {
			//	replacer := captureTokens(k, req.URL.Path)
			//	if replacer != nil {
			//		req.URL.Path = replacer.Replace(v)
			//	}
			//}

			// Fix header
			if req.Header.Get(echo.HeaderXRealIP) == "" {
				req.Header.Set(echo.HeaderXRealIP, c.RealIP())
			}
			if req.Header.Get(echo.HeaderXForwardedProto) == "" {
				req.Header.Set(echo.HeaderXForwardedProto, c.Scheme())
			}
			if c.IsWebSocket() && req.Header.Get(echo.HeaderXForwardedFor) == "" { // For HTTP, it is automatically set by Go HTTP reverse proxy.
				req.Header.Set(echo.HeaderXForwardedFor, c.RealIP())
			}

			// Proxy
			switch {
			case c.IsWebSocket():
				proxyRaw(tgt, c).ServeHTTP(res, req)
			case req.Header.Get(echo.HeaderAccept) == "text/event-stream":
			default:
				proxyHTTP(tgt, c, config).ServeHTTP(res, req)
			}

			return
		}
	}
}

func proxyHTTP(tgt *middleware.ProxyTarget, c echo.Context, config middleware.ProxyConfig) http.Handler {
	proxy := httputil.NewSingleHostReverseProxy(tgt.URL)
	proxy.ErrorHandler = func(resp http.ResponseWriter, req *http.Request, err error) {
		desc := tgt.URL.String()
		if tgt.Name != "" {
			desc = fmt.Sprintf("%s(%s)", tgt.Name, tgt.URL.String())
		}
		c.Logger().Errorf("remote %s unreachable, could not forward: %v", desc, err)
		c.Error(echo.NewHTTPError(http.StatusServiceUnavailable))
	}
	proxy.Transport = config.Transport
	return proxy
}

func proxyRaw(t *middleware.ProxyTarget, c echo.Context) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		in, _, err := c.Response().Hijack()
		if err != nil {
			c.Error(fmt.Errorf("proxy raw, hijack error=%v, url=%s", t.URL, err))
			return
		}
		defer in.Close()

		out, err := net.Dial("tcp", t.URL.Host)
		if err != nil {
			he := echo.NewHTTPError(http.StatusBadGateway, fmt.Sprintf("proxy raw, dial error=%v, url=%s", t.URL, err))
			c.Error(he)
			return
		}
		defer out.Close()

		// Write header
		err = r.Write(out)
		if err != nil {
			he := echo.NewHTTPError(http.StatusBadGateway, fmt.Sprintf("proxy raw, request header copy error=%v, url=%s", t.URL, err))
			c.Error(he)
			return
		}

		errCh := make(chan error, 2)
		cp := func(dst io.Writer, src io.Reader) {
			_, err = io.Copy(dst, src)
			errCh <- err
		}

		go cp(out, in)
		go cp(in, out)
		err = <-errCh
		if err != nil && err != io.EOF {
			c.Logger().Errorf("proxy raw, copy body error=%v, url=%s", t.URL, err)
		}
	})
}

func newMultiTargetProxyBalancer(service []Service) (*multiTargetProxyBalancer, error) {
	matchTargets := []*matchTarget{}
	for _, ser := range service {
		var tgts []*middleware.ProxyTarget
		for _, tgt := range ser.Targets {
			u, err := url.Parse(tgt)
			if err != nil {
				return nil, err
			}
			tgts = append(tgts, &middleware.ProxyTarget{
				URL: u,
			})
		}

		var matchFuncs []matchFunc
		for _, rule := range ser.MatchRules {
			matchFunName := strings.ToLower(rule.Name)

			matchFun, ok := matchFunMap[matchFunName]
			if !ok {
				return nil, &common.BaseError{Code: 400, Message: fmt.Sprintf("rule name not support : %s",
					matchFunName)}
			}
			fn, err := matchFun(rule.Param...)
			if err != nil {
				return nil, &common.BaseError{Code: 400, Message: fmt.Sprintf("规则分析失败： %s , 失败原因: %s", rule,
					err.Error())}
			}
			matchFuncs = append(matchFuncs, fn)
		}
		mt := &matchTarget{
			targets:  tgts,
			matchFns: matchFuncs,
		}

		matchTargets = append(matchTargets, mt)
	}
	return &multiTargetProxyBalancer{
		matchTargets: matchTargets,
	}, nil
}

func (mt *matchTarget) getTgt() *middleware.ProxyTarget {
	if mt.random == nil {
		mt.random = rand.New(rand.NewSource(int64(time.Now().Nanosecond())))
	}
	return mt.targets[mt.random.Intn(len(mt.targets))]
}

func (m *matchTarget) match(req *http.Request) *middleware.ProxyTarget {
	if len(m.matchFns) == 0 {
		return m.getTgt()
	}

	for _, fn := range m.matchFns {
		if !fn(req) {
			return nil
		}
	}
	return m.getTgt()
}

// AddTarget adds an upstream target to the list.
func (b *multiTargetProxyBalancer) AddTarget(target *middleware.ProxyTarget) bool {
	return false
}

// RemoveTarget removes an upstream target from the list.
func (b *multiTargetProxyBalancer) RemoveTarget(name string) bool {
	return false
}

// Next randomly returns an upstream target.
func (b *multiTargetProxyBalancer) Next(c echo.Context) *middleware.ProxyTarget {
	if len(b.matchTargets) == 0 {
		noTgt := common.Payload{
			Code:    503,
			Message: "no target",
		}
		log.LOGGER.Error(noTgt)
		//TODO no target
	}
	request := c.Request()
	for _, matchTarget := range b.matchTargets {
		match := matchTarget.match(request)
		if match != nil {
			//request.Host = match.URL.Host
			return match
		}
	}
	//TODO don't match any target
	return nil
}
