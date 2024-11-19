package usecase

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"time"

	"github.com/Code-Hex/synchro"
	"github.com/Code-Hex/synchro/tz"
)

type encodeRequest struct {
	mediaID          string
	uploadedFilePath string
}

func (u *Usecase) Run(ctx context.Context) {
	go func(ctx context.Context) {
		for {
			select {
			case req := <-u.encodeQueue:
				if err := u.encode(ctx, req); err != nil {
					slog.Error("failed to encode", slog.Any("error", err))
					// returnしない
				}
			case <-ctx.Done():
				return
			}
		}
	}(ctx)

	tickerFunc := func() {
		slog.Debug("start to download uploaded files")
		req, err := u.downloadUploadedFiles(ctx)
		if err != nil {
			slog.Error("failed to download uploaded files", slog.Any("error", err))
			return
		}

		if req == nil {
			slog.Debug("no uploaded files")
			return
		}
		slog.Debug("downloaded uploaded files", slog.Any("mediaID", req.mediaID), slog.Any("uploadedFilePath", req.uploadedFilePath))

		u.encodeQueue <- *req
	}
	tickerFunc()

	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			tickerFunc()
		case <-ctx.Done():
			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()
			if err := u.shutdown(ctx); err != nil {
				slog.Error("failed to shutdown", slog.Any("error", err))
			}
			return
		}
	}
}

func (u *Usecase) encode(ctx context.Context, req encodeRequest) error {
	slog.Debug("start to encode", slog.Any("mediaID", req.mediaID), slog.Any("uploadedFilePath", req.uploadedFilePath))
	encodedDir, err := u.encoder.Encode(req.mediaID, req.uploadedFilePath, false)
	if err != nil {
		return fmt.Errorf("failed to encode: %w", err)
	}

	go func(ctx context.Context) {
		if err := u.encodedRepo.Upload(ctx, req.mediaID, encodedDir); err != nil {
			slog.Error("failed to upload", slog.Any("error", err))
			return
		}
		if err := u.sourceRepo.DeleteSourceContent(ctx, req.mediaID); err != nil {
			slog.Error("failed to delete source content", slog.Any("error", err))
			// returnしない
		}
		if err := os.RemoveAll(encodedDir); err != nil {
			slog.Error("failed to remove encoded dir", slog.Any("error", err))
			// returnしない
		}

		if err := os.Remove(req.uploadedFilePath); err != nil {
			slog.Error("failed to remove uploaded file", slog.Any("error", err))
			// returnしない
		}
	}(ctx)

	return nil
}

func (u *Usecase) downloadUploadedFiles(ctx context.Context) (*encodeRequest, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	objectInfos := u.sourceRepo.ListUploadedFiles(ctx)
	for objectInfo, err := range objectInfos {
		if err != nil {
			return nil, fmt.Errorf("failed to list objects: %w", err)
		}

		if objectInfo.Tags != nil {
			startAt, ok := objectInfo.Tags["startAt"]
			if ok {
				startAtTime, err := synchro.ParseISO[tz.AsiaTokyo](startAt)
				if err != nil {
					return nil, fmt.Errorf("failed to parse startAt tag: %w", err)
				}
				slog.Debug("startAt", slog.Any("startAt", startAtTime))

				if !synchro.Now[tz.AsiaTokyo]().After(startAtTime.Add(u.encodeTimeout)) {
					continue
				}
			}
		} else {
			slog.Debug("tags not found")
		}

		if err := u.sourceRepo.SetObjectTags(ctx, objectInfo.ID, map[string]string{"startAt": synchro.Now[tz.AsiaTokyo]().Format(time.RFC3339), "hostname": u.hostname}); err != nil {
			return nil, fmt.Errorf("failed to set tags: %w", err)
		}

		object, err := u.sourceRepo.GetSourceContent(ctx, objectInfo.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to get object: %w", err)
		}
		defer object.Close()

		file, err := os.CreateTemp("", "mpeg-dash-encoder-downloaded")
		if err != nil {
			return nil, fmt.Errorf("failed to create temp file: %w", err)
		}
		defer file.Close()

		if _, err := io.Copy(file, object); err != nil {
			return nil, fmt.Errorf("failed to copy object: %w", err)
		}

		return &encodeRequest{
			mediaID:          objectInfo.ID,
			uploadedFilePath: file.Name(),
		}, nil
	}
	return nil, nil
}

func (u *Usecase) shutdown(ctx context.Context) error {
	close(u.encodeQueue)

	objectInfos := u.sourceRepo.ListUploadedFiles(ctx)
	for objectInfo, err := range objectInfos {
		if err != nil {
			slog.Error("failed to list objects", slog.Any("error", err))
			continue
		}

		if objectInfo.Tags != nil {
			hostname, ok := objectInfo.Tags["hostname"]
			if !ok {
				continue
			}
			if hostname != u.hostname {
				continue
			}
		}

		if err := u.sourceRepo.RemoveObjectTags(ctx, objectInfo.ID); err != nil {
			slog.Error("failed to remove tags", slog.Any("error", err))
			continue
		}
	}
	return nil
}
