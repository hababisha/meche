package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/araddon/dateparse"
)

func getExifTimestamp(imagePath string) (string, error) {
	cmd := exec.Command("exiftool", "-DateTimeOriginal", imagePath)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", err
	}

	lines := strings.Split(out.String(), ": ")
	if len(lines) < 2 {
		return "", fmt.Errorf("could not parse timestamp")
	}
	return strings.TrimSpace(lines[1]), nil
}

func humanReadableTime(timestamp string) (string, error) {
	parsedTime, err := dateparse.ParseLocal(timestamp)
	if err != nil {
		return "", err
	}

	duration := time.Since(parsedTime)

	if duration < time.Minute {
		return "Just now", nil
	} else if duration < time.Hour {
		return fmt.Sprintf("%d minutes ago", int(duration.Minutes())), nil
	} else if duration < 24*time.Hour {
		return fmt.Sprintf("%d hours ago", int(duration.Hours())), nil
	} else if duration < 30*24*time.Hour {
		return fmt.Sprintf("%d days ago", int(duration.Hours()/24)), nil
	} else if duration < 12*30*24*time.Hour {
		return fmt.Sprintf("%d months ago", int(duration.Hours()/(24*30))), nil
	} else {
		return fmt.Sprintf("%d years ago", int(duration.Hours()/(24*365))), nil
	}
}

func main() {
	humanReadable := flag.Bool("h", false, "Display the timestamp in a human-readable format")
	flag.Parse()

	if flag.NArg() < 1 {
		fmt.Println("Usage: meche [-h] <image-file>")
		os.Exit(1)
	}

	imagePath := flag.Arg(0)

	timestamp, err := getExifTimestamp(imagePath)
	if err != nil {
		fmt.Println("Error extracting timestamp:", err)
		os.Exit(1)
	}

	if *humanReadable {
		humanTime, err := humanReadableTime(timestamp)
		if err != nil {
			fmt.Println("Error formatting timestamp:", err)
			os.Exit(1)
		}
		fmt.Println(humanTime)
	} else {
		fmt.Println("Photo taken on:", timestamp)
	}
}
