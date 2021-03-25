import './App.css';
import Home from './Components/Pages/Home';
import About from './Components/Pages/About';
import Login from './Components/Pages/Login';
import Navigation from './Components/Pages/Navigation';
import { Route, Switch, HashRouter } from 'react-router-dom';
import '../node_modules/bootstrap/dist/css/bootstrap.min.css';

function App() {

  if (!true) {
    return <Login />
  }

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
      // <div className="App">
      //   <Router>
      //     <Navigation/>
      //     <Switch>
      //       <Route exact path='/' component={Home}/>
      //       <Route path='/about' component={About}/>
      //     </Switch>
      //   </Router>
      // </div>
  );
}

export default App;
