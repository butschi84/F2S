import React from 'react';
import logo from '../../images/f2s-logo.png';
import { NavLink, useParams } from 'react-router-dom';
import { connect, useDispatch } from 'react-redux';

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
                        
                        <div class="navbar-item has-dropdown is-hoverable">
                            <NavLink to="/f2sfunctions" className="navbar-item">
                                Functions
                            </NavLink>

                            <div class="navbar-dropdown">
                                <NavLink to={`/f2sfunctions`} className="navbar-item">
                                    List
                                </NavLink>
                                <NavLink to={`/f2sfunctions/create`} className="navbar-item">
                                    Create
                                </NavLink>
                            </div>
                        </div>
                        <div class="navbar-item has-dropdown is-hoverable">
                            <a class="navbar-link">
                            More
                            </a>

                            <div class="navbar-dropdown">
                                <NavLink to={`/settings`} className="navbar-item">
                                    Settings
                                </NavLink>
                                <a href={`${props.apiURL}/docs/`} className="navbar-item">
                                    API Docs
                                </a>
                            </div>
                        </div>
                    </div>
                </div>
            </nav>
            <div className="navbar-wrap"></div>
        </div>
    );
}

function mapStateToProps(state) {
    return { 
        apiURL: state.connectivitySlice.apiURL,
    };
  }
  
export default connect(mapStateToProps)(Navbar)
