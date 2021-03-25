import logo from '../../logo.svg';
import '../../App.css';

function Home() {
  return (
    <div className="App-header">
      <img src={logo} className="App-logo" alt="logo" />
      <p>
        goscrobble.com
      </p>
      <a
        className="App-link"
        href="https://git.m2.nz/idanoo/go-scrobble"
        target="_blank"
        rel="noopener noreferrer"
      >
        git.m2.nz/idanoo/go-scrobble
      </a>
    </div>
  );
}

export default Home;
