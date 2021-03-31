import React, { useContext } from 'react';
import '../App.css';
import './Dashboard.css';
import { useHistory } from "react-router";
import AuthContext from '../Contexts/AuthContext';

const Profile = () => {
  const history = useHistory();
  const { user } = useContext(AuthContext);

  if (!user) {
    history.push("/login");
  }

  return (
    <div className="pageWrapper">
      <h1>
        Welcome {user.username}!
      </h1>
    </div>
  );

}

export default Profile;