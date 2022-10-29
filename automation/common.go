package main

import (
	"archive/zip"
	"crypto/sha256"
	"fmt"
	"golang.org/x/mod/semver"
	"io"
	"os"
	"os/exec"
)

func isNewVersionGreater(newVersion string, prevVersion string) bool {
	if semver.IsValid(newVersion) && semver.IsValid(prevVersion) {
		return semver.Compare(newVersion, prevVersion) > 0
	} else {
		return newVersion > prevVersion
	}
}

func runCmd(command string) error {
	// get `pwsh` executable path
	pwsh, err := exec.LookPath("pwsh")
	if err != nil {
		return fmt.Errorf("error finding pwsh: %v\n", err)
	}

	// create a new process
	cmd := &exec.Cmd{
		Path:   pwsh,
		Args:   []string{pwsh, "-Command", command},
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}

	// run the command
	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("error running command: %v\n", err)
	}
	return nil
}

func getFileSha256Hash(filePath string) (string, error) {
	// open the file
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("error opening file for hashing: %v\n", err)
	}
	defer file.Close()

	// get the hash of the file using golang crypto
	hasher := sha256.New()
	_, err = io.Copy(hasher, file)
	if err != nil {
		return "", fmt.Errorf("error hashing file: %v\n", err)
	}

	return fmt.Sprintf("%x", hasher.Sum(nil)), nil
}

func extractFileFromZip(zipPath, filePathInsideZip string) error {
	zipReader, err := zip.OpenReader(zipPath)
	if err != nil {
		return fmt.Errorf("error opening zip file: %v\n", err)
	}
	defer zipReader.Close()

	for _, file := range zipReader.File {
		if file.Name != filePathInsideZip {
			continue
		}

		// copy the file to the current directory using io.Copy
		rc, err := file.Open()
		if err != nil {
			return fmt.Errorf("error opening file inside zip: %v\n", err)
		}
		defer rc.Close()

		// create a new file
		newFile, err := os.Create(file.Name)
		if err != nil {
			return fmt.Errorf("error creating file for zip extraction: %v\n", err)
		}
		defer newFile.Close()

		// copy the file
		_, err = io.Copy(newFile, rc)
		if err != nil {
			return fmt.Errorf("error copying file from zip: %v\n", err)
		}
	}
	return nil
}
