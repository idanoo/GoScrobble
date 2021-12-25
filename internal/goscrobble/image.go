package goscrobble

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func importImage(uuid string, url string) error {
	// Create image
	out, err := os.Create(DataDirectory + string(os.PathSeparator) + "img" + string(os.PathSeparator) + uuid + "_full.jpg")
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Bad response status: %s", resp.Status)
	}

	// Writer the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	// Goroutine the resize to keep it _faaaast_
	go resizeImage(uuid)
	return nil
}

func resizeImage(uuid string) {
	// resize to 300x300 and maybe smaller?

	return
}
