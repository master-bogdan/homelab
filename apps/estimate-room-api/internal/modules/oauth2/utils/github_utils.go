package oauth2utils

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	oauth2dto "github.com/master-bogdan/estimate-room-api/internal/modules/oauth2/dto"
	"github.com/master-bogdan/estimate-room-api/internal/pkg/utils"
)

const (
	defaultGithubAuthURL  = "https://github.com/login/oauth/authorize"
	defaultGithubTokenURL = "https://github.com/login/oauth/access_token"
	defaultGithubAPIBase  = "https://api.github.com"
)

type GithubConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
	StateSecret  string
	Scopes       []string
}

func (c GithubConfig) IsConfigured() bool {
	return c.ClientID != "" && c.ClientSecret != "" && c.RedirectURL != "" && c.StateSecret != ""
}

type GithubState struct {
	ExpiresAt           int64  `json:"exp"`
	ClientID            string `json:"client_id"`
	RedirectURI         string `json:"redirect_uri"`
	ResponseType        string `json:"response_type"`
	Scopes              string `json:"scopes"`
	State               string `json:"state"`
	CodeChallenge       string `json:"code_challenge"`
	CodeChallengeMethod string `json:"code_challenge_method"`
	Nonce               string `json:"nonce"`
}

type GithubProfile struct {
	ID          string
	Email       *string
	DisplayName string
	AvatarURL   *string
}

type githubTokenResponse struct {
	AccessToken      string `json:"access_token"`
	TokenType        string `json:"token_type"`
	Scope            string `json:"scope"`
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

type githubUserResponse struct {
	ID        int64   `json:"id"`
	Login     string  `json:"login"`
	Name      string  `json:"name"`
	AvatarURL string  `json:"avatar_url"`
	Email     *string `json:"email"`
}

type githubEmailResponse struct {
	Email    string `json:"email"`
	Primary  bool   `json:"primary"`
	Verified bool   `json:"verified"`
}

func BuildGithubAuthorizeURL(cfg GithubConfig, state string) string {
	scopes := cfg.Scopes
	if len(scopes) == 0 {
		scopes = []string{"read:user", "user:email"}
	}

	q := url.Values{}
	q.Set("client_id", cfg.ClientID)
	q.Set("redirect_uri", cfg.RedirectURL)
	q.Set("scope", strings.Join(scopes, " "))
	q.Set("state", state)
	q.Set("allow_signup", "true")

	return defaultGithubAuthURL + "?" + q.Encode()
}

func CreateGithubState(key []byte, query *AuthorizeQuery) (string, error) {
	if len(key) == 0 {
		return "", errors.New("github state secret is required")
	}

	state := GithubState{
		ExpiresAt:           time.Now().Add(10 * time.Minute).Unix(),
		ClientID:            query.ClientID,
		RedirectURI:         query.RedirectURI,
		ResponseType:        query.ResponseType,
		Scopes:              query.Scopes,
		State:               query.State,
		CodeChallenge:       query.CodeChallenge,
		CodeChallengeMethod: query.CodeChallengeMethod,
		Nonce:               query.Nonce,
	}

	return utils.GenerateToken(key, state)
}

func ParseGithubState(key []byte, token string) (*GithubState, error) {
	if token == "" {
		return nil, errors.New("state is required")
	}
	state, err := utils.ParseToken[GithubState](key, token)
	if err != nil {
		return nil, errors.New("invalid state")
	}
	if state.ExpiresAt > 0 && time.Now().Unix() > state.ExpiresAt {
		return nil, errors.New("state expired")
	}
	return state, nil
}

func ExchangeGithubCode(ctx context.Context, client *http.Client, cfg GithubConfig, code string) (string, error) {
	if code == "" {
		return "", errors.New("code is required")
	}

	body := url.Values{}
	body.Set("client_id", cfg.ClientID)
	body.Set("client_secret", cfg.ClientSecret)
	body.Set("code", code)
	body.Set("redirect_uri", cfg.RedirectURL)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, defaultGithubTokenURL, strings.NewReader(body.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var tokenResp githubTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return "", err
	}

	if tokenResp.Error != "" {
		return "", errors.New(tokenResp.ErrorDescription)
	}
	if tokenResp.AccessToken == "" {
		return "", errors.New("missing access token")
	}

	return tokenResp.AccessToken, nil
}

func FetchGithubUser(ctx context.Context, client *http.Client, token string) (*githubUserResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, defaultGithubAPIBase+"/user", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("User-Agent", "estimate-room-api")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, errors.New("failed to fetch github user")
	}

	var user githubUserResponse
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, err
	}

	return &user, nil
}

func FetchGithubEmails(ctx context.Context, client *http.Client, token string) ([]githubEmailResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, defaultGithubAPIBase+"/user/emails", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("User-Agent", "estimate-room-api")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, errors.New("failed to fetch github emails")
	}

	var emails []githubEmailResponse
	if err := json.NewDecoder(resp.Body).Decode(&emails); err != nil {
		return nil, err
	}

	return emails, nil
}

func selectGithubEmail(emails []githubEmailResponse) *string {
	for _, email := range emails {
		if email.Primary && email.Verified {
			return &email.Email
		}
	}
	for _, email := range emails {
		if email.Verified {
			return &email.Email
		}
	}
	return nil
}

func BuildGithubProfile(user *githubUserResponse, emails []githubEmailResponse) GithubProfile {
	displayName := user.Name
	if displayName == "" {
		displayName = user.Login
	}

	email := selectGithubEmail(emails)

	var avatarURL *string
	if user.AvatarURL != "" {
		avatarURL = &user.AvatarURL
	}

	return GithubProfile{
		ID:          strconv.FormatInt(user.ID, 10),
		Email:       email,
		DisplayName: displayName,
		AvatarURL:   avatarURL,
	}
}

type AuthorizeQuery struct {
	ClientID            string
	RedirectURI         string
	ResponseType        string
	Scopes              string
	State               string
	CodeChallenge       string
	CodeChallengeMethod string
	Nonce               string
}

func AuthorizeQueryFromState(state *GithubState) *AuthorizeQuery {
	return &AuthorizeQuery{
		ClientID:            state.ClientID,
		RedirectURI:         state.RedirectURI,
		ResponseType:        state.ResponseType,
		Scopes:              state.Scopes,
		State:               state.State,
		CodeChallenge:       state.CodeChallenge,
		CodeChallengeMethod: state.CodeChallengeMethod,
		Nonce:               state.Nonce,
	}
}

func AuthorizeQueryFromDTO(dto *oauth2dto.AuthorizeQueryDTO) *AuthorizeQuery {
	return &AuthorizeQuery{
		ClientID:            dto.ClientID,
		RedirectURI:         dto.RedirectURI,
		ResponseType:        dto.ResponseType,
		Scopes:              dto.Scopes,
		State:               dto.State,
		CodeChallenge:       dto.CodeChallenge,
		CodeChallengeMethod: dto.CodeChallengeMethod,
		Nonce:               dto.Nonce,
	}
}

func (q *AuthorizeQuery) ToDTO() *oauth2dto.AuthorizeQueryDTO {
	return &oauth2dto.AuthorizeQueryDTO{
		ClientID:            q.ClientID,
		RedirectURI:         q.RedirectURI,
		ResponseType:        q.ResponseType,
		Scopes:              q.Scopes,
		State:               q.State,
		CodeChallenge:       q.CodeChallenge,
		CodeChallengeMethod: q.CodeChallengeMethod,
		Nonce:               q.Nonce,
	}
}
