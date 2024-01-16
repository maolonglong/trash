// Copyright 2024 Shaolong Chen. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/pflag"
)

var (
	recursive = pflag.BoolP(
		"recursive",
		"r",
		false,
		"remove directories and their contents recursively",
	)
	force = pflag.BoolP(
		"force",
		"f",
		false,
		"ignore nonexistent files and arguments, never prompt",
	)
	protectedPaths = pflag.StringSliceP(
		"protect",
		"p",
		nil,
		"protected `paths`, separated by comma",
	)
	help = pflag.BoolP("help", "h", false, "display this help and exit")
)

var exitCode = 0

func rm(name string, m map[string]struct{}) error {
	path, err := filepath.Abs(name)
	if err != nil {
		return err
	}

	_, protected := m[path]
	if protected {
		log.Printf("skipping %q", name)
		return nil
	}

	stat, err := os.Stat(name)
	if err != nil {
		if !*force {
			log.Printf("cannot remove %q: %v", name, err)
			exitCode = 1
		}
		return nil
	}
	if stat.IsDir() && !*recursive {
		log.Printf("cannot remove %q: is a directory", name)
		exitCode = 1
		return nil
	}

	script := fmt.Sprintf(`tell application "Finder"
  move POSIX file %q to the trash
end tell
`, path)

	cmd := exec.Command("osascript", "-e", script)
	cmd.Stdout = nil
	var errbuf bytes.Buffer
	cmd.Stderr = &errbuf
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("%w: %s", err, errbuf.String())
	}
	return nil
}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [OPTION]... [FILE]...\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "Move FILE(s) to the trash.\n\n")
	pflag.PrintDefaults()
	fmt.Fprintf(os.Stderr, `
By default, trash does not remove directories.  Use the --recursive (-r)
option to remove each listed directory, too, along with all of its contents.

To remove a file whose name starts with a '-', for example '-foo',
use one of these commands:
  %[1]s -- -foo

  %[1]s ./-foo
`, os.Args[0])
}

func checkPrivileges() {
	if os.Getuid() == 0 {
		fmt.Fprintf(os.Stderr, "This program must not be run as root!!!\n")
		os.Exit(1)
	}
}

func main() {
	checkPrivileges()

	log.SetPrefix("trash: ")
	log.SetFlags(0)
	pflag.Parse()

	if *help {
		usage()
		return
	}

	if pflag.NArg() == 0 {
		fmt.Fprintf(os.Stderr, "rm: missing operand\n")
		fmt.Fprintf(os.Stderr, "Try '%s --help' for more information.\n", os.Args[0])
		os.Exit(2)
	}

	m := make(map[string]struct{}, len(*protectedPaths)+1)
	for _, path := range *protectedPaths {
		m[path] = struct{}{}
	}
	home, _ := os.UserHomeDir()
	if home != "" {
		m[home] = struct{}{}
	}

	for _, name := range pflag.Args() {
		if err := rm(name, m); err != nil {
			log.Fatal(err)
		}
	}
	os.Exit(exitCode)
}
