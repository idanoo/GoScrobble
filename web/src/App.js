import './App.css';
import Home from './Components/Pages/Home';
import About from './Components/Pages/About';
import Login from './Components/Pages/Login';
import Register from './Components/Pages/Register';
import Navigation from './Components/Pages/Navigation';

import { Route, Switch, HashRouter } from 'react-router-dom';
import { connect } from "react-redux";

import { ToastProvider } from 'react-toast-notifications';

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
    <HashRouter>
      <ToastProvider autoDismiss="true" autoDismissTimeout="5000" placement="bottom-right">
        <Navigation />
        <Switch>
          <Route exact={exact} path="/" component={Home} />
          <Route path="/about" component={About} />
          <Route path="/login" component={Login} />
          <Route path="/register" component={Register} />
        </Switch>
      </ToastProvider>
    </HashRouter>
  );
}

export default connect(mapStateToProps, mapDispatchToProps)(App);
