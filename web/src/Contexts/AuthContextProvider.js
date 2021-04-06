import React, { useState, useEffect } from 'react';
import { toast } from 'react-toastify';
import AuthContext from './AuthContext';
import { PostLogin, PostRegister, PostResetPassword, PostRefreshToken } from '../Api/index';

const AuthContextProvider = ({ children }) => {
  const [user, setUser] = useState();
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    setLoading(true)
    let curTime = Math.round((new Date()).getTime() / 1000);
    let user = JSON.parse(localStorage.getItem('user'));

    // Confirm JWT is set.
    if (user && user.jwt) {
      // Check refresh expiry is valid.
      if (user.refresh_exp && (user.refresh_exp > curTime)) {
        // Check if JWT is still valid
        if (user.exp < curTime) {
          // Refresh if not
          user = RefreshToken(user.refresh_token)
          localStorage.setItem('user', JSON.stringify(user));
        }

        setUser(user)
      }
    }
    setLoading(false)
  }, []);

  const Login = (formValues) => {
    setLoading(true);
    PostLogin(formValues).then(user => {
      if (user) {
        setUser(user);
        localStorage.setItem('user', JSON.stringify(user));
      }
      setLoading(false);
    })
  }

  const Register = (formValues) => {
    setLoading(true);
    return PostRegister(formValues).then(response => {
      setLoading(false);
    });
  };

  const ResetPassword = (formValues) => {
    return PostResetPassword(formValues);
  }

  const RefreshToken = (refreshToken) => {
    return PostRefreshToken(refreshToken);
  }

  const Logout = () => {
    localStorage.removeItem("user");
    setUser(null)
    toast.success('Successfully logged out.');
  };

  return (
    <AuthContext.Provider
      value={{
        Logout,
        Login,
        Register,
        ResetPassword,
        RefreshToken,
        loading,
        user,
      }}
    >
      {children}
    </AuthContext.Provider>
  );
};

export default AuthContextProvider;