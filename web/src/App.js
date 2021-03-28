import './App.css';
import Home from './Components/Pages/Home';
import About from './Components/Pages/About';
import Help from './Components/Pages/Help';
import Login from './Components/Pages/Login';
import Settings from './Components/Pages/Settings';
import Register from './Components/Pages/Register';
import Navigation from './Components/Navigation';

import { Route, Switch, withRouter } from 'react-router-dom';
import { connect } from "react-redux";
import '../node_modules/bootstrap/dist/css/bootstrap.min.css';

function mapStateToProps(state) {
  return {
    isLoggedIn: state
  };
}

function mapDispatchToProps(dispatch) {
  return {
    logIn: () => dispatch({type: true}),
    logOut: () => dispatch({type: false})
  };
}

const App = () => {
  let exact = true
  return (
    <div>
      <Navigation />
      <Switch>
        <Route exact={exact} path="/" component={Home} />
        <Route path="/about" component={About} />
        <Route path="/settings" component={Settings} />
        <Route path="/help" component={Help} />
        <Route path="/login" component={Login} />
        <Route path="/register" component={Register} />
      </Switch>
    </div>
  );
}

export default withRouter(connect(mapStateToProps, mapDispatchToProps)(App));
