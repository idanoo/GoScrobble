import React from 'react';
import '../../App.css';
import './Login.css';
import { Button } from 'reactstrap';

import { ToastProvider, useToasts } from 'react-toast-notifications';

// const FormWithToasts = () => {
//   const { addToast } = useToasts();

//   const onSubmit = async value => {
//     // const { error } = await dataPersistenceLayer(value);

//     if (error) {
//       addToast(error.message, { appearance: 'error' });
//     } else {
//       addToast('Saved Successfully', { appearance: 'success' });
//     }
//   };

//   return <form onSubmit={this.handleSubmit}>...</form>;
// };
// const { addToast } = useToasts();


class About extends React.Component {
  constructor(props) {
    super(props);
    this.state = {username: '', password: ''};
    this.handleUsernameChange = this.handleUsernameChange.bind(this);
    this.handlePasswordChange = this.handlePasswordChange.bind(this);
    this.handleSubmit = this.handleSubmit.bind(this);
  }
  handleUsernameChange(event) {
    this.setState({username: event.target.value});
    // addToast(error.message, { appearance: 'error' });

  }

  handlePasswordChange(event) {
    this.setState({password: event.target.value});
  }


  handleSubmit(event) {
    event.preventDefault();
    const requestOptions = {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        username: this.state.username,
        password: this.state.password,
      })
  };
    const apiUrl = 'http://127.0.0.1:42069/api/v1/login';
    fetch(apiUrl, requestOptions)
      .then((response) => response.json())
      .then((data) => function() {






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
                value={this.state.username}
                onChange={this.handleUsernameChange}
              />
            </label>
            <br/>
            <label>
              Password<br/>
              <input
                type="password"
                value={this.state.password}
                onChange={this.handlePasswordChange}
              />
            </label>
            <br/><br/>
            <Button
              color="primary"
              type="submit"
              className="loginButton"
            >Login</Button>
          </form>
        </div>
      </div>
    );
  }
}

export default About;
