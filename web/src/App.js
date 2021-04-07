import { Route, Switch, withRouter } from 'react-router-dom';
import Home from './Pages/Home';
import About from './Pages/About';
import Profile from './Pages/Profile';
import Artist from './Pages/Artist';
import Album from './Pages/Album';
import Track from './Pages/Track';
import User from './Pages/User';
import Admin from './Pages/Admin';
import Login from './Pages/Login';
import Register from './Pages/Register';
import Reset from './Pages/Reset';

import Navigation from './Components/Navigation';

import 'bootstrap/dist/css/bootstrap.min.css';
import './App.css';

const App = () => {
    let boolTrue = true;

    // Remove loading spinner on load
    const el = document.querySelector(".loader-container");
    if (el) {
      el.remove();
    }

    return (
      <div>
        <Navigation />
        <Switch>
          <Route exact={boolTrue} path={["/", "/home"]} component={Home} />
          <Route path="/about" component={About} />

          <Route path="/user" component={User} />
          <Route path="/u/:uuid" component={Profile} />
          <Route path="/artist/:uuid" component={Artist} />
          <Route path="/album/:uuid" component={Album} />
          <Route path="/track/:uuid" component={Track} />

          <Route path="/admin" component={Admin} />

          <Route path="/login" component={Login} />
          <Route path="/register" component={Register} />

          <Route path="/reset/:token" component={Reset} />
          <Route path="/reset" component={Reset} />

        </Switch>
      </div>
    );
}


export default withRouter(App);