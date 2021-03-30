import logo from '../logo.png';
import '../App.css';
import HomeBanner from '../Components/HomeBanner';
import React from 'react';

class Home extends React.Component {
  render() {
    return (
    <div className="App-header">
      <img src={logo} className="App-logo" alt="logo" />
      <p>
        goscrobble.com
      </p>
      <a
        className="App-link"
        href="https://gitlab.com/idanoo/go-scrobble"
        target="_blank"
        rel="noopener noreferrer"
      >
        gitlab.com/idanoo/go-scrobble
      </a>
      <HomeBanner />
    </div>
    );
  }
}

export default Home;
