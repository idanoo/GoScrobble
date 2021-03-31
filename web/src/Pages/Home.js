import logo from '../logo.png';
import '../App.css';
import './Home.css';
import HomeBanner from '../Components/HomeBanner';
import React from 'react';

const Home = () => {
  return (
    <div className="pageWrapper">
      <img src={logo} className="App-logo" alt="logo" />
      <p className="homeText">Go-Scrobble is an open source music scrobbling service written in Go and React.</p>
      <HomeBanner />
    </div>
    );
}

export default Home;
