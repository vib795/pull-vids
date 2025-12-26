class PullVids < Formula
  desc "Universal video downloader CLI supporting 1000+ websites"
  homepage "https://github.com/vib795/pull-vids"
  license "MIT"

  on_macos do
    on_arm do
      url "https://github.com/vib795/pull-vids/releases/download/v0.2.0/pull-vids-darwin-arm64.tar.gz"
      sha256 "a6d7c093bdd5cf387b5651f36dcdb92fec924f4ba6f2307383c0d15834bd5165"
    end
    on_intel do
      url "https://github.com/vib795/pull-vids/releases/download/v0.2.0/pull-vids-darwin-amd64.tar.gz"
      sha256 "596ac58187b127a549d51bf079c4404837b39550b0fd9a816f18d2bbf3eb6af5"
    end
  end

  on_linux do
    on_arm do
      url "https://github.com/vib795/pull-vids/releases/download/v0.2.0/pull-vids-linux-arm64.tar.gz"
      sha256 "f171eadb1884da3b4687a424c01fef56d2ea4394458319221b5c1d0e6e5bac23"
    end
    on_intel do
      url "https://github.com/vib795/pull-vids/releases/download/v0.2.0/pull-vids-linux-amd64.tar.gz"
      sha256 "2146527a437c5dc56ae37ef004a302000a1290eb2affa5611bc86a9b1ac92839"
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
