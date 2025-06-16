package disk

import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

type Disk struct {
	Device         string  `json:"device"`
	Used           float64 `json:"used"`
	UsedPercentage float64 `json:"usedPercentage"`
	Size           float64 `json:"size"`
	AvailableKB    float64 `json:"availableKB"`
	Mount          string  `json:"mount"`
}

func LoadDisks() ([]Disk, error) {
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
			disks = append(disks, disk)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading df output: %w", err)
	}

	return disks, nil
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

	disk := Disk{
		Device:         fields[0],
		Used:           usedKB,
		UsedPercentage: usedPercent,
		Size:           sizeKB,
		AvailableKB:    availKB,
		Mount:          fields[5],
	}

	return disk, nil
}
