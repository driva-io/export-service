package gateways

import (
	"encoding/json"
	"errors"
	"export-service/internal/server"
	"log"
	"net/url"
	"os"
	"time"
)

type HTTPAuthService struct {
	HttpClient      server.HttpClient
	adminToken      string
	tokenExpiration time.Time
}

func (s *HTTPAuthService) Login(email string, password string, company string) (response map[string]any, err error) {
	authUrl := os.Getenv("AUTH_API_URL") + "/authenticate"
	u, err := url.Parse(authUrl)
	if err != nil {
		log.Printf("Failed to parse URL: %v\n", err)
		return nil, err
	}

	query := u.Query()
	query.Set("company", company)
	u.RawQuery = query.Encode()

	body := map[string]any{"email": email, "password": password}
	loginResponse, err := s.HttpClient.Post(u.String(), body, nil)
	if err != nil {
		log.Printf("Error logging in: %v\n", err)
		return nil, err
	}

	if loginResponse.StatusCode != 201 {
		log.Println("Error logging in")
		return nil, errors.New("error logging in")
	}

	response, ok := loginResponse.Body.(map[string]any)
	if !ok {
		log.Println("Error logging in")
		return nil, errors.New("error deserializing login response")
	}
	return response, nil
}

func (s *HTTPAuthService) GetUserByToken(headers map[string]any) (AuthUser, error) {
	authUrl := os.Getenv("AUTH_API_URL") + "/current-user"
	u, err := url.Parse(authUrl)
	if err != nil {
		log.Printf("Failed to parse URL: %v\n", err)
		return AuthUser{}, err
	}

	response, err := s.HttpClient.Get(u.String(), headers)
	if err != nil {
		log.Printf("Error validating token: %v\n", err)
		return AuthUser{}, err
	}

	if response.StatusCode != 200 {
		log.Printf("Error validating token\n")
		return AuthUser{}, errors.New("error validating token")
	}

	var user AuthUser
	jsonData, err := json.Marshal(response.Body)
	if err != nil {
		return AuthUser{}, err
	}

	err = json.Unmarshal(jsonData, &user)
	if err != nil {
		return AuthUser{}, err
	}

	if user.ID == "" {
		return AuthUser{}, errors.New("invalid user")
	}
	return user, err
}

func (s *HTTPAuthService) GetAdminToken() string {
	if s.adminToken != "" && s.tokenExpiration.After(time.Now()) {
		return s.adminToken
	}

	adminEmail := os.Getenv("ADMIN_EMAIL")
	adminPassword := os.Getenv("ADMIN_PASSWORD")

	response, err := s.Login(adminEmail, adminPassword, "Driva")
	if err != nil {
		log.Printf("Error logging as admin: %v\n", err)
		return ""
	}

	s.adminToken = response["token"].(string)
	s.tokenExpiration = time.Now().Add(1 * time.Hour)

	return s.adminToken
}

func (s *HTTPAuthService) HasCredits(companyId string, amount int) (hasCredits bool, err error) {
	hasCredits = false
	adminToken := s.GetAdminToken()
	authUrl := os.Getenv("AUTH_API_URL") + "/credits/" + companyId

	var headers = map[string]any{
		"Authorization": "Bearer " + adminToken,
	}

	response, err := s.HttpClient.Get(authUrl, headers)
	if err != nil {
		log.Printf("Error obtaining credits: %v\n", err)
		return false, err
	}

	hasCredits = int(response.Body.(map[string]any)["data"].(map[string]any)["solicitation"].(map[string]any)["total"].(float64)) >= amount

	return
}

func (s *HTTPAuthService) TakeCredits() {
	panic("TakeCredits not implemented")
}

func (s *HTTPAuthService) RefundCredits() {
	panic("RefundCredits not implemented")
}

func (s *HTTPAuthService) GetBlacklist() {
	panic("GetBlacklist not implemented")
}
