import React, { useState, useEffect } from 'react';
import '../App.css';
import './Album.css';
import ScaleLoader from 'react-spinners/ScaleLoader';
import { getAlbum } from '../Api/index'

const Album = (route) => {
  const [loading, setLoading] = useState(true);
  const [album, setAlbum] = useState({});

  let albumUUID = false;
  if (route && route.match && route.match.params && route.match.params.uuid) {
    albumUUID = route.match.params.uuid;
  } else {
    albumUUID = false;
  }

  useEffect(() => {
    if (!albumUUID) {
      return false;
    }

    getAlbum(albumUUID)
      .then(data => {
        setAlbum(data);
        setLoading(false);
      })
  }, [albumUUID])

  if (loading) {
    return (
      <div className="pageWrapper">
        <ScaleLoader color="#6AD7E5" />
      </div>
    )
  }

  if (!albumUUID || !album) {
    return (
      <div className="pageWrapper">
        Unable to fetch user
      </div>
    )
  }

  return (
    <div className="pageWrapper">
      <h1>
        {album.name}
      </h1>
      <div className="pageBody">
        MusicBrainzId: {album.mbid}<br/>
        SpotifyID: {album.spotify_id}
      </div>
    </div>
  );
}

export default Album;