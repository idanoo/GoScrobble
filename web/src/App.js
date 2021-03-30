import './App.css';
import Home from './Pages/Home';
import About from './Pages/About';

import Dashboard from './Pages/Dashboard';
import Admin from './Pages/Admin';
import Profile from './Pages/Profile';
import Login from './Pages/Login';
import Settings from './Pages/Settings';
import Register from './Pages/Register';
import Navigation from './Components/Navigation';

import { logout } from './Actions/auth';
import { Route, Switch, withRouter } from 'react-router-dom';
import { connect } from 'react-redux';
import { Component } from 'react';
import 'bootstrap/dist/css/bootstrap.min.css';

function mapStateToProps(state) {
  const { user } = state.auth;
  return {
    user,
  };
}

class App extends Component {
  constructor(props) {
    super(props);
    this.logOut = this.logOut.bind(this);

    this.state = {
      // showAdminBoard: false,
      currentUser: undefined,
      // Don't even ask.. apparently you can't pass
      // exact="true".. it has to be a bool :|
      true: true,
    };
  }

  componentDidMount() {
    const user = this.props.user;

    if (user) {
      this.setState({
        currentUser: user,
        // showAdminBoard: user.roles.includes("ROLE_ADMIN"),
      });
    }
  }

  logOut() {
    this.props.dispatch(logout());
  }

  render() {
    // const { currentUser, showAdminBoard } = this.state;
    return (
      <div>
        <Navigation />
        <Switch>
          <Route exact={this.state.true} path="/" component={Home} />
          <Route path="/about" component={About} />

          <Route path="/dashboard" component={Dashboard} />
          <Route path="/profile" component={Profile} />
          <Route path="/admin" component={Admin} />

          <Route path="/settings" component={Settings} />
          <Route path="/login" component={Login} />
          <Route path="/register" component={Register} />
        </Switch>
      </div>
    );
  }
}
export default withRouter(connect(mapStateToProps)(App));