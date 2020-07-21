package main

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	help := hasArg("-h") || hasArg("--help")

	if help {
		printHelp()
		return
	}

	filePath := os.Args[1]
	targetFolder := os.Args[2]

	fmt.Printf("Unzipping %s to %s\r\n", filePath, targetFolder)

	err := unzip(filePath, targetFolder)
	if err != nil {
		fmt.Printf("ERROR: %s\r\n", err)
		os.Exit(1)
	}

	fmt.Printf("Successfully uncompressed all files.\r\n")
}

func hasArg(arg string) bool {
	for _, a := range os.Args {
		if a == arg {
			return true
		}
	}

	return false
}

func unzip(filePath string, targetFolder string) error {
	reader, err := zip.OpenReader(filePath)
	if err != nil {
		return err
	}
	defer reader.Close()

	os.MkdirAll(targetFolder, os.ModePerm)

	for _, file := range reader.File {
		targetFilePath := fmt.Sprintf("%s%s%s", targetFolder, string(os.PathSeparator), file.Name)

		fmt.Printf("Unzip compressed file %s to %s\r\n", file.Name, targetFilePath)

		// Check for ZipSlip. More Info: http://bit.ly/2MsjAWE
		requiredPrefix := filepath.Clean(targetFolder) + string(os.PathSeparator)
		if !strings.HasPrefix(targetFilePath, requiredPrefix) {
			return fmt.Errorf("%s: illegal file path due to ZipSlip, see http://bit.ly/2MsjAWE, required prefix missing %s", targetFilePath, requiredPrefix)
		}

		if file.FileInfo().IsDir() {
			// Make Folder
			os.MkdirAll(targetFilePath, os.ModePerm)
			continue
		}

		outFile, err := os.OpenFile(targetFilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return err
		}

		inputReader, err := file.Open()
		if err != nil {
			return err
		}

		_, err = io.Copy(outFile, inputReader)

		// Close the file without defer to close before next iteration of loop
		outFile.Close()
		inputReader.Close()

		if err != nil {
			return err
		}
	}

	return nil
}

func printHelp() {
	fmt.Printf("Unzip a file into a folder.\r\n")
	fmt.Printf("Usage: gounzip <archive_file> <target_folder>\r\n")
}
