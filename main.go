package main

import (
	"bufio"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
)

var cFiles = make(chan string)
var dir string // directory to search
var args []string

var ignore []string

// isText determines whether this is the kind of file we
// want to open in vim.
func isText(filename string) bool {
	file, err := os.Open(filename)
	if err != nil {
		log.Println("Unable to check type of", filename, ":", err)
		return false
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Scan()
	fileType := http.DetectContentType(scanner.Bytes())
	if fileType[:4] == "text" {
		return true
	}
	return false
}

// isDir accepts a string (file path) and returns
// a boolean which indicates if the path is
// a valid directory.
func isDir(path string) bool {
	stat, err := os.Stat(path)
	if err != nil {
		log.Println("Error checking whether", path, "was a directory:", err)
		return false
	}
	return stat.IsDir()
}

func init() {
	// Set up command-line flags.
	flag.StringVar(&dir, "p", ".", "path")
	flag.Parse()
	// Validate flags.
	if !isDir(dir) {
		log.Fatal(dir, "is not a valid path.")
	}
	args = flag.Args()
	if len(args) == 0 {
		log.Fatal("no arguments passed")
	}
	// Read rc file if available.
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	rc, err := ioutil.ReadFile(path.Join(cwd, ".foerc"))
	ignore = []string{}
	if err != nil {
		log.Println("no rc")
		return
	}
	for _, bad := range strings.Split(string(rc), "\n") {
		if bad != "" {
			ignore = append(ignore, bad)
		}
	}
}

// walker implements filepath.WalkFunc.
// Search filename for arguments passed in.
func walker(path string, info os.FileInfo, err error) error {
	if err != nil {
		log.Println(err, "in file", path)
		return nil
	}
	if isDir(path) {
		return nil
	}
	for _, arg := range args {
		if !strings.Contains(strings.ToLower(path), strings.ToLower(arg)) {
			return nil
		}
		for _, bad := range ignore {
			if strings.Contains(path, bad) {
				return nil
			}
		}
	}
	if !isText(path) {
		return nil
	}
	cFiles <- path
	return nil
}

func main() {
	go func() {
		filepath.Walk(dir, walker)
		close(cFiles)
	}()
	paths := make([]string, 0, 50)
	for path := range cFiles {
		paths = append(paths, path)
	}
	if len(paths) < 1 {
		log.Println("No results.")
		return
	}
	if len(paths) > 3 {
		log.Println("Too many results:")
		for _, x := range paths {
			log.Println(x)
		}
		return
	}
	cmd := exec.Command("vim", paths...)
	cmd.Stdout = os.Stdout
	cmd.Stdin = os.Stdin
	cmd.Run()
}
