import '../App.css';
import './About.css';

const About = () => {
  return (
    <div className="pageWrapper">
      <h1>
        About GoScrobble.com
      </h1>
      <p className="aboutBody">
        Go-Scrobble is an open source music scrobbling service written in Go and React.<br/>
        Used to track your listening history and build a profile to discover new music.
      </p>
      <a
        className="pageBody"
        href="https://gitlab.com/idanoo/goscrobble"
        target="_blank"
        rel="noopener noreferrer"
      >gitlab.com/idanoo/goscrobble
      </a>
    </div>
  );
}

export default About;
