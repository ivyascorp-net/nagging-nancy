class NaggingNancy < Formula
  desc "Terminal-based reminders and task management with TUI interface"
  homepage "https://github.com/ivyascorp-net/nagging-nancy"
  url "https://github.com/ivyascorp-net/nagging-nancy/archive/refs/tags/v1.0.0.tar.gz"
  sha256 "INSERT_TARBALL_SHA256_HERE"
  license "MIT"
  head "https://github.com/ivyascorp-net/nagging-nancy.git", branch: "main"

  depends_on "go" => :build

  def install
    ldflags = %W[
      -s -w
      -X main.version=#{version}
      -X main.buildTime=#{time.utc.strftime("%Y-%m-%d_%H:%M:%S")}
      -X main.gitCommit=#{tap.user}
    ]

    system "go", "build", *std_go_args(ldflags:), "./cmd/nancy"
  end

  def caveats
    <<~EOS
      nagging-nancy stores data in ~/.config/nancy/
      
      To start the notification daemon:
        nancy daemon start
      
      To add a reminder:
        nancy add "Call mom"
        nancy add "Meeting tomorrow at 2pm" --priority high
      
      To start the interactive TUI:
        nancy
    EOS
  end

  test do
    assert_match version.to_s, shell_output("#{bin}/nancy version")
    
    # Test basic functionality
    system bin/"nancy", "add", "test reminder"
    output = shell_output("#{bin}/nancy list")
    assert_match "test reminder", output
    
    # Test daemon commands (without actually starting)
    assert_match "daemon is not running", shell_output("#{bin}/nancy daemon status", 1)
  end
end