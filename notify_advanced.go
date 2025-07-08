package main

import (
	"fmt"
	"os/exec"
	"runtime"
)

// NotifyMac sends various types of notifications on macOS
type NotifyMac struct{}

// Alert shows a popup dialog (always visible)
func (n NotifyMac) Alert(title, message string) error {
	script := fmt.Sprintf(`display alert "%s" message "%s" as informational`, title, message)
	cmd := exec.Command("osascript", "-e", script)
	return cmd.Run()
}

// Sound plays a system sound
func (n NotifyMac) Sound(soundName string) error {
	// Available sounds: Glass, Basso, Blow, Bottle, Frog, Funk, Hero, Morse, Ping, Pop, Purr, Sosumi, Submarine, Tink
	soundPath := fmt.Sprintf("/System/Library/Sounds/%s.aiff", soundName)
	cmd := exec.Command("afplay", soundPath)
	return cmd.Run()
}

// Banner tries to show a notification banner
func (n NotifyMac) Banner(title, subtitle, message string) error {
	// First try notification center
	script := fmt.Sprintf(`display notification "%s" with title "%s" subtitle "%s"`, message, title, subtitle)
	cmd := exec.Command("osascript", "-e", script)
	return cmd.Run()
}

// SpeakText uses text-to-speech
func (n NotifyMac) SpeakText(text string) error {
	cmd := exec.Command("say", text)
	return cmd.Run()
}

// ShowNotification tries multiple methods to ensure user sees the notification
func ShowNotification(title, message string) {
	if runtime.GOOS != "darwin" {
		fmt.Println("この機能はmacOSでのみ動作します")
		return
	}

	n := NotifyMac{}
	
	// Method 1: Try notification banner
	err := n.Banner(title, "", message)
	if err != nil {
		fmt.Printf("通知バナーエラー: %v\n", err)
	}
	
	// Method 2: Play sound
	_ = n.Sound("Glass")
	
	// Method 3: If critical, show alert
	// n.Alert(title, message) // This will block until user clicks OK
}

// Example usage
func ExampleAdvancedNotification() {
	n := NotifyMac{}
	
	// Play completion sound
	n.Sound("Glass")
	
	// Show alert dialog (always visible)
	n.Alert("Claude Code", "すべての作業が完了しました！")
	
	// Speak completion message
	n.SpeakText("作業が完了しました")
}