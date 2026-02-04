package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/ellebam/hireme/api/internal/config"
	"github.com/ellebam/hireme/api/pkg/httputil"
)

// Test constants
const (
	testSecret     = "test-jwt-secret-key-for-testing-purposes"
	testUserID     = "user-123"
	bypassUserID   = "bypass-user-001"
)

// createTestConfig creates a config for testing
func createTestConfig(bypassEnabled bool, jwtSecret string) *config.Config {
	return &config.Config{
		Auth: config.AuthConfig{
			BypassEnabled: bypassEnabled,
			BypassUserID:  bypassUserID,
			JWTSecret:     jwtSecret,
		},
	}
}

// createValidJWT creates a valid JWT signed with HS256
func createValidJWT(userID string, secret string, expiresAt time.Time) string {
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": expiresAt.Unix(),
		"iat": time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte(secret))
	return tokenString
}

// createExpiredJWT creates an expired JWT
func createExpiredJWT(userID string, secret string) string {
	return createValidJWT(userID, secret, time.Now().Add(-time.Hour))
}

// createJWTWithWrongSigningMethod creates a JWT signed with RS256 (wrong method)
func createJWTWithWrongSigningMethod(userID string) string {
	// Create a JWT with HS256 header but we'll manually construct one that claims RS256
	// Since we don't have RSA keys, we'll create a token that claims to be RS256
	// but actually we'll use "none" algorithm to simulate wrong signing method detection
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(time.Hour).Unix(),
		"iat": time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodNone, claims)
	tokenString, _ := token.SignedString(jwt.UnsafeAllowNoneSignatureType)
	return tokenString
}

// parseErrorResponse parses an error response from the response body
func parseErrorResponse(rr *httptest.ResponseRecorder) (*httputil.Response, error) {
	var resp httputil.Response
	if err := json.NewDecoder(rr.Body).Decode(&resp); err != nil {
		return nil, err
	}
	return &resp, nil
}

// nextHandler is a simple handler that records whether it was called
func nextHandler(called *bool, capturedUserID *string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		*called = true
		if capturedUserID != nil {
			*capturedUserID = GetUserID(r.Context())
		}
		w.WriteHeader(http.StatusOK)
	})
}

func TestAuth_BypassEnabled(t *testing.T) {
	cfg := createTestConfig(true, "")
	middleware := Auth(cfg)

	var called bool
	var capturedUserID string
	handler := middleware(nextHandler(&called, &capturedUserID))

	req := httptest.NewRequest("GET", "/protected", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if !called {
		t.Error("expected next handler to be called")
	}

	if capturedUserID != bypassUserID {
		t.Errorf("expected userID '%s', got '%s'", bypassUserID, capturedUserID)
	}

	if rr.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rr.Code)
	}
}

func TestAuth_ValidJWT(t *testing.T) {
	cfg := createTestConfig(false, testSecret)
	middleware := Auth(cfg)

	var called bool
	var capturedUserID string
	handler := middleware(nextHandler(&called, &capturedUserID))

	token := createValidJWT(testUserID, testSecret, time.Now().Add(time.Hour))

	req := httptest.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if !called {
		t.Error("expected next handler to be called")
	}

	if capturedUserID != testUserID {
		t.Errorf("expected userID '%s', got '%s'", testUserID, capturedUserID)
	}

	if rr.Code != http.StatusOK {
		t.Errorf("expected status %d, got %d", http.StatusOK, rr.Code)
	}
}

func TestAuth_ValidJWT_CaseInsensitiveBearer(t *testing.T) {
	cfg := createTestConfig(false, testSecret)
	middleware := Auth(cfg)

	var called bool
	var capturedUserID string
	handler := middleware(nextHandler(&called, &capturedUserID))

	token := createValidJWT(testUserID, testSecret, time.Now().Add(time.Hour))

	// Test with lowercase "bearer"
	req := httptest.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "bearer "+token)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if !called {
		t.Error("expected next handler to be called with lowercase 'bearer'")
	}

	if capturedUserID != testUserID {
		t.Errorf("expected userID '%s', got '%s'", testUserID, capturedUserID)
	}
}

