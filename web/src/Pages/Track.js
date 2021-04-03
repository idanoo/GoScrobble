import React, { useState, useEffect } from 'react';
import '../App.css';
import './Track.css';
import ScaleLoader from 'react-spinners/ScaleLoader';
import ScrobbleTable from '../Components/ScrobbleTable'

const Track = (route) => {
  const [loading, setLoading] = useState(true);
  const [profile, setProfile] = useState({});

  let artist = false;
  if (route && route.match && route.match.params && route.match.params.uuid) {
    artist = route.match.params.uuid;
  } else {
    artist = false;
  }

  useEffect(() => {
    if (!artist) {
      return false;
    }

    // getProfile(username)
    //   .then(data => {
    //     setProfile(data);
    //     console.log(data)
    //     setLoading(false);
    //   })
  }, [artist])

  if (loading) {
    return (
      <div className="pageWrapper">
        <ScaleLoader color="#6AD7E5" />
      </div>
    )
  }

  if (!artist || !artist) {
    return (
      <div className="pageWrapper">
        Unable to fetch user
      </div>
    )
  }

  return (
    <div className="pageWrapper">
      <h1>
        {artist}
      </h1>
      <div className="pageBody">
        Artist
      </div>
    </div>
  );
}

export default Track;