import React from 'react';
import '../App.css';
import './Dashboard.css';

import { getRecentScrobbles } from '../Api/index';
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

export default Dashboard;
