# Homebrew Formula for pull-vids
class PullVids < Formula
  desc "Universal video downloader CLI supporting 1000+ websites"
  homepage "https://github.com/vib795/pull-vids"
  version "0.2.0"
  license "MIT"

  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/vib795/pull-vids/releases/download/v0.2.0/pull-vids-darwin-arm64.tar.gz"
      sha256 "REPLACE_WITH_ACTUAL_SHA256_ARM64" # Will be updated on release
    else
      url "https://github.com/vib795/pull-vids/releases/download/v0.2.0/pull-vids-darwin-amd64.tar.gz"
      sha256 "REPLACE_WITH_ACTUAL_SHA256_AMD64" # Will be updated on release
    end
  end

  on_linux do
    if Hardware::CPU.arm?
      url "https://github.com/vib795/pull-vids/releases/download/v0.2.0/pull-vids-linux-arm64.tar.gz"
      sha256 "REPLACE_WITH_ACTUAL_SHA256_LINUX_ARM64" # Will be updated on release
    else
      url "https://github.com/vib795/pull-vids/releases/download/v0.2.0/pull-vids-linux-amd64.tar.gz"
      sha256 "REPLACE_WITH_ACTUAL_SHA256_LINUX_AMD64" # Will be updated on release
    end
  end

  depends_on "ffmpeg"
  depends_on "yt-dlp"

  def install
    if OS.mac?
      if Hardware::CPU.arm?
        bin.install "pull-vids-darwin-arm64" => "pull-vids"
      else
        bin.install "pull-vids-darwin-amd64" => "pull-vids"
      end
    elsif OS.linux?
      if Hardware::CPU.arm?
        bin.install "pull-vids-linux-arm64" => "pull-vids"
      else
        bin.install "pull-vids-linux-amd64" => "pull-vids"
      end
    end
  end

  test do
    assert_match "pull-vids", shell_output("#{bin}/pull-vids --version")
  end
end
