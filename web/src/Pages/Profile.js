import React from 'react';
import '../App.css';
import './Dashboard.css';
import { connect } from 'react-redux';

class Profile extends React.Component {
  componentDidMount() {
    const { history } = this.props;
    const isLoggedIn = this.props.isLoggedIn;

    if (!isLoggedIn) {
      history.push("/login")
      window.location.reload()
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

function mapStateToProps(state) {
  const { isLoggedIn } = state.auth;
  return {
    isLoggedIn,
  };
}

export default connect(mapStateToProps)(Profile);