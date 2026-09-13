package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// youtubeAutoCaptions is the start of a real auto-generated YouTube track,
// kept verbatim, including the whitespace-only lines inside cues. Each cue
// repeats the previous cue's last line before adding new words.
const youtubeAutoCaptions = `WEBVTT
Kind: captions
Language: en

00:00:12.559 --> 00:00:14.350 align:start position:0%
 
so<00:00:12.759><c> in</c>

00:00:14.350 --> 00:00:14.360 align:start position:0%
so in
 

00:00:14.360 --> 00:00:17.029 align:start position:0%
so in
college<00:00:15.360><c> I</c><00:00:15.440><c> was</c><00:00:15.599><c> a</c><00:00:15.719><c> government</c><00:00:16.160><c> major</c><00:00:16.920><c> which</c>

00:00:17.029 --> 00:00:17.039 align:start position:0%
college I was a government major which
 

00:00:17.039 --> 00:00:19.630 align:start position:0%
college I was a government major which
means<00:00:17.320><c> I</c><00:00:17.439><c> had</c><00:00:17.560><c> to</c><00:00:17.720><c> write</c><00:00:17.880><c> a</c><00:00:18.000><c> lot</c><00:00:18.119><c> of</c><00:00:18.439><c> papers</c><00:00:19.439><c> now</c>

00:00:19.630 --> 00:00:19.640 align:start position:0%
means I had to write a lot of papers now
 
`

func TestCaptionText(t *testing.T) {
	// Editors and tools that strip trailing whitespace would silently turn the
	// fixture into one YouTube never produces, so refuse to run without it.
	if !strings.Contains(youtubeAutoCaptions, "\n \n") {
		t.Fatal("youtubeAutoCaptions lost its whitespace-only lines; restore them")
	}

	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "youtube auto captions lose their scrolling repeats",
			input: youtubeAutoCaptions,
			want: "so in\n" +
				"college I was a government major which\n" +
				"means I had to write a lot of papers now\n",
		},
		{
			name: "srt with BOM, CRLF, indices, markup and entities",
			input: "\uFEFF" + strings.ReplaceAll(`1
00:00:01,360 --> 00:00:03,040
{\an8}<i>Tom &amp; Jerry</i>

2
00:00:04,000 --> 00:00:06,000
Two lines
in one cue

3
00:00:07,000 --> 00:00:08,000
Escaped &lt;i&gt; stays literal
`, "\n", "\r\n"),
			want: "Tom & Jerry\n" +
				"Two lines\n" +
				"in one cue\n" +
				"Escaped <i> stays literal\n",
		},
		{
			name: "webvtt header, style, named cue ids and notes are dropped",
			input: `WEBVTT - with a title

STYLE
::cue { color: yellow }

intro
00:01.000 --> 00:04.000
<v Roger>Hello there

NOTE this comment
spans two lines

00:05.000 --> 00:06.000
General Kenobi
`,
			want: "Hello there\nGeneral Kenobi\n",
		},
		{
			// Real timings from an uploaded music video track: the chorus line
			// is sung twice, 760ms apart, and both must survive.
			name: "a line genuinely repeated after a gap is kept",
			input: `WEBVTT

00:01:59.840 --> 00:02:02.960
♪ (Ooh, give you up) ♪

00:02:03.720 --> 00:02:07.360
♪ (Ooh, give you up) ♪
`,
			want: "♪ (Ooh, give you up) ♪\n♪ (Ooh, give you up) ♪\n",
		},
		{
			name: "only the opening line of a contiguous cue can be a scrolling repeat",
			input: `WEBVTT

00:00:01.000 --> 00:00:02.000
go

00:00:02.000 --> 00:00:03.000
go
go
`,
			want: "go\ngo\n",
		},
		{
			name:  "no cues yields nothing",
			input: "WEBVTT\n\n",
			want:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := captionText([]byte(tt.input)); got != tt.want {
				t.Errorf("captionText() =\n%q\nwant\n%q", got, tt.want)
			}
		})
	}
}

func TestCueMillis(t *testing.T) {
	for in, want := range map[string]int64{
		"00:00:14.350": 14350,   // WebVTT
		"01:02:03,004": 3723004, // SRT
		"02:05.500":    125500,  // WebVTT without hours
	} {
		if got, ok := cueMillis(in); !ok || got != want {
			t.Errorf("cueMillis(%q) = %d, %v; want %d", in, got, ok, want)
		}
	}
	for _, bad := range []string{"14.350", "1:2:3:4", "aa:00.000"} {
		if _, ok := cueMillis(bad); ok {
			t.Errorf("cueMillis(%q) should fail", bad)
		}
	}
}

func TestTranscriptFormat(t *testing.T) {
	for in, want := range map[string]string{"": "txt", "txt": "txt", "SRT": "srt", "vtt": "vtt"} {
		if got, err := transcriptFormat(in); err != nil || got != want {
			t.Errorf("transcriptFormat(%q) = %q, %v; want %q", in, got, err, want)
		}
	}
	if _, err := transcriptFormat("mp4"); err == nil {
		t.Error("transcriptFormat(\"mp4\") should be rejected")
	}
}

func TestCollectTranscripts(t *testing.T) {
	write := func(path, content string) {
		t.Helper()
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
	}
	read := func(path string) string {
		t.Helper()
		data, err := os.ReadFile(path)
		if err != nil {
			t.Fatal(err)
		}
		return string(data)
	}

	t.Run("txt flattens, renames and keeps playlist folders", func(t *testing.T) {
		work, out := t.TempDir(), t.TempDir()
		write(filepath.Join(work, "My List", "01 - Talk.en.vtt"), youtubeAutoCaptions)
		write(filepath.Join(work, "stray.part"), "not a caption")

		written, err := collectTranscripts(work, out, "txt")
		if err != nil {
			t.Fatal(err)
		}
		want := filepath.Join(out, "My List", "01 - Talk.en.txt")
		if len(written) != 1 || written[0] != want {
			t.Fatalf("written = %v, want [%s]", written, want)
		}
		if got := read(want); !strings.HasPrefix(got, "so in\ncollege") {
			t.Errorf("transcript not flattened: %q", got)
		}
	})

	t.Run("srt is copied untouched and other formats ignored", func(t *testing.T) {
		work, out := t.TempDir(), t.TempDir()
		srt := "1\n00:00:01,000 --> 00:00:02,000\nHi\n"
		write(filepath.Join(work, "Clip.en.srt"), srt)
		write(filepath.Join(work, "Clip.en.vtt"), youtubeAutoCaptions)

		written, err := collectTranscripts(work, out, "srt")
		if err != nil {
			t.Fatal(err)
		}
		if len(written) != 1 || read(written[0]) != srt {
			t.Errorf("written = %v, want only Clip.en.srt copied verbatim", written)
		}
	})

	t.Run("nothing to collect is not an error", func(t *testing.T) {
		written, err := collectTranscripts(t.TempDir(), t.TempDir(), "txt")
		if err != nil || len(written) != 0 {
			t.Errorf("got %v, %v; want no files and no error", written, err)
		}
	})
}
