package global

import (
	"mxshop-api/user-web/config"
	"mxshop-api/user-web/proto"
	"strings"

	ut "github.com/go-playground/universal-translator"
)

var (
	ServerConfig  *config.ServerConfig
	Trans         ut.Translator
	UserSrvClient proto.UserClient
)

func RemoveTopStruct(fields map[string]string) map[string]string {
	res := map[string]string{}
	for field, err := range fields {
		res[field[strings.Index(field, ".")+1:]] = err
	}
	return res
}
