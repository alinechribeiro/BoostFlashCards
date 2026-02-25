package handlers

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/boostflashcards/backend/internal/auth"
	"github.com/boostflashcards/backend/internal/models"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/oauth2"
	"github.com/gorilla/mux"
)

const bcryptCost = 12

// AuthConfig is passed from main and holds JWT secret, cookie name, server URL, and OAuth configs.
type AuthConfig struct {
	JWTSecret  string
	CookieName  string
	ServerURL   string
	FrontendURL string
	OAuth       *auth.OAuthConfig
}

func (h *Handlers) Signup(w http.ResponseWriter, r *http.Request) {
	if h.Auth == nil {
		writeJSON(w, http.StatusServiceUnavailable, map[string]string{"error": "auth not configured"})
		return
	}
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Name     string `json:"name"`
		Role     string `json:"role"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid body"})
		return
	}
	body.Email = strings.TrimSpace(strings.ToLower(body.Email))
	body.Name = strings.TrimSpace(body.Name)
	if body.Email == "" || len(body.Password) < 8 {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "email required and password at least 8 characters"})
		return
	}
	if body.Role != "student" && body.Role != "tutor" {
		body.Role = "student"
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcryptCost)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to hash password"})
		return
	}
	user, err := h.Repo.CreateUser(r.Context(), body.Email, string(hash), body.Name, body.Role)
	if err != nil {
		if strings.Contains(err.Error(), "Duplicate") {
			writeJSON(w, http.StatusConflict, map[string]string{"error": "email already registered"})
			return
		}
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	h.setSessionAndRespond(w, r, user)
}

func (h *Handlers) Login(w http.ResponseWriter, r *http.Request) {
	if h.Auth == nil {
		writeJSON(w, http.StatusServiceUnavailable, map[string]string{"error": "auth not configured"})
		return
	}
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid body"})
		return
	}
	body.Email = strings.TrimSpace(strings.ToLower(body.Email))
	if body.Email == "" || body.Password == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "email and password required"})
		return
	}
	user, err := h.Repo.GetUserByEmail(r.Context(), body.Email)
	if err != nil || user == nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "invalid email or password"})
		return
	}
	if user.PasswordHash == "" {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "use social login for this account"})
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(body.Password)); err != nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "invalid email or password"})
		return
	}
	h.setSessionAndRespond(w, r, user)
}

func (h *Handlers) Logout(w http.ResponseWriter, r *http.Request) {
	if h.Auth == nil {
		writeJSON(w, http.StatusOK, map[string]string{"ok": "true"})
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     h.Auth.CookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
	writeJSON(w, http.StatusOK, map[string]string{"ok": "true"})
}

func (h *Handlers) Me(w http.ResponseWriter, r *http.Request) {
	if h.Auth == nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}
	cookie, err := r.Cookie(h.Auth.CookieName)
	if err != nil || cookie == nil || cookie.Value == "" {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}
	claims, err := auth.ParseToken(h.Auth.JWTSecret, cookie.Value)
	if err != nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}
	user, err := h.Repo.GetUserByID(r.Context(), claims.UserID)
	if err != nil || user == nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}
	writeJSON(w, http.StatusOK, userResponse(user))
}

func (h *Handlers) OAuthRedirect(w http.ResponseWriter, r *http.Request) {
	if h.Auth == nil || h.Auth.OAuth == nil {
		http.Error(w, "auth not configured", http.StatusServiceUnavailable)
		return
	}
	provider := mux.Vars(r)["provider"]
	oc := h.Auth.OAuth.Config(provider)
	if oc == nil {
		http.Error(w, "unknown provider", http.StatusBadRequest)
		return
	}
	state, err := oauthState(provider)
	if err != nil {
		http.Error(w, "state error", http.StatusInternalServerError)
		return
	}
	url := oc.AuthCodeURL(state, oauth2.AccessTypeOffline)
	http.Redirect(w, r, url, http.StatusFound)
}

func (h *Handlers) OAuthCallback(w http.ResponseWriter, r *http.Request) {
	if h.Auth == nil || h.Auth.OAuth == nil {
		http.Redirect(w, r, h.frontURL()+"/login?error=config", http.StatusFound)
		return
	}
	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")
	if code == "" || state == "" {
		http.Redirect(w, r, h.frontURL()+"/login?error=callback", http.StatusFound)
		return
	}
	provider, err := parseOAuthState(state)
	if err != nil {
		http.Redirect(w, r, h.frontURL()+"/login?error=state", http.StatusFound)
		return
	}
	oc := h.Auth.OAuth.Config(provider)
	if oc == nil {
		http.Redirect(w, r, h.frontURL()+"/login?error=provider", http.StatusFound)
		return
	}
	token, err := oc.Exchange(r.Context(), code)
	if err != nil {
		http.Redirect(w, r, h.frontURL()+"/login?error=exchange", http.StatusFound)
		return
	}
	client := oc.Client(r.Context(), token)
	var profile *auth.OAuthProfile
	switch provider {
	case auth.ProviderGoogle:
		profile, err = auth.FetchGoogleProfile(r.Context(), client, token)
	case auth.ProviderFacebook:
		profile, err = auth.FetchFacebookProfile(r.Context(), client, token)
	case auth.ProviderLinkedIn:
		profile, err = auth.FetchLinkedInProfile(r.Context(), client, token)
	case auth.ProviderInstagram:
		profile, err = auth.FetchInstagramProfile(r.Context(), client, token)
	default:
		http.Redirect(w, r, h.frontURL()+"/login?error=provider", http.StatusFound)
		return
	}
	if err != nil || profile == nil {
		http.Redirect(w, r, h.frontURL()+"/login?error=profile", http.StatusFound)
		return
	}
	user, err := h.Repo.GetUserByOAuth(r.Context(), provider, profile.ProviderID)
	if err != nil {
		http.Redirect(w, r, h.frontURL()+"/login?error=db", http.StatusFound)
		return
	}
	if user != nil {
		if profile.AvatarURL != "" && user.AvatarURL != profile.AvatarURL {
			_ = h.Repo.UpdateUserAvatar(r.Context(), user.ID, profile.AvatarURL)
			user.AvatarURL = profile.AvatarURL
		}
		h.setSessionCookie(w, user)
		http.Redirect(w, r, h.frontURL()+"/", http.StatusFound)
		return
	}
	// New user: redirect to frontend to choose role, with pending token
	pendingToken, err := h.issuePendingOAuthToken(provider, profile)
	if err != nil {
		http.Redirect(w, r, h.frontURL()+"/login?error=token", http.StatusFound)
		return
	}
	http.Redirect(w, r, h.frontURL()+"/signup/complete?pending="+pendingToken, http.StatusFound)
}

func (h *Handlers) CompleteSignup(w http.ResponseWriter, r *http.Request) {
	if h.Auth == nil {
		writeJSON(w, http.StatusServiceUnavailable, map[string]string{"error": "auth not configured"})
		return
	}
	var body struct {
		PendingToken string `json:"pending_token"`
		Role         string `json:"role"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid body"})
		return
	}
	if body.Role != "student" && body.Role != "tutor" {
		body.Role = "student"
	}
	provider, profile, err := h.parsePendingOAuthToken(body.PendingToken)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid or expired token"})
		return
	}
	// Ensure email is unique; if already taken, link OAuth to existing user or return conflict
	user, _ := h.Repo.GetUserByEmail(r.Context(), profile.Email)
	if user != nil {
		// Link this OAuth identity to existing user
		_, err = h.Repo.CreateOAuthIdentity(r.Context(), user.ID, provider, profile.ProviderID)
		if err != nil {
			writeJSON(w, http.StatusConflict, map[string]string{"error": "email already registered"})
			return
		}
		if profile.AvatarURL != "" {
			_ = h.Repo.UpdateUserAvatar(r.Context(), user.ID, profile.AvatarURL)
			user.AvatarURL = profile.AvatarURL
		}
		h.setSessionAndRespond(w, r, user)
		return
	}
	user, err = h.Repo.CreateUser(r.Context(), profile.Email, "", profile.Name, body.Role)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	_, err = h.Repo.CreateOAuthIdentity(r.Context(), user.ID, provider, profile.ProviderID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	if profile.AvatarURL != "" {
		_ = h.Repo.UpdateUserAvatar(r.Context(), user.ID, profile.AvatarURL)
		user.AvatarURL = profile.AvatarURL
	}
	h.setSessionAndRespond(w, r, user)
}

