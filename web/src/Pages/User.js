import React, { useContext, useState, useEffect } from 'react';
import '../App.css';
import './User.css';
import { useHistory } from "react-router";
import AuthContext from '../Contexts/AuthContext';
import ScaleLoader from 'react-spinners/ScaleLoader';
import { getUser, patchUser } from '../Api/index'
import { Button } from 'reactstrap';
import { confirmAlert } from 'react-confirm-alert';
import 'react-confirm-alert/src/react-confirm-alert.css';
import { spotifyConnectionRequest, spotifyDisonnectionRequest, resetScrobbleToken } from '../Api/index'
import TimezoneSelect from 'react-timezone-select'

const User = () => {
  const history = useHistory();
  const { user } = useContext(AuthContext);
  const [loading, setLoading] = useState(true);
  const [userdata, setUserdata] = useState({});

  const updateTimezone = (vals) => {
    setUserdata({...userdata, timezone: vals});
    patchUser({timezone: vals.value})
  }

  const resetTokenPopup = () => {
    confirmAlert({
      title: 'Reset token',
      message: 'Resetting your token will require you to update your sources with the new token. Continue?',
      buttons: [
        {
          label: 'Reset',
          onClick: () => resetToken()
        },
        {
          label: 'No',
        }
      ]
    });
  };

  const disconnectSpotifyPopup = () => {
    confirmAlert({
      title: 'Disconnect Spotify',
      message: 'Are you sure you want to disconnect your spotify account?',
      buttons: [
        {
          label: 'Disconnect',
          onClick: () => spotifyDisonnectionRequest()
        },
        {
          label: 'No',
        }
      ]
    });
  };

  const resetToken = () => {
    setLoading(true);
    resetScrobbleToken(user.uuid)
      .then(() => {
        getUser()
        .then(data => {
          setUserdata(data);
          setLoading(false);
        })
      })
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
        Token: {userdata.token}<br/>
        <Button
              color="primary"
              type="button"
              className="userButton"
              onClick={resetTokenPopup}
            >Reset Token</Button><br/><br/>
        Created At: {userdata.created_at}<br/>
        Email: {userdata.email}<br/>
        Verified: {userdata.verified ? '✓' : '✖'}<br/>

        {userdata.spotify_username
          ? <div>Spotify Account: {userdata.spotify_username}<br/><br/>
          <Button
            color="secondary"
            type="button"
            className="userButton"
            onClick={disconnectSpotifyPopup}
          >Disconnect Spotify</Button></div>
          : <div>
            <br/>
            <Button
              color="primary"
              type="button"
              className="userButton"
              onClick={spotifyConnectionRequest}
            >Connect To Spotify</Button>
          </div>
        }
      </p>
    </div>
  );
}

export default User;