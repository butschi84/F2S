import React, { useEffect } from 'react';
import { getF2SConfig } from '../../store/configSlice';
import { connect, useDispatch } from 'react-redux';


function Settings(props) {
    const dispaatch = useDispatch();

    useEffect(() => {
        dispaatch(getF2SConfig())
    })

    return (
        <React.Fragment>
            <h1 className='title'>Settings</h1>
            { props.config && props.config.Config.F2S &&
                <div className="card">
                    <div className="card-content">
                        <div className="content">
                            Request Timeout
                            <input 
                            className="input"
                            readOnly
                            value={props.config.Config.F2S.Timeouts.RequestTimeout} />
                            HTTP Timeout
                            <input 
                            className="input"
                            readOnly
                            value={props.config.Config.F2S.Timeouts.HttpTimeout} />
                            Scaling Timeout
                            <input 
                            className="input"
                            readOnly
                            value={props.config.Config.F2S.Timeouts.ScalingTimeout} />
                        </div>
                    </div>
                </div>
            }
        </React.Fragment>
    )
}


function mapStateToProps(state) {
    return { 
        config: state.configSlice.config,
    };
  }
  
export default connect(mapStateToProps)(Settings)
