import React from 'react';
import '../App.css';
import './Login.css';
import { Button } from 'reactstrap';
import ScaleLoader from "react-spinners/ScaleLoader";
import { withRouter } from 'react-router-dom'

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

    // if (this.state.password.len < 8) {
    //   this.props.addToast('Password must be at least 8 characters', { appearance: 'error' });
    //   return
    // }

    this.setState({loading: true});
    const requestOptions = {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      timeout: 5000,
      body: JSON.stringify({
        username: this.state.username,
        email: this.state.email,
        password: this.state.password,
      })
    };

    const apiUrl = process.env.REACT_APP_API_URL + '/api/v1/register';
    console.log(apiUrl);
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
        console.log(data);
        if (data.error) {
          this.props.addToast(data.error, { appearance: 'error' });
        } else if (data.message) {
          this.props.addToast(data.message, { appearance: 'success' });
          this.props.history.push('/login')
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
        {
          // TODO: Move to DB:config REGISTRATION_DISABLED=1|0
          process.env.REACT_APP_REGISTRATION_DISABLED === "true" ?
          <p>Registration is temporarily disabled. Please try again soon!</p>
          :
          <div>
            <h1>
            Register
            </h1>
            <div className="loginBody">
              <form onSubmit={this.handleSubmit}>
                <label>
                  Username*<br/>
                  <input
                    type="text"
                    required={trueBool}
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
                  Password*<br/>
                  <input
                    type="password"
                    required={trueBool}
                    className="loginFields"
                    value={this.state.password}
                    onChange={this.handlePasswordChange}
                  />
                </label>
                <br/>
                <label>
                  Password*<br/>
                  <input
                    type="password"
                    required={trueBool}
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
                  >{this.state.loading ? <ScaleLoader color="#FFF" size={35} /> : "Register"}</Button>
              </form>
            </div>
          </div>
        }
      </div>
    );
  }
}

export default withRouter(Register);
