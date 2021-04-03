import React, { useContext, useState, useEffect } from 'react';
import '../App.css';
import './Register.css';
import { Button } from 'reactstrap';
import ScaleLoader from "react-spinners/ScaleLoader";
import AuthContext from '../Contexts/AuthContext';
import { Formik, Field, Form } from 'formik';
import { useHistory } from 'react-router';
import { getServerInfo } from '../Api/index';
const Register = () => {
  const history = useHistory();
  let boolTrue = true;
  let { Register, user, loading } = useContext(AuthContext);
  let [serverInfo, setServerInfo] = useState({ registration_enabled: true });
  let [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    if (user) {
      return
    }
    getServerInfo()
      .then(data => {
        setServerInfo(data);
        setIsLoading(false);
      })
  }, [user])

  if (isLoading) {
    return (
      <div className="pageWrapper">
        <ScaleLoader color="#6AD7E5" />
      </div>
    )
  }

  if (user) {
    history.push("/dashboard");
  }

  return (
    <div className="pageWrapper">
      {
        serverInfo.registration_enabled !== "1" ?
        <p>Registration is temporarily disabled. Please try again soon!</p>
        :
        <div>
          <h1>
          Register
          </h1>
          <div className="registerBody">
            <Formik
              initialValues={{ username: '', email: '', password: '', passwordconfirm: '' }}
              onSubmit={async values => Register(values)}
            >
              <Form>
                <label>
                  Username*<br/>
                  <Field
                    name="username"
                    type="text"
                    required={boolTrue}
                    className="registerFields"
                  />
                </label>
                <br/>
                <label>
                  Email<br/>
                  <Field
                    name="email"
                    type="email"
                    className="registerFields"
                  />
                </label>
                <br/>
                <label>
                  Password*<br/>
                  <Field
                    name="password"
                    type="password"
                    required={boolTrue}
                    className="registerFields"
                  />
                </label>
                <br/>
                <label>
                  Confirm Password*<br/>
                  <Field
                    name="passwordconfirm"
                    type="password"
                    required={boolTrue}
                    className="registerFields"
                  />
                </label>
                <br/><br/>
                <Button
                  color="primary"
                  type="submit"
                  className="registerButton"
                  disabled={loading}
                  >{loading ? <ScaleLoader color="#FFF" size={35} /> : "Register"}</Button>
              </Form>
            </Formik>
          </div>
        </div>
      }
    </div>
  );
}

export default Register;
