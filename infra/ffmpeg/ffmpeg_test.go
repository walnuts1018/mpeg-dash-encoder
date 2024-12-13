package ffmpeg

import (
	"fmt"
	"log/slog"
	"os"
	"path"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/walnuts1018/mpeg-dash-encoder/config"
	"github.com/walnuts1018/mpeg-dash-encoder/util/random"
)

const testFilesDir = "./test_files"

var testfiles = map[string]map[string]string{
	"video": {
		"avi":  "file_example_AVI_480_750kB.avi",
		"mov":  "file_example_MOV_480_700kB.mov",
		"mp4":  "file_example_MP4_480_1_5MG.mp4",
		"webm": "file_example_WEBM_480_900KB.webm",
	},
	"audio": {
		"mp3": "file_example_MP3_700KB.mp3",
		"wav": "file_example_WAV_1MG.wav",
	},
}

func init() {
	slog.SetLogLoggerLevel(slog.LevelDebug)
}

func TestFFMPEG_CreateArgs(t *testing.T) {
	type args struct {
		inputFileName   string
		outputDirectory string
		audioOnly       bool
	}
	tests := []struct {
		name   string
		ffmpeg config.FFmpegConfig
		args   args
		want   []string
	}{
		{
			name: "with video, no hwAccel",
			ffmpeg: config.FFmpegConfig{
				LogDir:     "./log",
				FPS:        30,
				Preset:     config.Veryslow,
				HWAccel:    config.FFmpegHWAccelNone,
				AudioCodec: "aac",
			},
			args: args{
				inputFileName:   "input.mp4",
				outputDirectory: "Dash",
				audioOnly:       false,
			},
			want: []string{
				"-i", "input.mp4",
				"-y",
				"-hide_banner",
				"-progress", "-",
				"-preset", "veryslow",
				"-keyint_min", "100",
				"-g", "100",
				"-sc_threshold", "0",
				"-r", "30",
				"-c:v", "libx264",
				"-c:a", "aac",
				"-pix_fmt", "yuv420p",

				// 360p
				"-map", "v:0?",
				"-filter:v:0", "scale=-1:360",
				"-b:v:0", "365k",
				"-maxrate:0", "390k",
				"-bufsize:0", "640k",

				// 720p
				"-map", "v:0?",
				"-filter:v:1", "scale=-1:720",
				"-b:v:1", "4.5M",
				"-maxrate:1", "4.8M",
				"-bufsize:1", "8M",

				// 1080p
				"-map", "v:0?",
				"-filter:v:2", "scale=-1:1080",
				"-b:v:2", "7.8M",
				"-maxrate:2", "8.3M",
				"-bufsize:2", "14M",

				"-map", "0:a",
				"-init_seg_name", `init$RepresentationID$.$ext$`,
				"-media_seg_name", `chunk$RepresentationID$-$Number%05d$.$ext$`,
				"-use_template", "1",
				"-use_timeline", "1",
				"-seg_duration", "4",
				"-adaptation_sets", `id=0,streams=a id=1,streams=v`,
				"-f", "dash",
				filepath.Join("Dash", "dash.mpd"),
			},
		},
		{
			name: "with video, qsv",
			ffmpeg: config.FFmpegConfig{
				LogDir:     "./log",
				FPS:        30,
				Preset:     config.Veryslow,
				HWAccel:    config.FFmpegHWAccelQSV,
				AudioCodec: "aac",
			},
			args: args{
				inputFileName:   "input.mp4",
				outputDirectory: "Dash",
				audioOnly:       false,
			},
			want: []string{
				"-hwaccel", "qsv",
				"-hwaccel_output_format", "qsv",
				"-i", "input.mp4",
				"-y",
				"-hide_banner",
				"-progress", "-",
				"-preset", "veryslow",
				"-keyint_min", "100",
				"-g", "100",
				"-sc_threshold", "0",
				"-r", "30",
				"-c:v", "h264_qsv",
				"-c:a", "aac",

				// 360p
				"-map", "v:0?",
				"-filter:v:0", "scale_qsv=-1:360",
				"-b:v:0", "365k",
				"-maxrate:0", "390k",
				"-bufsize:0", "640k",

				// 720p
				"-map", "v:0?",
				"-filter:v:1", "scale_qsv=-1:720",
				"-b:v:1", "4.5M",
				"-maxrate:1", "4.8M",
				"-bufsize:1", "8M",

				// 1080p
				"-map", "v:0?",
				"-filter:v:2", "scale_qsv=-1:1080",
				"-b:v:2", "7.8M",
				"-maxrate:2", "8.3M",
				"-bufsize:2", "14M",

				"-map", "0:a",
				"-init_seg_name", `init$RepresentationID$.$ext$`,
				"-media_seg_name", `chunk$RepresentationID$-$Number%05d$.$ext$`,
				"-use_template", "1",
				"-use_timeline", "1",
				"-seg_duration", "4",
				"-adaptation_sets", `id=0,streams=a id=1,streams=v`,
				"-f", "dash",
				filepath.Join("Dash", "dash.mpd"),
			},
		},
		{
			name: "audioOnly, no hwAccel",
			ffmpeg: config.FFmpegConfig{
				LogDir:     "./log",
				FPS:        30,
				Preset:     config.Veryslow,
				HWAccel:    config.FFmpegHWAccelNone,
				AudioCodec: "aac",
			},
			args: args{
				inputFileName:   "input.mp4",
				outputDirectory: "Dash",
				audioOnly:       true,
			},
			want: []string{
				"-i", "input.mp4",
				"-y",
				"-hide_banner",
				"-progress", "-",
				"-r", "30",
				"-c:v", "libx264",
				"-c:a", "aac",
				"-pix_fmt", "yuv420p",
				"-map", "0:a",
				"-init_seg_name", `init$RepresentationID$.$ext$`,
				"-media_seg_name", `chunk$RepresentationID$-$Number%05d$.$ext$`,
				"-use_template", "1",
				"-use_timeline", "1",
				"-seg_duration", "4",
				"-adaptation_sets", `id=0,streams=a`,
				"-f", "dash",
				filepath.Join("Dash", "dash.mpd"),
			},
		},
		{
			name: "audioOnly, qsv",
			ffmpeg: config.FFmpegConfig{
				LogDir:     "./log",
				FPS:        30,
				Preset:     config.Veryslow,
				HWAccel:    config.FFmpegHWAccelQSV,
				AudioCodec: "aac",
			},
			args: args{
				inputFileName:   "input.mp4",
				outputDirectory: "Dash",
				audioOnly:       true,
			},
			want: []string{
				"-hwaccel", "qsv",
				"-hwaccel_output_format", "qsv",
				"-i", "input.mp4",
				"-y",
				"-hide_banner",
				"-progress", "-",
				"-r", "30",
				"-c:v", "h264_qsv",
				"-c:a", "aac",
				"-map", "0:a",
				"-init_seg_name", `init$RepresentationID$.$ext$`,
				"-media_seg_name", `chunk$RepresentationID$-$Number%05d$.$ext$`,
				"-use_template", "1",
				"-use_timeline", "1",
				"-seg_duration", "4",
				"-adaptation_sets", `id=0,streams=a`,
				"-f", "dash",
				filepath.Join("Dash", "dash.mpd"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f, err := NewFFMPEG(tt.ffmpeg)
			assert.NoError(t, err)
			assert.NotNil(t, f)

			got, err := f.createArgs(tt.args.inputFileName, tt.args.outputDirectory, tt.args.audioOnly)
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestFFMPEG_Encode(t *testing.T) {
	f, err := NewFFMPEG(config.FFmpegConfig{
		LogDir:     "./log",
		FPS:        30,
		Preset:     config.Veryfast,
		HWAccel:    config.FFmpegHWAccelNone,
		AudioCodec: "aac",
	})
	if err != nil {
		t.Errorf("failed to create ffmpeg: %v", err)
		return
	}

	type args struct {
		id        string
		path      string
		audioOnly bool
	}
	type test struct {
		name    string
		args    args
		wantErr bool
	}

	tests := make([]test, 0)

	// video
	for k, v := range testfiles["video"] {
		id, err := random.String(32, random.Alphanumeric)
		if err != nil {
			t.Errorf("failed to gen id: %s", err)
		}
		tests = append(tests, test{
			name: k,
			args: args{
				id:        id,
				path:      path.Join(testFilesDir, v),
				audioOnly: false,
			},
			wantErr: false,
		})
	}

	// audio
	for k, v := range testfiles["audio"] {
		id, err := random.String(32, random.Alphanumeric)
		if err != nil {
			t.Errorf("failed to gen id: %s", err)
		}
		tests = append(tests, test{
			name: k,
			args: args{
				id:        id,
				path:      path.Join(testFilesDir, v),
				audioOnly: true,
			},
			wantErr: false,
		})
	}

	workdir, err := os.Getwd()
	if err != nil {
		t.Errorf("failed to get workdir: %v", err)
		return
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hlsDir, err := f.Encode(tt.args.id, filepath.Join(workdir, tt.args.path), tt.args.audioOnly)
			if (err != nil) != tt.wantErr {
				t.Errorf("FFMPEG.Encode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			fmt.Println(hlsDir)
		})
	}
}
