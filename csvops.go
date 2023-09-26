package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"strings"
)

func deleteCsvFile() {
	err := os.Remove("data.csv")
	if err != nil {
		log.Fatal(err)
	}
}

func openFile() *os.File {
	filePath := "data.csv"

	f, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		log.Fatal(err)
	}
	return f
}

// Dosya içeriğini satır satır alıp qaMappingArr içerisinde barındırıp döndürüyor
func getAllData() []map[string]string {
	f := openFile()
	scanner := bufio.NewScanner(f)
	qaMappingArr := []map[string]string{}
	for scanner.Scan() {
		data := scanner.Text()
		if data == "Question,Answer,TrueAnswerCount,WrongAnswerCount,AskCount" {
			continue
		}
		splittedData := strings.Split(data, ",")
		qaMapping := map[string]string{}
		qaMapping["Question"] = splittedData[0]
		qaMapping["Answer"] = splittedData[1]
		qaMapping["TrueAnswerCount"] = splittedData[2]
		qaMapping["WrongAnswerCount"] = splittedData[3]
		qaMapping["AskCount"] = splittedData[4]
		qaMappingArr = append(qaMappingArr, qaMapping)
	}
	f.Close()
	return qaMappingArr
}

func readFirstLineOfFile(f io.Reader) string {
	scanner := bufio.NewScanner(f)
	firstLine := ""
	for scanner.Scan() {
		firstLine = scanner.Text()
		break
	}
	return firstLine
}

func writeCsvHeaderToFile(f *os.File) {
	if _, err := f.WriteString("Question,Answer,TrueAnswerCount,WrongAnswerCount,AskCount\n"); err != nil {
		log.Fatal(err)
	}
}

func appendDataToFile(f *os.File, data map[string]string) {
	lineData := ""
	lineData += data["Question"] + ","
	lineData += data["Answer"] + ","
	lineData += data["TrueAnswerCount"] + ","
	lineData += data["WrongAnswerCount"] + ","
	lineData += data["AskCount"] + "\n"
	if _, err := f.WriteString(lineData); err != nil {
		log.Fatal(err)
	}
}

func saveData(data map[string]string) {
	f := openFile()
	firstLine := readFirstLineOfFile(f)

	if firstLine == "" {
		writeCsvHeaderToFile(f)
	}

	appendDataToFile(f, data)

	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
	f.Close()
}
