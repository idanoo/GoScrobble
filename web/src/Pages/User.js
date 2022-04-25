import React, { useContext, useState, useEffect } from 'react';
import '../App.css';
import './User.css';
import { useHistory } from "react-router";
import AuthContext from '../Contexts/AuthContext';
import ScaleLoader from 'react-spinners/ScaleLoader';
import { Button } from 'reactstrap';
import { Formik, Form, Field } from 'formik';
import { confirmAlert } from 'react-confirm-alert';
import 'react-confirm-alert/src/react-confirm-alert.css';
import {
  getUser,
  patchUser,
  spotifyConnectionRequest,
  spotifyDisonnectionRequest,
  navidromeDisonnectionRequest,
  navidromeConnectionRequest,
} from '../Api/index'
import TimezoneSelect from 'react-timezone-select'

const User = () => {
  const history = useHistory();
  const { user, Logout } = useContext(AuthContext);
  const [loading, setLoading] = useState(true);
  const [userdata, setUserdata] = useState({});

  const updateTimezone = (vals) => {
    setUserdata({...userdata, timezone: vals});
    patchUser({timezone: vals.value})
  }

  const resetTokenPopup = () => {
    confirmAlert({
      title: 'Reset token',
      message: 'Resetting your token will require you to update your Jellyfin server / custom scroblers with the new token. Continue?',
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

  const deleteAccountPopup = () => {
    confirmAlert({
      title: 'Delete Account',
      message: 'This will disable your account and queue it for deletion. Are you sure?',
      buttons: [
        {
          label: 'Yes',
          onClick: () => deleteAccount()
        },
        {
          label: 'No',
        }
      ]
    });
  };

  const connectNavidromePopup = () => {
    confirmAlert({
      title: 'Connect Navidrome',
      buttons: [
        {
          label: 'Close',
        }
      ],
      childrenElement: () =>  <Formik
            initialValues={{ url: '', username: '', password: '' }}
            onSubmit={values => navidromeConnectionRequest(values)}
        >
          <Form>
          <label>
            Server URL<br/>
            <Field
              name="url"
              type="text"
            />
          </label>
          <br/>
          <label>
            Username<br/>
            <Field
              name="username"
              type="text"
            />
          </label>
          <br/>
          <label>
            Password<br/>
            <Field
              name="password"
              type="password"
            />
          </label>
          <br/><br/>
          <Button
            color="primary"
            type="submit"
            className="loginButton"
          >Connect</Button>
        </Form>
      </Formik>,
    });
  };

  const disconnectNavidromePopup = () => {
    confirmAlert({
      title: 'Disconnect Navidrome',
      message: 'Are you sure you want to disconnect your Navidrome connection?',
      buttons: [
        {
          label: 'Disconnect',
          onClick: () => navidromeDisonnectionRequest()
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
      message: 'Are you sure you want to disconnect your Spotify account?',
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

  const connectJellyfinPopup = () => {
    confirmAlert({
      title: 'Connect Jellyfin',
      message: 'Install the webhook plugin. Add a webhook to ' + process.env.REACT_APP_API_URL + '/api/v1/ingress/jellyfin?key='+userdata.token
        +'\nSet it to only send "Playback Start" and "Songs/Albums"',
      buttons: [
        {
          label: 'Close',
        }
      ]
    });
  }

  const connectOtherPopup = () => {
    confirmAlert({
      title: 'Connect Jellyfin',
      message: 'Endpoint: ' + process.env.REACT_APP_API_URL + '/api/v1/ingress/multiscrobbler?key='+userdata.token
        +'\nNeed to send JSON body with a string array for artists names, album:string, track:string, playDate:timestamp of scrobble, duration:tracklength in seconds',
      buttons: [
        {
          label: 'Close',
        }
      ]
    });
  }

  const resetToken = () => {
    setLoading(true);
    patchUser({ token: '' })
      .then(() => {
        getUser()
        .then(data => {
          setUserdata(data);
          setLoading(false);
        })
      })
  }

  const deleteAccount = () => {
    setLoading(true);
    patchUser({ active: 0 })
      .then(() => {
        getUser()
        .then(data => {
          setUserdata(data);
          Logout();
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
      <div style={{display: `flex`, flexWrap: `wrap`, textAlign: `center`}}>
        <div style={{width: `300px`, padding: `0 10px 10px 10px`, textAlign: `left`}}>
          <h3 style={{textAlign: `center`}}>Profile</h3>
          Timezone<br/>
          <TimezoneSelect
              className="userDropdown"
              value={userdata.timezone}
              onChange={updateTimezone}
          /><br/>
          Created At:<br/>{userdata.created_at}<br/>
          Email:<br/>{userdata.email}<br/>
          Verified: {userdata.verified ? '✓' : '✖'}
        </div>
        <div style={{width: `300px`, padding: `0 10px 10px 10px`}}>
          <h3>Scrobblers</h3>
          <br/>
          {userdata.spotify_username
          ? <div>
              <Button
                color="secondary"
                type="button"
                className="userButton"
                onClick={disconnectSpotifyPopup}
              >Disconnect Spotify ({userdata.spotify_username})</Button>
            </div>
          : <div>
              <Button
                color="primary"
                type="button"
                className="userButton"
                onClick={spotifyConnectionRequest}
              >Connect To Spotify</Button>
            </div>
          }
          <br/>
          {userdata.navidrome_server
            ? <Button
              color="secondary"
              type="button"
              className="userButton"
              onClick={disconnectNavidromePopup}
            >Disconnect Navidrome ({userdata.navidrome_server})</Button>
            : <Button
                color="primary"
                type="button"
                className="userButton"
                onClick={connectNavidromePopup}
              >Connect Navidrome</Button>
          }
          <br/><br/>
          <Button
            color="primary"
            type="button"
            className="userButton"
            onClick={connectJellyfinPopup}
          >Connect Jellyfin</Button>
          <br/><br/>
          <Button
            color="primary"
            type="button"
            className="userButton"
            onClick={connectOtherPopup}
          >Other Scrobblers</Button>
        </div>
        <div style={{width: `300px`, padding: `0 10px 10px 10px`}}>
          <h3>Sad Settings</h3>
          <br/>
          <Button
            color="secondary"
            type="button"
            className="userButton"
            onClick={deleteAccountPopup}
          >Disable Account</Button>
          <br/><br/>
          <Button
            color="secondary"
            type="button"
            className="userButton"
            onClick={resetTokenPopup}
          >Reset Scrobbler Token</Button>
        </div>
      </div>
    </div>
  );
}

export default User;