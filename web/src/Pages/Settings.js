import React from 'react';
import '../App.css';
import './Settings.css';

class Settings extends React.Component {
  constructor(props) {
    super(props);
    this.state = {username: '', password: '', loading: false};
  }

  render() {
    return (
      <div className="pageWrapper">
        <h1>
          Settings
        </h1>
        <div className="loginBody">
        <p>
          All the settings
        </p>
        </div>
      </div>
    );
  }
}

export default Settings;
