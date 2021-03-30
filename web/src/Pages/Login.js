import React from 'react';
import '../App.css';
import './Login.css';
import { Button } from 'reactstrap';
import { Formik, Form, Field } from 'formik';
import ScaleLoader from 'react-spinners/ScaleLoader';
import { connect } from 'react-redux';
import { login } from '../Actions/auth';
import eventBus from "../Actions/eventBus";
import { LOGIN_SUCCESS } from '../Actions/types';

class Login extends React.Component {
  constructor(props) {
    super(props);
    this.state = {username: '', password: '', loading: false};
  }

  componentDidMount() {
    const { history, isLoggedIn } = this.props;

    if (isLoggedIn) {
      history.push("/dashboard")
    }
  }

  handleLogin(values) {
    this.setState({loading: true});

    const { dispatch, history } = this.props;

    dispatch(login(values.username, values.password))
      .then(() => {
        this.setState({
          loading: false,
          isLoggedIn: true
        });

        eventBus.dispatch(LOGIN_SUCCESS, { isLoggedIn: true });
        history.push("/dashboard");
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
        <h1>
          Login
        </h1>
        <div className="loginBody">
          <Formik
            initialValues={{ username: '', password: '' }}
            onSubmit={async values => this.handleLogin(values)}
          >
            <Form>
            <label>
              Email / Username<br/>
              <Field
                name="username"
                type="text"
                required={trueBool}
                className="loginFields"
              />
            </label>
            <br/>
            <label>
              Password<br/>
              <Field
                name="password"
                type="password"
                required={trueBool}
                className="loginFields"
              />
            </label>
            <br/><br/>
            <Button
              color="primary"
              type="submit"
              className="loginButton"
              disabled={this.state.loading}
            >{this.state.loading ? <ScaleLoader color="#FFF" size={35} /> : "Login"}</Button>
          </Form>
          </Formik>
        </div>
      </div>
    );
  }
}

function mapStateToProps(state) {
  const { isLoggedIn } = state.auth;
  return {
    isLoggedIn
  };
}

export default connect(mapStateToProps)(Login);
