import React from 'react';
import { Navbar, NavbarBrand } from 'reactstrap';
import { Link } from 'react-router-dom';
import './Navigation.css';

const Navigation = () => {
    return (
    <div>
      <Navbar color="dark" dark fixed="top">
        <NavbarBrand exact href="/" className="mr-auto">GoScrobble</NavbarBrand>
        <Link class="navLink" to="/">Home</Link>
        <Link class="navLink" to="/about">About</Link>
      </Navbar>
    </div>
    );
}

export default Navigation;