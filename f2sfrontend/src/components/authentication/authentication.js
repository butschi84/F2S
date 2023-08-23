import React, { useState } from 'react';
import { connect, useDispatch } from 'react-redux';
import { 
    signinWithToken,
    signinWithUsernamePassword } from '../../store/connectivitySlice';

function Authentication(props) {
    const dispatch = useDispatch()

    const [jwttoken, setToken] = useState("")
    const [basicAuthUsername, setBasicAuthUsername] = useState("")
    const [basicAuthPassword, setBasicAuthPassword] = useState("")

    const { authenticated, authenticationType } = props;

    function loginWithToken() {
        dispatch(signinWithToken(jwttoken))
    }

    function loginWithUsernamePassword() {
        dispatch(signinWithUsernamePassword(basicAuthUsername, basicAuthPassword))
    }

    return (
        <React.Fragment>
            {/* TOKEN AUTH */}
            {!authenticated && authenticationType === "token" &&
                <div className='container' style={{width: "500px", marginTop: "300px"}}>
                    <React.Fragment>
                        <h2 className='title'>F2S Authentication</h2>
                        JWT Token:
                    <input 
                        value={jwttoken}
                        onChange={(e)=>setToken(e.target.value)}
                        className='input' />
                    <button 
                        onClick={loginWithToken}
                        className='button is-primaty'>Login</button>
                    </React.Fragment>
                </div>
            }

            {/* BASIC AUTH */}
            {!authenticated && authenticationType === "basic" &&
                <div className='container' style={{width: "500px", marginTop: "300px"}}>
                <React.Fragment>
                    <h2 className='title'>F2S Authentication</h2>
                    Username:
                    <input 
                        value={basicAuthUsername}
                        onChange={(e)=>setBasicAuthUsername(e.target.value)}
                        className='input' />
                    Password:
                    <input 
                        value={basicAuthPassword}
                        onChange={(e)=>setBasicAuthPassword(e.target.value)}
                        type="password"
                        className='input' />
                    <button 
                        onClick={loginWithUsernamePassword}
                        className='button is-primaty'>Login</button>
                </React.Fragment>
            </div>
            }

            { (authenticated || authenticationType === "none")  &&
                props.children
            }
        </React.Fragment>
    )
}

function mapStateToProps(state) {
    return { 
        apiURL: state.connectivitySlice.apiURL,
        authenticated: state.connectivitySlice.authenticated,
        authenticationType: state.connectivitySlice.authenticationType
    };
  }
  
export default connect(mapStateToProps)(Authentication)
