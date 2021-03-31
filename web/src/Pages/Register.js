import React, { useContext } from 'react';
import '../App.css';
import './Register.css';
import { Button } from 'reactstrap';
import ScaleLoader from "react-spinners/ScaleLoader";
import AuthContext from '../Contexts/AuthContext';
import { Formik, Field, Form } from 'formik';
import { useHistory } from "react-router";

const Register = () => {
  const history = useHistory();
  let boolTrue = true;
  let { Register, user, loading } = useContext(AuthContext);

  if (user) {
    history.push("/dashboard");
  }

  return (
    <div className="pageWrapper">
      {
        // TODO: Move to DB:config REGISTRATION_DISABLED=1|0 :upsidedownsmile:
        process.env.REACT_APP_REGISTRATION_DISABLED === "true" ?
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
