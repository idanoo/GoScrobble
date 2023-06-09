package goscrobble

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
)

type Track struct {
	UUID          string   `json:"uuid"`
	Name          string   `json:"name"`
	Length        int      `json:"length"`
	Desc          string   `json:"desc"`
	Img           string   `json:"img"`
	MusicBrainzID string   `json:"mbid"`
	SpotifyID     string   `json:"spotify_id"`
	Artists       []Artist `json:"artists"` // Optional Artist
	Albums        []Album  `json:"albums"`  // Optional Albums
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
		`SELECT uuid, name, COALESCE(desc,''), COALESCE(img,''), mbid FROM tracks WHERE "`+col+`" = $1 LIMIT 1`,
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
		`SELECT uuid, name, COALESCE(desc,''), COALESCE(img,''), mbid FROM tracks `+
			`LEFT JOIN track_artist ON tracks.uuid = track_artist.track `+
			`LEFT JOIN track_album ON tracks.uuid = track_album.track `+
			`WHERE name = $1 AND track_artist.artistIN ('`+artistString+`') `+
			`AND track_album.album = $2 LIMIT 1`,
		name, album).Scan(&track.UUID, &track.Name, &track.Desc, &track.Img, &track.MusicBrainzID)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Printf("Error fetching tracks: %+v", err)
		}
	}

	return track
}

