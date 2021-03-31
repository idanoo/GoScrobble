import { React, Component } from 'react';
import { Navbar, NavbarBrand, Collapse, Nav, NavbarToggler, NavItem } from 'reactstrap';
import { Link } from 'react-router-dom';
import logo from '../logo.png';
import './Navigation.css';

import logout from '../Contexts/AuthContextProvider';

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
    this.handleLogout = this.handleLogout.bind(this);

    // Yeah I know you might not hit home.. but I can't get the
    // path based finder thing working on initial load :sweatsmile:
    this.state = { active: "Home", collapsed: true };
  }

  componentDidMount() {
    const { isLoggedIn } = this.props;
    if (isLoggedIn) {
      this.setState({
        isLoggedIn: true,
      });
    }

  }

  _handleClick(menuItem) {
    this.setState({ active: menuItem, collapsed: !this.state.collapsed });
  }

  handleLogout() {
    this.dispatch(logout());
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
            <Link to="/" className="navLink" onClick={this.handleLogout}>Logout</Link>
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
              <Link to="/" className="navLink" onClick={this.handleLogout}>Logout</Link>
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

export default Navigation;