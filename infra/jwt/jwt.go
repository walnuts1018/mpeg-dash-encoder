package jwt

import (
	"fmt"
	"log/slog"

	"github.com/golang-jwt/jwt/v5"
	"github.com/walnuts1018/mpeg-dash-encoder/config"
	"github.com/walnuts1018/mpeg-dash-encoder/consts"
	"github.com/walnuts1018/mpeg-dash-encoder/util/anyslice"
)

type Manager struct {
	JwtSigningKey []byte
}

const media_ids = "media_ids"

func NewManager(jwtSigningKey config.JWTSigningKey) *Manager {
	return &Manager{
		JwtSigningKey: []byte(jwtSigningKey),
	}
}

func (m *Manager) CreateUserToken(
	mediaIDs []string,
) (string, error) {
	claims := jwt.MapClaims{
		"iss":     consts.ApplicationName,
		media_ids: mediaIDs,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	slog.Debug("token created", slog.Any("token", token))

	signed, err := token.SignedString(m.JwtSigningKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return signed, nil
}

func (m *Manager) GetMediaIDsFromToken(token string) ([]string, error) {
	t, err := jwt.Parse(
		token,
		func(t *jwt.Token) (interface{}, error) {
			return m.JwtSigningKey, nil
		},
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}),
		jwt.WithIssuer(consts.ApplicationName),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	claims, ok := t.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("failed to parse claims")
	}
	slog.Debug("claims parsed", slog.Any("claims", fmt.Sprintf("%#v", claims)))

	IDsAny, ok := claims[media_ids].([]any)
	if !ok {
		return nil, fmt.Errorf("failed to parse media_ids: %#v", claims[media_ids])
	}

	parsedIDs, err := anyslice.FromAny[string](IDsAny)
	if err != nil {
		return nil, fmt.Errorf("failed to parse media_ids: %w", err)
	}
	return parsedIDs, nil
}
