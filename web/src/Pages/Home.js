import logo from '../logo.png';
import '../App.css';
import './Home.css';
import HomeBanner from '../Components/HomeBanner';
import TopTable from '../Components/TopTable';
import React, { useState, useEffect } from 'react';
import { getTopTracks, getTopArtists } from '../Api/index'

const Home = () => {
  const [topArtists, setTopArtists] = useState({})
  const [topTracks, setTopTracks] = useState({})
  const [tableLoading, setTableLoading] = useState(true);

  useEffect(() => {
    // Fetch top tracks
    // if (topTracks && Object.keys(topTracks).length === 0) {
    //   getTopTracks("0", 7)
    //     .then(data => {
    //       setTopTracks(data.tracks)
    //   })
    // }

    // Fetch top artists
    if (topArtists && Object.keys(topArtists).length === 0) {
      getTopArtists("0", 7)
        .then(data => {
          setTopArtists(data.artists)
      })
    }

    setTableLoading(false);
  }, [topTracks, topArtists])

  return (
    <div className="pageWrapper">
      <div className="homeContainer">
        <div>
          <img src={logo} className="App-logo" alt="logo" />
        </div>
        <div className="homeItem">
          <p className="homeText">GoScrobble is an open source music scrobbling service.</p>
          <p className="subHomeText">Supports Spotify, Jellyfin, Navidrome / Subsonic / Airsonic.</p>
        </div>
      </div>
      <HomeBanner />
      {/* <TopTable type="track" items={topTracks} loading={tableLoading} /> */}
      <TopTable type="artist" items={topArtists} loading={tableLoading} extraText="this week" />
    </div>
    );
}

export default Home;
