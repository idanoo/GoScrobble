package goscrobble

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/disintegration/imaging"
)

func importUploadedImage(file multipart.File, uuid string, recordType string) error {
	// Create image
	out, err := os.Create(DataDirectory + string(os.PathSeparator) + "img" + string(os.PathSeparator) + uuid + "_full.jpg")
	if err != nil {
		return err
	}
	defer out.Close()

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	// Write image
	_, err = out.Write(fileBytes)
	if err == nil {
		// Make sure we queue it to process!
		_, err = db.Exec("UPDATE `"+recordType+"` SET `img` = 'pending' WHERE `uuid` = UUID_TO_BIN(?,true)", uuid)
	}

	return err
}

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

	return nil
}

func resizeImages() {
	resizeBulk("artists")
	resizeBulk("albums")
}

func resizeBulk(recordType string) {
	// Fetch pending 500 at a time cause we do it every minute anyway
	rows, err := db.Query("SELECT BIN_TO_UUID(`uuid`, true) FROM `" + recordType + "` WHERE `img` = 'pending' LIMIT 500")
	if err != nil {
		log.Printf("Failed to get pending images: %+v", err)
		return
	}

	// Fetch pending 100 at a time
	for rows.Next() {
		var uuid string
		err := rows.Scan(&uuid)
		if err != nil {
			log.Printf("Failed to fetch record image resize: %+v", err)
			rows.Close()
			return
		}

		// Run the resize to 300px
		success := resizeImage(uuid, 300)
		if !success {
			// If we get an error.. lets just remove the image link for now so it can reimport
			_, err = db.Exec("UPDATE `"+recordType+"` SET `img` = NULL WHERE `uuid` = UUID_TO_BIN(?,true)", uuid)
		} else {
			// Update DB to reflect complete
			_, err = db.Exec("UPDATE `"+recordType+"` SET `img` = 'complete' WHERE `uuid` = UUID_TO_BIN(?,true)", uuid)
		}
	}

	rows.Close()
}

func resizeAlbums() {

}

func resizeImage(uuid string, size int) bool {
	// Open source image
	src, err := imaging.Open(DataDirectory + string(os.PathSeparator) + "img" + string(os.PathSeparator) + uuid + "_full.jpg")
	if err != nil {
		log.Printf("Failed to open image: %+v", err)
		return false
	}

	// Resize image to specified size
	resizedImage := imaging.Resize(src, size, 0, imaging.Lanczos)

	// Save resized image
	err = imaging.Save(resizedImage, DataDirectory+string(os.PathSeparator)+"img"+string(os.PathSeparator)+uuid+"_300px.jpg")
	if err != nil {
		log.Printf("failed to save image: %v", err)
		return false
	}

	return true
}
