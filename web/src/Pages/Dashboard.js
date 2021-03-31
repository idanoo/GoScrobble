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

  if (!user) {
    history.push("/login");
  }

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

  return (
    <div className="pageWrapper">
      <h1>
        Dashboard!
      </h1>
      {loading
        ? <ScaleLoader color="#6AD7E5" size={60} />
        : <ScrobbleTable data={dashboardData} />
      }
    </div>
  );
}

export default Dashboard;
