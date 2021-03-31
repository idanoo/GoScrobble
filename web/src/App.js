import { Route, Switch, withRouter } from 'react-router-dom';

import Home from './Pages/Home';
import About from './Pages/About';
import Dashboard from './Pages/Dashboard';
import Profile from './Pages/Profile';
import Admin from './Pages/Admin';
import Login from './Pages/Login';
import Register from './Pages/Register';

import Navigation from './Components/Navigation';

import 'bootstrap/dist/css/bootstrap.min.css';
import './App.css';

const App = () => {
    let boolTrue = true

    return (
      <div>
        <Navigation />
        <Switch>
          <Route exact={boolTrue} path={["/", "/home"]} component={Home} />
          <Route path="/about" component={About} />

          <Route path="/dashboard" component={Dashboard} />
          <Route path="/profile" component={Profile} />

          <Route path="/admin" component={Admin} />

          <Route path="/login" component={Login} />
          <Route path="/register" component={Register} />

        </Switch>
      </div>
    );
}


export default withRouter(App);