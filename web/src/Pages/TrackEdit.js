import React, { useContext, useState, useEffect } from 'react';
import '../App.css';
import './TrackEdit.css';
import { useHistory } from 'react-router-dom';
import ScaleLoader from 'react-spinners/ScaleLoader';
import { getTrack } from '../Api/index'
import { Link } from 'react-router-dom';
import AuthContext from '../Contexts/AuthContext';

const TrackEdit = (route) => {
  const history = useHistory();
  const { user } = useContext(AuthContext);

  const [loading, setLoading] = useState(true);
  const [track, setTrack] = useState({});

  let trackUUID = false;
  if (route && route.match && route.match.params && route.match.params.uuid) {
    trackUUID = route.match.params.uuid;
  } else {
    trackUUID = false;
  }

  useEffect(() => {
    if (!trackUUID) {
      return false;
    }

    getTrack(trackUUID)
      .then(data => {
        setTrack(data);
        setLoading(false);
      })
  }, [trackUUID])

  if (!user) {
    history.push("/login")
  }

  if (user && !user.mod) {
    history.push("/Dashboard")
  }

  if (loading) {
    return (
      <div className="pageWrapper">
        <ScaleLoader color="#6AD7E5" />
      </div>
    )
  }

  if (!trackUUID || !track) {
    return (
      <div className="pageWrapper">
        Unable to fetch user
      </div>
    )
  }

  console.log(track)
  let length = "0";
  if (track.length && track.length !== '') {
    length = new Date(track.length * 1000).toISOString().substr(11, 8)
  }

  let artists = [];
  for (let artist of track.artists) {
    const row = (
      <Link
        key={artist.uuid}
        to={"/artist/" + artist.uuid}
      >{artist.name} </Link>
    );
    artists.push(row);
  }

  let albums = [];
  for (let album of track.albums) {
    const row = (
      <Link
        key={album.uuid}
        to={"/album/" + album.uuid}
      >{album.name} </Link>
    );
    albums.push(row);
  }

  return (
    <div className="pageWrapper">
      <h1 style={{margin: 0}}>
      {track.name} {<Link
            key="editbuttonomg"
            to={"/track/" + trackUUID}
          >unedit</Link>}
      </h1>
      <div className="pageBody" style={{width: `900px`, textAlign: `center`}}>
        <img src={process.env.REACT_APP_API_URL + "/img/" + track.img + "_full.jpg"} alt={track.name} style={{maxWidth: `300px`, maxHeight: `300px`}}/>
        <br/>
        <label>Primary Artist ({track.artists[0].name}):</label><br/>
        <input type="text" value={track.artists[0].uuid} style={{width: `420px`}} disabled="true"/><br/>
        <label>Primary Album ({track.albums[0].name})</label><br/>
        <input type="text" value={track.albums[0].uuid} style={{width: `420px`}} disabled="true"/><br/>
        <br/>
        <label>MBID</label><br/>
        <input type="text" value={track.mbid} style={{width: `420px`}} /><br/>
        <label>Spotify ID</label><br/>
        <input type="text" value={track.spotify_id} style={{width: `420px`}} /><br/>
      </div>
    </div>
  );
}

export default TrackEdit;