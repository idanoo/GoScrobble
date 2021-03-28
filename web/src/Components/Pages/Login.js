import React from 'react';
import '../../App.css';
import './Login.css';
import { Button } from 'reactstrap';
import { Formik, Form, Field } from 'formik';
import { useToasts } from 'react-toast-notifications';
import ScaleLoader from "react-spinners/ScaleLoader";

function withToast(Component) {
  return function WrappedComponent(props) {
    const toastFuncs = useToasts()
    return <Component {...props} {...toastFuncs} />;
  }
}

class Login extends React.Component {
  constructor(props) {
    super(props);
    this.state = {username: '', password: '', loading: false};
    this.handleUsernameChange = this.handleUsernameChange.bind(this);
    this.handlePasswordChange = this.handlePasswordChange.bind(this);
    this.handleSubmit = this.handleSubmit.bind(this);
  }

  handleUsernameChange(event) {
    this.setState({username: event.target.value});
  }

  handlePasswordChange(event) {
    this.setState({password: event.target.value});
  }

  handleSubmit(values) {
    this.setState({loading: true});
    const requestOptions = {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      timeout: 5000,
      body: JSON.stringify({
        username: values.username,
        password: values.password,
      })
  };
    const apiUrl = process.env.REACT_APP_API_URL + '/api/v1/login';
    fetch(apiUrl, requestOptions)
      .then((response) => {
        if (response.status === 429) {
          this.props.addToast("Rate limited. Please try again soon", { appearance: 'error' });
          return "{}"
        } else {
          return response.json()
        }
      })
      .then((function(data) {
        if (data.error) {
          this.props.addToast(data.error, { appearance: 'error' });
        } else if (data.token) {
          this.props.addToast(data.token, { appearance: 'success' });
        }
        this.setState({loading: false});
      }).bind(this))
      .catch(() => {
          this.props.addToast('Error submitting form. Please try again', { appearance: 'error' });
          this.setState({loading: false});
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
            onSubmit={async values => this.handleSubmit(values)}
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

export default withToast(Login);
