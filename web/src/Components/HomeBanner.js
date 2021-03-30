import React from 'react';
import '../App.css';
import './HomeBanner.css';
import { getStats } from '../Actions/api';
import ClipLoader from 'react-spinners/ClipLoader'

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
    .then((data) => {
      this.setState({
        isLoading: false,
        userCount: data.users,
        scrobbleCount: data.scrobbles,
        trackCount: data.tracks,
        artistCount: data.artists,
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
      <div className="homeBanner">
        <div className="homeBannerItem">
          {this.state.isLoading
            ? <ClipLoader color="#6AD7E5" size={36} />
            : <span className="homeBannerItemCount">{this.state.scrobbleCount}</span>}<br/>Scrobbles
        </div>
        <div className="homeBannerItem">
          {this.state.isLoading
            ? <ClipLoader color="#6AD7E5" size={36} />
            : <span className="homeBannerItemCount">{this.state.userCount}</span>}<br/>Users
        </div>
        <div className="homeBannerItem">
          {this.state.isLoading
            ? <ClipLoader color="#6AD7E5" size={36} />
            : <span className="homeBannerItemCount">{this.state.trackCount}</span>}<br/>Tracks
        </div>
        <div className="homeBannerItem">
          {this.state.isLoading
            ? <ClipLoader color="#6AD7E5" size={36} />
            : <span className="homeBannerItemCount">{this.state.artistCount}</span>}<br/>Artists
        </div>
      </div>
    );
  }
}

export default HomeBanner;
