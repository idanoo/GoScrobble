import React, { useEffect, useState } from 'react';
import '../App.css';
import './HomeBanner.css';
import { getStats } from '../Api/index';
import ClipLoader from 'react-spinners/ClipLoader'

const HomeBanner = () => {
  let [bannerData, setBannerData] = useState({});
  let [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    getStats()
      .then(data => {
        if (data.users) {
          setBannerData(data);
          setIsLoading(false);
        }
      })
  }, [])

  return (
    <div className="homeBanner">
      <div className="homeBannerItem">
        {isLoading
          ? <ClipLoader color="#6AD7E5" size={34} />
          : <span className="homeBannerItemCount">{bannerData.scrobbles}</span>
        }
        <br/>Scrobbles
      </div>
      <div className="homeBannerItem">
        {isLoading
          ? <ClipLoader color="#6AD7E5" size={34} />
          : <span className="homeBannerItemCount">{bannerData.users}</span>
        }
        <br/>Users
      </div>
      <div className="homeBannerItem">
        {isLoading
          ? <ClipLoader color="#6AD7E5" size={34} />
          : <span className="homeBannerItemCount">{bannerData.tracks}</span>
        }
        <br/>Tracks
      </div>
      <div className="homeBannerItem">
        {isLoading
          ? <ClipLoader color="#6AD7E5" size={34} />
          : <span className="homeBannerItemCount">{bannerData.artists}</span>
        }
        <br/>Artists
      </div>
    </div>
  );
}

export default HomeBanner;
