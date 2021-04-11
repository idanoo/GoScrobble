package goscrobble

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
)

type Track struct {
	UUID          string `json:"uuid"`
	Name          string `json:"name"`
	Length        int    `json:"length"`
	Desc          string `json:"desc"`
	Img           string `json:"img"`
	MusicBrainzID string `json:"mbid"`
	SpotifyID     string `json:"spotify_id"`
}

type TopTrack struct {
	UUID  string `json:"uuid"`
	Name  string `json:"name"`
	Img   string `json:"img"`
	Plays int    `json:"plays"`
}
type TopTracks struct {
	Tracks map[int]TopTrack `json:"tracks"`
}

// insertTrack - This will return if it exists or create it based on MBID > Name
func insertTrack(name string, legnth int, mbid string, spotifyId string, album string, artists []string, tx *sql.Tx) (Track, error) {
	track := Track{}

	// Try locate our track
	if mbid != "" {
		track = getTrackByCol("mbid", mbid, tx)
	} else if spotifyId != "" {
		track = getTrackByCol("spotify_id", spotifyId, tx)
	}

	// If it didn't match above, lookup by name
	if track.UUID == "" {
		track = getTrackWithArtists(name, artists, album, tx)
	}

	// If we can't find it. Lets add it!
	if track.UUID == "" {
		err := insertNewTrack(&track, name, legnth, mbid, spotifyId, tx)

		if err != nil {
			return track, errors.New("Failed to insert track")
		}

		if track.UUID == "" {
			return track, errors.New("Unable to fetch track!")
		}

		err = track.linkTrack(album, artists, tx)
		if err != nil {
			return track, errors.New("Unable to link tracks!")
		}
	}

	// Updates these values if we match earlier!
	if track.MusicBrainzID != mbid {
		track.MusicBrainzID = mbid
		track.updateTrack("mbid", mbid, tx)
	}

	if track.SpotifyID != spotifyId {
		track.SpotifyID = spotifyId
		track.updateTrack("spotify_id", spotifyId, tx)
	}

	return track, nil
}

func getTrackByCol(col string, val string, tx *sql.Tx) Track {
	var track Track
	err := tx.QueryRow(
		"SELECT BIN_TO_UUID(`uuid`, true), `name`, IFNULL(`desc`,''), IFNULL(`img`,''), `mbid` FROM `tracks` WHERE `"+col+"` = ? LIMIT 1",
		val).Scan(&track.UUID, &track.Name, &track.Desc, &track.Img, &track.MusicBrainzID)

	if err != nil {
		if err != sql.ErrNoRows {
			log.Printf("Error fetching tracks: %+v", err)
		}
	}

	return track
}

func getTrackWithArtists(name string, artists []string, album string, tx *sql.Tx) Track {
	var track Track
	artistString := strings.Join(artists, "','")
	err := tx.QueryRow(
		"SELECT BIN_TO_UUID(`uuid`, true), `name`, IFNULL(`desc`,''), IFNULL(`img`,''), `mbid` FROM `tracks` "+
			"LEFT JOIN `track_artist` ON `tracks`.`uuid` = `track_artist`.`track` "+
			"LEFT JOIN `track_album` ON `tracks`.`uuid` = `track_album`.`track` "+
			"WHERE `name` = ? AND BIN_TO_UUID(`track_artist`.`artist`, true) IN ('"+artistString+"') "+
			"AND BIN_TO_UUID(`track_album`.`album`,true) = ? LIMIT 1",
		name, album).Scan(&track.UUID, &track.Name, &track.Desc, &track.Img, &track.MusicBrainzID)

	if err != nil {
		if err != sql.ErrNoRows {
			log.Printf("Error fetching tracks: %+v", err)
		}
	}

	return track
}

