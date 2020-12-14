package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateDateFormat(t *testing.T) {

	assert := assert.New(t)

	err := validateDateFormat("2sdaf")
	assert.Error(err)

	err = validateDateFormat("2020-8-6")
	assert.NoError(err)

}

func TestMakeDir(t *testing.T) {

	assert := assert.New(t)

	tempDir := t.TempDir()

	path, err := makeDir(tempDir, "subfolder")
	assert.NoError(err)
	assert.DirExists(path)

	os.Remove(tempDir)
	os.Remove(path)

	path, err = makeDir("", "subfolder")
	assert.NoError(err)
	assert.DirExists(path)

	os.Remove("subfolder")
}

func TestGetDateImagesURL(t *testing.T) {

	assert := assert.New(t)

	myAPIKey := "myapikey"
	fileURL := "https://foo.com/image.jpg"
	photo := Photo{
		ImgSrc: fileURL,
	}
	d := Data{
		Photos: []*Photo{&photo},
	}

	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		page := req.URL.Query().Get("page")
		fmt.Println("page", page)
		if page == "1" {
			bts, _ := json.Marshal(d)
			rw.Header().Set("Content-Type", "application/json")
			rw.Write(bts)
		} else {
			assert.Equal("2", page)
			rw.Header().Set("Content-Type", "application/json")
			rw.Write([]byte("{}"))
		}

	}))
	defer server.Close()

	images := make(chan string, 1)

	count, err := getDateImagesURL("2020-8-6", myAPIKey, server.URL, images)
	assert.NoError(err)
	assert.Equal(1, count)

	img := <-images
	assert.Equal(fileURL, img)

}

func TestDownloadImages(t *testing.T) {

	assert := assert.New(t)

	date := "2020-8-3"

	server := httptest.NewServer(http.FileServer(http.Dir("/test_data")))
	defer server.Close()

	tempDir := t.TempDir()

	path, err := makeDir(tempDir, date)
	assert.NoError(err)

	images := make(chan string, 1)
	images <- server.URL + "/image.jpg"
	close(images)

	done := make(chan error, 1)

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	downloadImages(path, images, done)
	e := <-done
	assert.Equal("", e.Error())

	assert.FileExists(tempDir + "/" + date + "/" + "image.jpg")

	os.Remove(tempDir)

}

func TestDownloadDateImages(t *testing.T) {

	assert := assert.New(t)
	date := "2020-8-6"
	myAPIKey := "myapikey"
	dir := t.TempDir()

	fileServer := httptest.NewServer(http.FileServer(http.Dir("/test_data")))
	defer fileServer.Close()

	fileURL := "http://foo.com/image.jpg"
	photo := Photo{
		ImgSrc: fileURL,
	}
	data := Data{
		Photos: []*Photo{&photo},
	}

	callCount := 0
	roverAPIserver := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Header().Set("Content-Type", "application/json")
		page := req.URL.Query().Get("page")
		apiKey := req.URL.Query().Get("api_key")
		earthDate := req.URL.Query().Get("earth_date")
		assert.Equal(myAPIKey, apiKey)
		assert.Equal(date, earthDate)
		if callCount == 0 {
			assert.Equal("1", page)
			bts, _ := json.Marshal(data)
			rw.Write(bts)
		} else {
			assert.Equal("2", page)
			rw.Write([]byte("{}"))
		}
		callCount++
	}))
	defer roverAPIserver.Close()

	// Execution
	count, err := downloadDateImages(dir, myAPIKey, roverAPIserver.URL, []string{"", date})
	assert.NoError(err)
	assert.Equal(1, count)

	assert.FileExists(dir + "/" + date + "/" + "image.jpg")

	os.Remove(dir)

	count, err = downloadDateImages(dir, myAPIKey, roverAPIserver.URL, []string{"", "7634rs"})
	assert.Error(err)
	assert.Equal(0, count)

	count, err = downloadDateImages(dir, myAPIKey, roverAPIserver.URL, []string{""})
	assert.Error(err)
	assert.Equal(0, count)

}

func TestDownloadImage(t *testing.T) {

	assert := assert.New(t)

	server := httptest.NewServer(http.FileServer(http.Dir("/test_data")))
	defer server.Close()

	tempDir := t.TempDir()

	err := downloadImage(tempDir, server.URL+"image2.jpg")
	assert.Error(err)

	assert.NoFileExists(tempDir + "/" + "image2.jpg")

	os.Remove(tempDir)

}