func userResponse(u *models.User) map[string]interface{} {
	return map[string]interface{}{
		"id":         u.ID,
		"email":      u.Email,
		"name":       u.Name,
		"role":       u.Role,
		"avatar_url": u.AvatarURL,
		"created_at": u.CreatedAt,
	}
}

func (h *Handlers) setSessionAndRespond(w http.ResponseWriter, r *http.Request, user *models.User) {
	h.setSessionCookie(w, user)
	writeJSON(w, http.StatusOK, userResponse(user))
}

func (h *Handlers) setSessionCookie(w http.ResponseWriter, user *models.User) {
	tok, _ := auth.NewToken(h.Auth.JWTSecret, user.ID, user.Email, user.Role, 7*24*time.Hour)
	http.SetCookie(w, &http.Cookie{
		Name:     h.Auth.CookieName,
		Value:    tok,
		Path:     "/",
		MaxAge:   7 * 24 * 3600,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   false,
	})
}

func (h *Handlers) frontURL() string {
	if h.Auth != nil && h.Auth.FrontendURL != "" {
		return strings.TrimSuffix(h.Auth.FrontendURL, "/")
	}
	return "http://localhost:3000"
}

func oauthState(provider string) (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	payload := provider + ":" + base64.URLEncoding.EncodeToString(b)
	return base64.URLEncoding.EncodeToString([]byte(payload)), nil
}

