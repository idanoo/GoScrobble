import { React, Component } from 'react';
import { Navbar, NavbarBrand } from 'reactstrap';
import { Link } from 'react-router-dom';
import './Navigation.css';

const menuItems = [
  'Home',
  'About',
];

class Navigation extends Component {
    constructor(props) {
      super(props);
      // Yeah I know you might not hit home.. but I can't get the
      // path based finder thing working on initial load :sweatsmile:
      console.log(this.props.initLocation)
      this.state = { isLoggedIn: false, active: "Home" };
    }

    toggleLogin() {
      this.setState({ isLoggedIn: !this.state.isLoggedIn })
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
                  <Link to="/" className="navLink" onClick={this.toggleLogin.bind(this)}>Logout</Link>
                </div>;
        } else {
          return <div className="navLinkLogin">
                  <Link to="/login" className="navLink">Login</Link>
                  <Link to="/register" className="navLink">Register</Link>
                </div>;
        }
      }

      return (
        <div>
          <Navbar color="dark" dark fixed="top">
            <NavbarBrand href="/" className="mr-auto">GoScrobble</NavbarBrand>
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
          {renderAuthButtons()}
          </Navbar>
        </div>
      );
    }
  }

export default Navigation;