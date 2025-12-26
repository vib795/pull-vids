# Creating Demo GIF

This guide shows how to create the demo GIF for pull-vids using VHS.

## Install VHS

```bash
# macOS
brew install vhs

# Or with Go
go install github.com/charmbracelet/vhs@latest
```

## Option 1: Quick Demo (Recommended)

Creates a shorter, focused demo perfect for README:

```bash
vhs demo-quick.tape
```

This creates `demo-quick.gif` showing:
- Version check
- YouTube download
- Audio extraction
- Quality selection

**Duration:** ~15 seconds
**File size:** ~2-3 MB

## Option 2: Full Demo

Creates a comprehensive demo with installation:

```bash
vhs demo.tape
```

This creates `demo.gif` showing:
- Homebrew installation
- Help output
- Multiple download examples
- Different use cases

**Duration:** ~30-40 seconds
**File size:** ~5-8 MB

## Tips for Best Results

### 1. Clean Terminal
```bash
# Clear history
clear
history -c

# Use a clean shell
bash --noprofile --norc
```

### 2. Customize Theme

Edit the `.tape` file to change the theme:

```tape
Set Theme "Dracula"          # Dark purple theme
Set Theme "Catppuccin Mocha" # Modern dark theme
Set Theme "Nord"             # Cool blue theme
Set Theme "Tokyo Night"      # Dark theme
Set Theme "One Dark"         # Atom-inspired
```

### 3. Adjust Timing

Make it faster or slower:

```tape
# Faster
Sleep 500ms  # Change to 300ms

# Slower
Sleep 1s     # Change to 2s
```

### 4. Change Size

Make it bigger or smaller:

```tape
Set Width 1200   # Increase for wider
Set Height 600   # Increase for taller
Set FontSize 16  # Increase for bigger text
```

## Recording a Real Demo (Alternative)

If you prefer recording a real terminal session:

### Using Terminalizer

```bash
# Install
npm install -g terminalizer

# Record
terminalizer record demo -c terminalizer.yml

# Render
terminalizer render demo
```

### Using Asciinema + agg

```bash
# Install
brew install asciinema agg

# Record
asciinema rec demo.cast

# Convert to GIF
agg demo.cast demo.gif
```

## Example Real Recording Script

Create `record-demo.sh`:

```bash
#!/bin/bash

echo "# pull-vids - Universal Video Downloader"
sleep 2

echo "# Install via Homebrew"
sleep 1
echo "$ brew tap vib795/tap && brew install pull-vids"
sleep 2

echo ""
echo "# Download a video"
sleep 1
pull-vids --version
sleep 1

pull-vids --help | head -20
sleep 3

echo ""
echo "# Example: Download from YouTube"
sleep 1
echo "$ pull-vids 'https://youtu.be/VIDEO_ID'"
sleep 2

echo ""
echo "# That's it! Star on GitHub: github.com/vib795/pull-vids ⭐"
sleep 2
```

Then record:

```bash
chmod +x record-demo.sh
asciinema rec demo.cast -c ./record-demo.sh
agg demo.cast demo.gif --theme monokai
```

## Optimizing GIF Size

If the GIF is too large:

### Using gifsicle

```bash
# Install
brew install gifsicle

# Optimize
gifsicle -O3 --colors 256 demo.gif -o demo-optimized.gif

# Further compress (lossy)
gifsicle -O3 --lossy=80 --colors 128 demo.gif -o demo-small.gif
```

### Using ffmpeg

```bash
# Convert to optimized GIF
ffmpeg -i demo.gif -vf "fps=10,scale=1200:-1:flags=lanczos" -c:v gif demo-optimized.gif
```

## Adding to README

Once you have the GIF:

```markdown
## Demo

![pull-vids demo](demo.gif)

*Quick demo showing installation and basic usage*
```

Or with a link:

```markdown
## Demo

[![Demo](demo-quick.gif)](https://github.com/vib795/pull-vids)

*Click to see more examples*
```

## Troubleshooting

### VHS not found
```bash
# Check installation
which vhs

# Reinstall if needed
brew reinstall vhs
```

### Font issues
```bash
# Install common fonts
brew tap homebrew/cask-fonts
brew install font-fira-code
brew install font-jetbrains-mono

# Update tape file
Set FontFamily "JetBrains Mono"
```

### GIF too large
- Reduce dimensions: `Set Width 1000` and `Set Height 500`
- Reduce font size: `Set FontSize 14`
- Use fewer colors in theme
- Shorten sleep times
- Optimize with gifsicle (see above)

## Best Practices

1. **Keep it short** - Under 20 seconds is ideal
2. **Show key features** - Installation + 2-3 examples
3. **Use real commands** - Don't fake the output
4. **Good contrast** - Use a theme with clear text
5. **Reasonable size** - Under 5 MB for GitHub
6. **Test first** - Generate and review before committing

Happy recording! 🎬
