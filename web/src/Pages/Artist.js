import React, { useContext, useState, useEffect } from 'react';
import '../App.css';
import './Artist.css';
import ScaleLoader from 'react-spinners/ScaleLoader';
import { getArtist } from '../Api/index'
import TracksForRecordTable from '../Components/TracksForRecordTable';
import TopUserTable from '../Components/TopUserTable';
import AuthContext from '../Contexts/AuthContext';
import { Link } from 'react-router-dom';

const Artist = (route) => {
  const [loading, setLoading] = useState(true);
  const [artist, setArtist] = useState({});
  const { user } = useContext(AuthContext);

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
        Unable to fetch artist
      </div>
    )
  }

  return (
    <div className="pageWrapper">
      <h1 style={{margin: 0}}>
        {artist.name} {user && <Link
            key="editbuttonomg"
            to={"/artist/" + artist.uuid + "/edit"}
          >edit</Link>}
      </h1>
      <div className="pageBody">
        <div style={{display: `flex`, flexWrap: `wrap`, textAlign: `center`}}>
          <div style={{width: `300px`, padding: `0 10px 10px 10px`, textAlign: `left`}}>
            <img src={process.env.REACT_APP_API_URL + "/img/" + artist.uuid + "_full.jpg"} alt={artist.name} style={{maxWidth: `300px`, maxHeight: `300px`}}/>
          </div>
          <div style={{width: `290px`, padding: `0 10px 10px 10px`, margin: `0 5px 0 5px`, textAlign: `left`}}>
            <span style={{fontSize: '14pt'}}>
              {artist.mbid && <a rel="noreferrer" target="_blank" href={"https://musicbrainz.org/artist/" + artist.mbid}>Open on MusicBrainz<br/></a>}
              {artist.spotify_id && <a rel="noreferrer" target="_blank" href={"https://open.spotify.com/artist/" + artist.spotify_id}>Open on Spotify<br/></a>}
            </span>
          </div>
          <div style={{width: `290px`, padding: `0 10px 10px 10px`}}>
            <h3>Top 10 Scrobblers</h3>
            <TopUserTable artistuuid={artist.uuid}/>
          </div>
        </div>
        <TracksForRecordTable artistuuid={artist.uuid}/>
      </div>
    </div>
  );
}

export default Artist;