func parseOAuthState(state string) (string, error) {
	dec, err := base64.URLEncoding.DecodeString(state)
	if err != nil {
		return "", err
	}
	p := strings.SplitN(string(dec), ":", 2)
	if len(p) < 1 {
		return "", nil
	}
	return p[0], nil
}

// issuePendingOAuthToken returns a short-lived JWT that encodes provider + profile for completion.
func (h *Handlers) issuePendingOAuthToken(provider string, profile *auth.OAuthProfile) (string, error) {
	claims := pendingClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		Purpose:        "oauth_complete",
		Provider:       provider,
		ProviderUserID: profile.ProviderID,
		Email:          profile.Email,
		Name:           profile.Name,
		AvatarURL:      profile.AvatarURL,
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString([]byte(h.Auth.JWTSecret))
}

func (h *Handlers) parsePendingOAuthToken(tokenString string) (provider string, profile *auth.OAuthProfile, err error) {
	t, err := jwt.ParseWithClaims(tokenString, &pendingClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, auth.ErrInvalidToken
		}
		return []byte(h.Auth.JWTSecret), nil
	})
	if err != nil {
		return "", nil, err
	}
	pc, ok := t.Claims.(*pendingClaims)
	if !ok || !t.Valid || pc.Purpose != "oauth_complete" {
		return "", nil, auth.ErrInvalidToken
	}
	return pc.Provider, &auth.OAuthProfile{
		Email:      pc.Email,
		Name:       pc.Name,
		AvatarURL:  pc.AvatarURL,
		ProviderID: pc.ProviderUserID,
	}, nil
}

type pendingClaims struct {
	jwt.RegisteredClaims
	Purpose        string `json:"purpose"`
	Provider       string `json:"provider"`
	ProviderUserID string `json:"provider_user_id"`
	Email          string `json:"email"`
	Name           string `json:"name"`
	AvatarURL      string `json:"avatar_url"`
}
