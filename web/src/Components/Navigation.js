import { React, useState, useContext } from 'react';
import { Navbar, NavbarBrand, Collapse, Nav, NavbarToggler, NavItem } from 'reactstrap';
import { Link, useLocation } from 'react-router-dom';
import logo from '../logo.png';
import './Navigation.css';

import AuthContext from '../Contexts/AuthContext';

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

const Navigation = () => {
  const location = useLocation();

  // Lovely hack to highlight the current page (:
  let active = "Home"
  if (location && location.pathname && location.pathname.length > 1) {
    active = location.pathname.replace(/\//, "");
  }

  let activeStyle = { color: '#FFFFFF' };
  let { user, Logout } = useContext(AuthContext);
  let [collapsed, setCollapsed] = useState(true);

  const renderMobileNav = () => {
    return <Navbar color="dark" dark fixed="top">
      <NavbarBrand className="mr-auto"><img src={logo} className="nav-logo" alt="logo" /> GoScrobble</NavbarBrand>
      <NavbarToggler onClick={setCollapsed(!collapsed)} className="mr-2" />
      <Collapse isOpen={!collapsed} navbar>
        {user ?
        <Nav className="navLinkLoginMobile" navbar>
          {loggedInMenuItems.map(menuItem =>
            <NavItem>
              <Link
                key={menuItem}
                className="navLinkMobile"
                style={active === menuItem ? activeStyle : {}}
                to={menuItem}
              >{menuItem}</Link>
            </NavItem>
          )}
          <Link
            to="/profile"
            style={active === "profile" ? activeStyle : {}}
            className="navLinkMobile"
            >Profile</Link>
          <Link to="/" className="navLink" onClick={Logout}>Logout</Link>
        </Nav>
      : <Nav className="navLinkLoginMobile" navbar>
              {menuItems.map(menuItem =>
                <NavItem>
                  <Link
                    key={menuItem}
                    className="navLinkMobile"
                    style={active === menuItem ? activeStyle : {}}
                    to={menuItem === "Home" ? "/" : menuItem}
                  >
                    {menuItem}
                  </Link>
                </NavItem>
              )}
              <NavItem>
                <Link
                  to="/Login"
                  style={active === "Login" ? activeStyle : {}}
                  className="navLinkMobile"
                  >Login</Link>
              </NavItem>
              <NavItem>
                <Link
                  to="/Register"
                  className="navLinkMobile"
                  style={active === "Register" ? activeStyle : {}}
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
      {user ?
        <div>
        {loggedInMenuItems.map(menuItem =>
          <Link
            key={menuItem}
            className="navLink"
            style={active === menuItem ? activeStyle : {}}
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
            style={active === menuItem ? activeStyle : {}}
            to={menuItem === "Home" ? "/" : menuItem}
          >
            {menuItem}
          </Link>
        )}
      </div>
      }
      {user ?
        <div className="navLinkLogin">
            <Link
              to="/profile"
              style={active === "profile" ? activeStyle : {}}
              className="navLink"
            >{user.username}</Link>
            <Link to="/" className="navLink" onClick={Logout}>Logout</Link>
          </div>
      :
      <div className="navLinkLogin">
        <Link
          to="/login"
          style={active === "login" ? activeStyle : {}}
          className="navLink"
        >Login</Link>
        <Link
          to="/register"
          className="navLink"
          style={active === "register" ? activeStyle : {}}
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

export default Navigation;