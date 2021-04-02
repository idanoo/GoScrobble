import React, { useContext, useState, useEffect } from 'react';
import '../App.css';
import './User.css';
import { useHistory } from "react-router";
import AuthContext from '../Contexts/AuthContext';
import ScaleLoader from 'react-spinners/ScaleLoader';
import { getUser } from '../Api/index'
import { Button } from 'reactstrap';
import { spotifyConnectionRequest, spotifyDisonnectionRequest } from '../Api/index'

const User = () => {
  const history = useHistory();
  const { user } = useContext(AuthContext);
  const [loading, setLoading] = useState(true);
  const [userdata, setUserdata] = useState({});


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