import React from 'react';
import {Link} from 'react-router-dom';
import './NavBar.css';

function NavBar() {
  return (
    <nav className="navbar navbar-inverse bg-primary fixed-top">
      <Link className="navbar-brand" to="/">
        Media Player <h6> by namo</h6>
      </Link>
    </nav>
  );
}

export default NavBar;