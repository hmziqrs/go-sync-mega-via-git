package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"time"
)

func main() {
	start := time.Now()

	agrsLen := len(os.Args)

	if agrsLen < 3 {
		fmt.Println("Pass source and destination dumbass")
		return
	}
	source := strings.TrimSpace(os.Args[1])
	destination := strings.TrimSpace(os.Args[2])

	if agrsLen == 4 {
		if strings.TrimSpace(os.Args[3]) != "init" {
			fmt.Println("What's wrong with you only \"init\" is allowed. go read the readme or a english book")
			return
		}

		files, err := ioutil.ReadDir(source)
		if err != nil {
			log.Fatal(err)
		}

		// fmt.Println(source)
		// fmt.Println(destination)

		syncList(files, source, destination, nil)

		// return
	}
	fmt.Println(time.Since(start))

}

func syncList(files []os.FileInfo, source string, destination string, baseIgnores []string) {
	var ignores []string

	if baseIgnores != nil {
		ignores = append([]string{}, baseIgnores...)
		println(ignores[0])
		println(source)
		println(destination)
	}

	gitIgnoreBytes, _ := ioutil.ReadFile(filepath.Join(source, ".gitignore"))
	if gitIgnoreBytes != nil {
		gitIgnoreRaw := strings.Split(string(gitIgnoreBytes), "\n")
		for _, str := range gitIgnoreRaw {
			trimed := strings.TrimSpace(str)
			if trimed != "" && !strings.Contains(trimed, "#") {
				ignores = append(ignores, slash(str))
			}
		}
	}

	for _, f := range files {
		path, _ := filepath.Abs(filepath.Join(source, f.Name()))
		if f.IsDir() {
			path = slash(path + "/")
		}
		ignore := shouldIgnore(path, ignores)

		if !ignore {
			newPath, _ := filepath.Abs(filepath.Join(destination, f.Name()))
			if f.IsDir() {
				_ = os.Mkdir(newPath, os.ModeDir)
				newFiles, _ := ioutil.ReadDir(path)
				println("NESTING PASSED")
				newSource := slash(source + "/" + f.Name())
				newDestination := slash(destination + "/" + f.Name())
				syncList(newFiles, newSource, newDestination, ignores)
			} else {
				content, _ := ioutil.ReadFile(path)
				ioutil.WriteFile(newPath, content, 0777)
			}
		}
	}
}

func shouldIgnore(path string, patterns []string) bool {
	for _, pattern := range patterns {
		match, _ := regexp.MatchString(pattern, path)
		contains := strings.Contains(path, pattern)

		if match || contains {
			return true
		}

	}
	return false
}

func slash(str string) string {
	if runtime.GOOS == "windows" {
		return strings.ReplaceAll(str, "/", "\\")
	}
	return strings.ReplaceAll(str, "\\", "/")

}
