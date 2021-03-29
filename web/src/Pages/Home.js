import logo from '../logo.png';
import '../App.css';

function Home() {
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
    </div>
  );
}

export default Home;
