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

  let length = "0";
  if (track.length !== '') {
    length = new Date(track.length * 1000).toISOString().substr(11, 8)
  }


  return (
    <div className="pageWrapper">
      <h1>
        {track.name}
      </h1>
      <div className="pageBody">
        MusicBrainzId: {track.mbid}<br/>
        SpotifyID: {track.spotify_id}<br/>
        Track Length: {length}
      </div>
    </div>
  );
}

export default Track;