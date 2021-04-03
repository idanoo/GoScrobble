import React, { useContext } from 'react';
import '../App.css';
import './Login.css';
import { Button } from 'reactstrap';
import { Formik, Form, Field } from 'formik';
import ScaleLoader from 'react-spinners/ScaleLoader';
import AuthContext from '../Contexts/AuthContext';
import { useHistory } from "react-router";

const Login = () => {
  const history = useHistory();
  let boolTrue = true;
  let { Login, loading, user } = useContext(AuthContext);

  if (user) {
    history.push("/dashboard");
  }

  const redirectReset = () => {
    history.push("/reset")
  }

  return (
    <div className="pageWrapper">
      <h1>
        Login
      </h1>
      <div className="pageBody">
        <Formik
          initialValues={{ username: '', password: '' }}
          onSubmit={values => Login(values)}
        >
          <Form>
          <label>
            Email / Username<br/>
            <Field
              name="username"
              type="text"
              required={boolTrue}
              className="loginFields"
            />
          </label>
          <br/>
          <label>
            Password<br/>
            <Field
              name="password"
              type="password"
              className="loginFields"
            />
          </label>
          <br/><br/>
          <Button
            color="primary"
            type="submit"
            className="loginButton"
            disabled={loading}
          >{loading ? <ScaleLoader color="#FFF" size={35} /> : "Login"}</Button>
          <br/><br/>
            <Button
            color="secondary"
            type="button"
            className="loginButton"
            onClick={redirectReset}
            disabled={loading}
          >Reset Password</Button>
        </Form>
        </Formik>
      </div>
    </div>
  );
}

export default Login;
