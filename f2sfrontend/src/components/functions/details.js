import React, { useState, useEffect } from "react";
import { NavLink, useParams } from 'react-router-dom';
import { connect, useDispatch } from 'react-redux';
import * as _ from 'lodash';
import axios from 'axios';
import spinner from '../../images/spinner2.gif';
import yaml from 'js-yaml';

function getF2SFunctionYAMLDefinition(f2sfunction) {
    let clone = {...f2sfunction}
    clone["metadata"] = {
        name: clone.name
    }
    delete clone.name
    delete clone.uid
    clone["kind"] = "Function"
    clone["apiVersion"] = "f2s.opensight.ch/v1alpha1"
    return yaml.dump(clone)
}

function F2SFunctionDetails(props) {
    const routeParams = useParams();
    const [f2sfunction, setF2SFunction] = useState()
    const [tab, setTab] = useState("metadata")

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

            <div class="tabs">
                <ul>
                    <li className={tab=="metadata" ? "is-active" : ""}><a onClick={()=>setTab("metadata")}>Metadata</a></li>
                    <li className={tab=="specification" ? "is-active" : ""}><a onClick={()=>setTab("specification")}>Specification</a></li>
                    <li className={tab=="target" ? "is-active" : ""}><a onClick={()=>setTab("target")}>Target</a></li>
                    <li className={tab=="yaml" ? "is-active" : ""}><a onClick={()=>setTab("yaml")}>YAML Definition</a></li>
                </ul>
            </div>

            {/* Metadata */}
            {tab == "metadata" &&
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
            }

            {/* Specification */}
            {tab == "specification" &&
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
            }

            {/* Target */}
            {tab == "target" &&
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
            }

            {/* YAML definition */}
            {tab == "yaml" &&
            <div className="card">
            <div className="card-content">
                <div className="content">
                    <textarea className="input" style={{height: "700px"}}>
                    {getF2SFunctionYAMLDefinition(f2sfunction)}
                    </textarea>
                </div>
            </div>
        </div>
            }
        </React.Fragment>
    )
}


function mapStateToProps(state) {
    return { 
        functions: state.functionsSlice.functions
    };
  }
  
export default connect(mapStateToProps)(F2SFunctionDetails)
