package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/schollz/progressbar/v3"
)

const version = "0.2.0"

var (
	cyan    = color.New(color.FgCyan)
	green   = color.New(color.FgGreen)
	yellow  = color.New(color.FgYellow)
	red     = color.New(color.FgRed)
	magenta = color.New(color.FgMagenta)
)

type Config struct {
	URL        string
	Output     string
	Quality    string
	AudioOnly  bool
	Playlist   bool
	Format     string
	NoBanner   bool
	ShowVersion bool
}

type VideoInfo struct {
	Title    string  `json:"title"`
	Duration float64 `json:"duration"`
	Uploader string  `json:"uploader"`
}

func printBanner() {
	cyan.Println("╔═══════════════════════════════════════╗")
	cyan.Printf("║          pull-vids v%-8s       ║\n", version)
	cyan.Println("║   Free YouTube Video Downloader CLI   ║")
	cyan.Println("║         Built with Go - Fast! 🚀      ║")
	cyan.Println("╚═══════════════════════════════════════╝")
	fmt.Println()
}

func checkYtDlp() error {
	cmd := exec.Command("yt-dlp", "--version")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("yt-dlp not found. Install it with: pip install yt-dlp")
	}
	return nil
}

func getFormatString(quality string, audioOnly bool) string {
	if audioOnly {
		return "bestaudio/best"
	}

	qualityMap := map[string]string{
		"best":   "bestvideo+bestaudio/best",
		"high":   "bestvideo[height<=1080]+bestaudio/best[height<=1080]",
		"medium": "bestvideo[height<=720]+bestaudio/best[height<=720]",
		"low":    "bestvideo[height<=480]+bestaudio/best[height<=480]",
		"2160p":  "bestvideo[height<=2160]+bestaudio/best[height<=2160]",
		"1440p":  "bestvideo[height<=1440]+bestaudio/best[height<=1440]",
		"1080p":  "bestvideo[height<=1080]+bestaudio/best[height<=1080]",
		"720p":   "bestvideo[height<=720]+bestaudio/best[height<=720]",
		"480p":   "bestvideo[height<=480]+bestaudio/best[height<=480]",
		"360p":   "bestvideo[height<=360]+bestaudio/best[height<=360]",
	}

	if format, ok := qualityMap[quality]; ok {
		return format
	}
	return "bestvideo+bestaudio/best"
}

func getVideoInfo(url string) (*VideoInfo, error) {
	cmd := exec.Command("yt-dlp", "--dump-json", "--no-playlist", url)
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var info VideoInfo
	if err := json.Unmarshal(output, &info); err != nil {
		return nil, err
	}

	return &info, nil
}

func downloadVideo(config *Config) error {
	// Ensure output directory exists
	outputDir := os.ExpandEnv(config.Output)
	if strings.HasPrefix(config.Output, "~") {
		home, _ := os.UserHomeDir()
		outputDir = filepath.Join(home, config.Output[1:])
	}

	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Build output template
	outputTemplate := filepath.Join(outputDir, "%(title)s.%(ext)s")
	if config.Playlist {
		outputTemplate = filepath.Join(outputDir, "%(playlist)s", "%(playlist_index)s - %(title)s.%(ext)s")
	}

	cyan.Printf("Starting download from: %s\n", config.URL)
	cyan.Printf("Output directory: %s\n", outputDir)

	if config.AudioOnly {
		yellow.Println("Mode: Audio only")
	} else {
		yellow.Printf("Quality: %s\n", config.Quality)
	}
	fmt.Println()

	// Get video info first
	if !config.Playlist {
		info, err := getVideoInfo(config.URL)
		if err == nil {
			mins := int(info.Duration) / 60
			secs := int(info.Duration) % 60
			magenta.Printf("Title: %s\n", info.Title)
			magenta.Printf("Duration: %dm %ds\n", mins, secs)
			fmt.Println()
		}
	}

	// Build yt-dlp command
	args := []string{
		"--newline",
		"--progress",
		"-f", getFormatString(config.Quality, config.AudioOnly),
		"-o", outputTemplate,
	}

	if !config.Playlist {
		args = append(args, "--no-playlist")
	}

	if config.AudioOnly {
		format := "mp3"
		if config.Format != "" {
			format = config.Format
		}
		args = append(args, "-x", "--audio-format", format, "--audio-quality", "192K")
	} else if config.Format != "" {
		args = append(args, "--merge-output-format", config.Format)
	}

	args = append(args, config.URL)

	// Execute download
	cmd := exec.Command("yt-dlp", args...)

	// Create pipe for stdout
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	// Create progress bar
	bar := progressbar.NewOptions(-1,
		progressbar.OptionSetDescription(green.Sprint("Downloading...")),
		progressbar.OptionSetWidth(50),
		progressbar.OptionShowBytes(true),
		progressbar.OptionShowCount(),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        green.Sprint("="),
			SaucerHead:    green.Sprint(">"),
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}),
	)

	// Read output
	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			line := scanner.Text()
			if strings.Contains(line, "[download]") {
				if strings.Contains(line, "ETA") {
					bar.Add(1)
				} else if strings.Contains(line, "has already been downloaded") {
					yellow.Println("Already downloaded, skipping...")
				} else if strings.Contains(line, "Downloading") {
					fmt.Println(line)
				}
			}
		}
	}()

	// Read stderr for errors
	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			line := scanner.Text()
			if strings.Contains(line, "ERROR") {
				red.Println(line)
			}
		}
	}()

	// Wait for completion
	err = cmd.Wait()

	if err != nil {
		return fmt.Errorf("download failed: %w", err)
	}

	bar.Finish()
	fmt.Println()

	green.Println("==================================================")
	green.Println("✓ Download completed successfully!")
	green.Println("==================================================")
	fmt.Println()

	return nil
}

