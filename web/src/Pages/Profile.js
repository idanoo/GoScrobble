import React, { useState, useEffect } from 'react';
import '../App.css';
import './Profile.css';
import ScaleLoader from 'react-spinners/ScaleLoader';
import { getProfile } from '../Api/index'
import ScrobbleTable from '../Components/ScrobbleTable'

const Profile = (route) => {
  const [loading, setLoading] = useState(true);
  const [profile, setProfile] = useState({});

  let username = false;
  if (route && route.match && route.match.params && route.match.params.uuid) {
    username = route.match.params.uuid
  }

  useEffect(() => {
    if (!username) {
      return false;
    }

    getProfile(username)
      .then(data => {
        setProfile(data);
        console.log(data)
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

  if (!username || Object.keys(profile).length === 0) {
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
      <div className="profileBody">
        Last 10 scrobbles...<br/>
      <ScrobbleTable data={profile.scrobbles}/>
      </div>
    </div>
  );
}

export default Profile;