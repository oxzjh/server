package webdav

import (
	"fmt"
	"net/http"
	"time"

	"golang.org/x/net/webdav"
)

type webDAV struct {
	auths    map[string]string
	readonly bool
	cert     string
	key      string
	timeout  time.Duration
	prefix   string
	dir      string
	logger   func(*http.Request, error)
	handler  *webdav.Handler
}

func (wd *webDAV) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if len(wd.auths) > 0 {
		username, password, ok := r.BasicAuth()
		if !ok {
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		if p, ok := wd.auths[username]; !ok || p != password {
			http.Error(w, "WebDAV: authorization fail!", http.StatusUnauthorized)
			return
		}
	}
	if wd.readonly && (r.Method == http.MethodPut || r.Method == http.MethodDelete || r.Method == "PROPPATCH" || r.Method == "MKCOL" || r.Method == "COPY" || r.Method == "MOVE") {
		http.Error(w, "WebDAV: readonly!", http.StatusForbidden)
		return
	}
	wd.handler.ServeHTTP(w, r)
}

func Serve(addr string, opts ...Option) error {
	wd := &webDAV{timeout: 5 * time.Second}
	for _, opt := range opts {
		opt(wd)
	}
	wd.handler = &webdav.Handler{
		Prefix:     wd.prefix,
		FileSystem: webdav.Dir(wd.dir),
		LockSystem: webdav.NewMemLS(),
		Logger:     wd.logger,
	}

	s := &http.Server{Addr: addr, Handler: wd, ReadHeaderTimeout: wd.timeout}
	if wd.cert != "" && wd.key != "" {
		fmt.Println("Serve WebDAVS on", addr)
		return s.ListenAndServeTLS(wd.cert, wd.key)
	}
	fmt.Println("Serve WebDAV on", addr)
	return s.ListenAndServe()
}
