package configs

import (
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	GoogleOauthConfig = &oauth2.Config{
		RedirectURL: "http://127.0.0.1:5500/assets/html/profile.html",
		ClientID: "896520400467-dn4s2ei25ibjqtbutmini99djvchgao0.apps.googleusercontent.com",
		ClientSecret: "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",
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}
)