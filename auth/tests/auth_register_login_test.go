package tests

import (
	"auth/tests/suite"
	"crypto/rand"
	"math/big"
	"testing"
	"time"

	auth1 "github.com/MrKrik/protos/gen/go/auth"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	emptyAppID = 0
	appID      = 1
	appSecret  = "test-secret"

	passDefaultLen = 10
)

func TestRegisterLogin_Login_HappyPath(t *testing.T) {
	ctx, st := suite.New(t)

	login := generateLogin(7)
	password := generatePassword(passDefaultLen)

	respReg, err := st.AuthClient.Register(ctx, &auth1.RegisterRequest{
		Login:    login,
		Password: password,
	})
	require.NoError(t, err)

	respLogin, err := st.AuthClient.Login(ctx, &auth1.LoginRequest{
		Login:    login,
		Password: password,
		AppId:    appID,
	})

	token := respLogin.GetToken()
	require.NotEmpty(t, token)

	loginTime := time.Now()

	tokenParsed, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(appSecret), nil
	})
	require.NoError(t, err)

	claims, ok := tokenParsed.Claims.(jwt.MapClaims)
	require.True(t, ok)

	assert.Equal(t, respReg.GetUserId(), int64(claims["uid"].(float64)))
	assert.Equal(t, login, claims["login"].(string))
	assert.Equal(t, appID, int(claims["app_id"].(float64)))

	const deltaSeconds = 1

	// check if exp of token is in correct range, ttl get from st.Cfg.TokenTTL
	assert.InDelta(t, loginTime.Add(st.Cfg.TokenTTL).Unix(), claims["exp"].(float64), deltaSeconds)
}

func TestRegisterLogin_DuplicatedRegistration(t *testing.T) {
	ctx, st := suite.New(t)

	login := generateLogin(7)
	pass := generatePassword(passDefaultLen)

	respReg, err := st.AuthClient.Register(ctx, &auth1.RegisterRequest{
		Login:    login,
		Password: pass,
	})
	require.NoError(t, err)
	require.NotEmpty(t, respReg.GetUserId())

	respReg, err = st.AuthClient.Register(ctx, &auth1.RegisterRequest{
		Login:    login,
		Password: pass,
	})
	require.Error(t, err)
	assert.Empty(t, respReg.GetUserId())
	assert.ErrorContains(t, err, "user already exists")
}

func TestLogin_FailCases(t *testing.T) {
	ctx, st := suite.New(t)

	tests := []struct {
		name        string
		login       string
		password    string
		appID       int32
		expectedErr string
	}{
		{
			name:        "Login with Empty Password",
			login:       generateLogin(passDefaultLen),
			password:    "",
			appID:       appID,
			expectedErr: "password is required",
		},
		{
			name:        "Login with Empty login",
			login:       "",
			password:    generatePassword(passDefaultLen),
			appID:       appID,
			expectedErr: "login is required",
		},
		{
			name:        "Login with Both Empty login and Password",
			login:       "",
			password:    "",
			appID:       appID,
			expectedErr: "login is required",
		},
		{
			name:        "Login with Non-Matching Password",
			login:       generateLogin(passDefaultLen),
			password:    generatePassword(passDefaultLen),
			appID:       appID,
			expectedErr: "invalid login or password",
		},
		{
			name:        "Login without AppID",
			login:       generateLogin(passDefaultLen),
			password:    generatePassword(passDefaultLen),
			appID:       emptyAppID,
			expectedErr: "app_id is required",
		},
		{
			name:        "Invisible rune",
			login:       generateLogin(passDefaultLen),
			password:    generatePassword(passDefaultLen) + "ㅤㅤ",
			appID:       emptyAppID,
			expectedErr: "invalid login or password",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := st.AuthClient.Register(ctx, &auth1.RegisterRequest{
				Login:    generateLogin(passDefaultLen),
				Password: generatePassword(passDefaultLen),
			})
			require.NoError(t, err)

			_, err = st.AuthClient.Login(ctx, &auth1.LoginRequest{
				Login:    tt.login,
				Password: tt.password,
				AppId:    tt.appID,
			})
			require.Error(t, err)
			require.Contains(t, err.Error(), tt.expectedErr)
		})
	}
}

func generatePassword(lenght int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*"
	password := make([]byte, lenght)
	for i := range password {
		charIndex, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		password[i] = charset[charIndex.Int64()]
	}

	return string(password)
}

func generateLogin(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, length)
	for i := range b {
		index, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		b[i] = charset[index.Int64()]
	}
	return "user_" + string(b)
}
