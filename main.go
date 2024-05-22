package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"
)

const (
	windowsEOL           = "\r\n"
	windowsPathSeparator = "\\"
	linuxEOL             = "\n"
	eolDetection         = '\n'
)

var filePath string
var replaceFormat string
var replaceTo string

func main() {
	var eol string
	if runtime.GOOS == "windows" {
		fmt.Println("*** This service is running on Windows base os ***")
		eol = windowsEOL
	} else {
		fmt.Println("*** This service is running on Linux base os ***")
		eol = linuxEOL
	}
	displayInstruction()
	if err := prepareData(eol); err != nil {
		fmt.Println("error when preparing data", err)
		time.Sleep(5 * time.Second)
		os.Exit(1)
	}
	if err := changeFilename(); err != nil {
		fmt.Println("error when change filename", err)
		time.Sleep(5 * time.Second)
		os.Exit(1)
	}
	fmt.Println("complete to change all files")
	time.Sleep(5 * time.Second)
	os.Exit(0)
}

func displayInstruction() {
	fmt.Println("-----------------------------------HOW TO USE-----------------------------------")
	fmt.Println("INPUT file path in format> C:/pathA/pathB/... (for windows; \\ is support too)")
	fmt.Println("INPUT file replace format> -Copy")
	fmt.Println("INPUT file replace to    > -TEST (No input if you want to remove)")
	fmt.Println("For example, File-Copy.pdf -> File-TEST.pdf")
	fmt.Println("-------------------------------------WARNING-------------------------------------")
	fmt.Println("1.) If the file prefixes on the same path are similar, you may receive an error.")
	fmt.Println("---------------------------------------------------------------------------------")
	fmt.Println("")
}

func prepareData(eol string) (err error) {
	inputReader := bufio.NewReader(os.Stdin)
	fmt.Printf("please input file path to replace: ")
	filePath, err = inputReader.ReadString(eolDetection)
	if err != nil {
		return
	}
	fmt.Printf("please input file replace format : ")
	replaceFormat, err = inputReader.ReadString(eolDetection)
	if err != nil {
		return
	}
	fmt.Printf("please input file replace to     : ")
	replaceTo, err = inputReader.ReadString(eolDetection)
	if err != nil {
		return
	}
	if eol == windowsEOL {
		filePath = strings.TrimSuffix(filePath, windowsEOL)
		filePath = strings.Replace(filePath, windowsPathSeparator, "/", -1)
		replaceFormat = strings.TrimSuffix(replaceFormat, windowsEOL)
		replaceTo = strings.TrimSuffix(replaceTo, windowsEOL)
	} else {
		filePath = strings.TrimSuffix(filePath, linuxEOL)
		replaceFormat = strings.TrimSuffix(replaceFormat, linuxEOL)
		replaceTo = strings.TrimSuffix(replaceTo, linuxEOL)
	}
	return nil
}

func changeFilename() (err error) {
	if _, err = os.Stat(filePath); errors.Is(os.ErrNotExist, err) {
		fmt.Println("path is not exist", filePath)
		return err
	}
	dirEntries, err := os.ReadDir(filePath)
	if err != nil {
		return err
	}
	var count uint16 = 0
	errorFile := make([]string, 0)
	for _, entry := range dirEntries {
		if !entry.IsDir() && strings.Contains(entry.Name(), replaceFormat) {
			editedName := strings.Replace(entry.Name(), replaceFormat, replaceTo, -1)
			fmt.Println("change filename from", entry.Name(), "to", editedName)
			if err = os.Rename(filePath+"/"+entry.Name(), filePath+"/"+editedName); err != nil {
				fmt.Println("error occur:", err)
				errorFile = append(errorFile, entry.Name())
				continue
			}
			count++
		}
	}
	fmt.Println("total filename was changed:", count)
	fmt.Println("total filename get an error:", len(errorFile))
	for index, file := range errorFile {
		fmt.Printf("%d.) %v", index+1, file)
	}
	return nil
}
