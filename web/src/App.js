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

// import { logout } from './Actions/auth';
import { Route, Switch, withRouter } from 'react-router-dom';
import { Component } from 'react';
import 'bootstrap/dist/css/bootstrap.min.css';

class App extends Component {
  constructor(props) {
    super(props);
    this.state = {
      true: true,
    };
  }

  render() {
    return (
      <div>
        <Navigation />
        <Switch>
          <Route exact={this.state.true} path={["/", "/home"]} component={Home} />
          <Route path="/about" component={About} />

          <Route path="/dashboard" component={Dashboard} />
          <Route path="/profile" component={Profile} />
          <Route path="/settings" component={Settings} />

          <Route path="/admin" component={Admin} />

          <Route path="/login" component={Login} />
          <Route path="/register" component={Register} />

        </Switch>
      </div>
    );
  }
}

export default withRouter(App);