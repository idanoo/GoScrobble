import React, { useState, useEffect } from 'react';
import '../App.css';
import './Profile.css';
import ScaleLoader from 'react-spinners/ScaleLoader';
import { getProfile, getTopTracks, getTopArtists } from '../Api/index'
import ScrobbleTable from '../Components/ScrobbleTable'
import TopTable from '../Components/TopTable'

const profileDateRanges = {
  'All time': false,
  'Last year': '365',
  'Last month': '30',
  'Last week': '7',
};

const defaultDateRange = 'Last month';
let activeStyle = { color: '#FFFFFF' };

const Profile = (route) => {
  const [loading, setLoading] = useState(true);
  const [tableLoading, setTableLoading] = useState(false);
  const [active, setActive] = useState(defaultDateRange);
  const [profile, setProfile] = useState({});
  const [topTracks, setTopTracks] = useState({})
  const [topArtists, setTopArtists] = useState({})
  
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
        getTopTracks(data.uuid, profileDateRanges[defaultDateRange])
          .then(data => {
            setTopTracks(data.tracks)
        })

        // Fetch top artists
        getTopArtists(data.uuid, profileDateRanges[defaultDateRange])
          .then(data => {
            setTopArtists(data.artists)
        })

        setLoading(false);
      })

  }, [username])

  const reloadScrobblesForDate = (username, days, name) => {
    setActive(name);
    setTableLoading(true);
    getProfile(username)
    .then(data => {
      setProfile(data);

      // Fetch top tracks
      getTopTracks(data.uuid, days)
        .then(data => {
          setTopTracks(data.tracks);
      })

      // Fetch top artists
      getTopArtists(data.uuid, days)
        .then(data => {
          setTopArtists(data.artists);
      })

      setTableLoading(false)
    })
  };

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
        <div className="profileDateRange">
          {
            Object.entries(profileDateRanges).map((t,k) => <span>
              <span onClick={() => reloadScrobblesForDate(username,t[1], t[0])}
                style={active === t[0] ? activeStyle : {}}
                className="profileDateRangeText">
                  {t[0]}
              </span>
            </span>)
          }
        </div>
        <TopTable type="track" items={topTracks} loading={tableLoading} />
        <br/>
        <TopTable type="artist" items={topArtists} loading={tableLoading} />
        <br/>
        Last 10 scrobbles<br/>
        <ScrobbleTable data={profile.scrobbles}/>
      </div>
    </div>
  );
}

export default Profile;