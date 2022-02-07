package apis

import (
	"net/http"

	"github.com/bugfan/de"
	"github.com/bugfan/trojan-auth/env"
	"github.com/bugfan/trojan-auth/srv/services"
	"github.com/bugfan/trojan-auth/utils"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var (
	AuthName         string = env.GetDefault("AUTH_KEY", "TrojanAuth")
	ApiName          string = env.GetDefault("API_KEY", "ApiAuth")
	AuthRemoteIP     string = env.GetDefault("REMOTE_IP", "RemoteIp")
	TK               string = "TrojanHash"
	trojanApiCryptor        = de.New(env.Get("trojan_api_secret"))
)

func AuthMiddleware(ctx *gin.Context) {
	token := ctx.Request.Header.Get(ApiName)
	if token == "" {
		ctx.Writer.WriteHeader(http.StatusForbidden)
		ctx.Abort()
		return
	}
	_, err := trojanApiCryptor.DecodeEx([]byte(token))
	if err != nil {
		ctx.Writer.WriteHeader(http.StatusForbidden)
		ctx.Abort()
		return
	}
	ctx.Next()
}
func VerifyVPNRequest(ctx *gin.Context) {
	authToken := ctx.Request.Header.Get(AuthName)
	if authToken == "" {
		ctx.Writer.WriteHeader(http.StatusForbidden)
		ctx.Abort()
		return
	}
	hash, err := trojanApiCryptor.Decode([]byte(authToken))
	if err != nil {
		ctx.Writer.WriteHeader(http.StatusForbidden)
		ctx.Abort()
		return
	}
	remoteIP := ctx.Request.Header.Get(AuthRemoteIP)
	logrus.Info("verify from vpn,ip is %v,hash is:%v\n", remoteIP, hash)
	if _, has := services.GetPassByHash(string(hash)); has {
		ctx.JSON(http.StatusOK, nil)
		return
	}
	ctx.Writer.WriteHeader(http.StatusForbidden)
	return
}

type packet struct {
	Name string `json:"name"`
	Pass string `json:"pass"`
	Hash string `json:"hash"`
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
	ctx.JSON(http.StatusCreated, nil)
	return
}
