import React from 'react';
import '../App.css';
import './Dashboard.css';
import { connect } from 'react-redux';
import { getRecentScrobbles } from '../Actions/api';
import ScaleLoader from 'react-spinners/ScaleLoader';
import ScrobbleTable from "../Components/ScrobbleTable";

class Dashboard extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      isLoading: true,
      scrobbleData: [],
      uuid: null,
    };
  }

  componentDidMount() {
    const { history, uuid } = this.props;
    const isLoggedIn = this.props.isLoggedIn;

    if (!isLoggedIn) {
      history.push("/login")
    }

    getRecentScrobbles(uuid)
    .then((data) => {
      this.setState({
        isLoading: false,
        data: data
      });
    })
    .catch(() => {
      this.setState({
        isLoading: false
      });
    });
  }

  render() {
    return (
      <div className="pageWrapper">
        <h1>
          Dashboard!
        </h1>
        {this.state.isLoading
          ? <ScaleLoader color="#FFF" size={60} />
          : <ScrobbleTable data={this.state.data} />
        }
      </div>
    );
  }
}

function mapStateToProps(state) {
  const { isLoggedIn } = state.auth;
  const { uuid } = state.auth.user;

  return {
    isLoggedIn,
    uuid,
  };
}

export default connect(mapStateToProps)(Dashboard);
