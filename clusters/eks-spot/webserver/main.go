package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"syscall"
	"time"
)

var srcData = make([]byte, 10*1024*1024)

func init() {
	for i := 0; i < len(srcData); i++ {
		srcData[i] = byte(i % 256)
	}
}

func openFileDirect(filePath string) (*os.File, error) {
	file, err := os.OpenFile(filePath, syscall.O_DIRECT|os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}
	return file, nil
}

func readLargeFile(filePath string) ([]byte, error) {
	startTime := time.Now()
	file, err := openFileDirect(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	data := make([]byte, 10*1024*1024)
	_, err = file.Read(data)
	if err != nil {
		return nil, err
	}

	fmt.Printf("File read in %v\n", time.Since(startTime))
	return data, nil
}

func writeLargeFile(filePath string, data []byte) error {
	startTime := time.Now()
	file, err := openFileDirect(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(data)
	if err != nil {
		return err
	}

	fmt.Printf("File written in %v\n", time.Since(startTime))
	return nil
}

func testReadWriteHandler(w http.ResponseWriter, r *http.Request) {
	filePath := "largefile.txt"
	// Add timestamp
	// filePath = fmt.Sprintf("%s-%d", filePath, time.Now().Unix())

	// Copy 1MB
	src := make([]byte, 10*1024*1024)
	for i := 0; i < len(src); i++ {
		src[i] = byte(i % 256)
		// random calculation
		src[i] = src[i] + 1
		src[i] = src[i] * 2
		src[i] = src[i] / 2
		src[i] = src[i] - 1
	}

	// memory copy here
	content := make([]byte, len(src))
	startTime := time.Now()
	copy(src, content)
	fmt.Printf("Memory copy 10*1024*1024 using copy() took %v\n", time.Since(startTime))

	err := writeLargeFile(filePath, srcData[:1024])
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to write large file, %s", err.Error()), http.StatusInternalServerError)
		return
	}

	_, err = readLargeFile(filePath)
	if err != nil {
		http.Error(w, "Failed to read large file", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "Test completed successfully")
}

func startServer() {
	http.HandleFunc("/test", testReadWriteHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func main() {
	fmt.Println("Starting server on port 8080...")
	startServer()
}
