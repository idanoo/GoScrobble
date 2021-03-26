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

class Register extends React.Component {
  constructor(props) {
    super(props);
    this.state = {username: '', email: '', password: '', passwordconfirm: '', loading: false};
    this.handleUsernameChange = this.handleUsernameChange.bind(this);
    this.handleEmailChange = this.handleEmailChange.bind(this);
    this.handlePasswordChange = this.handlePasswordChange.bind(this);
    this.handlePasswordConfirmChange = this.handlePasswordConfirmChange.bind(this);
    this.handleSubmit = this.handleSubmit.bind(this);
  }

  handleUsernameChange(event) {
    this.setState({username: event.target.value});
  }

  handleEmailChange(event) {
    this.setState({email: event.target.value});
  }

  handlePasswordChange(event) {
    this.setState({password: event.target.value});
  }

  handlePasswordConfirmChange(event) {
    this.setState({passwordconfirm: event.target.value});
  }

  handleSubmit(event) {
    event.preventDefault();

    if (this.state.password !== this.state.passwordconfirm) {
      this.props.addToast('Passwords do not match', { appearance: 'error' });
      return
    }

    if (this.state.password.len < 8) {
      this.props.addToast('Passwords do not match', { appearance: 'error' });
      return
    }

    this.setState({loading: true});
    const requestOptions = {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        username: this.state.username,
        email: this.state.email,
        password: this.state.password,
      })
    };

    const apiUrl = 'http://127.0.0.1:42069/api/v1/register';
    fetch(apiUrl, requestOptions)
      .then((response) => response.json())
      .then((function(data) {
        if (data.error) {
          this.props.addToast(data.error, { appearance: 'error' });
        } else {
          this.props.addToast(data.message, { appearance: 'success' });
        }
        this.setState({loading: false});
    }).bind(this));
  }

  render() {
    return (
      <div className="pageWrapper">
        <h1>
        Register
        </h1>
        <div className="loginBody">
          <form onSubmit={this.handleSubmit}>
            <label>
              Username*<br/>
              <input
                type="text"
                required="true"
                className="loginFields"
                value={this.state.username}
                onChange={this.handleUsernameChange}
              />
            </label>
            <br/>
            <label>
              Email<br/>
              <input
                type="email"
                className="loginFields"
                value={this.state.email}
                onChange={this.handleEmailChange}
              />
            </label>
            <br/>
            <label>
              Password<br/>
              <input
                type="password"
                required="true"
                className="loginFields"
                value={this.state.password}
                onChange={this.handlePasswordChange}
              />
            </label>
            <br/>
            <label>
              Password<br/>
              <input
                type="password"
                required="true"
                className="loginFields"
                value={this.state.passwordconfirm}
                onChange={this.handlePasswordConfirmChange}
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

export default withToast(Register);
