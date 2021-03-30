import { React, Component } from 'react';
import { Navbar, NavbarBrand, Collapse, Nav, NavbarToggler, NavItem } from 'reactstrap';
import { Link } from 'react-router-dom';
import logo from '../logo.png';
import './Navigation.css';
import { connect } from 'react-redux';
import { logout } from '../Actions/auth';
import eventBus from "../Actions/eventBus";

import {
  LOGIN_SUCCESS,
  LOGOUT,
} from "../Actions/types";

const menuItems = [
  'Home',
  'About',
];

const loggedInMenuItems = [
  'Dashboard',
  'About',
]

const isMobile = () => {
  return (/Android|webOS|iPhone|iPad|iPod|BlackBerry|IEMobile|Opera Mini/i.test(navigator.userAgent))
};

class Navigation extends Component {
  constructor(props) {
    super(props);
    this.toggleNavbar = this.toggleNavbar.bind(this);

    // Yeah I know you might not hit home.. but I can't get the
    // path based finder thing working on initial load :sweatsmile:
    this.state = { active: "Home", collapsed: true};
  }

  componentDidMount() {
    const { isLoggedIn } = this.props;
    if (isLoggedIn) {
      this.setState({
        isLoggedIn: true,
      });
    }

    eventBus.on(LOGIN_SUCCESS, () =>
      this.setState({ isLoggedIn: true })
    );

    eventBus.on(LOGOUT, () =>
      this.setState({ isLoggedIn: false })
    );
  }

  componentWillUnmount() {
    eventBus.remove(LOGIN_SUCCESS);
    eventBus.remove(LOGOUT);
  }

  _handleClick(menuItem) {
    this.setState({ active: menuItem, collapsed: !this.state.collapsed });
  }

  toggleNavbar() {
    this.setState({ collapsed: !this.state.collapsed });
  }

  // This is a real mess. TO CLEAN UP.
  render() {
    const activeStyle = { color: '#FFFFFF' };

    const renderMobileNav = () => {
      return <Navbar color="dark" dark fixed="top">
        <NavbarBrand className="mr-auto"><img src={logo} className="nav-logo" alt="logo" /> GoScrobble</NavbarBrand>
        <NavbarToggler onClick={this.toggleNavbar} className="mr-2" />
        <Collapse isOpen={!this.state.collapsed} navbar>
          {this.state.isLoggedIn ?
          <Nav className="navLinkLoginMobile" navbar>
            {loggedInMenuItems.map(menuItem =>
              <NavItem>
                <Link
                  key={menuItem}
                  className="navLinkMobile"
                  style={this.state.active === menuItem ? activeStyle : {}}
                  onClick={this._handleClick.bind(this, menuItem)}
                  to={menuItem}
                >{menuItem}</Link>
              </NavItem>
            )}
            <Link
              to="/profile"
              style={this.state.active === "profile" ? activeStyle : {}}
              onClick={this._handleClick.bind(this, "profile")}
              className="navLinkMobile"
              >Profile</Link>
            <Link to="/" className="navLink" onClick={logout()}>Logout</Link>
          </Nav>
        : <Nav className="navLinkLoginMobile" navbar>
                {menuItems.map(menuItem =>
                  <NavItem>
                    <Link
                      key={menuItem}
                      className="navLinkMobile"
                      style={this.state.active === menuItem ? activeStyle : {}}
                      onClick={this._handleClick.bind(this, menuItem)}
                      to={menuItem === "Home" ? "/" : menuItem}
                    >
                      {menuItem}
                    </Link>
                  </NavItem>
                )}
                <NavItem>
                  <Link
                    to="/login"
                    style={this.state.active === "login" ? activeStyle : {}}
                    onClick={this._handleClick.bind(this, "login")}
                    className="navLinkMobile"
                    >Login</Link>
                </NavItem>
                <NavItem>
                  <Link
                    to="/register"
                    className="navLinkMobile"
                    style={this.state.active === "register" ? activeStyle : {}}
                    onClick={this._handleClick.bind(this, "register")}
                    history={this.props.history}
                  >Register</Link>
                </NavItem>
              </Nav>
        }
        </Collapse>
      </Navbar>
    }

    const renderDesktopNav = () => {
      return <Navbar color="dark" dark fixed="top">
        <NavbarBrand className="mr-auto"><img src={logo} className="nav-logo" alt="logo" /> GoScrobble</NavbarBrand>
        {this.state.isLoggedIn ?
          <div>
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
        </div>
        : <div>
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
        </div>
        }
        {this.state.isLoggedIn ?
          <div className="navLinkLogin">
              <Link
                to="/profile"
                style={this.state.active === "profile" ? activeStyle : {}}
                onClick={this._handleClick.bind(this, "profile")}
                className="navLink"
              >Profile</Link>
              <Link to="/" className="navLink" onClick={logout()}>Logout</Link>
            </div>
        :
        <div className="navLinkLogin">
          <Link
            to="/login"
            style={this.state.active === "login" ? activeStyle : {}}
            onClick={this._handleClick.bind(this, "login")}
            className="navLink"
          >Login</Link>
          <Link
            to="/register"
            className="navLink"
            style={this.state.active === "register" ? activeStyle : {}}
            onClick={this._handleClick.bind(this, "register")}
            history={this.props.history}
          >Register</Link>
      </div>

        }
      </Navbar>
    }

    return (
      <div>
        {
          isMobile()
            ? renderMobileNav()
            : renderDesktopNav()
        }
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