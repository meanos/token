package TokenMaster

import (
	"net/http"
	"time"
)

func Verify(id, ip string) (int, string) {
	code, t := GetToken(id)
	if code != 200 {
		return code, ""
	} else {
		if t.Expires <= time.Now().Unix() || t.Ip != ip {
			go t.Remove()
			return http.StatusForbidden, ""
		} else {
			return http.StatusOK, t.Uid
		}
	}
}

func (t Token) Remove() {
	RemoveToken(t.TokenId)
}
