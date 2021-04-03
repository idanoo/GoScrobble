package goscrobble

import (
	"database/sql"
	"errors"
	"log"
	"strings"
)

type Track struct {
	Uuid          string         `json:"uuid"`
	Name          string         `json:"name"`
	Length        int            `json:"length"`
	Desc          sql.NullString `json:"desc"`
	Img           sql.NullString `json:"img"`
	MusicBrainzID sql.NullString `json:"mbid"`
	SpotifyID     sql.NullString `json:"spotify_id"`
}

// insertTrack - This will return if it exists or create it based on MBID > Name
func insertTrack(name string, legnth int, mbid string, spotifyId string, album string, artists []string, tx *sql.Tx) (Track, error) {
	track := Track{}

	// Try locate our track
	if mbid != "" {
		track = fetchTrack("mbid", mbid, tx)
	} else if spotifyId != "" {
		track = fetchTrack("spotify_id", spotifyId, tx)
	}

	// If it didn't match above, lookup by name
	if track.Uuid == "" {
		// TODO: add artist check here too
		track = fetchTrackWithArtists(name, artists, tx)
	}

	// If we can't find it. Lets add it!
	if track.Uuid == "" {
		err := insertNewTrack(name, legnth, mbid, spotifyId, tx)
		if err != nil {
			return track, errors.New("Failed to insert track")
		}

		// Fetch the recently inserted track to get the UUID
		if mbid != "" {
			track = fetchTrack("mbid", mbid, tx)
		} else if spotifyId != "" {
			track = fetchTrack("spotify_id", spotifyId, tx)
		}

		if track.Uuid == "" {
			track = fetchTrackWithArtists(name, artists, tx)
		}

		err = track.linkTrack(album, artists, tx)
		if err != nil {
			return track, errors.New("Unable to link tracks!")
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

func fetchTrackWithArtists(name string, artists []string, tx *sql.Tx) Track {
	var track Track
	artistString := strings.Join(artists, "','")

	err := tx.QueryRow(
		"SELECT BIN_TO_UUID(`uuid`, true), `name`, `desc`, `img`, `mbid` FROM `tracks` "+
			"LEFT JOIN `track_artist` ON `tracks`.`uuid` = `track_artist`.`track` "+
			"WHERE `name` = ? AND BIN_TO_UUID(`track_artist`.`artist`, true) IN ('`"+artistString+"`')",
		name).Scan(&track.Uuid, &track.Name, &track.Desc, &track.Img, &track.MusicBrainzID)

	if err != nil {
		if err != sql.ErrNoRows {
			log.Printf("Error fetching tracks: %+v", err)
		}
	}

	return track
}

func insertNewTrack(name string, length int, mbid string, spotifyId string, tx *sql.Tx) error {
	_, err := tx.Exec("INSERT INTO `tracks` (`uuid`, `name`, `length`, `mbid`, `spotify_id`) "+
		"VALUES (UUID_TO_BIN(UUID(), true),?,?,?,?)", name, length, mbid, spotifyId)

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

func updateTrack(uuid string, col string, val string, tx *sql.Tx) error {
	_, err := tx.Exec("UPDATE `tracks` SET `"+col+"` = ? WHERE `uuid` = UUID_TO_BIN(?,true)", val, uuid)

	return err
}
