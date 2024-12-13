package ffmpeg

import (
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"

	"github.com/walnuts1018/mpeg-dash-encoder/config"
	"github.com/walnuts1018/mpeg-dash-encoder/util/fileutil"
)

const (
	outDirPrefix = "mpeg-dash-encoder-outdir"
)

type FFmpeg struct {
	fps              string
	preset           config.FFmpegPreset
	audioCodec       string
	videoQualityKeys []VideoQualityKey
	logFileDir       string
	hwAccel          config.FFmpegHWAccel
}

func NewFFMPEG(cfg config.FFmpegConfig) (*FFmpeg, error) {
	if cfg.HWAccel == config.FFmpegHWAccelNone {
		slog.Warn("QSV is not available, using software codec")
	}

	return &FFmpeg{
		fps:        strconv.Itoa(cfg.FPS),
		preset:     cfg.Preset,
		audioCodec: cfg.AudioCodec,
		videoQualityKeys: []VideoQualityKey{
			VideoQualityKey360P,
			VideoQualityKey720P,
			VideoQualityKey1080P,
		},
		logFileDir: cfg.LogDir,
		hwAccel:    cfg.HWAccel,
	}, nil
}

func (f *FFmpeg) GetOutDirPrefix() string {
	return outDirPrefix
}

func (f *FFmpeg) createArgs(inputFileName, outputDirectory string, audioOnly bool) ([]string, error) {
	args := make([]string, 0, 65)

	// hwaccel option
	switch f.hwAccel {
	case config.FFmpegHWAccelQSV:
		args = append(args,
			"-hwaccel", "qsv",
			"-hwaccel_output_format", "qsv",
		)
	}

	args = append(args,
		"-i", inputFileName,
		"-y",
		"-hide_banner",
		"-progress", "-",
	)

	if !audioOnly {
		args = append(args,
			"-preset", string(f.preset),
			"-keyint_min", "100",
			"-g", "100",
			"-sc_threshold", "0",
		)
	}

	var videoCodec string
	switch f.hwAccel {
	case config.FFmpegHWAccelQSV:
		videoCodec = "h264_qsv"
	case config.FFmpegHWAccelNone:
		videoCodec = "libx264"
	}

	args = append(args,
		"-r", f.fps,
		"-c:v", videoCodec,
		"-c:a", f.audioCodec,
	)

	switch f.hwAccel {
	case config.FFmpegHWAccelNone:
		args = append(args, "-pix_fmt", "yuv420p")
	}

	if !audioOnly {
		for i, quality := range f.videoQualityKeys {
			var videoFilter string
			switch f.hwAccel {
			case config.FFmpegHWAccelQSV:
				videoFilter = "scale_qsv=" + videoQualities[quality].Scale
			case config.FFmpegHWAccelNone:
				videoFilter = "scale=" + videoQualities[quality].Scale
			}

			args = append(args,
				"-map", "v:0?",
				fmt.Sprintf("-filter:v:%d", i), videoFilter,
				fmt.Sprintf("-b:v:%d", i), videoQualities[quality].Bitrate,
				fmt.Sprintf("-maxrate:%d", i), videoQualities[quality].MaxBitrate,
				fmt.Sprintf("-bufsize:%d", i), videoQualities[quality].Bufsize,
			)
		}
	}

	args = append(args,
		"-map", "0:a",
		"-init_seg_name", `init$RepresentationID$.$ext$`,
		"-media_seg_name", `chunk$RepresentationID$-$Number%05d$.$ext$`,
		"-use_template", "1",
		"-use_timeline", "1",
		"-seg_duration", "4",
	)

	if audioOnly {
		args = append(args,
			"-adaptation_sets", `id=0,streams=a`,
		)
	} else {
		args = append(args,
			"-adaptation_sets", `id=0,streams=a id=1,streams=v`,
		)
	}

	args = append(args, "-f", "dash", filepath.Join(outputDirectory, "dash.mpd"))
	return args, nil
}

func (f *FFmpeg) Encode(mediaID string, sourceFilePath string, audioOnly bool) (string, error) {
	outDir, err := os.MkdirTemp("", outDirPrefix)
	if err != nil {
		return "", err
	}

	args, err := f.createArgs(sourceFilePath, outDir, audioOnly)
	if err != nil {
		return "", err
	}
	slog.Debug("ffmpeg args", slog.Any("args", args))

	cmd := exec.Command("ffmpeg", args...)
	cmd.Dir = outDir

	var stdout bytes.Buffer
	var stderr bytes.Buffer

	logfile, err := fileutil.CreateFileRecursive(filepath.Join(f.logFileDir, mediaID+".log"))
	if err != nil {
		return "", fmt.Errorf("failed to create log file: %w", err)
	}

	cmd.Stdout = io.MultiWriter(logfile, &stdout)
	cmd.Stderr = io.MultiWriter(logfile, &stderr)

	if err := cmd.Run(); err != nil {
		slog.Error("ffmpeg error",
			slog.String("stdout", stdout.String()),
			slog.String("stderr", stderr.String()),
		)
		return "", fmt.Errorf("failed to run ffmpeg: %w", err)
	}

	return outDir, nil
}
