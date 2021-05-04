package providers

import (
	"context"
	"fmt"
	"net/url"

	"github.com/oauth2-proxy/oauth2-proxy/v7/pkg/apis/sessions"
	"github.com/oauth2-proxy/oauth2-proxy/v7/pkg/logger"
	"github.com/oauth2-proxy/oauth2-proxy/v7/pkg/requests"
)

// System76Provider represents a System76 based Identity Provider
type System76Provider struct {
	*ProviderData
}

var _ Provider = (*System76Provider)(nil)

const (
	system76ProviderName  = "System76"
	system76DefaultScope = "profile:read"
)

var (
	// Default Login URL for System76.
	// Pre-parsed URL of https://accounts.system76.com/oauth/authorize.
	system76DefaultLoginURL = &url.URL{
		Scheme: "https",
		Host:   "account.system76.com",
		Path:   "/oauth/authorize",
		RawQuery: "",
	}

	// Default Redeem URL for System76.
	// Pre-parsed URL of https://accounts.system76.com/oauth/token.
	system76DefaultRedeemURL = &url.URL{
		Scheme: "https",
		Host:   "account.system76.com",
		Path:   "/oauth/token",
	}

	// Default Profile URL for System76.
	// Pre-parsed URL of https://accounts.system76.com/api/settings.
	system76DefaultProfileURL = &url.URL{
		Scheme: "https",
		Host:   "account.system76.com",
		Path:   "/api/settings",
	}
)

// NewNextcloudProvider initiates a new NextcloudProvider
func NewSystem76Provider(p *ProviderData) *System76Provider {
	p.setProviderDefaults(providerDefaults{
		name:        system76ProviderName,
		loginURL:    system76DefaultLoginURL,
		redeemURL:   system76DefaultRedeemURL,
		profileURL:  system76DefaultProfileURL,
		validateURL: nil,
		scope:       system76DefaultScope,
	})
	return &System76Provider{
		ProviderData: p,
	}
}

// EnrichSession uses the Recognizer settings endpoint to populate the session's email.
func (p *System76Provider) EnrichSession(ctx context.Context, s *sessions.SessionState) error {
	json, err := requests.New(p.ProfileURL.String()).
		WithContext(ctx).
		SetHeader("Authorization", "Bearer "+s.AccessToken).
		Do().
		UnmarshalJSON()
	if err != nil {
		logger.Errorf("failed making request %v", err)
		return err
	}

	email, err := json.Get("user").Get("email").String()
	if err != nil {
		return fmt.Errorf("unable to extract email from settings endpoint: %v", err)
	}
	s.Email = email

	staff, err := json.Get("user").Get("staff").Bool()
	if err != nil {
		return fmt.Errorf("unable to extract staff bool from settings endpoint: %v", err)
	}
	if staff {
		s.Groups = append(s.Groups, "staff")
	}

	return nil
}
