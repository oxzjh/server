package acme

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"

	"github.com/oxzjh/server/http/proxy"
	"golang.org/x/crypto/acme/autocert"
)

func Serve(targets map[string]string, opts ...Option) error {
	hosts := make([]string, 0, len(targets))
	for host := range targets {
		hosts = append(hosts, host)
	}
	ops := &options{}
	for _, opt := range opts {
		opt(ops)
	}
	m := &autocert.Manager{
		Prompt:      autocert.AcceptTOS,
		HostPolicy:  autocert.HostWhitelist(hosts...),
		RenewBefore: ops.renewBefore,
	}
	if ops.cacheDir == "" {
		ops.cacheDir = cacheDir()
	}
	if err := os.MkdirAll(ops.cacheDir, 0700); err != nil {
		log.Printf("warning: acme not using a cache: %v", err)
	} else {
		m.Cache = autocert.DirCache(ops.cacheDir)
	}
	return http.Serve(m.Listener(), proxy.NewReverse(targets, ops.errorHandler))
}

func homeDir() string {
	if runtime.GOOS == "windows" {
		return os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
	}
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return "/"
}

func cacheDir() string {
	const base = "acme"
	switch runtime.GOOS {
	case "darwin":
		return filepath.Join(homeDir(), "Library", "Caches", base)
	case "windows":
		for _, ev := range []string{"APPDATA", "CSIDL_APPDATA", "TEMP", "TMP"} {
			if v := os.Getenv(ev); v != "" {
				return filepath.Join(v, base)
			}
		}
		// Worst case:
		return filepath.Join(homeDir(), base)
	}
	if xdg := os.Getenv("XDG_CACHE_HOME"); xdg != "" {
		return filepath.Join(xdg, base)
	}
	return filepath.Join(homeDir(), ".cache", base)
}
