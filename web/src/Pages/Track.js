import React, { useState, useEffect } from 'react';
import '../App.css';
import './Track.css';
import ScaleLoader from 'react-spinners/ScaleLoader';
import { getTrack } from '../Api/index'

const Track = (route) => {
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


  return (
    <div className="pageWrapper">
      <h1>
        {track.name}
      </h1>

      <div className="pageBody" style={{textAlign: `center`}}>
        <img src={process.env.REACT_APP_API_URL + "/img/" + track.img + "_full.jpg"} alt={track.name} style={{maxWidth: `300px`, maxHeight: `300px`}}/><br/><br/>
        {track.mbid && <a rel="noreferrer" target="_blank" href={"https://musicbrainz.org/track/" + track.mbid}>Open on MusicBrainz<br/></a>}
        {track.spotify_id && <a rel="noreferrer" target="_blank" href={"https://open.spotify.com/track/" + track.spotify_id}>Open on Spotify<br/></a>}
        Track Length: {length && length}
      </div>
    </div>
  );
}

export default Track;