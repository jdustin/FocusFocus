// Copyright 2016 James Dustin. All rights reserved.
// license that can be found in the LICENSE file.

// A simple tool to help focus on getting stuff done and not
// procrastinating looking at sites that waste too much time.
//
package focusfocus

import (
	//"errors"
	"bufio"
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"
)

const StartMarker = "FocusFocusStart"
const EndMarker = "FocusFocusStop"
const CommentMarker = "#"
const DefaultRedirect = "127.0.0.1"

type Hostfile struct {
	lines []string
}

func (hf Hostfile) Lines() []string {
	return hf.lines
}

func GetLineEnding() string {
	if runtime.GOOS == "windows" {
		return "\r\n"
	}
	return "\n"
}

func GetHostPath() string {
	if runtime.GOOS == "windows" {
		return "${SystemRoot}/System32/drivers/etc/hosts"
	}
	return "/etc/hosts"
}

func CheckComment(line string, commentID string) bool {
	trimmed := strings.TrimSpace(line)
	return strings.HasPrefix(trimmed, CommentMarker)
}

func IsFocusStart(line string) bool {
	if CheckComment(line, CommentMarker) {
		return strings.Contains(line, StartMarker)
	}
	return false
}

func IsFocusEnd(line string) bool {
	if CheckComment(line, CommentMarker) {
		return strings.Contains(line, EndMarker)
	}
	return false
}

func (hf *Hostfile) RemoveFocusLines() (bool, error) {
	startFound := false
	endFound := false
	startPos := 0
	endPos := 0
	linecount := len(hf.lines)

	// try to find the start of the focusfocus block of lines
	for i := 0; i < linecount; i++ {
		if IsFocusStart(hf.lines[i]) {
			startFound = true
			startPos = i
			break
		}
	}
	if !startFound {
		//nothing to do
		return false, nil
	}
	// now try to find the end of the block
	for i := startPos; i < linecount; i++ {
		if IsFocusEnd(hf.lines[i]) {
			endFound = true
			endPos = i
			break
		}
	}
	// problem. the start line was found but not the end marker.
	// don't proceed and raise an error
	if !endFound {
		return false, fmt.Errorf("Invalid file contents. Start of focusfocus lines found (at line %d but the end marker was not found. The hosts file may need manually edited")
	}
	newlines := make([]string, len(hf.lines)-(endPos-startPos+1))
	j := 0
	for i, line := range hf.lines {
		if i >= startPos && i <= endPos {
			continue
		}
		newlines[j] = line
		j++
	}
	hf.lines = newlines
	return true, nil
}

// check that the file is writeable for later saving
func WriteAccessAvailable(filename string) error {
	filew, err := os.OpenFile(filename, os.O_WRONLY, 0660)
	if err == nil {
		filew.Close()
	}
	return err
}

func (hf *Hostfile) AddFocusLines(newHosts []string, redirect string) (bool, error) {
	newLineCount := len(newHosts)
	if newLineCount < 1 {
		//nothing to do
		return false, nil
	}
	existingCount := len(hf.lines)
	totalLineCount := existingCount + newLineCount + 2
	focusStrs := make([]string, totalLineCount)

	t := time.Now()
	startStr := fmt.Sprintf("%s %s lines:%d Added:%s ",
		CommentMarker, StartMarker, newLineCount, t.Format(time.RFC1123))
	endStr := fmt.Sprintf("%s %s", CommentMarker, EndMarker)

	src := 0
	dest := 0
	focusStrs[dest] = startStr
	dest++
	for src = 0; src < newLineCount; src++ {
		focusStrs[dest] = fmt.Sprintf("%s %s", redirect, newHosts[src])
		dest++
	}
	focusStrs[dest] = endStr
	dest++
	for src = 0; src < existingCount; src++ {
		focusStrs[dest] = hf.lines[src]
		dest++
	}
	hf.lines = focusStrs
	return true, nil
}

//todo: change to a stream
func (hf *Hostfile) WriteFile(filename string) error {
	if err := WriteAccessAvailable(filename); err != nil {
		return fmt.Errorf("Failed write access. Err:%s", err)
	}
	//file, err := os.OpenFile(filename, os.O_WRONLY, 0660)
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("Failed open. Err:%s", err)
	}
	defer file.Close()

	w := bufio.NewWriter(file)
	eol := GetLineEnding()
	for _, line := range hf.lines {
		fmt.Fprintf(w, "%s%s", line, eol)
	}
	err = w.Flush()
	if err != nil {
		return fmt.Errorf("Failed to Flush write buffer. Err:%s", err)
	}
	return nil
}

//todo: change to a stream
func (hf *Hostfile) ReadFile(filename string) error {
	if err := WriteAccessAvailable(filename); err != nil {
		return err
	}

	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		hf.lines = append(hf.lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return err
	}
	return nil
}
