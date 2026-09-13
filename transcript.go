package main

import (
	"fmt"
	"html"
	"io/fs"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

var (
	// captionTiming matches a cue timing line in either format:
	//   00:00:12.559 --> 00:00:14.350 align:start position:0%   (WebVTT)
	//   00:00:01,360 --> 00:00:03,040                            (SRT)
	captionTiming = regexp.MustCompile(`^([\d:.,]+)\s*-->\s*([\d:.,]+)`)

	// captionTag matches inline markup that carries no words: <i>, <v Name>,
	// and the <00:00:15.360><c> I</c> per-word timing YouTube emits.
	captionTag = regexp.MustCompile(`<[^>]*>`)

	// captionOverride matches the {\an8} positioning codes SRT inherits from ASS.
	captionOverride = regexp.MustCompile(`\{\\[^}]*\}`)
)

// transcriptFormats maps each -f value accepted with --transcript to the
// subtitle format yt-dlp is asked to produce. Plain text has no yt-dlp
// equivalent, so it is fetched as WebVTT and flattened by captionText.
var transcriptFormats = map[string]string{
	"txt": "vtt",
	"srt": "srt",
	"vtt": "vtt",
}

// transcriptFormat normalises the -f value for transcript mode, where an
// unset format means plain text.
func transcriptFormat(format string) (string, error) {
	format = strings.ToLower(format)
	if format == "" {
		return "txt", nil
	}
	if _, ok := transcriptFormats[format]; !ok {
		return "", fmt.Errorf("transcript format must be txt, srt or vtt (got %q)", format)
	}
	return format, nil
}

// subtitleArgs returns the yt-dlp flags that fetch captions for langs as
// subFormat ("srt" or "vtt").
//
// Both --write-subs and --write-auto-subs are passed because most YouTube
// videos only have auto-generated captions. yt-dlp still prefers an uploaded
// track, and uses an auto-generated one only for a language that has none.
// --convert-subs is a no-op when the site already serves subFormat, so ffmpeg
// is only invoked when a conversion is really needed.
func subtitleArgs(langs, subFormat string) []string {
	return []string{
		"--write-subs", "--write-auto-subs",
		"--sub-langs", langs,
		"--sub-format", subFormat + "/best",
		"--convert-subs", subFormat,
	}
}

// collectTranscripts copies the captions yt-dlp wrote into workDir over to
// outputDir, keeping any playlist subdirectories, and returns the paths it
// wrote. Plain-text transcripts are flattened on the way.
//
// Files are copied rather than renamed because workDir lives in the system
// temp directory, which is often a different filesystem from outputDir.
func collectTranscripts(workDir, outputDir, format string) ([]string, error) {
	wantExt := "." + transcriptFormats[format]
	var written []string

	err := filepath.WalkDir(workDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || !strings.EqualFold(filepath.Ext(path), wantExt) {
			return err
		}

		rel, err := filepath.Rel(workDir, path)
		if err != nil {
			return err
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		dest := filepath.Join(outputDir, rel)
		if format == "txt" {
			dest = strings.TrimSuffix(dest, filepath.Ext(dest)) + ".txt"
			data = []byte(captionText(data))
		}

		if err := os.MkdirAll(filepath.Dir(dest), 0755); err != nil {
			return err
		}
		if err := os.WriteFile(dest, data, 0644); err != nil {
			return err
		}
		written = append(written, dest)
		return nil
	})

	return written, err
}

// captionText flattens a WebVTT or SRT file into plain text, one caption line
// per output line.
//
// Only the words survive. Everything before the first cue (the WEBVTT header,
// metadata, STYLE and REGION blocks) is dropped, along with timing lines, NOTE
// blocks and inline markup. Cue identifiers, including SRT's numeric indices,
// are recognised by position rather than content: both formats require a
// blank line between cues, so the line directly above a timing line can never
// be caption text.
//
// YouTube's auto-generated captions scroll: each cue opens by repeating the
// previous cue's last line, so read naively every line appears two or three
// times. A repeat is dropped only when it opens a cue starting at the exact
// millisecond the previous cue ended, because that is how scrolling repeats
// are laid out, while a line genuinely said twice, such as a chorus in
// uploaded captions, follows a gap and is kept. Timestamps that fail to parse
// never count as contiguous, so unclear input keeps its text.
func captionText(data []byte) string {
	src := strings.TrimPrefix(string(data), "\uFEFF")
	lines := strings.Split(strings.ReplaceAll(src, "\r\n", "\n"), "\n")

	var out []string
	seenCue, inNote := false, false
	// contiguous: the current cue starts where the previous one ended.
	// openingLine: no text from the current cue has been seen yet.
	contiguous, openingLine := false, false
	prevEnd := int64(-1)

	for i, raw := range lines {
		line := strings.TrimSpace(raw)

		if m := captionTiming.FindStringSubmatch(line); m != nil {
			start, startOK := cueMillis(m[1])
			end, endOK := cueMillis(m[2])
			contiguous = startOK && prevEnd >= 0 && start == prevEnd
			prevEnd = -1
			if endOK {
				prevEnd = end
			}
			seenCue, inNote, openingLine = true, false, true
			continue
		}

		switch {
		case line == "":
			inNote = false
			continue
		case !seenCue || inNote:
			continue
		case i+1 < len(lines) && captionTiming.MatchString(strings.TrimSpace(lines[i+1])):
			continue // cue identifier
		case (line == "NOTE" || strings.HasPrefix(line, "NOTE ")) &&
			(i == 0 || strings.TrimSpace(lines[i-1]) == ""):
			inNote = true
			continue
		}

		// Strip markup before unescaping, so an escaped "&lt;i&gt;" in the
		// source stays literal text instead of becoming a tag and vanishing.
		text := captionOverride.ReplaceAllString(captionTag.ReplaceAllString(line, ""), "")
		text = strings.Join(strings.Fields(html.UnescapeString(text)), " ")

		if text == "" {
			continue
		}
		repeat := openingLine && contiguous && len(out) > 0 && out[len(out)-1] == text
		openingLine = false
		if repeat {
			continue
		}
		out = append(out, text)
	}

	if len(out) == 0 {
		return ""
	}
	return strings.Join(out, "\n") + "\n"
}

// cueMillis parses a caption timestamp (hh:mm:ss.ttt, mm:ss.ttt, or SRT's
// hh:mm:ss,ttt) into milliseconds. Whole milliseconds are compared rather
// than float seconds so contiguous cues match exactly.
func cueMillis(ts string) (int64, bool) {
	parts := strings.Split(strings.Replace(ts, ",", ".", 1), ":")
	if len(parts) < 2 || len(parts) > 3 {
		return 0, false
	}

	seconds, err := strconv.ParseFloat(parts[len(parts)-1], 64)
	if err != nil {
		return 0, false
	}
	total := seconds
	for i, unit := 0, 60.0; i < len(parts)-1; i, unit = i+1, unit*60 {
		n, err := strconv.Atoi(parts[len(parts)-2-i])
		if err != nil {
			return 0, false
		}
		total += float64(n) * unit
	}
	return int64(math.Round(total * 1000)), true
}
