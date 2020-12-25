package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strings"
)

func uploadFile(w http.ResponseWriter, r *http.Request) {
	fmt.Println("File Upload Endpoint Hit")

	// Parse our multipart form, 10 << 20 specifies a maximum
	// upload of 10 MB files.
	// FormFile returns the first file for the given key `myFile`
	// it also returns the FileHeader so we can get the Filename,
	// the Header and the size of the file
	file, handler, err := r.FormFile("myFile")
	if err != nil {
		fmt.Println("Error Retrieving the File")
		fmt.Println(err)
		return
	}
	fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	fmt.Printf("File Size: %+v\n", handler.Size)
	fmt.Printf("MIME Header: %+v\n", handler.Header)

	// Create a temporary file within our temp-images directory that follows
	// a particular naming pattern
	if _, err := os.Stat("rawVideo/" + handler.Filename); !os.IsNotExist(err) {
		os.Remove("rawVideo/" + handler.Filename)
		// path/to/whatever does not exist
	}
	tempFile, err := os.Create("rawVideo/" + handler.Filename)
	if err != nil {
		fmt.Println(err)
	}

	// read all of the contents of our uploaded file into a
	// byte array
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
	}
	// write this byte array to our temporary file
	tempFile.Write(fileBytes)
	// return that we have successfully uploaded our file!

	fmt.Println("Successfully Uploaded File")
	fmt.Println("Starting FFmpeg Pass 1")

	file.Close()
	tempFile.Close()
	file = nil
	tempFile = nil
	ffmpegprocess := exec.Command("ffmpeg", "-i", "rawVideo/"+handler.Filename, "-crf", "15", "--pass", "1", "webmVid"+FilenameWithoutExtension(handler.Filename)+".webm")
	erro := ffmpegprocess.Start().Error()
	if erro != "" {
		fmt.Println("Error on pass 1: " + erro)
	}
	ffmpegprocess = exec.Command("ffmpeg", "-i", "rawVideo/"+handler.Filename, "-crf", "15", "--pass", "2", "webmVid"+FilenameWithoutExtension(handler.Filename)+".webm")
	if erro != "" {
		fmt.Println("Error on pass 1: " + erro)
	}
	uploadhtm, err := os.Open("upload.html")
	if err != nil {
		panic(err)
	}
	io.Copy(w, uploadhtm)

}

func setupRoutes() {
	http.HandleFunc("/upload", uploadFile)
	port := os.Getenv("PORT")
	fmt.Println(http.ListenAndServe(port, nil))
	fmt.Println("hello")
}

func main() {
	fmt.Println("Hello World")
	setupRoutes()
}

func FilenameWithoutExtension(fn string) string {
	return strings.TrimSuffix(fn, path.Ext(fn))
}
