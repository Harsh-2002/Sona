class Sona < Formula
  desc "Audio transcription tool using AssemblyAI"
  homepage "https://github.com/Harsh-2002/Sona"
  license "MIT"

  on_macos do
    if Hardware::CPU.intel?
      url "https://s3.srvr.site/artifact/sona/sona-darwin-amd64"
      sha256 "436f472752a4f62dd4add087cb5f96012517563bd2987048746211338e1bc8c1"
    else
      url "https://s3.srvr.site/artifact/sona/sona-darwin-arm64"
      sha256 "b673f73d73d8b58e6f73c958255f8096493ee7bf28ac49f30bcddc9528ac062d"
    end
  end

  on_linux do
    if Hardware::CPU.intel?
      url "https://s3.srvr.site/artifact/sona/sona-linux-amd64"
      sha256 "6ea4926e89357690fc77db97551ed60b0a51762193246be2b468fcb57f04c7ef"
    else
      url "https://s3.srvr.site/artifact/sona/sona-linux-arm64"
      sha256 "151f774ac18a39694326375b5ecf7545c1f5025a5db1f7a9be147f044d2b5544"
    end
  end

  def install
    bin.install "sona-darwin-amd64" => "sona" if OS.mac? && Hardware::CPU.intel?
    bin.install "sona-darwin-arm64" => "sona" if OS.mac? && Hardware::CPU.arm?
    bin.install "sona-linux-amd64" => "sona" if OS.linux? && Hardware::CPU.intel?
    bin.install "sona-linux-arm64" => "sona" if OS.linux? && Hardware::CPU.arm?
  end

  test do
    system "#{bin}/sona", "--version"
  end
end
