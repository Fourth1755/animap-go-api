package adapters

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"net/url"

	"github.com/Fourth1755/animap-go-api/internal/core/services"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

type HttpAuthHandler struct {
	auth *services.Authenticator
}

func NewHttpAuthHandler(auth *services.Authenticator) *HttpAuthHandler {
	return &HttpAuthHandler{auth: auth}
}

// Handler for our login.
func (h *HttpAuthHandler) Login() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		state, err := generateRandomState()
		if err != nil {
			ctx.String(http.StatusInternalServerError, err.Error())
			return
		}

		// Save the state inside the session.
		session := sessions.Default(ctx)
		session.Set("state", state)
		if err := session.Save(); err != nil {
			ctx.String(http.StatusInternalServerError, err.Error())
			return
		}

		ctx.Redirect(http.StatusTemporaryRedirect, h.auth.AuthCodeURL(state))
	}
}

func generateRandomState() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	state := base64.StdEncoding.EncodeToString(b)

	return state, nil
}

// Handler for our callback.
func (h *HttpAuthHandler) Callback(ctx *gin.Context) {
	session := sessions.Default(ctx)
	if ctx.Query("state") != session.Get("state") {
		ctx.String(http.StatusBadRequest, "Invalid state parameter.")
		return
	}

	// Exchange an authorization code for a token.
	token, err := h.auth.Exchange(ctx.Request.Context(), ctx.Query("code"))
	if err != nil {
		ctx.String(http.StatusUnauthorized, "Failed to exchange an authorization code for a token.")
		return
	}

	idToken, err := h.auth.VerifyIDToken(ctx.Request.Context(), token)
	if err != nil {
		ctx.String(http.StatusInternalServerError, "Failed to verify ID Token.")
		return
	}

	var profile map[string]interface{}
	if err := idToken.Claims(&profile); err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	session.Set("access_token", token.AccessToken)
	session.Set("profile", profile)
	sid := profile["sid"].(string)
	ctx.SetCookie("access_token", token.AccessToken, 3600, "/", "localhost", false, true)
	ctx.SetCookie("user_id", sid, 3600, "/", "localhost", false, true)
	if err := session.Save(); err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	// Redirect to logged in page.
	ctx.Redirect(http.StatusTemporaryRedirect, "http://localhost:3000/")

}

// Handler for our logout.
func (h *HttpAuthHandler) Logout(ctx *gin.Context) {
	logoutUrl, err := url.Parse("https://" + viper.GetString("auth0.domain") + "/v2/logout")
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	scheme := "http"
	if ctx.Request.TLS != nil {
		scheme = "https"
	}

	returnTo, err := url.Parse(scheme + "://" + ctx.Request.Host)
	if err != nil {
		ctx.String(http.StatusInternalServerError, err.Error())
		return
	}

	parameters := url.Values{}
	parameters.Add("returnTo", returnTo.String())
	parameters.Add("client_id", viper.GetString("auth0.clientID"))
	logoutUrl.RawQuery = parameters.Encode()

	//ctx.Redirect(http.StatusTemporaryRedirect, logoutUrl.String())
	ctx.Redirect(http.StatusTemporaryRedirect, "http://localhost:3000/")
}
