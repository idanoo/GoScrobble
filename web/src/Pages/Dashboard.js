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
  let [isLoading, setIsLoading] = useState(true);
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
        setIsLoading(false);
      })
  }, [user])

  return (
    <div className="pageWrapper">
      <h1>
        Dashboard!
      </h1>
      {isLoading
        ? <ScaleLoader color="#FFF" size={60} />
        : <ScrobbleTable data={dashboardData} />
      }
    </div>
  );
}

export default Dashboard;
