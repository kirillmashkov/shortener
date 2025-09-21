package net

import (
	"net/http"
	"net/netip"

	"github.com/kirillmashkov/shortener.git/internal/app"
	"go.uber.org/zap"
)

const realIPHeader = "X-Real-IP"

// IsFromTrustSubnet - проверка вхождения ip в подсеть
func IsFromTrustSubnet(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		trustedIPNet, err := netip.ParsePrefix(app.ServerConf.TrustedSubnet)
		if err != nil {
			app.Log.Error("Can't parse trustedIP from config", zap.Error(err))
			http.Error(w, "Something went wrong", http.StatusForbidden)
			return
		}

		ipStr := r.Header.Get(realIPHeader)
		ip, err := netip.ParseAddr(ipStr)
		if err != nil {
			app.Log.Error("Can't parse trustedIP from header", zap.Error(err))
			http.Error(w, "Something went wrong", http.StatusForbidden)
			return
		}

		if !trustedIPNet.Contains(ip) {
			app.Log.Error("subnet in config doesn't contain ip in header", zap.Error(err))
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
