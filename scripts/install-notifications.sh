#!/bin/bash

# Nagging Nancy - Notification Dependencies Installer
# This script installs platform-specific notification dependencies

set -e

echo "🔔 Setting up notification dependencies for Nagging Nancy..."
echo

# Detect platform and install appropriate dependencies
case "$(uname -s)" in
    Linux*)
        echo "Detected Linux platform"
        
        # Detect package manager and install libnotify
        if command -v apt &> /dev/null; then
            echo "Using apt package manager..."
            if ! command -v notify-send &> /dev/null; then
                echo "Installing libnotify-bin..."
                sudo apt update && sudo apt install -y libnotify-bin
            else
                echo "✅ notify-send already installed"
            fi
            
        elif command -v dnf &> /dev/null; then
            echo "Using dnf package manager..."
            if ! command -v notify-send &> /dev/null; then
                echo "Installing libnotify..."
                sudo dnf install -y libnotify
            else
                echo "✅ notify-send already installed"
            fi
            
        elif command -v yum &> /dev/null; then
            echo "Using yum package manager..."
            if ! command -v notify-send &> /dev/null; then
                echo "Installing libnotify..."
                sudo yum install -y libnotify
            else
                echo "✅ notify-send already installed"
            fi
            
        elif command -v pacman &> /dev/null; then
            echo "Using pacman package manager..."
            if ! command -v notify-send &> /dev/null; then
                echo "Installing libnotify..."
                sudo pacman -S --noconfirm libnotify
            else
                echo "✅ notify-send already installed"
            fi
            
        elif command -v zypper &> /dev/null; then
            echo "Using zypper package manager..."
            if ! command -v notify-send &> /dev/null; then
                echo "Installing libnotify-tools..."
                sudo zypper install -y libnotify-tools
            else
                echo "✅ notify-send already installed"
            fi
            
        else
            echo "❌ Unsupported package manager. Please install libnotify manually:"
            echo "   - Ubuntu/Debian: sudo apt install libnotify-bin"
            echo "   - RHEL/CentOS: sudo yum install libnotify"
            echo "   - Fedora: sudo dnf install libnotify"
            echo "   - Arch: sudo pacman -S libnotify"
            echo "   - openSUSE: sudo zypper install libnotify-tools"
            exit 1
        fi
        ;;
        
    Darwin*)
        echo "Detected macOS platform"
        
        # Check if Homebrew is installed
        if command -v brew &> /dev/null; then
            echo "Using Homebrew package manager..."
            if ! command -v terminal-notifier &> /dev/null; then
                echo "Installing terminal-notifier..."
                brew install terminal-notifier
            else
                echo "✅ terminal-notifier already installed"
            fi
        else
            echo "⚠️  Homebrew not found. terminal-notifier is recommended but optional."
            echo "   macOS has built-in notification support via AppleScript."
            echo "   To install Homebrew: /bin/bash -c \"\$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)\""
            echo "   Then run: brew install terminal-notifier"
        fi
        ;;
        
    CYGWIN*|MINGW*|MSYS*)
        echo "Detected Windows platform"
        echo "✅ Windows has built-in notification support via PowerShell"
        echo "   No additional dependencies required."
        ;;
        
    *)
        echo "❌ Unsupported platform: $(uname -s)"
        echo "   Nagging Nancy will fall back to terminal bell notifications."
        exit 1
        ;;
esac

echo
echo "🎉 Notification setup complete!"
echo
echo "To test notifications, run:"
echo "   nancy test notification"
echo
echo "To view available notification methods:"
echo "   nancy test notification --verbose"