// getAllTracksForAlbum - return all album tracks!
func getAllTracksForAlbum(uuid string) (TopTracks, error) {
	var topTracks TopTracks

	rows, err := db.Query("SELECT BIN_TO_UUID(`tracks`.`uuid`, true), `tracks`.`name` "+
		"FROM `albums` "+
		"JOIN `track_album` ON `track_album`.`album` = `albums`.`uuid` "+
		"JOIN `tracks` ON `tracks`.`uuid` = `track_album`.`track` "+
		"WHERE albums.uuid = UUID_TO_BIN(?, true)", uuid)

	if err != nil {
		log.Printf("Failed to fetch top tracks: %+v", err)
		return topTracks, errors.New("Failed to fetch top tracks")
	}
	defer rows.Close()

	i := 1
	tempTracks := make(map[int]TopTrack)

	for rows.Next() {
		var track TopTrack
		err := rows.Scan(&track.UUID, &track.Name)
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

// getAllTracksForArtist - return all album tracks!
func getAllTracksForArtist(uuid string) (TopTracks, error) {
	var topTracks TopTracks

	rows, err := db.Query("SELECT BIN_TO_UUID(`tracks`.`uuid`, true), `tracks`.`name` "+
		"FROM `artists` "+
		"JOIN `track_artist` ON `track_artist`.`artist` = `artists`.`uuid` "+
		"JOIN `tracks` ON `tracks`.`uuid` = `track_artist`.`track` "+
		"WHERE artists.uuid = UUID_TO_BIN(?,true)", uuid)

	if err != nil {
		log.Printf("Failed to fetch top tracks: %+v", err)
		return topTracks, errors.New("Failed to fetch top tracks")
	}
	defer rows.Close()

	i := 1
	tempTracks := make(map[int]TopTrack)

	for rows.Next() {
		var track TopTrack
		err := rows.Scan(&track.UUID, &track.Name)
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

func insertNewTrack(track *Track, name string, length int, mbid string, spotifyId string, tx *sql.Tx) error {
	track.UUID = newUUID()
	track.Name = name
	track.Length = length
	track.MusicBrainzID = mbid
	track.SpotifyID = spotifyId

	_, err := tx.Exec(`INSERT INTO tracks (uuid, name, length, mbid, spotify_id) `+
		`VALUES ($1,$2,$3,$4,$5)`, track.UUID, track.Name, track.Length, track.MusicBrainzID, track.SpotifyID)

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
	_, err := tx.Exec(`INSERT INTO track_album (track, album) `+
		`VALUES ($1, $2)`, track.UUID, albumUuid)

	return err
}

func (track *Track) linkTrackToArtists(artists []string, tx *sql.Tx) error {
	var err error
	for _, artist := range artists {
		_, err = tx.Exec(`INSERT INTO track_artist (track, artist) `+
			`VALUES ($1,$2)`, track.UUID, artist)
		if err != nil {
			return err
		}
	}

	return nil
}

func (track *Track) updateTrack(col string, val string, tx *sql.Tx) error {
	_, err := tx.Exec(`UPDATE tracks SET "`+col+`" = $1 WHERE uuid = $2`, val, track.UUID)

	return err
}

func getTrackByUUID(uuid string) (Track, error) {
	var track Track
	err := db.QueryRow(`SELECT tracks.uuid, tracks.name, COALESCE(albums.desc,''), COALESCE(albums.uuid,''), tracks.length, tracks.mbid, tracks.spotify_id `+
		`FROM tracks `+
		`LEFT JOIN track_album ON track_album.track = tracks.uuid `+
		`LEFT JOIN albums ON track_album.album = albums.uuid `+
		`WHERE tracks.uuid = $1`,
		uuid).Scan(&track.UUID, &track.Name, &track.Desc, &track.Img, &track.Length, &track.MusicBrainzID, &track.SpotifyID)

	if err != nil {
		fmt.Println(err)
		return track, errors.New("Invalid UUID")
	}

	return track, nil
}

func getTopTracks(userUuid string, dayRange string) (TopTracks, error) {
	var topTracks TopTracks

	// dateClause := ""
	// if dayRange != "" {
	// 	dateClause = " AND DATE(created_at) > SUBDATE(CURRENT_DATE, " + dayRange + ") "
	// }

	// whereClause := ""
	// if userUuid != "0" {
	// 	whereClause = "WHERE `scrobbles`.`user` = UUID_TO_BIN('" + userUuid + "', true) "
	// }

	// rows, err := db.Query("SELECT BIN_TO_UUID(`tracks`.`uuid`, true), `tracks`.`name`, IFNULL(BIN_TO_UUID(`albums`.`uuid`, true),''), count(*) " +
	// 	"FROM `scrobbles` " +
	// 	"JOIN `tracks` ON `tracks`.`uuid` = `scrobbles`.`track` " +
	// 	"JOIN track_album ON track_album.track = tracks.uuid " +
	// 	"JOIN albums ON track_album.album = albums.uuid " +
	// 	whereClause +
	// 	dateClause +
	// 	"GROUP BY `scrobbles`.`track` " +
	// 	"ORDER BY count(*) DESC " +
	// 	"LIMIT 14")

	rows, err := db.Query(`SELECT tracks.uuid, tracks.name, COALESCE(albums.uuid,''), count(*) `+
		`FROM scrobbles `+
		`JOIN tracks ON tracks.uuid = scrobbles.track `+
		`JOIN track_album ON track_album.track = tracks.uuid `+
		`JOIN albums ON track_album.album = albums.uuid `+
		`WHERE "user" = $1 `+
		`GROUP BY scrobbles.track `+
		`ORDER BY count(*) DESC `+
		`LIMIT 14`,
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

func (track *Track) loadExtraTrackInfo() error {
	err := track.getArtistsForTrack()
	if err != nil {
		return err
	}
	err = track.getAlbumsForTrack()
	return err
}

func (track *Track) getArtistsForTrack() error {
	artists := []Artist{}

	rows, err := db.Query(`SELECT track_artist.artist `+
		`FROM track_artist `+
		`WHERE track_artist.track = $1`,
		track.UUID)
	if err != nil {
		log.Printf("Failed to fetch artists for track: %+v", err)
		return errors.New("Failed to fetch top tracks")
	}
	defer rows.Close()

	for rows.Next() {
		var artistUUID string
		err := rows.Scan(&artistUUID)
		if err != nil {
			log.Printf("Failed to fetch track_artist: %+v", err)
			return errors.New("Failed to fetch track")
		}

		artist, err := getArtistByUUID(artistUUID)
		if err != nil {
			return err
		}
		artists = append(artists, artist)
	}

	track.Artists = artists
	return nil
}

func (track *Track) getAlbumsForTrack() error {
	albums := []Album{}

	rows, err := db.Query(`SELECT track_album.album `+
		`FROM track_album `+
		`WHERE track_album.track = $1`,
		track.UUID)
	if err != nil {
		log.Printf("Failed to fetch album for track: %+v", err)
		return errors.New("Failed to fetch top tracks")
	}
	defer rows.Close()

	for rows.Next() {
		var albumUUID string
		err := rows.Scan(&albumUUID)
		if err != nil {
			log.Printf("Failed to fetch track_album: %+v", err)
			return errors.New("Failed to fetch track")
		}

		album, err := getAlbumByUUID(albumUUID)
		if err != nil {
			return err
		}
		albums = append(albums, album)
	}

	track.Albums = albums
	return nil
}

// getTopUsersForTrackUUID  - Returns list of top users for a track
func getTopUsersForTrackUUID(trackUUID string, limit int, page int) (TopUserResponse, error) {
	response := TopUserResponse{}
	var count int

	total, err := getDbCount(
		`SELECT COUNT(*) FROM scrobbles WHERE track = $1 GROUP BY track, "user"`, trackUUID)

	if err != nil {
		log.Printf("Failed to fetch scrobble count: %+v", err)
		return response, errors.New("Failed to fetch combined scrobbles")
	}

	rows, err := db.Query(
		`SELECT scrobbles.user, users.username, COUNT(*) `+
			`FROM scrobbles `+
			`JOIN users ON scrobbles."user" = users.uuid `+
			`WHERE track = $1 `+
			`GROUP BY scrobbles."user" `+
			`ORDER BY COUNT(*) DESC LIMIT $2`,
		trackUUID, limit)

	if err != nil {
		log.Printf("Failed to fetch scrobbles: %+v", err)
		return response, errors.New("Failed to fetch combined scrobbles")
	}
	defer rows.Close()

	for rows.Next() {
		item := TopUserResponseItem{}
		err := rows.Scan(&item.UserUUID, &item.UserName, &item.Count)
		if err != nil {
			log.Printf("Failed to fetch scrobbles: %+v", err)
			return response, errors.New("Failed to fetch combined scrobbles")
		}
		count++
		response.Items = append(response.Items, item)
	}

	err = rows.Err()
	if err != nil {
		log.Printf("Failed to fetch scrobbles: %+v", err)
		return response, errors.New("Failed to fetch scrobbles")
	}

	response.Meta.Count = count
	response.Meta.Total = total
	response.Meta.Page = page

	return response, nil
}
