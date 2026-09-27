package main

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"
)

// useFakeYtDlp puts a shell script named yt-dlp, running body, at the front of
// PATH.
func useFakeYtDlp(t *testing.T, body string) {
	t.Helper()
	if runtime.GOOS == "windows" {
		t.Skip("the fake yt-dlp is a POSIX shell script")
	}
	dir := t.TempDir()
	script := "#!/bin/sh\n" + body
	if err := os.WriteFile(filepath.Join(dir, "yt-dlp"), []byte(script), 0755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", dir+string(os.PathListSeparator)+os.Getenv("PATH"))
}

// TestExecuteDownloadKeepsFinalError covers what the retry loop depends on:
// yt-dlp's error text reaching the returned error, so downloadVideo can see a
// 403 and back off instead of giving up.
//
// exec.Cmd.Wait closes the output pipes as soon as the process exits, which
// discards anything not yet read. With a real process that lost the error in
// about 1 of 300 failures, too rarely to catch by repetition. Here a
// background child keeps stderr open and writes the error after the fake
// yt-dlp has exited, so the output is certain to still be unread when the
// exit is seen, and only code that reads to EOF before calling Wait passes.
func TestExecuteDownloadKeepsFinalError(t *testing.T) {
	useFakeYtDlp(t, `
( sleep 0.3; echo "ERROR: unable to download video data: HTTP Error 403: Forbidden" >&2 ) &
exit 1
`)

	config := &Config{
		URL:         "https://example.com/watch?v=test",
		Output:      t.TempDir(),
		Quality:     "best",
		Connections: 1,
		Downloader:  "native",
	}

	err := executeDownload(config)
	if err == nil || !strings.Contains(err.Error(), "HTTP Error 403") {
		t.Fatalf("yt-dlp's final 403 was lost: %v", err)
	}
}

func TestRetryReason(t *testing.T) {
	tests := []struct{ msg, want string }{
		{"ERROR: unable to download video data: HTTP Error 403: Forbidden", "forbidden"},
		{"ERROR: Unable to download video subtitles for 'en': HTTP Error 429: Too Many Requests", "rate-limited"},
		{"ERROR: [youtube] abc: Sign in to confirm you're not a bot. Use --cookies-from-browser", "rate-limited"},
		// The old check matched "bot" anywhere, so any URL or title containing
		// it was retried five times before failing.
		{`ERROR: [generic] "https://example.com/robot-wars" is not a valid URL`, ""},
		{"ERROR: Video unavailable. This video is private", ""},
	}
	for _, tt := range tests {
		if got := retryReason(tt.msg); got != tt.want {
			t.Errorf("retryReason(%q) = %q, want %q", tt.msg, got, tt.want)
		}
	}
}

func TestParseInfoLine(t *testing.T) {
	tests := []struct {
		line, title, duration string
		ok                    bool
	}{
		{infoLinePrefix + "212.0\tRick Astley - Never Gonna Give You Up", "Rick Astley - Never Gonna Give You Up", "3m 32s", true},
		{infoLinePrefix + "19\tMe at the zoo", "Me at the zoo", "0m 19s", true},
		// Livestreams and some extractors have no duration.
		{infoLinePrefix + "NA\tLive now", "Live now", "", true},
		// A tab in the title must survive; only the first one separates.
		{infoLinePrefix + "5\ta\tb", "a\tb", "0m 5s", true},
		{"[download]  45.3% of 12.34MiB at 1.23MiB/s ETA 00:10", "", "", false},
	}
	for _, tt := range tests {
		title, duration, ok := parseInfoLine(tt.line)
		if title != tt.title || duration != tt.duration || ok != tt.ok {
			t.Errorf("parseInfoLine(%q) = %q, %q, %v; want %q, %q, %v",
				tt.line, title, duration, ok, tt.title, tt.duration, tt.ok)
		}
	}
}

func TestCookieHint(t *testing.T) {
	const chromeErr = `ERROR: could not find chrome cookies database in "/Users/u/Library/Application Support/Google/Chrome"`

	mac := cookieHint(chromeErr, "darwin")
	if !strings.Contains(mac, "Full Disk Access") || !strings.Contains(mac, "chrome") {
		t.Errorf("macOS hint should name the browser and Full Disk Access, got %q", mac)
	}

	// Elsewhere the message really does mean the profile isn't there, so the
	// hint must not send people to a macOS setting they don't have.
	other := cookieHint(chromeErr, "windows")
	if strings.Contains(other, "Full Disk Access") {
		t.Errorf("non-macOS hint should not mention Full Disk Access, got %q", other)
	}
	for _, want := range []string{"chrome", "--cookies"} {
		if !strings.Contains(other, want) {
			t.Errorf("non-macOS hint should mention %q, got %q", want, other)
		}
	}

	if got := cookieHint(`ERROR: could not find safari cookies database in "/Users/u/Library/Cookies"`, "darwin"); !strings.Contains(got, "safari") {
		t.Errorf("hint should name safari, got %q", got)
	}

	for _, msg := range []string{
		"ERROR: unable to download video data: HTTP Error 403: Forbidden",
		"ERROR: Video unavailable",
		"",
	} {
		if got := cookieHint(msg, "darwin"); got != "" {
			t.Errorf("cookieHint(%q) = %q, want no hint", msg, got)
		}
	}
}

// TestDownloadVideoAddsCookieHint covers the wiring: a cookie failure is not
// retryable, so it must come back from the first attempt with the hint
// attached rather than after four rounds of backoff.
func TestDownloadVideoAddsCookieHint(t *testing.T) {
	useFakeYtDlp(t, `
echo 'ERROR: could not find chrome cookies database in "/Users/u/Library/Application Support/Google/Chrome"' >&2
exit 1
`)

	config := &Config{
		URL:                "https://example.com/watch?v=test",
		Output:             t.TempDir(),
		Quality:            "best",
		Connections:        1,
		Downloader:         "native",
		CookiesFromBrowser: "chrome",
	}

	done := make(chan error, 1)
	go func() { done <- downloadVideo(config) }()

	select {
	case err := <-done:
		if err == nil || !strings.Contains(err.Error(), "--cookies") {
			t.Fatalf("cookie failure lost its hint: %v", err)
		}
		if strings.Contains(err.Error(), "after 3 retries") {
			t.Errorf("cookie failure should not be retried: %v", err)
		}
	case <-time.After(10 * time.Second):
		t.Fatal("downloadVideo retried a cookie failure instead of returning at once")
	}
}
