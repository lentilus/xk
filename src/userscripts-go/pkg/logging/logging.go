package logging

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

const LogMaxLines = 100 // Set the maximum number of log lines to keep

// TruncateLogFile truncates the log file to the last N lines.
func TruncateLogFile(logFilePath string, maxLines int) error {
	file, err := os.Open(logFilePath)
	if err != nil {
		return fmt.Errorf("error opening log file for truncation: %v", err)
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading log file: %v", err)
	}

	// Keep only the last `maxLines` lines
	if len(lines) > maxLines {
		lines = lines[len(lines)-maxLines:]
	}

	// Write the truncated lines back to the log file
	err = ioutil.WriteFile(logFilePath, []byte(strings.Join(lines, "\n")+"\n"), 0644)
	if err != nil {
		return fmt.Errorf("error writing to log file: %v", err)
	}
	return nil
}

// SetLogOutput sets the log output to the specified log file.
func SetLogOutput(logFilePath string) (*os.File, error) {
	logFile, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("error opening log file: %v", err)
	}
	log.SetOutput(logFile)
	return logFile, nil
}

// PanicWithLog logs the message and then panics.
func PanicWithLog(msg string, args ...interface{}) {
	log.Printf(msg, args...)
	panic(fmt.Sprintf(msg, args...))
}
