import React, { useContext, useState, useEffect } from 'react';
import '../App.css';
import './Track.css';
import TopUserTable from '../Components/TopUserTable';
import ScaleLoader from 'react-spinners/ScaleLoader';
import { getTrack } from '../Api/index'
import { Link } from 'react-router-dom';
import AuthContext from '../Contexts/AuthContext';

const Track = (route) => {
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
        {track.name} {user && false && <Link
            key="editbuttonomg"
            to={"/track/" + track.uuid + "/edit"}
          >edit</Link>}
      </h1>
      <div className="pageBody">
        <div style={{display: `flex`, flexWrap: `wrap`, textAlign: `center`}}>
          <div style={{width: `300px`, padding: `0 10px 10px 10px`, textAlign: `left`}}>
            <img src={process.env.REACT_APP_API_URL + "/img/" + track.img + "_full.jpg"} alt={track.name} style={{maxWidth: `300px`, maxHeight: `300px`}}/>
          </div>
          <div style={{width: `290px`, padding: `0 10px 10px 10px`, margin: `0 5px 0 5px`, textAlign: `left`}}>
            <span style={{fontSize: '14pt'}}>
              {artists}
            </span>
            <br/>
            <span style={{fontSize: '14pt', textDecoration: 'none'}}>
              {albums}
            </span>
            <br/><br/>
            {track.mbid && <a rel="noreferrer" target="_blank" href={"https://musicbrainz.org/track/" + track.mbid}>Open on MusicBrainz<br/></a>}
            {track.spotify_id && <a rel="noreferrer" target="_blank" href={"https://open.spotify.com/track/" + track.spotify_id}>Open on Spotify<br/></a>}
            {length && <span>Track Length: {length}</span>}
          </div>
          <div style={{width: `290px`, padding: `0 10px 10px 10px`}}>
            <h3>Top 10 Scrobblers</h3>
            <TopUserTable trackuuid={track.uuid}/>
          </div>
        </div>
      </div>
    </div>
  );
}

export default Track;