package main

import (
	"flag"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func walkFnClosure(path string) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		if strings.Contains(path, "converted-480") {
			log.Printf("\n\n\t already converted %s, exiting:: ", path)
			return nil
		}

		//convert
		fileToConvert := path
		convert(fileToConvert)

		return nil
	}
}

func convert(filename string) {
	log.Println("\n\n\t converting:: ", filename)
	cmd := exec.Command("ffmpeg", "-y", "-i", filename, "-vf", "scale=640:480", "-movflags", "+faststart", "-tune", "zerolatency", "-crf", "23", "-maxrate", "600k", "-bufsize", "600k", filename+"-converted-480"+".mp4")
	// cmd.Stdin = strings.NewReader("some input")
	// var out bytes.Buffer
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		log.Printf("\n\n\t conversion of %s failed. err=%v", filename, err)
	}
	log.Printf("\n\n\t conversion of %s succeded:: ", filename)

	// delete old file
	deleteFile(filename)

}

func deleteFile(filename string) {
	err := os.Remove(filename)
	if err != nil {
		log.Printf("\n\n\t unable to delete %s. err=%v", filename, err)
	}
}

func main() {
	// 1. cd to directory
	// 2. check if dir has done.txt
	// 3. if it has, exit
	// 4. if it hasn't;
	// 5. check if there are any other dirs inside
	// 6. if there are dirs, goto 1
	// 7. else run ffmpeg
	// 8. create done.txt
	// 9. exit

	var root = "~/Downloads/test"
	flag.StringVar(
		&root,
		"r",
		root,
		"path to dir with item.")
	flag.Parse()

	err := filepath.Walk(root, walkFnClosure(root))
	if err != nil {
		log.Printf("unable to walk filepath. err=%#+v", err)
	}
}
