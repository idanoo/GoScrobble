import React, { useState, useEffect } from 'react';
import '../App.css';
import './Profile.css';
import ScaleLoader from 'react-spinners/ScaleLoader';
import { getProfile, getTopTracks } from '../Api/index'
import ScrobbleTable from '../Components/ScrobbleTable'
import TopTable from '../Components/TopTable'

const Profile = (route) => {
  const [loading, setLoading] = useState(true);
  const [profile, setProfile] = useState({});
  const [topTracks, setTopTracks] = useState({})

  let username = false;
  if (route && route.match && route.match.params && route.match.params.uuid) {
    username = route.match.params.uuid;
  } else {
    username = false;
  }

  useEffect(() => {
    if (!username) {
      return false;
    }

    getProfile(username)
      .then(data => {
        setProfile(data);

        // Fetch top tracks
        getTopTracks(data.uuid)
          .then(data => {
            setTopTracks(data)
          }
        )

        setLoading(false);
      })

  }, [username])

  if (loading) {
    return (
      <div className="pageWrapper">
        <ScaleLoader color="#6AD7E5" />
      </div>
    )
  }

  if (!username || !profile.username) {
    return (
      <div className="pageWrapper">
        Unable to fetch user
      </div>
    )
  }

  return (
    <div className="pageWrapper">
      <h1>
        {profile.username}'s Profile
      </h1>
      <div className="pageBody">
        <TopTable type="track" items={topTracks} />
        <br/>
        Last 10 scrobbles...<br/>
        <ScrobbleTable data={profile.scrobbles}/>
      </div>
    </div>
  );
}

export default Profile;