import React from 'react';
import logo from '../../images/f2s-logo.png';

function Navbar(props) {
    return (
        <div>
            <nav className="navbar" role="navigation" aria-label="main navigation" >
                <div className="navbar-brand">
                    <a className="navbar-item" href="/">
                    <img src={logo} width="112" height="28" />
                    </a>

                    <a role="button" className="navbar-burger" aria-label="menu" aria-expanded="false" data-target="f2sNavbar">
                    <span aria-hidden="true"></span>
                    <span aria-hidden="true"></span>
                    <span aria-hidden="true"></span>
                    </a>
                </div>

                <div id="f2sNavbar" className="navbar-menu">
                    <div className="navbar-start">
                    <a className="navbar-item">
                        Functions
                    </a>
                    </div>
                </div>
            </nav>
        </div>
    );
}

export default Navbar;