import React, { useContext, useState, useEffect } from 'react';
import '../App.css';
import './Album.css';
import ScaleLoader from 'react-spinners/ScaleLoader';
import { getAlbum } from '../Api/index'
import TopUserTable from '../Components/TopUserTable';
import TracksForRecordTable from '../Components/TracksForRecordTable';
import AuthContext from '../Contexts/AuthContext';
import { Link } from 'react-router-dom';

const Album = (route) => {
  const [loading, setLoading] = useState(true);
  const [album, setAlbum] = useState({});
  const { user } = useContext(AuthContext);

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
        Unable to fetch album
      </div>
    )
  }

  return (
    <div className="pageWrapper">
      <h1 style={{margin: 0}}>
        {album.name} {user && user.mod && <Link
            key="editbuttonomg"
            to={"/album/" + album.uuid + "/edit"}
          >edit</Link>}
      </h1>
      <div className="pageBody">
        <div style={{display: `flex`, flexWrap: `wrap`, textAlign: `center`}}>
          <div style={{width: `300px`, padding: `0 10px 10px 10px`, textAlign: `left`}}>
            <img src={process.env.REACT_APP_API_URL + "/img/" + album.uuid + "_full.jpg"} alt={album.name} style={{maxWidth: `300px`, maxHeight: `300px`}}/>
          </div>
          <div style={{width: `290px`, padding: `0 10px 10px 10px`, margin: `0 5px 0 5px`, textAlign: `left`}}>
            <span style={{fontSize: '14pt'}}>
              {album.mbid && <a rel="noreferrer" target="_blank" href={"https://musicbrainz.org/album/" + album.mbid}>Open on MusicBrainz<br/></a>}
              {album.spotify_id && <a rel="noreferrer" target="_blank" href={"https://open.spotify.com/album/" + album.spotify_id}>Open on Spotify<br/></a>}
            </span>
          </div>
          <div style={{width: `290px`, padding: `0 10px 10px 10px`}}>
            <h3>Top 10 Scrobblers</h3>
            <TopUserTable albumuuid={album.uuid}/>
          </div>
        </div>
        <TracksForRecordTable albumuuid={album.uuid}/>
      </div>
    </div>
  );
}

export default Album;