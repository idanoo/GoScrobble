import React from 'react';
import '../../App.css';
import './Login.css';
import { Button } from 'reactstrap';

import { useToasts } from 'react-toast-notifications';

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

  handleSubmit(event) {
    event.preventDefault();
    this.setState({loading: true});
    const requestOptions = {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      timeout: 5000,
      body: JSON.stringify({
        username: this.state.username,
        password: this.state.password,
      })
  };
    const apiUrl = process.env.REACT_APP_API_URL + '/api/v1/login';
    fetch(apiUrl, requestOptions)
      .then((response) => response.json())
      .then((function(data) {
        if (data.error) {
          this.props.addToast(data.error, { appearance: 'error' });
        } else {
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
    return (
      <div className="pageWrapper">
        <h1>
          Login
        </h1>
        <div className="loginBody">
          <form onSubmit={this.handleSubmit}>
            <label>
              Email / Username<br/>
              <input
                type="text"
                className="loginFields"
                value={this.state.username}
                onChange={this.handleUsernameChange}
              />
            </label>
            <br/>
            <label>
              Password<br/>
              <input
                type="password"
                className="loginFields"
                value={this.state.password}
                onChange={this.handlePasswordChange}
              />
            </label>
            <br/><br/>
            <Button
              color="primary"
              type="submit"
              className="loginButton"
              disabled={this.state.loading}
            >Login</Button>
          </form>
        </div>
      </div>
    );
  }
}

export default withToast(Login);
