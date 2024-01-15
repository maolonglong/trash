// Copyright 2024 Shaolong Chen. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package main

import (
	"os"
	"testing"
)

// Generated by GitHub Copilot
func TestRm(t *testing.T) {
	// Create a temporary file for testing
	file, err := os.CreateTemp("", "trash_test_*")
	if err != nil {
		t.Fatalf("Failed to create temporary file: %v", err)
	}
	defer os.Remove(file.Name())

	// Call the rm function to remove the temporary file
	err = rm(file.Name())
	if err != nil {
		t.Fatalf("Failed to remove file: %v", err)
	}

	// Check if the file still exists
	_, err = os.Stat(file.Name())
	if !os.IsNotExist(err) {
		t.Errorf("File was not removed: %s", file.Name())
	}
}