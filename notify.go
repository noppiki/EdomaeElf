package main

import (
	"fmt"
	"os/exec"
)

// SendMacNotification sends a notification on macOS
func SendMacNotification(title, subtitle, message string) error {
	script := fmt.Sprintf(`display notification "%s" with title "%s" subtitle "%s" sound name "Glass"`, 
		message, title, subtitle)
	
	cmd := exec.Command("osascript", "-e", script)
	return cmd.Run()
}

// SendMacNotificationWithSound sends a notification with custom sound
func SendMacNotificationWithSound(title, subtitle, message, sound string) error {
	script := fmt.Sprintf(`display notification "%s" with title "%s" subtitle "%s" sound name "%s"`, 
		message, title, subtitle, sound)
	
	cmd := exec.Command("osascript", "-e", script)
	return cmd.Run()
}

// Example usage:
func ExampleNotification() {
	// Simple notification
	err := SendMacNotification("Claude Code", "作業完了", "すべてのタスクが完了しました！")
	if err != nil {
		fmt.Printf("通知送信エラー: %v\n", err)
	}
	
	// Notification with custom sound
	// Available sounds: "Basso", "Blow", "Bottle", "Frog", "Funk", "Glass", "Hero", "Morse", "Ping", "Pop", "Purr", "Sosumi", "Submarine", "Tink"
	err = SendMacNotificationWithSound("Claude Code", "エラー発生", "ビルドに失敗しました", "Basso")
	if err != nil {
		fmt.Printf("通知送信エラー: %v\n", err)
	}
}