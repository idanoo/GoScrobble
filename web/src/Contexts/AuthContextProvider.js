import React, { useState, useEffect } from 'react';
import { toast } from 'react-toastify';
import AuthContext from './AuthContext';
import { PostLogin, PostRegister } from '../Api/index';

const AuthContextProvider = ({ children }) => {
  const [user, setUser] = useState();
  const [loading, setLoading] = useState(false);

  useEffect(() => {
    setLoading(true)
    const user = JSON.parse(localStorage.getItem('user'));
    if (user && user.jwt) {
      setUser(user)
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
      // Do stuff here?
      setLoading(false);
    });
  };

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
        loading,
        user,
      }}
    >
      {children}
    </AuthContext.Provider>
  );
};

export default AuthContextProvider;