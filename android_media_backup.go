package main

import (
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

//
// Simple script to download&remove photos from an Android phone which
// disconnects before you can get everything off (even using ADB).
// Last used on a Google Pixel 2 which had a bad USB-C port.
//

var MediaDirectory = flag.String("dir", "DCIM/Camera", "Folder to backup")

func main() {
	flag.Parse()

	// Step 1: is the device connected?
	b, err := exec.Command("adb", "devices").Output()
	if err != nil {
		log.Fatal(err)
	}

	out := string(bytes.TrimSpace(b))
	fmt.Printf("%q\n", out)

	parts := strings.Split(out, "\n")

	if len(parts) == 1 {
		log.Fatal("No devices found, connect phone")
	}

	if len(parts) != 2 {
		log.Fatal("Too many devices listed, only connect one at a time")
	}

	parts = strings.Split(parts[1], "\t")
	deviceID := parts[0]
	fmt.Printf("Downloading from %s\n", deviceID)

	// Step 2: what images have already been brought over?
	if ok, _ := exists(*MediaDirectory); !ok {
		err = os.MkdirAll(*MediaDirectory, 0755)
		if err != nil {
			log.Fatal(err)
		}
	}

	ComputerFiles, err := FilePathWalkDir(*MediaDirectory)
	fmt.Printf("Computer Files:\n%s\n", strings.Join(ComputerFiles, "\n"))

	if err != nil {
		log.Fatal("Error looking at files", err)
	}

	// Step 3: What files are on the device?
	b, err = exec.Command("adb", "shell", "ls", "/sdcard/"+*MediaDirectory).Output()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Phone Files:\n%s\n", string(b))

	PhoneFiles := strings.Split(string(b), "\n")

	// Step 4: Remove files on the phone we have saved locally
	if len(PhoneFiles) > 0 {
		for _, pf := range PhoneFiles {
			var found bool
			for _, cf := range ComputerFiles {
				if cf == pf {
					found = true
					break
				}
			}

			if found && len(pf) > 4 {
				// fmt.Printf("REMOVE: %s\n", pf)
				b, err = exec.Command("adb", "shell", "rm", filepath.Join("/sdcard/", *MediaDirectory, pf)).Output()
				if err != nil {
					log.Fatal(err)
				}
				fmt.Printf("REMOVE %s -> %q\n", pf, b)
			}
		}
	}

	// Step 4: Download new ones: adb pull -a /sdcard/DCIM/Camera/ ./
	b, err = exec.Command("adb", "pull", "-a", "/sdcard/"+*MediaDirectory, *MediaDirectory).Output()
	if err != nil {
		log.Fatal(err)
	}

	out = string(bytes.TrimSpace(b))
	fmt.Println(out)
	// Step 5: if at first you don't succeed, fail, fail again!

}

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

func FilePathWalkDir(root string) ([]string, error) {
	var files []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			files = append(files, filepath.Base(path))
		}
		return nil
	})
	return files, err
}

// 2019 http://davidpennington.me - Released under the MIT License
