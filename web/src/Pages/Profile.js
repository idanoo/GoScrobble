import React from 'react';
import '../App.css';
import './Dashboard.css';

class Profile extends React.Component {
  componentDidMount() {
    const { history, isLoggedIn } = this.props;
    
    if (!isLoggedIn) {
      history.push("/login")
    }
  }

  render() {
    return (
      <div className="pageWrapper">
        <h1>
          Hai User
        </h1>
      </div>
    );
  }
}

export default Profile;