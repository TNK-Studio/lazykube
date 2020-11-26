package utils

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
	"os/user"
	"runtime"
	"strings"
)

// FilePath replace ~ -> $HOME
func FilePath(path string) string {
	path = strings.Replace(path, "~", os.Getenv("HOME"), 1)
	return path
}

// FileExited check file exited
func FileExited(path string) bool {
	info, err := os.Stat(FilePath(path))
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// IsDirector IsDir
func IsDirector(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}

// Home returns the home directory for the executing user.
//
// This uses an OS-specific method for discovering the home directory.
// An error is returned if a home directory cannot be detected.
func Home() (string, error) {
	current, err := user.Current()
	if err == nil {
		return current.HomeDir, nil
	}

	// cross compile support

	if "windows" == runtime.GOOS {
		return HomeWindows()
	}

	// Unix-like system, so just assume Unix
	return HomeUnix()
}

func HomeUnix() (string, error) {
	// First prefer the HOME environmental variable
	if home := os.Getenv("HOME"); home != "" {
		return home, nil
	}

	// If that fails, try the shell
	var stdout bytes.Buffer
	cmd := exec.Command("sh", "-c", "eval echo ~$USER")
	cmd.Stdout = &stdout
	if err := cmd.Run(); err != nil {
		return "", err
	}

	result := strings.TrimSpace(stdout.String())
	if result == "" {
		return "", errors.New("blank output when reading home directory")
	}

	return result, nil
}

func HomeWindows() (string, error) {
	drive := os.Getenv("HOMEDRIVE")
	path := os.Getenv("HOMEPATH")
	home := drive + path
	if drive == "" || path == "" {
		home = os.Getenv("USERPROFILE")
	}
	if home == "" {
		return "", errors.New("HOMEDRIVE, HOMEPATH, and USERPROFILE are blank")
	}

	return home, nil
}