func TestAuth_MissingHeader(t *testing.T) {
	cfg := createTestConfig(false, testSecret)
	middleware := Auth(cfg)

	var called bool
	handler := middleware(nextHandler(&called, nil))

	req := httptest.NewRequest("GET", "/protected", nil)
	// No Authorization header
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if called {
		t.Error("expected next handler NOT to be called")
	}

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("expected status %d, got %d", http.StatusUnauthorized, rr.Code)
	}

	resp, err := parseErrorResponse(rr)
	if err != nil {
		t.Fatalf("failed to parse error response: %v", err)
	}

	if resp.Error == nil {
		t.Fatal("expected error in response")
	}

	if resp.Error.Message != "missing authorization header" {
		t.Errorf("expected message 'missing authorization header', got '%s'", resp.Error.Message)
	}
}

func TestAuth_InvalidFormat_Basic(t *testing.T) {
	cfg := createTestConfig(false, testSecret)
	middleware := Auth(cfg)

	var called bool
	handler := middleware(nextHandler(&called, nil))

	req := httptest.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Basic dXNlcjpwYXNz") // Base64 encoded "user:pass"
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if called {
		t.Error("expected next handler NOT to be called")
	}

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("expected status %d, got %d", http.StatusUnauthorized, rr.Code)
	}

	resp, err := parseErrorResponse(rr)
	if err != nil {
		t.Fatalf("failed to parse error response: %v", err)
	}

	if resp.Error == nil {
		t.Fatal("expected error in response")
	}

	if resp.Error.Message != "invalid authorization header format" {
		t.Errorf("expected message 'invalid authorization header format', got '%s'", resp.Error.Message)
	}
}

func TestAuth_InvalidFormat_NoToken(t *testing.T) {
	cfg := createTestConfig(false, testSecret)
	middleware := Auth(cfg)

	var called bool
	handler := middleware(nextHandler(&called, nil))

	req := httptest.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer") // Missing token
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if called {
		t.Error("expected next handler NOT to be called")
	}

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("expected status %d, got %d", http.StatusUnauthorized, rr.Code)
	}
}

func TestAuth_InvalidToken_Malformed(t *testing.T) {
	cfg := createTestConfig(false, testSecret)
	middleware := Auth(cfg)

	var called bool
	handler := middleware(nextHandler(&called, nil))

	req := httptest.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer not-a-valid-jwt")
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if called {
		t.Error("expected next handler NOT to be called")
	}

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("expected status %d, got %d", http.StatusUnauthorized, rr.Code)
	}

	resp, err := parseErrorResponse(rr)
	if err != nil {
		t.Fatalf("failed to parse error response: %v", err)
	}

	if resp.Error == nil {
		t.Fatal("expected error in response")
	}

	if resp.Error.Message != "invalid token" {
		t.Errorf("expected message 'invalid token', got '%s'", resp.Error.Message)
	}
}

func TestAuth_ExpiredToken(t *testing.T) {
	cfg := createTestConfig(false, testSecret)
	middleware := Auth(cfg)

	var called bool
	handler := middleware(nextHandler(&called, nil))

	token := createExpiredJWT(testUserID, testSecret)

	req := httptest.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if called {
		t.Error("expected next handler NOT to be called")
	}

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("expected status %d, got %d", http.StatusUnauthorized, rr.Code)
	}

	resp, err := parseErrorResponse(rr)
	if err != nil {
		t.Fatalf("failed to parse error response: %v", err)
	}

	if resp.Error == nil {
		t.Fatal("expected error in response")
	}

	if resp.Error.Message != "invalid token" {
		t.Errorf("expected message 'invalid token', got '%s'", resp.Error.Message)
	}
}

func TestAuth_WrongSigningMethod(t *testing.T) {
	cfg := createTestConfig(false, testSecret)
	middleware := Auth(cfg)

	var called bool
	handler := middleware(nextHandler(&called, nil))

	// Create a token with "none" signing method (wrong method)
	token := createJWTWithWrongSigningMethod(testUserID)

	req := httptest.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if called {
		t.Error("expected next handler NOT to be called")
	}

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("expected status %d, got %d", http.StatusUnauthorized, rr.Code)
	}
}

func TestAuth_WrongSecret(t *testing.T) {
	cfg := createTestConfig(false, testSecret)
	middleware := Auth(cfg)

	var called bool
	handler := middleware(nextHandler(&called, nil))

	// Create token with different secret
	token := createValidJWT(testUserID, "wrong-secret", time.Now().Add(time.Hour))

	req := httptest.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if called {
		t.Error("expected next handler NOT to be called")
	}

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("expected status %d, got %d", http.StatusUnauthorized, rr.Code)
	}
}

