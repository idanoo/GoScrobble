import React, { useState, useEffect } from 'react';
import '../App.css';
import './Artist.css';
import ScaleLoader from 'react-spinners/ScaleLoader';
import { getArtist } from '../Api/index'

const Artist = (route) => {
  const [loading, setLoading] = useState(true);
  const [artist, setArtist] = useState({});

  let artistUUID = false;
  if (route && route.match && route.match.params && route.match.params.uuid) {
    artistUUID = route.match.params.uuid;
  } else {
    artistUUID = false;
  }

  useEffect(() => {
    if (!artistUUID) {
      return false;
    }

    getArtist(artistUUID)
      .then(data => {
        setArtist(data);
        setLoading(false);
      })
  }, [artistUUID])

  if (loading) {
    return (
      <div className="pageWrapper">
        <ScaleLoader color="#6AD7E5" />
      </div>
    )
  }

  if (!artistUUID || !artist) {
    return (
      <div className="pageWrapper">
        Unable to fetch user
      </div>
    )
  }

  return (
    <div className="pageWrapper">
      <h1>
        {artist.name}
      </h1>
      <div className="pageBody">
        MusicBrainzId: {artist.mbid}<br/>
        SpotifyID: {artist.spotify_id}
      </div>
    </div>
  );
}

export default Artist;