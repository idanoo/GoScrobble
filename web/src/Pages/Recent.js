import React, { useState, useEffect } from 'react';
import '../App.css';
import './Recent.css';
import RecentScrobbleTable from '../Components/RecentScrobbleTable'
import { getRecentScrobbles } from '../Api/index'
import ScaleLoader from 'react-spinners/ScaleLoader';

const Recent = (route) => {
  const [loading, setLoading] = useState(true);
  const [data, setData] = useState({});

  useEffect(() => {
    if (loading) {
      getRecentScrobbles()
        .then(data => {
          setData(data);
          setLoading(false);
        })
      }

  }, [loading])

  if (loading) {
    return (
      <div className="pageWrapper">
        <ScaleLoader color="#6AD7E5" />
      </div>
    )
  }

  return (
    <div className="pageWrapper">
      Last 50 Scrobbles
      <RecentScrobbleTable data={data.items} />
    </div>
  );
}

export default Recent;