package toolkit

import (
	"fmt"
	"image"
	"image/png"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"sync"
	"testing"
)

func TestTools_RandomString(t *testing.T) {
	var testTools Tools

	s := testTools.RandomString(10)
	if len(s) != 10 {
		t.Errorf("wrong length random string length")

	}
}

var uploadTests = []struct {
	name          string
	allowedTypes  []string
	renameFile    bool
	errorExpected bool
}{
	{name: "allowed no rename", allowedTypes: []string{"image/png", "image/jpeg"}, renameFile: false, errorExpected: false},
	{name: "allowed rename", allowedTypes: []string{"image/png", "image/jpeg"}, renameFile: true, errorExpected: false},
	{name: "file not allowed", allowedTypes: []string{"image/jpeg"}, renameFile: false, errorExpected: true},
}

func TestTools_UploadFiles(t *testing.T) {
	for _, e := range uploadTests {
		pr, pw := io.Pipe()
		writer := multipart.NewWriter(pw)
		wg := sync.WaitGroup{}
		wg.Add(1)

		go func() {
			defer writer.Close()
			defer wg.Done()

			part, err := writer.CreateFormFile("file", "./testdata/gola.png")
			if err != nil {
				t.Error(err)
			}

			f, err := os.Open("./testdata/gola.png")
			if err != nil {
				t.Error(err)
			}

			img, _, err := image.Decode(f)
			if err != nil {
				t.Error("Error decoding the image")
			}

			err = png.Encode(part, img)
			if err != nil {
				t.Error(err)
			}

		}()

		// read fromt the pipe which recieves the data
		request := httptest.NewRequest("POST", "/", pr)
		request.Header.Add("Content-Type", writer.FormDataContentType())

		var testTools Tools
		testTools.AllowedFileTypes = e.allowedTypes
		UploadedFiles, err := testTools.UploadFiles(request, "./testdata/uploads/", e.renameFile)
		if err != nil && !e.errorExpected {
			t.Error(err)
		}

		if !e.errorExpected {
			if _, err := os.Stat(fmt.Sprintf("./testdata/uploads/%s", UploadedFiles[0].NewFileName)); os.IsNotExist(err) {
				t.Errorf("%s: expected file to exist: %s", e.name, err.Error())
			}
			//? clean Up
			_ = os.Remove(fmt.Sprintf("./testdata/uploads/%s", UploadedFiles[0].NewFileName))
		}

		if !e.errorExpected && err != nil {
			t.Errorf("%s: Expected error bit none recieved", e.name)
		}
		wg.Wait()
	}
}

func TestTools_UploadOneFile(t *testing.T) {
	for _, e := range uploadTests {
		pr, pw := io.Pipe()
		writer := multipart.NewWriter(pw)

		go func() {
			defer writer.Close()

			part, err := writer.CreateFormFile("file", "./testdata/gola.png")
			if err != nil {
				t.Error(err)
			}

			f, err := os.Open("./testdata/gola.png")
			if err != nil {
				t.Error(err)
			}

			img, _, err := image.Decode(f)
			if err != nil {
				t.Error("Error decoding the image")
			}

			err = png.Encode(part, img)
			if err != nil {
				t.Error(err)
			}

		}()

		// read fromt the pipe which recieves the data
		request := httptest.NewRequest("POST", "/", pr)
		request.Header.Add("Content-Type", writer.FormDataContentType())

		var testTools Tools

		UploadedFiles, err := testTools.UploadOneFile(request, "./testdata/uploads/", true)
		if err != nil && !e.errorExpected {
			t.Error(err)
		}

		if _, err := os.Stat(fmt.Sprintf("./testdata/uploads/%s", UploadedFiles.NewFileName)); os.IsNotExist(err) {
			t.Errorf("expected file to exist: %s", err.Error())
		}
		//? clean Up
		_ = os.Remove(fmt.Sprintf("./testdata/uploads/%s", UploadedFiles.NewFileName))

	}
}

func TestTools_CreateDirIfNotExist(t *testing.T) {
	var testTools Tools

	err := testTools.CreateDirIfNotExist("./testdata/myDir")
	if err != nil {
		t.Error(err)
	}

	err = testTools.CreateDirIfNotExist("./testdata/myDir")
	if err != nil {
		t.Error(err)
	}
	_ = os.Remove("./testdata/myDir")
}

var slugTests = []struct {
	name          string
	s             string
	expected      string
	errorExpected bool
}{
	{name: "valid string", s: "now is the time", expected: "now-is-the-time", errorExpected: false},
	{name: "empty string", s: "", expected: "", errorExpected: true},
	{name: "complex string", s: "now is the TIME? eh !! 123 &(* 4", expected: "now-is-the-time-eh-123-4", errorExpected: false},
	{name: "japanese string", s: "こんにちは", expected: "", errorExpected: true},
	{name: "japanese+english string", s: "こんにちは means hello there !!", expected: "means-hello-there", errorExpected: false},
}

func TestTools_Slugify(t *testing.T) {
	var tools Tools

	for _, e := range slugTests {
		slug, err := tools.Slugify(e.s)
		if err != nil && !e.errorExpected {
			t.Errorf("%s: error recieved when mone expected: %s", e.name, err.Error())
		}

		if !e.errorExpected && slug != e.expected {
			t.Errorf("%s: wrong slug return; expeected %s, got %s", e.name, e.expected, slug)
		}
	}
}

func TestTools_DownloadStaticFile(t *testing.T) {
	rr := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/", nil)

	var testTools Tools

	testTools.DownloadStaticFile(rr, req, "./testdata", "gola.png", "golang.png")

	res := rr.Result()
	defer res.Body.Close()
	if res.Header["Content-Length"][0] != "404763" {
		t.Error("Wrong content length of ", res.Header["Content-Length"][0])
	}

	if res.Header["Content-Disposition"][0] != "attachment; filename=\"golang.png\"" {
		t.Error("Wrong content disposition")
	}

	_, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Error(err)
	}
}
