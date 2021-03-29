import { React, Component } from 'react';
import { Navbar, NavbarBrand } from 'reactstrap';
import { Link } from 'react-router-dom';
import logo from '../logo.png';
import './Navigation.css';
import { connect } from 'react-redux';
import { logout } from '../Actions/auth';

const menuItems = [
  'Home',
  'About',
];

const loggedInMenuItems = [
  'Dashboard',
  'About',
]
class Navigation extends Component {
  constructor(props) {
    super(props);
    // Yeah I know you might not hit home.. but I can't get the
    // path based finder thing working on initial load :sweatsmile:
    this.state = { active: "Home" };
  }

  componentDidMount() {
    const isLoggedIn = this.props.isLoggedIn;

    if (isLoggedIn) {
      this.setState({
        isLoggedIn: true,
      });
    }
  }

  _handleClick(menuItem) {
    this.setState({ active: menuItem });
  }

  render() {
    const activeStyle = { color: '#FFFFFF' };

    const renderAuthButtons = () => {
      if (this.state.isLoggedIn) {
        return <div className="navLinkLogin">
                <Link to="/profile" className="navLink">Profile</Link>
                <Link to="/" className="navLink" onClick={logout()}>Logout</Link>
              </div>;
      } else {
        return <div className="navLinkLogin">
                <Link to="/login" className="navLink">Login</Link>
                <Link to="/register" className="navLink" history={this.props.history}>Register</Link>
              </div>;
      }
    }

    const renderMenuButtons = () => {
      if (this.state.isLoggedIn) {
        return <div>
                {loggedInMenuItems.map(menuItem =>
                  <Link
                    key={menuItem}
                    className="navLink"
                    style={this.state.active === menuItem ? activeStyle : {}}
                    onClick={this._handleClick.bind(this, menuItem)}
                    to={menuItem}
                  >
                    {menuItem}
                  </Link>
                )}
              </div>;
      } else {
        return <div>
                {menuItems.map(menuItem =>
                  <Link
                    key={menuItem}
                    className="navLink"
                    style={this.state.active === menuItem ? activeStyle : {}}
                    onClick={this._handleClick.bind(this, menuItem)}
                    to={menuItem === "Home" ? "/" : menuItem}
                  >
                    {menuItem}
                  </Link>
                )}
              </div>;
      }
    }

    return (
      <div>
        <Navbar color="dark" dark fixed="top">
          <NavbarBrand href="/" className="mr-auto"><img src={logo} className="nav-logo" alt="logo" /> GoScrobble</NavbarBrand>
          {renderMenuButtons()}
          {renderAuthButtons()}
        </Navbar>
      </div>
    );
  }
}

function mapStateToProps(state) {
  const { isLoggedIn } = state.auth;
  return {
    isLoggedIn,
  };
}

export default connect(mapStateToProps)(Navigation);