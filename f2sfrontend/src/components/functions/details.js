import React, { useState, useEffect } from "react";
import { NavLink, useParams } from 'react-router-dom';
import { connect, useDispatch } from 'react-redux';
import * as _ from 'lodash';
import axios from 'axios';
import spinner from '../../images/spinner2.gif';

function F2SFunctionDetails(props) {
    const routeParams = useParams();
    const [f2sfunction, setF2SFunction] = useState()

    // set current subscription as state
    useEffect(() => {
        // find specific subscription
        const functionId = routeParams.id;
        setF2SFunction(_.find(props.functions, f => { return f.uid == functionId }))
    }, [])

    if(!f2sfunction) return ""
    return (
        <React.Fragment>
            <h1 className='title'>F2S Function Details</h1>

            <div className="card">
                <div className="card-content">
                    <div className="content">
                        {/* UID */}
                        Function UID
                        <input 
                        className="input"
                        readOnly
                        value={f2sfunction.uid} />

                        {/* NAME */}
                        Function Name
                        <input 
                        className="input"
                        readOnly
                        value={f2sfunction.name} />
                    </div>
                </div>
            </div>

            {/* Specification */}
            <div className="card">
                <div className="card-content">
                    <div class="media">
                        <div class="media-content">
                            <p class="title is-4">Specification</p>
                        </div>
                    </div>
                    <div className="content">
                        {/* Endpoint */}
                        Endpoint
                        <input 
                        className="input"
                        readOnly
                        value={f2sfunction.spec.endpoint} />

                        {/* Method */}
                        Method
                        <input 
                        className="input"
                        readOnly
                        value={f2sfunction.spec.method} />
                    </div>
                </div>
            </div>

            {/* Target */}
            <div className="card">
                <div className="card-content">
                    <div class="media">
                        <div class="media-content">
                            <p class="title is-4">Target</p>
                        </div>
                    </div>
                    <div className="content">
                        {/* Container Image */}
                        Container Image
                        <input 
                        className="input"
                        readOnly
                        value={f2sfunction.target.containerImage} />

                        {/* Endpoint */}
                        Endpoint
                        <input 
                        className="input"
                        readOnly
                        value={f2sfunction.target.endpoint} />

                        {/* Port */}
                        Port
                        <input 
                        className="input"
                        readOnly
                        value={f2sfunction.target.port} />

                        {/* Maximum Replicas */}
                        Maximum Replicas
                        <input 
                        className="input"
                        readOnly
                        value={f2sfunction.target.maxReplicas} />

                        {/* Minimum Replicas */}
                        Minimum Replicas
                        <input 
                        className="input"
                        readOnly
                        value={f2sfunction.target.minReplicas} />
                    </div>
                </div>
            </div>
        </React.Fragment>
    )
}


function mapStateToProps(state) {
    return { 
        functions: state.functionsSlice.functions
    };
  }
  
export default connect(mapStateToProps)(F2SFunctionDetails)
