package goscrobble

import (
	"database/sql"
	"errors"
	"log"
)

type Track struct {
	Uuid          string         `json:"uuid"`
	Name          string         `json:"name"`
	Desc          sql.NullString `json:"desc"`
	Img           sql.NullString `json:"img"`
	MusicBrainzID sql.NullString `json:"mbid"`
}

// insertTrack - This will return if it exists or create it based on MBID > Name
func insertTrack(name string, mbid string, album string, artists []string, tx *sql.Tx) (Track, error) {
	track := Track{}

	if mbid != "" {
		track = fetchTrack("mbid", mbid, tx)
		if track.Uuid == "" {
			err := insertNewTrack(name, mbid, tx)
			if err != nil {
				log.Printf("Error inserting track via MBID %s  %+v", name, err)
				return track, errors.New("Failed to insert track")
			}

			track = fetchTrack("mbid", mbid, tx)
			err = track.linkTrack(album, artists, tx)
			if err != nil {
				return track, err
			}
		}
	} else {
		track = fetchTrack("name", name, tx)
		if track.Uuid == "" {
			err := insertNewTrack(name, mbid, tx)
			if err != nil {
				log.Printf("Error inserting track via Name %s %+v", name, err)
				return track, errors.New("Failed to insert track")
			}

			track = fetchTrack("name", name, tx)
			err = track.linkTrack(album, artists, tx)
			if err != nil {
				return track, err
			}
		}
	}

	if track.Uuid == "" {
		return track, errors.New("Unable to fetch track!")
	}

	return track, nil
}

func fetchTrack(col string, val string, tx *sql.Tx) Track {
	var track Track
	err := tx.QueryRow(
		"SELECT BIN_TO_UUID(`uuid`, true), `name`, `desc`, `img`, `mbid` FROM `tracks` WHERE `"+col+"` = ?",
		val).Scan(&track.Uuid, &track.Name, &track.Desc, &track.Img, &track.MusicBrainzID)

	if err != nil {
		if err != sql.ErrNoRows {
			log.Printf("Error fetching tracks: %+v", err)
		}
	}

	return track
}

func insertNewTrack(name string, mbid string, tx *sql.Tx) error {
	_, err := tx.Exec("INSERT INTO `tracks` (`uuid`, `name`, `mbid`) "+
		"VALUES (UUID_TO_BIN(UUID(), true),?,?)", name, mbid)

	return err
}

func (track *Track) linkTrack(album string, artists []string, tx *sql.Tx) error {
	err := track.linkTrackToAlbum(album, tx)
	if err != nil {
		return err
	}
	err = track.linkTrackToArtists(artists, tx)
	if err != nil {
		return err
	}
	return nil
}

func (track Track) linkTrackToAlbum(albumUuid string, tx *sql.Tx) error {
	_, err := tx.Exec("INSERT INTO `track_album` (`track`, `album`) "+
		"VALUES (UUID_TO_BIN(?, true), UUID_TO_BIN(?, true))", track.Uuid, albumUuid)

	return err
}

func (track Track) linkTrackToArtists(artists []string, tx *sql.Tx) error {
	var err error
	for _, artist := range artists {
		_, err = tx.Exec("INSERT INTO `track_artist` (`track`, `artist`) "+
			"VALUES (UUID_TO_BIN(?, true),UUID_TO_BIN(?, true))", track.Uuid, artist)
		if err != nil {
			return err
		}
	}

	return nil
}
