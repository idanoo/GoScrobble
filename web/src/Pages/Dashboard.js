import React, { useState, useEffect, useContext } from 'react';
import '../App.css';
import './Dashboard.css';
import { useHistory } from "react-router";
import { getRecentScrobbles } from '../Api/index';
import ScaleLoader from 'react-spinners/ScaleLoader';
import ScrobbleTable from "../Components/ScrobbleTable";
import AuthContext from '../Contexts/AuthContext';

const Dashboard = () => {
  const history = useHistory();
  let { user } = useContext(AuthContext);
  let [loading, setLoading] = useState(true);
  let [dashboardData, setDashboardData] = useState({});

  useEffect(() => {
    if (!user) {
      return
    }
    getRecentScrobbles(user.uuid)
      .then(data => {
        setDashboardData(data);
        setLoading(false);
      })
  }, [user])


  if (!user) {
    history.push("/login")
  }

  if (loading) {
    return (
      <div className="pageWrapper">
        <ScaleLoader color="#6AD7E5" />
      </div>
    )
  }

  return (
    <div className="pageWrapper">
      <h1>
        {user.username}'s Dashboard!
      </h1>
      <div className="dashboardBody">
      {loading
        ? <ScaleLoader color="#6AD7E5" size={60} />
        : <ScrobbleTable data={dashboardData.items} />
      }
      </div>
    </div>
  );
}

export default Dashboard;
