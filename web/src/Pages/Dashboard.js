import React from 'react';
import '../App.css';
import './Dashboard.css';
import { connect } from 'react-redux';

class Dashboard extends React.Component {
  componentDidMount() {
    const { history } = this.props;
    const isLoggedIn = this.props.isLoggedIn;

    if (!isLoggedIn) {
      history.push("/login")
    }
  }

  render() {
    return (
      <div className="pageWrapper">
        <h1>
          Dashboard!
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

export default connect(mapStateToProps)(Dashboard);
