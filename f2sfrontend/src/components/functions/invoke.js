import React, { useState, useEffect } from "react";
import { useParams } from 'react-router-dom';
import { connect } from 'react-redux';
import * as _ from 'lodash';
import spinner from '../../images/spinner2.gif';
import { get, post } from '../../services/common'
import ReactMarkdown from 'react-markdown';

function InvokeFunction(props) {
    const routeParams = useParams();
    const [f2sfunction, setF2SFunction] = useState({})
    const [invocationInProgress, setInvocationInProgress] = useState(false)
    const [invocationResult, setInvocationResult] = useState("")
    const [postData, setPostData] = useState("")

    // set current subscription as state
    useEffect(() => {
        // find specific subscription
        const functionId = routeParams.id;
        setF2SFunction(_.find(props.functions, f => { return f.uid === functionId }))
    }, [props.functions, routeParams.id])

    function invoke(f2sfunction, apiURL) {
        if(!f2sfunction.spec) return
        setInvocationInProgress(true)
        setInvocationResult("")
        switch(f2sfunction.spec.method){
            case "GET":
                get(`/invoke${f2sfunction.spec.endpoint}`).then((response) => {
                    setInvocationInProgress(false)
                    setInvocationResult(response)
                }).catch((error) => {
                    console.log(error)
                    setInvocationInProgress(false)
                })
                break;
            case "POST":
                post(`/invoke${f2sfunction.spec.endpoint}`, postData).then((response) => {
                    setInvocationInProgress(false)
                    setInvocationResult(response)
                }).catch((error) => {
                    console.log(error)
                    setInvocationInProgress(false)
                })
                break;
            default:
                break;
        }
    }

    if(!f2sfunction || !f2sfunction.spec) return ""
    return (
        <React.Fragment>
            <h1 className='title'>Invoke - {f2sfunction.name}</h1>

            <div className="card">
                <div className="card-content">
                    <div class="media">
                        <div class="media-content">
                            <p class="title is-4">Description</p>
                        </div>
                    </div>
                    <div className="content">
                        {/* Description */}
                        <ReactMarkdown>{f2sfunction.spec.description}</ReactMarkdown>
                    </div>
                </div>
            </div>

            <div className="card">
                <div className="card-content">
                <div class="media">
                    <div class="media-content">
                            { ["PUT", "POST"].includes(f2sfunction.spec.method) &&
                            <p class="title is-4">{f2sfunction.spec.method.toLowerCase()} data</p>
                            }
                        </div>
                    </div>
                    <div className="content">
                        { ["PUT", "POST"].includes(f2sfunction.spec.method) &&
                        <textarea
                        className="input"
                        style={{height:"150px"}}
                        readOnly={f2sfunction.spec.method === "GET" || f2sfunction.spec.method === "DELETE"}
                        onChange={(e)=>setPostData(e.target.value)}
                        value={postData}
                        rows="10"></textarea>
                        }
                        
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
                        { invocationResult !== "" &&
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
                            <img src={spinner} alt="in progress" />
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
