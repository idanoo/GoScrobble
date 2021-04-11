package goscrobble

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func importImage(uuid string, url string) error {
	// Create the file
	path, err := os.Getwd()
	if err != nil {
		return err
	}

	out, err := os.Create(path + string(os.PathSeparator) + StaticDirectory + string(os.PathSeparator) + "img" + string(os.PathSeparator) + uuid + "_full.jpg")
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

	return nil
}

func resizeImage(uuid string) error {
	return nil
}
