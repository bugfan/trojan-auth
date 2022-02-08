package apis

import (
	"errors"
	"net/http"

	"github.com/bugfan/de"
	"github.com/bugfan/trojan-auth/env"
	"github.com/bugfan/trojan-auth/srv/services"
	"github.com/bugfan/trojan-auth/utils"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var (
	TokenKey         string = env.GetDefault("AUTH_KEY", "AuthToken")
	AuthHash         string = env.GetDefault("TROJAN_HASH", "AuthHash")
	AuthRemoteIP     string = env.GetDefault("TROJAN_REMOTE_IP", "AuthRemoteIp")
	TrojanApiCryptor        = de.New(env.Get("trojan_api_secret"))
)

func AuthMiddleware(ctx *gin.Context) {
	token := ctx.Request.Header.Get(TokenKey)
	if token == "" {
		ctx.AbortWithError(http.StatusForbidden, errors.New("nil token"))
		return
	}
	_, err := TrojanApiCryptor.DecodeEx([]byte(token))
	if err != nil {
		ctx.AbortWithError(http.StatusForbidden, err)
		ctx.Abort()
		return
	}
	ctx.Next()
}
func VerifyVPNRequest(ctx *gin.Context) {
	hash := ctx.Request.Header.Get(AuthHash)
	remoteIP := ctx.Request.Header.Get(AuthRemoteIP)
	logrus.Info("verify from vpn,ip is %v,hash is:%v\n", remoteIP, hash)
	if _, has := services.GetPassByHash(hash); has {
		ctx.JSON(http.StatusOK, nil)
		return
	}
	// todo: verify if or not vip  by remoteip
	ctx.Writer.WriteHeader(http.StatusForbidden)
	return
}

type packet struct {
	Name string `json:"name"`
	Pass string `json:"pass"`
	Hash string `json:"hash,omitempty"`
}

func GenerateCredential(ctx *gin.Context) {
	p := &packet{}
	err := ctx.Bind(p)
	if err != nil {
		ctx.AbortWithError(http.StatusBadRequest, err)
		return
	}
	if p.Pass == "" || p.Hash == "" || p.Name == "" {
		ctx.AbortWithError(http.StatusUnprocessableEntity, err)
		return
	}
	if utils.SHA224String(p.Pass) != p.Hash {
		logrus.Error("given hash not match pass")
		ctx.AbortWithError(http.StatusUnprocessableEntity, err)
		return
	}
	if err = services.NewCredential(p.Name, p.Pass, p.Hash); err != nil {
		ctx.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	p.Hash = ""
	ctx.JSON(http.StatusCreated, p)
	return
}
