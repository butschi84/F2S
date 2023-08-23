import React, { useEffect, useState } from 'react';
import { connect, useDispatch } from 'react-redux';
import { checkConnectivity, setBackendURL, logout} from '../../store/connectivitySlice'

function ConnectivityCheck(props) {
    const dispatch = useDispatch()
    const [urlPath, setUrlPath] = useState("")

    useEffect(() => {
        dispatch(checkConnectivity())
    }, [dispatch])

    function logoff() {
        dispatch(logout())
    }

    function connect() {
        dispatch(setBackendURL(urlPath))
    }

    return (
        <React.Fragment>
            {
                !props.connectivity &&
                <div className='container' style={{width: "500px", marginTop: "300px"}}>
                    <div className="field">
                        <label className="label">F2S API Connection</label>
                        <div className="control">
                            <input 
                                className="input" 
                                type="text" 
                                placeholder="Text input" 
                                value={urlPath} 
                                onChange={(p) => setUrlPath(p.target.value)} />
                        </div>
                        <p className="help">Please specify address of f2s backend API</p>
                    </div>
                    <button className='button is-primary' onClick={connect}>Save</button>
                </div>
            }

            {
                props.connectivity && 
                <React.Fragment>
                {props.children}
                <button className='backendUrlButton button' onClick={logoff}>Logout</button>
                </React.Fragment>
            }
        </React.Fragment>
    )
}

function mapStateToProps(state) {
    return { 
        connectivity: state.connectivitySlice.ApiConnectionEstablished,
        apiURL: state.connectivitySlice.apiURL,
    };
  }
  
export default connect(mapStateToProps)(ConnectivityCheck)
