import './App.css';
import Home from './Components/Pages/Home';
import About from './Components/Pages/About';
import Navigation from './Components/Pages/Navigation';
import { Route, Switch, HashRouter } from 'react-router-dom';
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

const App = (props) => {
  return (
    <HashRouter>
      <div className="wrapper">
        <Navigation />
        <Switch>
          <Route exact path="/" component={Home} />
          <Route path="/about" component={About} />
        </Switch>
      </div>
      </HashRouter>
  );
}

export default connect(mapStateToProps, mapDispatchToProps)(App);
