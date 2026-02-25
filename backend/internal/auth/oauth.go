package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/linkedin"
)

const (
	ProviderGoogle   = "google"
	ProviderFacebook = "facebook"
	ProviderLinkedIn = "linkedin"
	ProviderInstagram = "instagram"
)

type OAuthConfig struct {
	Google   *oauth2.Config
	Facebook *oauth2.Config
	LinkedIn *oauth2.Config
	Instagram *oauth2.Config
}

func NewOAuthConfig(serverURL, googleID, googleSecret, facebookID, facebookSecret, linkedinID, linkedinSecret, instagramID, instagramSecret string) *OAuthConfig {
	redirectURL := strings.TrimSuffix(serverURL, "/") + "/api/auth/callback"
	cfg := &OAuthConfig{}
	if googleID != "" && googleSecret != "" {
		cfg.Google = &oauth2.Config{
			ClientID:     googleID,
			ClientSecret: googleSecret,
			RedirectURL:  redirectURL,
			Scopes:       []string{"openid", "email", "profile"},
			Endpoint:     google.Endpoint,
		}
	}
	if facebookID != "" && facebookSecret != "" {
		cfg.Facebook = &oauth2.Config{
			ClientID:     facebookID,
			ClientSecret: facebookSecret,
			RedirectURL:  redirectURL,
			// Some app types don't allow the "email" scope without extra setup.
			// We request only public_profile and try to read email if it's available.
			Scopes:       []string{"public_profile"},
			Endpoint:     facebook.Endpoint,
		}
	}
	if linkedinID != "" && linkedinSecret != "" {
		cfg.LinkedIn = &oauth2.Config{
			ClientID:     linkedinID,
			ClientSecret: linkedinSecret,
			RedirectURL:  redirectURL,
			Scopes:       []string{"openid", "profile", "email"},
			Endpoint:     linkedin.Endpoint,
		}
	}
	if instagramID != "" && instagramSecret != "" {
		cfg.Instagram = &oauth2.Config{
			ClientID:     instagramID,
			ClientSecret: instagramSecret,
			RedirectURL:  redirectURL,
			Scopes:       []string{"user_profile", "user_media"},
			Endpoint: oauth2.Endpoint{
				AuthURL:  "https://api.instagram.com/oauth/authorize",
				TokenURL: "https://api.instagram.com/oauth/access_token",
			},
		}
	}
	return cfg
}

func (c *OAuthConfig) Config(provider string) *oauth2.Config {
	switch provider {
	case ProviderGoogle:
		return c.Google
	case ProviderFacebook:
		return c.Facebook
	case ProviderLinkedIn:
		return c.LinkedIn
	case ProviderInstagram:
		return c.Instagram
	default:
		return nil
	}
}

// OAuthProfile holds common fields we extract from each provider.
type OAuthProfile struct {
	Email      string
	Name       string
	AvatarURL  string
	ProviderID string
}

func FetchGoogleProfile(ctx context.Context, client *http.Client, token *oauth2.Token) (*OAuthProfile, error) {
	req, _ := http.NewRequestWithContext(ctx, "GET", "https://www.googleapis.com/oauth2/v2/userinfo", nil)
	token.SetAuthHeader(req)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("google profile: %d", resp.StatusCode)
	}
	var v struct {
		ID      string `json:"id"`
		Email   string `json:"email"`
		Name    string `json:"name"`
		Picture string `json:"picture"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
		return nil, err
	}
	return &OAuthProfile{Email: v.Email, Name: v.Name, AvatarURL: v.Picture, ProviderID: v.ID}, nil
}

func FetchFacebookProfile(ctx context.Context, client *http.Client, token *oauth2.Token) (*OAuthProfile, error) {
	u := "https://graph.facebook.com/me?fields=id,name,email,picture.type(large)&access_token=" + url.QueryEscape(token.AccessToken)
	req, _ := http.NewRequestWithContext(ctx, "GET", u, nil)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("facebook profile: %d", resp.StatusCode)
	}
	var v struct {
		ID    string `json:"id"`
		Email string `json:"email"`
		Name  string `json:"name"`
		Picture struct {
			Data struct { URL string `json:"url"` }
		} `json:"picture"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
		return nil, err
	}
	return &OAuthProfile{Email: v.Email, Name: v.Name, AvatarURL: v.Picture.Data.URL, ProviderID: v.ID}, nil
}

func FetchLinkedInProfile(ctx context.Context, client *http.Client, token *oauth2.Token) (*OAuthProfile, error) {
	req, _ := http.NewRequestWithContext(ctx, "GET", "https://api.linkedin.com/v2/me", nil)
	token.SetAuthHeader(req)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("linkedin profile: %d", resp.StatusCode)
	}
	var me struct {
		ID        string `json:"id"`
		FirstName struct { Localized map[string]string `json:"localized"` } `json:"firstName"`
		LastName  struct { Localized map[string]string `json:"localized"` } `json:"lastName"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&me); err != nil {
		return nil, err
	}
	name := ""
	if n, ok := me.FirstName.Localized["en_US"]; ok {
		name = n
	}
	if n, ok := me.LastName.Localized["en_US"]; ok {
		name = strings.TrimSpace(name + " " + n)
	}
	// LinkedIn v2 doesn't return email in /v2/me; need email endpoint with r_emailaddress
	emailReq, _ := http.NewRequestWithContext(ctx, "GET", "https://api.linkedin.com/v2/emailAddress?q=members&projection=(elements*(handle~))", nil)
	token.SetAuthHeader(emailReq)
	emailResp, err := client.Do(emailReq)
	email := ""
	if err == nil && emailResp.StatusCode == http.StatusOK {
		var em struct {
			Elements []struct {
				Handle struct { Email string `json:"emailAddress"` } `json:"handle~"`
			} `json:"elements"`
		}
		_ = json.NewDecoder(emailResp.Body).Decode(&em)
		if len(em.Elements) > 0 {
			email = em.Elements[0].Handle.Email
		}
		emailResp.Body.Close()
	}
	return &OAuthProfile{Email: email, Name: name, ProviderID: me.ID}, nil
}

func FetchInstagramProfile(ctx context.Context, client *http.Client, token *oauth2.Token) (*OAuthProfile, error) {
	// Instagram Graph API (after token exchange) or Basic Display
	u := "https://graph.instagram.com/me?fields=id,username,account_type&access_token=" + url.QueryEscape(token.AccessToken)
	req, _ := http.NewRequestWithContext(ctx, "GET", u, nil)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("instagram profile: %d", resp.StatusCode)
	}
	var v struct {
		ID       string `json:"id"`
		Username string `json:"username"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&v); err != nil {
		return nil, err
	}
	// Instagram often doesn't give email; use username@instagram.placeholder or leave empty and require completion
	email := v.Username + "@instagram.user"
	if v.Username == "" {
		email = "ig-" + v.ID + "@instagram.user"
	}
	return &OAuthProfile{Email: email, Name: v.Username, ProviderID: v.ID}, nil
}
