import React, { useState, useEffect } from 'react';
import '../App.css';
import './Album.css';
import ScaleLoader from 'react-spinners/ScaleLoader';
import ScrobbleTable from '../Components/ScrobbleTable'

const Album = (route) => {
  const [loading, setLoading] = useState(true);
  const [profile, setProfile] = useState({});

  let album = false;
  if (route && route.match && route.match.params && route.match.params.uuid) {
    album = route.match.params.uuid;
  } else {
    album = false;
  }

  useEffect(() => {
    if (!album) {
      return false;
    }

    // getProfile(username)
    //   .then(data => {
    //     setProfile(data);
    //     console.log(data)
    //     setLoading(false);
    //   })
  }, [album])

  if (loading) {
    return (
      <div className="pageWrapper">
        <ScaleLoader color="#6AD7E5" />
      </div>
    )
  }

  if (!album || !album) {
    return (
      <div className="pageWrapper">
        Unable to fetch user
      </div>
    )
  }

  return (
    <div className="pageWrapper">
      <h1>
        {album}
      </h1>
      <div className="pageBody">
        Album
      </div>
    </div>
  );
}

export default Album;