func TestAuth_MissingSubClaim(t *testing.T) {
	cfg := createTestConfig(false, testSecret)
	middleware := Auth(cfg)

	var called bool
	handler := middleware(nextHandler(&called, nil))

	// Create token without sub claim
	claims := jwt.MapClaims{
		"exp": time.Now().Add(time.Hour).Unix(),
		"iat": time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte(testSecret))

	req := httptest.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+tokenString)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if called {
		t.Error("expected next handler NOT to be called")
	}

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("expected status %d, got %d", http.StatusUnauthorized, rr.Code)
	}

	resp, err := parseErrorResponse(rr)
	if err != nil {
		t.Fatalf("failed to parse error response: %v", err)
	}

	if resp.Error == nil {
		t.Fatal("expected error in response")
	}

	if resp.Error.Message != "missing user ID in token" {
		t.Errorf("expected message 'missing user ID in token', got '%s'", resp.Error.Message)
	}
}

func TestAuth_EmptySubClaim(t *testing.T) {
	cfg := createTestConfig(false, testSecret)
	middleware := Auth(cfg)

	var called bool
	handler := middleware(nextHandler(&called, nil))

	// Create token with empty sub claim
	claims := jwt.MapClaims{
		"sub": "",
		"exp": time.Now().Add(time.Hour).Unix(),
		"iat": time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte(testSecret))

	req := httptest.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+tokenString)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if called {
		t.Error("expected next handler NOT to be called")
	}

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("expected status %d, got %d", http.StatusUnauthorized, rr.Code)
	}

	resp, err := parseErrorResponse(rr)
	if err != nil {
		t.Fatalf("failed to parse error response: %v", err)
	}

	if resp.Error.Message != "missing user ID in token" {
		t.Errorf("expected message 'missing user ID in token', got '%s'", resp.Error.Message)
	}
}

func TestAuth_SubClaimWrongType(t *testing.T) {
	cfg := createTestConfig(false, testSecret)
	middleware := Auth(cfg)

	var called bool
	handler := middleware(nextHandler(&called, nil))

	// Create token with sub claim as number instead of string
	claims := jwt.MapClaims{
		"sub": 12345,
		"exp": time.Now().Add(time.Hour).Unix(),
		"iat": time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte(testSecret))

	req := httptest.NewRequest("GET", "/protected", nil)
	req.Header.Set("Authorization", "Bearer "+tokenString)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if called {
		t.Error("expected next handler NOT to be called")
	}

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("expected status %d, got %d", http.StatusUnauthorized, rr.Code)
	}
}

// Tests for GetUserID helper function

func TestGetUserID_Present(t *testing.T) {
	ctx := context.WithValue(context.Background(), UserIDKey, testUserID)

	userID := GetUserID(ctx)

	if userID != testUserID {
		t.Errorf("expected userID '%s', got '%s'", testUserID, userID)
	}
}

func TestGetUserID_Missing(t *testing.T) {
	ctx := context.Background()

	userID := GetUserID(ctx)

	if userID != "" {
		t.Errorf("expected empty userID, got '%s'", userID)
	}
}

func TestGetUserID_WrongType(t *testing.T) {
	// Set userID with wrong type (int instead of string)
	ctx := context.WithValue(context.Background(), UserIDKey, 12345)

	userID := GetUserID(ctx)

	if userID != "" {
		t.Errorf("expected empty userID when type is wrong, got '%s'", userID)
	}
}

// Tests for MustGetUserID helper function

func TestMustGetUserID_Present(t *testing.T) {
	ctx := context.WithValue(context.Background(), UserIDKey, testUserID)

	userID := MustGetUserID(ctx)

	if userID != testUserID {
		t.Errorf("expected userID '%s', got '%s'", testUserID, userID)
	}
}

func TestMustGetUserID_Missing(t *testing.T) {
	ctx := context.Background()

	defer func() {
		if r := recover(); r == nil {
			t.Error("expected MustGetUserID to panic when userID is missing")
		} else {
			// Verify panic message
			msg, ok := r.(string)
			if !ok {
				t.Error("expected panic message to be a string")
			} else if msg != "user ID not found in context - auth middleware not applied?" {
				t.Errorf("unexpected panic message: %s", msg)
			}
		}
	}()

	MustGetUserID(ctx)
}

func TestMustGetUserID_EmptyString(t *testing.T) {
	ctx := context.WithValue(context.Background(), UserIDKey, "")

	defer func() {
		if r := recover(); r == nil {
			t.Error("expected MustGetUserID to panic when userID is empty")
		}
	}()

	MustGetUserID(ctx)
}