func insertNewTrack(track *Track, name string, length int, mbid string, spotifyId string, tx *sql.Tx) error {
	track.UUID = newUUID()
	track.Name = name
	track.Length = length
	track.MusicBrainzID = mbid
	track.SpotifyID = spotifyId

	_, err := tx.Exec("INSERT INTO `tracks` (`uuid`, `name`, `length`, `mbid`, `spotify_id`) "+
		"VALUES (UUID_TO_BIN(?, true),?,?,?,?)", track.UUID, track.Name, track.Length, track.MusicBrainzID, track.SpotifyID)

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

func (track *Track) linkTrackToAlbum(albumUuid string, tx *sql.Tx) error {
	_, err := tx.Exec("INSERT INTO `track_album` (`track`, `album`) "+
		"VALUES (UUID_TO_BIN(?, true), UUID_TO_BIN(?, true))", track.UUID, albumUuid)

	return err
}

func (track *Track) linkTrackToArtists(artists []string, tx *sql.Tx) error {
	var err error
	for _, artist := range artists {
		_, err = tx.Exec("INSERT INTO `track_artist` (`track`, `artist`) "+
			"VALUES (UUID_TO_BIN(?, true),UUID_TO_BIN(?, true))", track.UUID, artist)
		if err != nil {
			return err
		}
	}

	return nil
}

func (track *Track) updateTrack(col string, val string, tx *sql.Tx) error {
	_, err := tx.Exec("UPDATE `tracks` SET `"+col+"` = ? WHERE `uuid` = UUID_TO_BIN(?,true)", val, track.UUID)

	return err
}

func getTrackByUUID(uuid string) (Track, error) {
	var track Track
	err := db.QueryRow("SELECT BIN_TO_UUID(`tracks`.`uuid`, true), `tracks`.`name`, IFNULL(`albums`.`desc`,''), IFNULL(BIN_TO_UUID(`albums`.`uuid`, true),''), `tracks`.`length`, `tracks`.`mbid`, `tracks`.`spotify_id` "+
		"FROM `tracks` "+
		"LEFT JOIN track_album ON track_album.track = tracks.uuid "+
		"LEFT JOIN albums ON track_album.album = albums.uuid "+
		"WHERE `tracks`.`uuid` = UUID_TO_BIN(?, true)",
		uuid).Scan(&track.UUID, &track.Name, &track.Desc, &track.Img, &track.Length, &track.MusicBrainzID, &track.SpotifyID)

	if err != nil {
		fmt.Println(err)
		return track, errors.New("Invalid UUID")
	}

	return track, nil
}

func getTopTracks(userUuid string) (TopTracks, error) {
	var topTracks TopTracks

	rows, err := db.Query("SELECT BIN_TO_UUID(`tracks`.`uuid`, true), `tracks`.`name`, IFNULL(BIN_TO_UUID(`albums`.`uuid`, true),''), count(*) "+
		"FROM `scrobbles` "+
		"JOIN `tracks` ON `tracks`.`uuid` = `scrobbles`.`track` "+
		"JOIN track_album ON track_album.track = tracks.uuid "+
		"JOIN albums ON track_album.album = albums.uuid "+
		"WHERE `user` = UUID_TO_BIN(?, true) "+
		"GROUP BY `scrobbles`.`track` "+
		"ORDER BY count(*) DESC "+
		"LIMIT 14",
		userUuid)
	if err != nil {
		log.Printf("Failed to fetch top tracks: %+v", err)
		return topTracks, errors.New("Failed to fetch top tracks")
	}
	defer rows.Close()

	i := 1
	tempTracks := make(map[int]TopTrack)

	for rows.Next() {
		var track TopTrack
		err := rows.Scan(&track.UUID, &track.Name, &track.Img, &track.Plays)
		if err != nil {
			log.Printf("Failed to fetch track: %+v", err)
			return topTracks, errors.New("Failed to fetch track")
		}

		tempTracks[i] = track
		i++
	}

	topTracks.Tracks = tempTracks

	return topTracks, nil
}
