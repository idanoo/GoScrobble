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
  'Home',
  'My Profile',
  // 'Docs',
]

const isMobile = () => {
  return (/Android|webOS|iPhone|iPad|iPod|BlackBerry|IEMobile|Opera Mini/i.test(navigator.userAgent))
};

const Navigation = () => {
  const location = useLocation();

  // Lovely hack to highlight the current page (:
  let active = "home"
  if (location && location.pathname && location.pathname.length > 1) {
    active = location.pathname.replace(/\//, "");
  }

  let activeStyle = { color: '#FFFFFF' };
  let { user, Logout } = useContext(AuthContext);
  let [collapsed, setCollapsed] = useState(true);

  const toggleCollapsed = () => {
    setCollapsed(!collapsed)
  }

  const renderMobileNav = () => {
    return <Navbar color="dark" dark fixed="top">
      <NavbarBrand className="mr-auto"><img src={logo} className="nav-logo" alt="logo" /> GoScrobble</NavbarBrand>
      <NavbarToggler onClick={toggleCollapsed} className="mr-2" />
      <Collapse isOpen={!collapsed} navbar>
        {user ?
        <Nav className="navLinkLoginMobile" navbar>
          {loggedInMenuItems.map(menuItem =>
            <NavItem key={menuItem}>
                <Link
                key={menuItem}
                className="navLinkMobile"
                style={active === menuItem.toLowerCase() ? activeStyle : {}}
                to={menuItem === "My Profile" ? "/u/" + user.username : "/" + menuItem.toLowerCase()}                onClick={toggleCollapsed}
              >{menuItem}</Link>
            </NavItem>
          )}
          <Link
            to="/user"
            style={active === "user" ? activeStyle : {}}
            className="navLinkMobile"
            onClick={toggleCollapsed}
            >Settings</Link>
          {user.admin &&
            <Link
              to="/admin"
              style={active === "admin" ? activeStyle : {}}
              className="navLink"
              onClick={toggleCollapsed}
            >Admin</Link>}
          <Link to="/" className="navLink" onClick={Logout}>Logout</Link>
        </Nav>
      : <Nav className="navLinkLoginMobile" navbar>
              {menuItems.map(menuItem =>
                <NavItem key={menuItem}>
                  <Link
                    key={menuItem}
                    className="navLinkMobile"
                    style={active === "home" && menuItem.toLowerCase() === "home" ? activeStyle : (active === menuItem.toLowerCase() ? activeStyle : {})}
                    to={menuItem.toLowerCase() === "home" ? "/" : "/" + menuItem.toLowerCase()}
                    onClick={toggleCollapsed}
                  >{menuItem}
                  </Link>
                </NavItem>
              )}
              <NavItem>
                <Link
                  to="/login"
                  style={active === "login" ? activeStyle : {}}
                  className="navLinkMobile"
                  onClick={toggleCollapsed}
                  >Login</Link>
              </NavItem>
              <NavItem>
                <Link
                  to="/register"
                  className="navLinkMobile"
                  style={active === "register" ? activeStyle : {}}
                  onClick={toggleCollapsed}
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
            style={active === menuItem.toLowerCase() ? activeStyle : {}}
            to={menuItem === "My Profile" ? "/u/" + user.username : "/" + menuItem.toLowerCase()}
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
            style={active === "home" && menuItem.toLowerCase() === "home" ? activeStyle : (active === menuItem.toLowerCase() ? activeStyle : {})}
            to={menuItem.toLowerCase() === "home" ? "/" : "/" + menuItem.toLowerCase()}
          >
            {menuItem}
          </Link>
        )}
      </div>
      }
      {user ?
        <div className="navLinkLogin">
            <Link
              to="/user"
              style={active === "user" ? activeStyle : {}}
              className="navLink"
            >Settings</Link>
            {user.admin &&
            <Link
              to="/admin"
              style={active === "admin" ? activeStyle : {}}
              className="navLink"
            >Admin</Link>}
            <Link to="/admin" className="navLink" onClick={Logout}>Logout</Link>
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