import React from 'react';
import { Navbar, NavbarBrand, NavLink } from 'reactstrap';

const Navigation = () => {
    return (
    <div>
      <Navbar color="dark" dark fixed="top">
        <NavbarBrand exact href="/" className="mr-auto">GoScrobble</NavbarBrand>
        <NavLink exact href="/">Home</NavLink>
        <NavLink href="/about">About</NavLink>
      </Navbar>
    </div>
    );
}

export default Navigation;