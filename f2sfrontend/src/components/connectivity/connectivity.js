import React, { useEffect } from 'react';
import { connect, useDispatch } from 'react-redux';
import { checkConnectivity} from '../../store/connectivitySlice'

function ConnectivityCheck(props) {
    const dispatch = useDispatch()

    useEffect(() => {
        dispatch(checkConnectivity())
    }, [])

    return (
        <React.Fragment>
            {JSON.stringify(props.connectivity)}

            <div className='container' style={{width: "500px"}}>
                <div className="field">
                    <label className="label">F2S API</label>
                    <div className="control">
                        <input className="input" type="text" placeholder="Text input" />
                    </div>
                    <p className="help">Please specify address of f2s backend API</p>
                </div>
                <button className='button is-primary'>Save</button>
            </div>
        </React.Fragment>
    )
}

function mapStateToProps(state) {
    return { 
        connectivity: state.connectivitySlice.ApiConnectionEstablished,
    };
  }
  
export default connect(mapStateToProps)(ConnectivityCheck)
