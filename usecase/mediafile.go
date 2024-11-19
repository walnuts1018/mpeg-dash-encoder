package usecase

import (
	"context"
	"io"
)

func (u *Usecase) GetMediaFile(ctx context.Context, mediaID string, fileName string) (io.ReadSeekCloser, error) {
	return u.encodedRepo.GetObject(ctx, mediaID, fileName)
}
