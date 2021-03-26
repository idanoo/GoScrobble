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
      this.state = { isLoggedIn: false };
    }

    toggleLogin() {
      this.setState({ isLoggedIn: !this.state.isLoggedIn })
    }

    _handleClick(menuItem) {
      this.setState({ active: menuItem });
    }

    render() {
      const activeStyle = { color: '#FFF' };

      const renderAuthButton = () => {
        if (this.state.isLoggedIn) {
          return <Link class="navLink" onClick={this.toggleLogin.bind(this)}>Logout</Link>;
        } else {
          return <Link class="navLink" onClick={this.toggleLogin.bind(this)}>Login</Link>;
        }
      }

      return (
        <div>
          <Navbar color="dark" dark fixed="top">
            <NavbarBrand exact href="/" className="mr-auto">GoScrobble</NavbarBrand>
            {menuItems.map(menuItem =>
            <Link
              class="navLink"
              style={this.state.active === menuItem ? activeStyle : {}}
             onClick={this._handleClick.bind(this, menuItem)}
            >
              {menuItem}
            </Link>
         )}

            <Link class="navLink" to="/">Home</Link>
            <Link class="navLink" to="/about">About</Link>
            {renderAuthButton()}
          </Navbar>
        </div>
      );
    }
  }

export default Navigation;