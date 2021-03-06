package middleware

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/facebook"
	"github.com/markbates/goth/providers/github"
	"github.com/markbates/goth/providers/gplus"
	"net/http"
)

// 0 => not in list, 1 => unique, > 1 => duplicate
func isUnique(list map[string][]string, name string) int {
	count := 0
	for key, _ := range list {
		if key == name {
			count++
		}
	}
	return count
}
func isInAuths(provider string, auths map[string][]string) []string {
	for key, value := range auths {
		if key == provider {
			return value
		}
	}
	return nil
}
func mapProvider(auths map[string][]string) error {
	useProviders := []goth.Provider{}
	for providerName, providerSetting := range auths {

		if isUnique(auths, providerName) != 1 {
			return errors.New("Providers duplicate") //Should unique this
		}
		switch providerName {
		case "facebook":
			useProviders = append(useProviders, facebook.New(providerSetting[0], providerSetting[1], providerSetting[2])) //(providerSetting...)
		case "google":
			useProviders = append(useProviders, gplus.New(providerSetting[0], providerSetting[1], providerSetting[2])) //(providerSetting...)
		case "github":
			useProviders = append(useProviders, github.New(providerSetting[0], providerSetting[1], providerSetting[2])) //(providerSetting...)
		default:
			return errors.New("Provider required")
		}
	}
	goth.UseProviders(useProviders...)
	gothic.GetState = func(req *http.Request) string {
		return req.URL.Query().Get("state")
	}

	return nil
}
func Goth(authConf map[string][]string) gin.HandlerFunc {
	if err := mapProvider(authConf); err != nil {
		panic(err)
	}

	return func(c *gin.Context) {
		c.Next()
	}
}
