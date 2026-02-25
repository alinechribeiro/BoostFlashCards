package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/boostflashcards/backend/internal/auth"
	"github.com/boostflashcards/backend/internal/models"
	"github.com/boostflashcards/backend/internal/repository"
	"github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/account"
	"github.com/stripe/stripe-go/v82/accountlink"
	"github.com/stripe/stripe-go/v82/paymentintent"
	"github.com/stripe/stripe-go/v82/webhook"
)

// BillingConnect creates or reuses a Stripe Connect Express account for the current tutor
// and returns an onboarding link URL.
func (h *Handlers) BillingConnect(w http.ResponseWriter, r *http.Request) {
	if h.Auth == nil {
		writeJSON(w, http.StatusServiceUnavailable, map[string]string{"error": "billing not configured"})
		return
	}

	user := auth.UserFromContext(r.Context())
	if user == nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}
	if user.Role != "tutor" {
		writeJSON(w, http.StatusForbidden, map[string]string{"error": "only tutors can connect billing"})
		return
	}

	secret := os.Getenv("STRIPE_SECRET_KEY")
	if secret == "" {
		writeJSON(w, http.StatusServiceUnavailable, map[string]string{"error": "stripe not configured"})
		return
	}
	stripe.Key = secret

	ctx := r.Context()
	tutor, err := h.Repo.GetTutorWithProfile(ctx, user.ID)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "tutor profile not found"})
		return
	}

	acctID := tutor.Profile.StripeConnectAccountID
	if acctID == "" {
		params := &stripe.AccountParams{
			Type:         stripe.String(string(stripe.AccountTypeExpress)),
			Email:        stripe.String(user.Email),
			BusinessType: stripe.String("individual"),
		}
		acct, err := account.New(params)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to create connect account"})
			return
		}
		acctID = acct.ID
		if err := h.Repo.UpdateTutorStripeAccountID(ctx, user.ID, acctID); err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to persist connect account"})
			return
		}
	}

	baseFront := strings.TrimSuffix(h.Auth.FrontendURL, "/")
	refreshURL := baseFront + "/billing/connect/refresh"
	returnURL := baseFront + "/billing/connect/return"

	link, err := accountlink.New(&stripe.AccountLinkParams{
		Account:    stripe.String(acctID),
		RefreshURL: stripe.String(refreshURL),
		ReturnURL:  stripe.String(returnURL),
		Type:       stripe.String("account_onboarding"),
	})
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to create connect link"})
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"url": link.URL})
}

type createBundleBody struct {
	TutorID     int64 `json:"tutor_id"`
	TotalEssays int   `json:"total_essays"`
}

// CreateBundle creates an essay bundle and its Stripe PaymentIntent.
// It enforces the minimum £20 rule and uses destination charges so the
// tutor's connected account receives their share while the platform keeps a fee.
func (h *Handlers) CreateBundle(w http.ResponseWriter, r *http.Request) {
	if h.Auth == nil {
		writeJSON(w, http.StatusServiceUnavailable, map[string]string{"error": "billing not configured"})
		return
	}
	user := auth.UserFromContext(r.Context())
	if user == nil {
		writeJSON(w, http.StatusUnauthorized, map[string]string{"error": "unauthorized"})
		return
	}
	if user.Role != "student" {
		writeJSON(w, http.StatusForbidden, map[string]string{"error": "only students can buy bundles"})
		return
	}

	var body createBundleBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid body"})
		return
	}
	if body.TutorID == 0 || body.TotalEssays <= 0 {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "tutor_id and total_essays are required"})
		return
	}
	if body.TotalEssays != 3 && body.TotalEssays != 5 {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "only 3 or 5‑essay bundles are supported"})
		return
	}

	secret := os.Getenv("STRIPE_SECRET_KEY")
	if secret == "" {
		writeJSON(w, http.StatusServiceUnavailable, map[string]string{"error": "stripe not configured"})
		return
	}
	stripe.Key = secret

	// Ensure tutor has a connected Stripe account
	tutor, err := h.Repo.GetTutorWithProfile(r.Context(), body.TutorID)
	if err != nil || tutor.Profile.StripeConnectAccountID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": "tutor is not ready for billing"})
		return
	}

	// Pricing: base it on premium tier, with bundle discounts from tutor share.
	premium := repository.TierPricing[models.EssayTierPremium]
	originalTotal := body.TotalEssays * premium.StudentPriceCents
	platformFee := body.TotalEssays * premium.PlatformCutCents

	discountFactor := 1.0
	if body.TotalEssays == 3 {
		discountFactor = 0.9 // 10% discount
	} else if body.TotalEssays == 5 {
		discountFactor = 0.85 // 15% discount
	}
	studentPays := int(float64(originalTotal) * discountFactor)
	if studentPays < repository.MinimumBundleCents {
		studentPays = repository.MinimumBundleCents
	}

	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(int64(studentPays)),
		Currency: stripe.String(string(stripe.CurrencyGBP)),
		PaymentMethodTypes: stripe.StringSlice([]string{
			"card",
		}),
		ApplicationFeeAmount: stripe.Int64(int64(platformFee)),
		TransferData: &stripe.PaymentIntentTransferDataParams{
			Destination: stripe.String(tutor.Profile.StripeConnectAccountID),
		},
	}
	pi, err := paymentintent.New(params)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to create payment intent"})
		return
	}

	bundle := &models.EssayBundle{
		StudentID:             user.ID,
		TutorID:               &tutor.User.ID,
		TotalEssays:           body.TotalEssays,
		UsedEssays:            0,
		PriceCents:            studentPays,
		Status:                models.EssayBundleStatusPendingPayment,
		StripePaymentIntentID: pi.ID,
	}
	created, err := h.Repo.CreateEssayBundle(r.Context(), bundle)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, http.StatusCreated, map[string]interface{}{
		"bundle":        created,
		"client_secret": pi.ClientSecret,
	})
}

// StripeWebhook listens for payment_intent.succeeded and activates bundles,
// then moves any pending essay requests for that bundle to awaiting_tutor.
func (h *Handlers) StripeWebhook(w http.ResponseWriter, r *http.Request) {
	payload, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "failed to read body", http.StatusBadRequest)
		return
	}
	sig := r.Header.Get("Stripe-Signature")
	secret := os.Getenv("STRIPE_WEBHOOK_SECRET")
	if secret == "" {
		http.Error(w, "webhook not configured", http.StatusServiceUnavailable)
		return
	}

	event, err := webhook.ConstructEvent(payload, sig, secret)
	if err != nil {
		http.Error(w, "invalid signature", http.StatusBadRequest)
		return
	}

	if event.Type == "payment_intent.succeeded" {
		var pi stripe.PaymentIntent
		if err := json.Unmarshal(event.Data.Raw, &pi); err != nil {
			http.Error(w, "bad payload", http.StatusBadRequest)
			return
		}
		ctx := r.Context()
		bundle, err := h.Repo.GetEssayBundleByPaymentIntent(ctx, pi.ID)
		if err == nil && bundle != nil {
			_ = h.Repo.UpdateEssayBundleStatus(ctx, bundle.ID, models.EssayBundleStatusActive)
			_ = h.Repo.SetRequestsAwaitingTutorForBundle(ctx, bundle.ID)
		}
	}

	w.WriteHeader(http.StatusOK)
}


