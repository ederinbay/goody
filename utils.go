package goody

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
)

func ReadFileToSlice(sourcePath string, separator string) []string {
	list := make([]string, 0)
	file, err := os.Open(sourcePath)
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Println(err)
		}
	}(file)
	if err != nil {
		log.Println(err)
	}
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		list = append(list, scanner.Text())
	}
	if err = scanner.Err(); err != nil {
		log.Println(err)
	}
	return list
}

func SliceToTxtFile(filePath []string, targetFilePath string) {
	var wg sync.WaitGroup
	var sem = make(chan int, 1000000)

	CreateDocumentIfNotExists(targetFilePath)
	file, err := os.OpenFile(targetFilePath, os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("failed opening file: %s", err)
	}
	defer file.Close()

	for i, a := range filePath {
		sem <- 1
		wg.Add(1)
		i := i
		go func(a string) {
			defer wg.Done()
			AppendFile(fmt.Sprintf(" %s\n", a), *file)
			_ = fmt.Sprintf("%d", i)
			<-sem
		}(a)
	}
	log.Println("Total number of lines: ", len(filePath))
	wg.Wait()
}
func CreateDocumentIfNotExists(path string) {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		log.Println("File does not exist. Creating one.")
		file, err1 := os.Create(path)
		if err1 != nil {
			log.Printf("File couldn't be created. %s", err1)
		}
		defer file.Close()
	}
}
func AppendFile(str string, file os.File) {
	trimmed := strings.Trim(str, " ")
	_, err := file.WriteString(trimmed)
	if err != nil {
		log.Fatalf("Failed writing to file: %s", err)
	}
}
