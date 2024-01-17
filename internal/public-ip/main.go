package public_ip

import (
	"encoding/json"
	"io"
	"net"
	"net/http"
)

type IP struct {
	Query string
}

func GetPublicIpAddr() string {
	req, err := http.Get("http://ip-api.com/json/")
	if err != nil {
		return err.Error()
	}
	defer req.Body.Close()

	body, err := io.ReadAll(req.Body)
	if err != nil {
		return err.Error()
	}
	var ip IP
	err = json.Unmarshal(body, &ip)
	if err != nil {
		return "0.0.0.0"
	}
	if net.ParseIP(ip.Query) == nil {
		return "0.0.0.0"
	}
	return ip.Query
}
