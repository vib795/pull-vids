# VHS documentation: https://github.com/charmbracelet/vhs

# Set up the terminal
Set FontSize 16
Set Width 1200
Set Height 600
Set Padding 20
Set Theme "Dracula"

# Show the banner/intro
Type "# pull-vids - Universal Video Downloader" Sleep 500ms Enter Sleep 1s

# Installation
Type "# Install via Homebrew" Sleep 500ms Enter
Type "brew tap vib795/tap" Sleep 500ms Enter
Sleep 3s
Type "brew install pull-vids" Sleep 500ms Enter
Sleep 4s

# Version check
Type "pull-vids --version" Sleep 500ms Enter
Sleep 2s

# Show help
Type "pull-vids --help" Sleep 500ms Enter
Sleep 4s

# Clear screen
Type "clear" Sleep 500ms Enter
Sleep 1s

# Demo 1: Download video
Type "# Download a YouTube video" Sleep 500ms Enter
Type "pull-vids 'https://www.youtube.com/watch?v=dQw4w9WgXcQ'" Sleep 500ms Enter
Sleep 6s
Ctrl+C
Sleep 1s

# Demo 2: Audio only
Type "clear" Sleep 500ms Enter
Type "# Extract audio only (MP3)" Sleep 500ms Enter
Type "pull-vids --audio 'https://www.youtube.com/watch?v=dQw4w9WgXcQ'" Sleep 500ms Enter
Sleep 6s
Ctrl+C
Sleep 1s

# Demo 3: Specific quality
Type "clear" Sleep 500ms Enter
Type "# Download in 720p" Sleep 500ms Enter
Type "pull-vids --quality 720p 'https://www.youtube.com/watch?v=dQw4w9WgXcQ'" Sleep 500ms Enter
Sleep 6s
Ctrl+C
Sleep 1s

# Final message
Type "clear" Sleep 500ms Enter
Type "# That's pull-vids! ✨" Sleep 1s Enter
Type "# Free, fast, and supports 1000+ websites" Sleep 1s Enter
Type "# Star on GitHub: github.com/vib795/pull-vids" Sleep 2s Enter

# Save output
Sleep 3s
Output demo.gif
