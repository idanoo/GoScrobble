import React from 'react';
import '../App.css';
import './Register.css';
import { Button } from 'reactstrap';
import ScaleLoader from "react-spinners/ScaleLoader";
import register from '../Contexts/AuthContextProvider';
import { Formik, Field, Form } from 'formik';

class Register extends React.Component {
  constructor(props) {
    super(props);
    this.state = {username: '', email: '', password: '', passwordconfirm: '', loading: false};
  }

  componentDidMount() {
    const { history, isLoggedIn } = this.props;

    if (isLoggedIn) {
      history.push("/dashboard")
    }
  }

  handleRegister(values) {
    console.log(values)
    this.setState({loading: true});

    const { dispatch, history } = this.props;

    dispatch(register(values.username, values.email, values.password))
      .then(() => {
        this.setState({
          loading: false,
        });
        history.push("/login");
      })
      .catch(() => {
        this.setState({
          loading: false
        });
      });
  }

  render() {
    let trueBool = true;
    return (
      <div className="pageWrapper">
        {
          // TODO: Move to DB:config REGISTRATION_DISABLED=1|0
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
                onSubmit={async values => this.handleRegister(values)}
              >
                <Form>
                  <label>
                    Username*<br/>
                    <Field
                      name="username"
                      type="text"
                      required={trueBool}
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
                      required={trueBool}
                      className="registerFields"
                    />
                  </label>
                  <br/>
                  <label>
                    Confirm Password*<br/>
                    <Field
                      name="passwordconfirm"
                      type="password"
                      required={trueBool}
                      className="registerFields"
                    />
                  </label>
                  <br/><br/>
                  <Button
                    color="primary"
                    type="submit"
                    className="registerButton"
                    disabled={this.state.loading}
                    >{this.state.loading ? <ScaleLoader color="#FFF" size={35} /> : "Register"}</Button>
                </Form>
              </Formik>
            </div>
          </div>
        }
      </div>
    );
  }
}

export default Register;
