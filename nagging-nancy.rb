class NaggingNancy < Formula
  desc "Terminal-based reminders and task management with TUI interface"
  homepage "https://github.com/ivyascorp-net/nagging-nancy"
  url "https://github.com/ivyascorp-net/nagging-nancy/archive/refs/tags/v1.0.0.tar.gz"
  sha256 "INSERT_TARBALL_SHA256_HERE"
  license "MIT"
  head "https://github.com/ivyascorp-net/nagging-nancy.git", branch: "main"

  depends_on "go" => :build
  
  # Platform-specific notification dependencies
  on_macos do
    depends_on "terminal-notifier" => :recommended
  end
  
  on_linux do
    depends_on "libnotify" => :recommended
  end

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
    caveats_text = <<~EOS
      nagging-nancy stores data in ~/.config/nancy/
      
      To start the notification daemon:
        nancy daemon start
      
      To add a reminder:
        nancy add "Call mom"
        nancy add "Meeting tomorrow at 2pm" --priority high
      
      To start the interactive TUI:
        nancy
        
      For optimal notification support:
    EOS
    
    case OS.mac?
    when true
      caveats_text += <<~EOS
        • terminal-notifier is recommended for enhanced macOS notifications
        • Install with: brew install terminal-notifier
        • Built-in AppleScript notifications are used as fallback
      EOS
    when false
      if OS.linux?
        caveats_text += <<~EOS
        • libnotify (notify-send) is recommended for desktop notifications
        • Install with your package manager (apt install libnotify-bin, etc.)
        • dunstify is also supported as an alternative
        EOS
      end
    end
    
    caveats_text += <<~EOS
      
      To test notifications:
        nancy test notification
    EOS
    
    caveats_text
  end

  test do
    assert_match version.to_s, shell_output("#{bin}/nancy version")
    
    # Test basic functionality
    system bin/"nancy", "add", "test reminder"
    output = shell_output("#{bin}/nancy list")
    assert_match "test reminder", output
    
    # Test daemon commands (without actually starting)
    assert_match "daemon is not running", shell_output("#{bin}/nancy daemon status", 1)
    
    # Test notification system (should not fail even without GUI)
    system bin/"nancy", "test", "notification"
  end
end