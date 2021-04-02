import React, { useContext, useState, useEffect } from 'react';
import '../App.css';
import './User.css';
import { useHistory } from "react-router";
import AuthContext from '../Contexts/AuthContext';
import ScaleLoader from 'react-spinners/ScaleLoader';
import { getUser, patchUser } from '../Api/index'
import { Button } from 'reactstrap';

import { spotifyConnectionRequest, spotifyDisonnectionRequest } from '../Api/index'
import TimezoneSelect from 'react-timezone-select'

const User = () => {
  const history = useHistory();
  const { user } = useContext(AuthContext);
  const [loading, setLoading] = useState(true);
  const [userdata, setUserdata] = useState({});

  const updateTimezone = (vals) => {
    console.log(vals)
    setUserdata({...userdata, timezone: vals});
    patchUser({timezone: vals.value})
  }

  useEffect(() => {
    if (!user) {
      return
    }

    getUser()
      .then(data => {
        setUserdata(data);
        setLoading(false);
      })
  }, [user])

  if (!user) {
    history.push("/login")
  }

  if (loading) {
    return (
      <div className="pageWrapper">
        <ScaleLoader color="#6AD7E5" />
      </div>
    )
  }

  return (
    <div className="pageWrapper">
      <h1>
        Welcome {userdata.username}
      </h1>
      <p className="userBody">
      Timezone<br/>
      <TimezoneSelect
          className="userDropdown"
          value={userdata.timezone}
          onChange={updateTimezone}
      /><br/>
        Created At: {userdata.created_at}<br/>
        Email: {userdata.email}<br/>
        Verified: {userdata.verified ? '✓' : '✖'}<br/>
        {userdata.spotify_username
          ? <div>Spotify Account: {userdata.spotify_username}<br/><br/>
          <Button
            color="secondary"
            type="button"
            className="loginButton"
            onClick={spotifyDisonnectionRequest}
          >Disconnect Spotify</Button></div>
          : <div>
            <br/>
            <Button
              color="primary"
              type="button"
              className="loginButton"
              onClick={spotifyConnectionRequest}
            >Connect To Spotify</Button>
          </div>
        }
      </p>
    </div>
  );
}

export default User;