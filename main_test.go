package main

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// useFakeYtDlp puts a shell script named yt-dlp at the front of PATH. The
// metadata lookup executeDownload makes first gets an empty success, and
// every other invocation runs body.
func useFakeYtDlp(t *testing.T, body string) {
	t.Helper()
	if runtime.GOOS == "windows" {
		t.Skip("the fake yt-dlp is a POSIX shell script")
	}
	dir := t.TempDir()
	script := "#!/bin/sh\nif [ \"$1\" = --dump-json ]; then exit 0; fi\n" + body
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
