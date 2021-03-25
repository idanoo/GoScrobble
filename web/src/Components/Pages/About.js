import logo from '../../logo.svg';
import '../../App.css';

function About() {
  return (
    <div>
      <img src={logo} className="App-logo" alt="logo" />
      <p>
        TEST!!!
      </p>
      <a
        className="App-link"
        href="https://git.m2.nz/idanoo/go-scrobble"
        target="_blank"
        rel="noopener noreferrer"
      >
        TEST@!@@
      </a>
    </div>
  );
}

export default About;
