import React from 'react';
import '../App.css';
import './HomeBanner.css';
import { getStats } from '../Actions/api';

class HomeBanner extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      isLoading: true,
      userCount: 0,
      scrobbleCount: 0,
      trackCount: 0,
      artistCount: 0,
    };
  }

  componentDidMount() {
    getStats()
    // .then((data) => {
    //   this.setState({
    //     loading: false,
    //     userCount: data.users,
    //     scrobbleCount: data.scrobbles,
    //     trackCount: data.tracks,
    //     artistCount: data.artists,
    //   });
    // })
    // .catch(() => {
    //   this.setState({
    //     loading: false
    //   });
    // });
  }

  render() {
    return (
      <div className="container">
      </div>
    );
  }
}

export default HomeBanner;
