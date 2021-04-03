import React, { useContext, useState, useEffect } from 'react';
import { useHistory } from 'react-router-dom';
import '../App.css';
import './Admin.css';
import { Button } from 'reactstrap';
import { Formik, Form, Field } from 'formik';
import ScaleLoader from 'react-spinners/ScaleLoader';
import AuthContext from '../Contexts/AuthContext';
import { Switch } from 'formik-material-ui';
import { getConfigs, postConfigs } from '../Api/index'

const Admin = () => {
  const history = useHistory();
  const { user } = useContext(AuthContext);
  const [loading, setLoading] = useState(true);
  const [configs, setConfigs] = useState({})
  const [toggle, setToggle] = useState(false);

  useEffect(() => {
    getConfigs()
      .then(data => {
        if (data.configs) {
          setConfigs(data.configs);
          setToggle(data.configs.REGISTRATION_ENABLED === "1")
          setLoading(false);
        }
      })
  }, [])

  const handleToggle = () => {
    setToggle(!toggle);
  };

  if (!user) {
    history.push("/login")
  }

  if (user && !user.admin) {
    history.push("/Dashboard")
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
        Admin Panel
      </h1>
      <div className="pageBody">
        <Formik
          initialValues={configs}
          onSubmit={(values) => postConfigs(values, toggle)}
        >
          <Form><br/>
          <label>
            <Field
              type="checkbox"
              name="REGISTRATION_ENABLED"
              onChange={handleToggle}
              component={Switch}
              checked={toggle}
              value={toggle}
            />
            Registration Enabled
          </label><br/><br/>
          <label>
            LastFM Api Key<br/>
            <Field
              name="LASTFM_API_KEY"
              type="text"
              className="loginFields"
            />
          </label>
          <br/>
          <label>
            Spotify App ID<br/>
            <Field
              name="SPOTIFY_APP_ID"
              type="text"
              className="loginFields"
            />
          </label>
          <label>
            Spotify App Secret<br/>
            <Field
              name="SPOTIFY_APP_SECRET"
              type="text"
              className="loginFields"
            />
          </label>
          <br/><br/>
          <Button
            color="primary"
            type="submit"
            className="loginButton"
            disabled={loading}
          >Update</Button>
        </Form>
        </Formik>
      </div>
    </div>
  );
}

export default Admin;
