package disk

import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

type Status int

const (
	Good Status = iota
	Warning
	Error
)

type Disk struct {
	Device         string  `json:"device"`
	Used           float64 `json:"used"`
	UsedPercentage float64 `json:"usedPercentage"`
	Size           float64 `json:"size"`
	AvailableKB    float64 `json:"availableKB"`
	Mount          string  `json:"mount"`
	Status         Status  `json:"status"`
}

func LoadDisks(excludePatterns ...regexp.Regexp) ([]Disk, error) {
	// Use -kP to ensure POSIX format and sizes in kilobytes
	cmd := exec.Command("df", "-kP")
	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("failed to run df: %w", err)
	}

	var disks []Disk
	scanner := bufio.NewScanner(&out)
	firstLine := true

	for scanner.Scan() {
		line := scanner.Text()

		if firstLine {
			firstLine = false // skip header
			continue
		}

		disk, err := parseLine(line)

		if err != nil {
			fmt.Printf("Error parsing line: '%s': %s", line, err)
		} else {
			if isRelevant(disk, excludePatterns) {
				disks = append(disks, disk)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading df output: %w", err)
	}

	return disks, nil
}

func isRelevant(disk Disk, excludePatterns []regexp.Regexp) bool {
	for _, pattern := range excludePatterns {
		if pattern.MatchString(disk.Mount) {
			return false
		}
	}

	return true
}

func parseLine(line string) (Disk, error) {
	fields := strings.Fields(line)
	if len(fields) < 6 {
		// skip malformed lines
		return Disk{}, fmt.Errorf("Expected 6 fields, got %d", len(fields))
	}

	sizeKB, _ := strconv.ParseFloat(fields[1], 64)
	usedKB, _ := strconv.ParseFloat(fields[2], 64)
	availKB, _ := strconv.ParseFloat(fields[3], 64)
	usedPercentStr := strings.TrimRight(fields[4], "%")
	usedPercent, _ := strconv.ParseFloat(usedPercentStr, 64)

	status := Good
	if usedPercent > 90 {
		status = Error
	} else if usedPercent > 70 {
		status = Warning
	}

	disk := Disk{
		Device:         fields[0],
		Used:           usedKB,
		UsedPercentage: usedPercent,
		Size:           sizeKB,
		AvailableKB:    availKB,
		Mount:          fields[5],
		Status:         status,
	}

	return disk, nil
}
