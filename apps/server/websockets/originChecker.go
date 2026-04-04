package websockets

import (
	"net/http"
)

// TODO: add cors implementation
func CheckOriginWithSettings(r *http.Request) bool {
	return true
	// origin := r.Header.Get("Origin")
	//
	// // no origin means same origin, browsers generally don't send origin to the same origin afaik
	// if origin == "" {
	// 	return true
	// }
	//
	// host := r.Host
	// hostWithProtocol := host
	// if !strings.HasPrefix(host, "http://") && !strings.HasPrefix(host, "https://") {
	// 	if r.TLS != nil {
	// 		hostWithProtocol = "https://" + host
	// 	} else {
	// 		hostWithProtocol = "http://" + host
	// 	}
	// }
	//
	// origin = strings.TrimSuffix(origin, "/")
	// hostWithProtocol = strings.TrimSuffix(hostWithProtocol, "/")
	//
	// if origin == hostWithProtocol {
	// 	return true
	// }
	//
	// settings, err := models.GetSystemSettings()
	// if err != nil {
	// 	return false
	// }
	//
	// if settings.AllowedOrigins == "" {
	// 	return false
	// }
	//
	// allowedOrigins := strings.Split(settings.AllowedOrigins, ",")
	// for _, allowed := range allowedOrigins {
	// 	allowed = strings.TrimSpace(allowed)
	// 	allowed = strings.TrimSuffix(allowed, "/")
	// 	if origin == allowed {
	// 		return true
	// 	}
	// }
	//
	// return false
}
