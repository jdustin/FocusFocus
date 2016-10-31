// Copyright 2016 James Dustin. All rights reserved.
// license that can be found in the LICENSE file.

// A simple tool to help focus on getting stuff done and not
// procrastinating looking at sites that waste too much time.
//
// Usage: (with sudo or a windows cmd with admin access)
//
//    To Focus and waste less time:
//
//       focuscmd focus sitestolimit.txt
//
//    To Relax and waste more time
//
//       focuscmd relax
package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/jdustin/focusfocus"
	"os"
	"path/filepath"
	"strings"
)

func Usage() {
	fmt.Fprintf(os.Stderr,
		"Usage:\n %s focus file - Apply restrictions\n %s relax - Remove restrictions\n", os.Args[0], os.Args[0])
}
func LoadFocusList(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func main() {
	flag.Parse()
	flag.Usage = Usage

	if len(flag.Args()) < 1 {
		fmt.Fprintln(os.Stderr, "No action specified")
		Usage()
		return
	}
	action := strings.ToLower(flag.Arg(0))
	if action != "relax" && action != "focus" {
		fmt.Fprintln(os.Stderr, "Invalid action (", action, ")")
		Usage()
		return
	}
	if action == "focus" && len(flag.Args()) <= 1 {
		fmt.Fprintln(os.Stderr, "Missing focus host list")
		Usage()
		return
	}
	filePath := focusfocus.GetHostPath()
	fullFilePath := os.ExpandEnv(filepath.FromSlash(filePath))
	hf := focusfocus.Hostfile{}

	if err := hf.ReadFile(fullFilePath); err != nil {
		fmt.Fprintln(os.Stderr, "Unable to proceed. Error:", err)
		os.Exit(1)
	}
	// remove any lines previously added by focusfocus
	if action == "relax" {
		Relax(&hf, fullFilePath)
		return
	}
	// add new lines to the host file (removing any previous focusfocus lines first)
	if action == "focus" {
		param := flag.Arg(1)
		hosts, err := LoadFocusList(param)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Unable to load focus list (", param, ") Err:", err)
			Usage()
			return
		}
		if !Relax(&hf, fullFilePath) {
			return
		}
		if !Focus(&hf, fullFilePath, hosts) {
			return
		}
	}
}

// remove lines from a host file that have been added previously by focusfocus
func Relax(hf *focusfocus.Hostfile, file string) bool {
	removed, err := hf.RemoveFocusLines()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Unable to remove focus host entries. Error:", err)
		return false
	}
	if removed {
		err = hf.WriteFile(file)
		if err != nil {
			fmt.Fprintln(os.Stderr,
				"Unable to write host file while trying to remove focus entries. Error:", err)
			return false
		}
		fmt.Fprintln(os.Stdout,
			"Focus lines removed from Hosts file (", file, ")")
	}
	return true
}

// add lines to a host file
func Focus(hf *focusfocus.Hostfile, file string, hosts []string) bool {
	added, err := hf.AddFocusLines(hosts, focusfocus.DefaultRedirect)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Unable to add focus host entries. Error:", err)
		return false
	}
	if added {
		err = hf.WriteFile(file)
		if err != nil {
			fmt.Fprintln(os.Stderr,
				"Unable to write host file while trying to add focus entries. Error:", err)
			return false
		}
		fmt.Fprintln(os.Stdout,
			"Focus list added to Hosts file (", file, ")")
	}
	return true
}
