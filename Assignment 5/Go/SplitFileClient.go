/**
 * SplitFileClient.go
 * ID : 20194094
 * Name : Yongmin Yoo
 **/

package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"strings"
)

const SERVER_ADDRESS_1 = "localhost:44094"
const SERVER_ADDRESS_2 = "localhost:54094"

func main() {
	if len(os.Args) != 3 {
		fmt.Println("Usage: go run client.go <put|get> <filename>")
		return
	}

	cmd := os.Args[1]
	filename := os.Args[2]

	if cmd == "put" {
		filePart1, filePart2, err := splitFile(filename)
		if err != nil {
			os.Remove(filename + "-part1.tmp")
			os.Remove(filename + "-part2.tmp")
			exitError(fmt.Sprintf("Failed to split file - %s", err.Error()))
		}

		err = sendPart(filename, filePart1, SERVER_ADDRESS_1)
		os.Remove(filePart1)
		if err != nil {
			exitError(fmt.Sprintf("Failed to send Part 1 - %s", err.Error()))
		}

		err = sendPart(filename, filePart2, SERVER_ADDRESS_2)
		os.Remove(filePart2)
		if err != nil {
			exitError(fmt.Sprintf("Failed to send Part 2 - %s", err.Error()))
		}

		fmt.Println("File successfully split and sent to servers.")
	} else if cmd == "get" {
		filePart1, err := getPart(filename, SERVER_ADDRESS_1, 1)
		if err != nil {
			os.Remove(filename + "-part1.tmp")
			exitError(fmt.Sprintf("Failed to get Part 1 - %s", err.Error()))
		}

		filePart2, err := getPart(filename, SERVER_ADDRESS_2, 2)
		if err != nil {
			os.Remove(filename + "-part1.tmp")
			os.Remove(filename + "-part2.tmp")
			exitError(fmt.Sprintf("Failed to get Part 2 - %s", err.Error()))
		}

		outputFilename := getMergedFilename(filename)
		err = mergeFiles(filePart1, filePart2, outputFilename)
		os.Remove(filePart1)
		os.Remove(filePart2)
		if err != nil {
			exitError(fmt.Sprintf("Failed to merge files - %s", err.Error()))
		}

		fmt.Println("File successfully retrieved and merged:", outputFilename)
	} else {
		fmt.Println("Usage: go run SplitFileClient.go <put|get> <filename>")
	}
}

func sendPart(filename, partFilename, serverAddress string) error {
	serverConn, err := net.Dial("tcp4", serverAddress)
	if err != nil {
		return err
	}
	defer serverConn.Close()

	partFile, err := os.Open(partFilename)
	if err != nil {
		return err
	}
	defer partFile.Close()

	serverConn.Write([]byte(fmt.Sprintf("PUT:%s", filename)))
	response, err := bufio.NewReader(serverConn).ReadString('\n')
	if err != nil {
		return err
	}
	if strings.TrimSpace(response) != "READY" {
		return fmt.Errorf("Server not ready for file content")
	}

	fileBuffer := make([]byte, 1024)
	for {
		fileLength, err := partFile.Read(fileBuffer)
		if err != nil && err != io.EOF {
			return err
		}
		if fileLength == 0 {
			break
		}
		serverConn.Write(fileBuffer[:fileLength])
	}
	return nil
}

func getPart(filename, serverAddress string, partNum int) (string, error) {
	serverConn, err := net.Dial("tcp4", serverAddress)
	if err != nil {
		return "", err
	}
	defer serverConn.Close()

	partFilename := filename + fmt.Sprintf("-part%d.tmp", partNum)
	serverConn.Write([]byte(fmt.Sprintf("GET:%s", filename)))

	partFile, err := os.Create(partFilename)
	if err != nil {
		return "", err
	}
	defer partFile.Close()

	partFileBuffer := make([]byte, 1025)
	for {
		partFileLength, err := serverConn.Read(partFileBuffer)
		if err != nil && err != io.EOF {
			return "", err
		}
		if partFileLength == 0 {
			break
		}

		if strings.Contains(string(partFileBuffer[:partFileLength]), "EOF") {
			eofIndex := strings.Index(string(partFileBuffer[:partFileLength]), "EOF")
			if eofIndex > 0 {
				partFile.Write(partFileBuffer[1:eofIndex])
			}
			break
		} else if string(partFileBuffer[:partFileLength][:6]) == "NOFILE" {
			return "", errors.New("Server has an error with file")
		} else if string(partFileBuffer[:partFileLength][:5]) == "ERROR" {
			return "", errors.New("Server Returned an Error")
		} else {
			partFile.Write(partFileBuffer[1:partFileLength])
		}
	}

	return partFilename, nil
}

func splitFile(filename string) (string, string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", "", err
	}
	defer file.Close()

	partFile1, err := os.Create(filename + "-part1.tmp")
	if err != nil {
		return "", "", err
	}
	defer partFile1.Close()

	partFile2, err := os.Create(filename + "-part2.tmp")
	if err != nil {
		return "", "", err
	}
	defer partFile2.Close()

	fileBuffer := make([]byte, 1)
	isPart1 := true
	for {
		fileLength, err := file.Read(fileBuffer)
		if err != nil && err != io.EOF {
			return "", "", err
		}
		if fileLength == 0 {
			break
		}

		if isPart1 {
			partFile1.Write(fileBuffer[:fileLength])
		} else {
			partFile2.Write(fileBuffer[:fileLength])
		}
		isPart1 = !isPart1
	}

	return filename + "-part1.tmp", filename + "-part2.tmp", nil
}

func mergeFiles(part1, part2, outputFile string) error {
	outFile, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer outFile.Close()

	partFile1, err := os.Open(part1)
	if err != nil {
		return err
	}
	defer partFile1.Close()

	partFile2, err := os.Open(part2)
	if err != nil {
		return err
	}
	defer partFile2.Close()

	partFile1Buffer := make([]byte, 1)
	partFile2Buffer := make([]byte, 1)
	for {
		partFile1Length, err1 := partFile1.Read(partFile1Buffer)
		partFile2Length, err2 := partFile2.Read(partFile2Buffer)

		if partFile1Length == 0 && partFile2Length == 0 {
			break
		}

		if partFile1Length > 0 {
			outFile.Write(partFile1Buffer[:partFile1Length])
		}
		if partFile2Length > 0 {
			outFile.Write(partFile2Buffer[:partFile2Length])
		}

		if err1 != nil && err1 != io.EOF {
			return err1
		}
		if err2 != nil && err2 != io.EOF {
			return err2
		}
	}

	return nil
}

func getMergedFilename(filename string) string {
	ext := filepath.Ext(filename)
	base := strings.TrimSuffix(filename, ext)
	return base + "-merged" + ext
}

func exitError(msg string) {
	_, err := fmt.Fprintf(os.Stderr, "Error: %s\n", msg)
	if err != nil {
		fmt.Printf("Error: %s\n", msg)
	}
	os.Exit(1)
}
