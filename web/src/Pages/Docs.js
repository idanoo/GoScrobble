import '../App.css';
import './Docs.css';

const Docs = () => {
  return (
    <div className="pageWrapper">
      <h1>
        Documentation
      </h1>
      <p className="aboutBody">
        Go-Scrobble is an open source music scorbbling service written in Go and React.<br/>
        Used to track your listening history and build a profile to discover new music.
      </p>
      <a
        className="pageBody"
        href="https://gitlab.com/idanoo/go-scrobble"
        target="_blank"
        rel="noopener noreferrer"
      >gitlab.com/idanoo/go-scrobble
      </a>
    </div>
  );
}

export default Docs;