func cleanURL(url string) string {
	// Remove backslash escapes that shells add to special characters
	// These are common when users don't quote URLs properly
	replacements := []string{
		`\?`, `?`,
		`\=`, `=`,
		`\&`, `&`,
		`\:`, `:`,
		`\/`, `/`,
	}

	cleaned := url
	for i := 0; i < len(replacements); i += 2 {
		cleaned = strings.ReplaceAll(cleaned, replacements[i], replacements[i+1])
	}

	return cleaned
}

func parseFlags() *Config {
	config := &Config{}

	flag.StringVar(&config.Output, "o", "~/Downloads/pull-vids", "Output directory")
	flag.StringVar(&config.Output, "output", "~/Downloads/pull-vids", "Output directory")
	flag.StringVar(&config.Quality, "q", "best", "Video quality (best, high, medium, low, 2160p, 1440p, 1080p, 720p, 480p, 360p)")
	flag.StringVar(&config.Quality, "quality", "best", "Video quality (best, high, medium, low, 2160p, 1440p, 1080p, 720p, 480p, 360p)")
	flag.BoolVar(&config.AudioOnly, "a", false, "Download audio only")
	flag.BoolVar(&config.AudioOnly, "audio-only", false, "Download audio only")
	flag.BoolVar(&config.Playlist, "p", false, "Download entire playlist")
	flag.BoolVar(&config.Playlist, "playlist", false, "Download entire playlist")
	flag.StringVar(&config.Format, "f", "", "Output format (mp4, mkv, mp3, m4a, etc.)")
	flag.StringVar(&config.Format, "format", "", "Output format (mp4, mkv, mp3, m4a, etc.)")
	flag.BoolVar(&config.NoBanner, "no-banner", false, "Don't show the banner")
	flag.BoolVar(&config.ShowVersion, "v", false, "Show version")
	flag.BoolVar(&config.ShowVersion, "version", false, "Show version")

	flag.Usage = func() {
		printBanner()
		fmt.Fprintf(os.Stderr, "Usage: pull-vids [options] <URL>\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  # Download video in best quality\n")
		fmt.Fprintf(os.Stderr, "  pull-vids \"https://www.youtube.com/watch?v=VIDEO_ID\"\n\n")
		fmt.Fprintf(os.Stderr, "  # Download audio only as MP3\n")
		fmt.Fprintf(os.Stderr, "  pull-vids -a \"https://www.youtube.com/watch?v=VIDEO_ID\"\n\n")
		fmt.Fprintf(os.Stderr, "  # Download in 720p quality to specific directory\n")
		fmt.Fprintf(os.Stderr, "  pull-vids -q 720p -o ~/Videos \"https://www.youtube.com/watch?v=VIDEO_ID\"\n\n")
		fmt.Fprintf(os.Stderr, "  # Download entire playlist\n")
		fmt.Fprintf(os.Stderr, "  pull-vids -p \"https://www.youtube.com/playlist?list=PLAYLIST_ID\"\n\n")
	}

	flag.Parse()

	return config
}

func main() {
	config := parseFlags()

	if config.ShowVersion {
		fmt.Printf("pull-vids %s\n", version)
		os.Exit(0)
	}

	// Show banner unless disabled
	if !config.NoBanner {
		printBanner()
	}

	// Get URL from remaining args
	args := flag.Args()
	if len(args) == 0 {
		red.Println("✗ Error: YouTube URL required")
		fmt.Println("\nRun 'pull-vids -h' for usage information")
		os.Exit(1)
	}

	config.URL = cleanURL(args[0])

	// Validate URL
	if !strings.Contains(config.URL, "youtube.com") && !strings.Contains(config.URL, "youtu.be") {
		red.Println("✗ Error: Invalid YouTube URL")
		os.Exit(1)
	}

	// Check if yt-dlp is installed
	if err := checkYtDlp(); err != nil {
		red.Printf("✗ Error: %v\n", err)
		os.Exit(1)
	}

	// Download the video
	start := time.Now()
	if err := downloadVideo(config); err != nil {
		red.Printf("✗ Error: %v\n", err)
		os.Exit(1)
	}

	elapsed := time.Since(start)
	cyan.Printf("Total time: %.2fs\n", elapsed.Seconds())
}
