package auth

import (
	"net/http"
	"encoding/json"
	"context"

	"remood/pkg/configs"
	"io/ioutil"

)

type GoogleUser struct {
	ID 				string 	`json:"id"`
	Email 			string 	`json:"email"`
	VerifiedEmail 	bool 	`json:"verified_email"`
	Name 			string 	`json:"name"`
	GivenName 		string 	`json:"given_name"`
	FamilyName 		string 	`json:"family_name"`
	Picture 		string 	`json:"picture"`
	Locale 			string 	`json:"locale"`
}

func GetGoogleUserInfo(code string)(GoogleUser, error) {
	var googleUser GoogleUser

	token, err := configs.GoogleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		return googleUser, err
	}

	res, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		return googleUser, err
	}
	defer res.Body.Close()


	content, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return googleUser, err
	}
	
	_ = json.Unmarshal(content, &googleUser)
	return googleUser, nil
}