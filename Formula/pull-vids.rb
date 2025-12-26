class PullVids < Formula
  desc "Universal video downloader CLI supporting 1000+ websites"
  homepage "https://github.com/vib795/pull-vids"
  license "MIT"

  on_macos do
    on_arm do
      url "https://github.com/vib795/pull-vids/releases/download/v0.2.0/pull-vids-darwin-arm64.tar.gz"
      sha256 "a114a82b6a149568877b7af70926c86787470a855a421b34006c2db33ecaad5f4"
    end
    on_intel do
      url "https://github.com/vib795/pull-vids/releases/download/v0.2.0/pull-vids-darwin-amd64.tar.gz"
      sha256 "1b447f79fca1d236e37988840b4c25bc1fbd0d326600ec8207a1b1a0d17f0b03"
    end
  end

  on_linux do
    on_arm do
      url "https://github.com/vib795/pull-vids/releases/download/v0.2.0/pull-vids-linux-arm64.tar.gz"
      sha256 "51a8d4829f84d37c56c0ee25cf0d6c48e7db83f4fda69a04e35b70e30b73fe89"
    end
    on_intel do
      url "https://github.com/vib795/pull-vids/releases/download/v0.2.0/pull-vids-linux-amd64.tar.gz"
      sha256 "1ef17442d890ba8fb748a1f5cbbd6622d80217ab9ebcc0c0bcc8e87928381a6c"
    end
  end

  depends_on "ffmpeg"
  depends_on "yt-dlp"

  def install
    # Determine the correct binary name based on platform and architecture
    binary_name = if OS.mac?
      Hardware::CPU.arm? ? "pull-vids-darwin-arm64" : "pull-vids-darwin-amd64"
    else
      Hardware::CPU.arm? ? "pull-vids-linux-arm64" : "pull-vids-linux-amd64"
    end

    bin.install binary_name => "pull-vids"
  end

  test do
    assert_match version.to_s, shell_output("#{bin}/pull-vids --version")

    # Test help output
    help_output = shell_output("#{bin}/pull-vids --help 2>&1")
    assert_match "Universal Video Downloader", help_output
    assert_match "Supports 1000+ sites", help_output
  end
end