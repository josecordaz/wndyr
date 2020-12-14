package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

const API_URL = "https://api.nasa.gov/mars-photos/api/v1/rovers/curiosity/photos"
const API_KEY = "MOPu6sylikQcHo64dlzdvi80h2dxH5pQaz7c5kWl"

func main() {
	count, err := downloadDateImages("", API_KEY, API_URL, os.Args)
	if err != nil {
		log.Fatal("request was unsuccessful ", err)
	}
	fmt.Printf("Request was successful, %d images downloaded \n", count)
}

func downloadDateImages(path, apiKey, apiURL string, args []string) (int, error) {

	if len(args) == 1 {
		return 0, fmt.Errorf("missing date parameter (go run . 2012-08-03)")
	}

	date := args[1]

	if err := validateDateFormat(date); err != nil {
		return 0, err
	}

	dir, err := makeDir(path, date)
	if err != nil {
		return 0, err
	}

	done := make(chan error, 1)
	imagesURL := make(chan string)

	go downloadImages(dir, imagesURL, done)

	count, err := getDateImagesURL(date, apiKey, apiURL, imagesURL)
	if err != nil {
		return 0, err
	}

	if err := <-done; err != nil && err.Error() != "" {
		return 0, err
	}

	return count, nil
}

func downloadImages(path string, images chan string, done chan error) {

	errs := make(chan error, len(images))

	var wg sync.WaitGroup
	for image := range images {
		wg.Add(1)
		go func(image string) {
			defer wg.Done()
			if err := downloadImage(path, image); err != nil {
				errs <- err
			}
		}(image)
	}
	wg.Wait()

	close(errs)

	var err string
	for e := range errs {
		err += e.Error()
	}

	done <- fmt.Errorf(err)

}

func makeDir(path, dir string) (finalPath string, rerr error) {
	if path != "" {
		finalPath = path + "/" + dir
	} else {
		finalPath = dir
	}
	if err := os.Mkdir(finalPath, 0755); err != nil && os.IsNotExist(err) {
		rerr = fmt.Errorf("error creating folder %s", err)
		return
	}
	return
}

func validateDateFormat(date string) error {
	_, err := time.Parse("2006-1-2", date)
	if err != nil {
		return fmt.Errorf("date in the wrong format(YYYY-M-D) %s", err)
	}
	return nil
}

func getDateImagesURL(date string, apiKey, apiURL string, imagesURL chan string) (int, error) {

	params := url.Values{}
	params.Set("api_key", apiKey)
	params.Set("earth_date", date)
	params.Set("page", "1")

	var page int64 = 1

	count := 0
	for {
		response, err := http.Get(apiURL + "?" + params.Encode())
		if err != nil {
			return 0, fmt.Errorf("error on get %s", err)
		}

		var data Data
		err = json.NewDecoder(response.Body).Decode(&data)
		if err != nil {
			return 0, fmt.Errorf("error decoding response %s", err)
		}
		err = response.Body.Close()
		if err != nil {
			return 0, fmt.Errorf("error closing response body %s", err)
		}
		if data.Err != nil {
			return 0, fmt.Errorf(data.Err.Message)
		}
		if len(data.Photos) == 0 {
			break
		}

		for _, photo := range data.Photos {
			imagesURL <- photo.ImgSrc
			count++
		}
		page++
		params.Set("page", strconv.FormatInt(page, 10))

	}

	close(imagesURL)

	return count, nil
}

func downloadImage(folder, url string) error {

	lastSlashIndex := strings.LastIndexByte(url, '/')
	fileName := url[lastSlashIndex+1:]
	filePath := folder + "/" + fileName

	response, err := http.Get(url)
	if err != nil && err != io.EOF {
		return fmt.Errorf("error downloading image %s ", err)
	}

	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("error creating file %s ", err)
	}

	if _, err := io.Copy(file, response.Body); err != nil {
		return fmt.Errorf("error writing file content %s", err)
	}

	response.Body.Close()

	return nil
}
