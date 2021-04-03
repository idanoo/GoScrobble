import React, { useState, useEffect, useContext } from 'react';
import '../App.css';
import './Reset.css';
import { Button } from 'reactstrap';
import { Formik, Form, Field } from 'formik';
import ScaleLoader from 'react-spinners/ScaleLoader';
import { validateResetPassword, sendPasswordReset } from '../Api/index';
import AuthContext from '../Contexts/AuthContext';

const Reset = (route) => {
  let boolTrue = true;
  const [loading, setLoading] = useState(true);
  const [reset, setReset] = useState({});
  const [sent, setSent] = useState(false);

  let { ResetPassword } = useContext(AuthContext);

  let reqToken = false;
  if (route && route.match && route.match.params && route.match.params.token) {
    reqToken = route.match.params.token
  }

  const sendReset = (values) => {
    sendPasswordReset(values).then(() => {
      setSent(true);
    });
  }

  useEffect(() => {
    if (!reqToken) {
      setLoading(false);
      return false;
    }

    validateResetPassword(reqToken)
      .then(data => {
        setReset(data);
        console.log(data)
        setLoading(false);
      })
  }, [reqToken])

  if (loading) {
    return (
      <div className="pageWrapper">
        <ScaleLoader color="#6AD7E5" />
      </div>
    )
  }

  if (sent) {
    return (
      <div className="pageWrapper">
        <h1>
          Check your email!
        </h1>
      </div>
    )
  }

  if (!reqToken) {
    return (
      <div className="pageWrapper">
        <h1>
          Reset Password
        </h1>
        <div className="pageBody">
          <Formik
            initialValues={{ email: '' }}
            onSubmit={values => sendReset(values)}
          >
            <Form>
            <label>
              Email<br/>
              <Field
                name="email"
                type="email"
                required={boolTrue}
                className="loginFields"
              />
            </label>
            <br/><br/>
            <Button
              color="primary"
              type="submit"
              className="loginButton"
              disabled={loading}
            >{loading ? <ScaleLoader color="#FFF" size={35} /> : "Reset"}</Button>
          </Form>
          </Formik>
        </div>
      </div>
    )
  }

  if (reqToken && !reset.valid) {
    return (
      <div className="pageWrapper">
        Invalid Reset Token or Token expired
      </div>
    )
  }

  return (
    <div className="pageWrapper">
      <h1>
        Reset Password
      </h1>
      <div className="resetBody">
        <Formik
          initialValues={{ password: '', comfirmpassword: '', token: reqToken }}
          onSubmit={values => ResetPassword(values)}
        >
          <Form>
          <label>
            New Password<br/>
            <Field
              name="password"
              type="password"
              required={boolTrue}
              className="resetFields"
            />
          </label>
          <br/>
          <label>
            Confirm New Password<br/>
            <Field
              name="comfirmpassword"
              type="password"
              required={boolTrue}
              className="resetFields"
            />
          </label>
          <Field
              name="token"
              type="hidden"
              className="resetFields"
            />
          <br/><br/>
          <Button
            color="primary"
            type="submit"
            className="loginButton"
            disabled={loading}
          >{loading ? <ScaleLoader color="#FFF" size={35} /> : "Reset"}</Button>
        </Form>
        </Formik>
      </div>
    </div>
  );
}

export default Reset;
