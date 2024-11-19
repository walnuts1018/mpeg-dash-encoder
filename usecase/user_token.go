package usecase

import (
	"errors"

	"github.com/walnuts1018/mpeg_dash-encoder/domain"
)

func (u *Usecase) CreateUserToken(
	mediaIDs []string,
) (string, error) {
	return u.tokenIssuer.CreateUserToken(mediaIDs)
}

func (u *Usecase) GetMediaIDsFromToken(token string) ([]string, error) {
	ids, err := u.tokenIssuer.GetMediaIDsFromToken(token)
	if err != nil {
		return nil, errors.Join(err, domain.ErrInvalidToken)
	}
	return ids, nil
}
