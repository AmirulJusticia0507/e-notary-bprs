package sso

import (
	"context"
	"errors"
	"fmt"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
)

// Client adalah klien OIDC generik untuk Keycloak self-host.
// Dua URL dibedakan karena topologi Docker:
//   - issuer: URL publik untuk browser (mis. http://localhost:8180/realms/enotary)
//   - internalURL: URL antar-container untuk backend
//     (mis. http://keycloak:8080/realms/enotary). Bila dikosongkan = issuer.
type Client struct {
	issuer        string
	discoveryBase string
	clientID      string
	clientSecret  string
	redirectURL   string
	provider      *oidc.Provider
	oauthCfg      oauth2.Config
}

type Claims struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Username string `json:"preferred_username"`
	Nonce    string `json:"nonce"`
	Roles    []string `json:"roles"`
	RealmAccess struct {
		Roles []string `json:"roles"`
	} `json:"realm_access"`
}

func NewClient(ctx context.Context, issuer, internalURL, clientID, clientSecret, redirectURL string) (*Client, error) {
	if issuer == "" || clientID == "" {
		return nil, errors.New("keycloak issuer and client id are required")
	}
	discoveryURL := internalURL
	if discoveryURL == "" {
		discoveryURL = issuer
	}
	provider, err := oidc.NewProvider(ctx, discoveryURL)
	if err != nil {
		return nil, fmt.Errorf("keycloak discovery: %w", err)
	}
	// Endpoint token/jwks pakai hasil discovery (bisa internal),
	// endpoint auth SELALU publik agar bisa dibuka browser.
	endpoint := provider.Endpoint()
	endpoint.AuthURL = issuer + "/protocol/openid-connect/auth"
	return &Client{
		issuer:        issuer,
		discoveryBase: discoveryURL,
		clientID:      clientID,
		clientSecret:  clientSecret,
		redirectURL:   redirectURL,
		provider:      provider,
		oauthCfg: oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			RedirectURL:  redirectURL,
			Endpoint:     endpoint,
			Scopes:       []string{oidc.ScopeOpenID, "email", "profile"},
		},
	}, nil
}

// IsConfigured true bila kredensial Keycloak diisi.
func (c *Client) IsConfigured() bool {
	return c != nil && c.issuer != "" && c.clientID != ""
}

// AuthURL membangun URL login Keycloak (dibuka browser).
func (c *Client) AuthURL(state, nonce string) string {
	return c.oauthCfg.AuthCodeURL(state, oidc.Nonce(nonce))
}

// Exchange menukar code dari callback menjadi token OIDC.
func (c *Client) Exchange(ctx context.Context, code string) (*oauth2.Token, error) {
	return c.oauthCfg.Exchange(ctx, code)
}

// Verify memvalidasi ID token: signature via JWKS + client + issuer publik + nonce.
func (c *Client) Verify(ctx context.Context, rawIDToken, expectedNonce string) (*Claims, error) {
	verifier := c.provider.Verifier(&oidc.Config{ClientID: c.clientID, SkipIssuerCheck: true})
	idToken, err := verifier.Verify(ctx, rawIDToken)
	if err != nil {
		return nil, fmt.Errorf("verify id token: %w", err)
	}
	// Issuer Keycloak mengikuti Host yang dipakai saat login browser
	// (publik) vs discovery backend (internal); terima keduanya.
	if idToken.Issuer != c.issuer && idToken.Issuer != c.discoveryBase {
		return nil, fmt.Errorf("unexpected token issuer %q", idToken.Issuer)
	}
	var claims Claims
	if err := idToken.Claims(&claims); err != nil {
		return nil, fmt.Errorf("parse claims: %w", err)
	}
	if expectedNonce != "" && claims.Nonce != expectedNonce {
		return nil, errors.New("nonce mismatch")
	}
	if len(claims.Roles) == 0 {
		claims.Roles = claims.RealmAccess.Roles
	}
	return &claims, nil
}

// MapRole memetakan role Keycloak ke role aplikasi (prioritas tertinggi menang).
func MapRole(roles []string) string {
	for _, want := range []string{"admin", "legal_officer", "notary", "nasabah"} {
		for _, got := range roles {
			if got == want {
				return want
			}
		}
	}
	return "nasabah"
}
