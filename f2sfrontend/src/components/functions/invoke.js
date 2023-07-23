import React, { useState, useEffect } from "react";
import { NavLink, useParams } from 'react-router-dom';
import { connect, useDispatch } from 'react-redux';
import * as _ from 'lodash';
import axios from 'axios';
import spinner from '../../images/spinner2.gif';

function InvokeFunction(props) {
    const routeParams = useParams();
    const [f2sfunction, setF2SFunction] = useState({})
    const [invocationInProgress, setInvocationInProgress] = useState(false)
    const [invocationResult, setInvocationResult] = useState("")

    // set current subscription as state
    useEffect(() => {
        // find specific subscription
        const functionId = routeParams.id;
        setF2SFunction(_.find(props.functions, f => { return f.uid == functionId }))
    }, [])

    function invoke(f2sfunction, apiURL) {
        if(!f2sfunction.spec) return
        setInvocationInProgress(true)
        setInvocationResult("")
        axios.get(`${apiURL}/invoke${f2sfunction.spec.endpoint}`).then((response) => {
            setInvocationInProgress(false)
            setInvocationResult(response.data)
        }).catch((error) => {
            console.log(error)
            setInvocationInProgress(false)
        })
    }

    return (
        <React.Fragment>
            <h1 className='title'>Invoke Function</h1>

            <div className="card">
                <div className="card-content">
                    <div className="content">
                        Function Name
                        <input 
                        className="input"
                        readOnly
                        value={f2sfunction.name} />
                        <br />
                        Method
                        <input 
                        className="input"
                        readOnly
                        value={f2sfunction.spec.method} />
                        <br />
                        Endpoint
                        <input 
                        className="input"
                        readOnly
                        value={`${props.apiURL}/invoke${f2sfunction.spec.endpoint}`} />
                        <br />
                        <br />
                        
                        <button 
                        className="button is-primary"
                        disabled={invocationInProgress}
                        onClick={()=>invoke(f2sfunction, props.apiURL)}>Invoke</button>
                        
                    </div>
                </div>
            </div>

            <div className="card">
                <div className="card-content">
                    <div className="content">
                        { invocationResult != "" &&
                            <React.Fragment>
                            Result
                            <textarea className="input" rows="40" cols="40" style={{height: "150px"}}>
                            {
                                JSON.stringify(invocationResult, null, 2)
                            }
                            </textarea>
                            </React.Fragment>
                        }

                        { invocationInProgress &&
                            <img src={spinner} />
                        }
                    </div>
                </div>
            </div>
        </React.Fragment>
    )
}


function mapStateToProps(state) {
    return { 
        functions: state.functionsSlice.functions,
        apiURL: state.connectivitySlice.apiURL
    };
  }
  
export default connect(mapStateToProps)(InvokeFunction)
