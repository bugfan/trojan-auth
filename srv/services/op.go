package services

import (
	"github.com/bugfan/trojan-auth/models"
)

func GetPassByHash(hash string) (string, bool) {
	credential := &models.Credential{}
	exist, err := models.GetEngine().Cols("pass").Where("hash = ?", hash).Get(credential)
	if err != nil || !exist {
		return "", false
	}
	return credential.Pass, true
}

func NewCredential(name, pass, hash string) error {
	credential := &models.Credential{
		Name: name,
		Pass: pass,
		Hash: hash,
	}
	user := &models.Credential{Name: name}
	_, err := models.GetEngine().Table(user).Get(user)
	if err != nil {
		return err
	}

	_, err = models.Insert(credential)
	return err
}
