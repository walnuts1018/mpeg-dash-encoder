package ffmpeg

import (
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/walnuts1018/mpeg-dash-encoder/config"
	"github.com/walnuts1018/mpeg-dash-encoder/util/fileutil"
)

var baseURL = url.URL{Scheme: "https", Host: "to-be-replaced.example.com"}
var OutDirPrefix = "mpeg-dash-encoder-outdir"

type FFMPEG struct {
	fps              int
	preset           Preset
	videoCodec       string
	audioCodec       string
	videoQualityKeys []VideoQualityKey
	logFileDir       string

	useQSV bool
}

func NewFFMPEG(cfg config.Config) (*FFMPEG, error) {
	useQSV := cfg.FFMPEGHwAccel == "qsv"

	var videoCodec string
	if useQSV {
		videoCodec = "h264_qsv"
	} else {
		slog.Warn("QSV is not available, using software codec")
		videoCodec = "libx264"
	}

	return &FFMPEG{
		fps:        30,
		preset:     Medium,
		videoCodec: videoCodec,
		audioCodec: "aac",
		videoQualityKeys: []VideoQualityKey{
			VideoQualityKey360P,
			VideoQualityKey720P,
			VideoQualityKey1080P,
		},
		logFileDir: filepath.Join(cfg.LogDir, "ffmpeg"),
		useQSV:     useQSV,
	}, nil
}

func (f *FFMPEG) GetOutDirPrefix() string {
	return OutDirPrefix
}

func (f *FFMPEG) createArgs(inputFileName, outputDirectory string, audioOnly bool) ([]string, error) {
	args := make([]string, 0, 65)

	if f.useQSV {
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

	args = append(args,
		"-r", fmt.Sprintf("%d", f.fps),
		"-c:v", f.videoCodec,
		"-c:a", f.audioCodec,
	)

	if !f.useQSV {
		args = append(args, "-pix_fmt", "yuv420p")
	}

	if !audioOnly {
		for i, quality := range f.videoQualityKeys {
			var videoFilter string
			if f.useQSV {
				videoFilter = fmt.Sprintf("scale_qsv=%s", videoQualities[quality].Scale)
			} else {
				videoFilter = fmt.Sprintf("scale=%s", videoQualities[quality].Scale)
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

func (f *FFMPEG) Encode(mediaID string, sourceFilePath string, audioOnly bool) (string, error) {
	outDir, err := os.MkdirTemp("", OutDirPrefix)
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
