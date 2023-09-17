import React, { useState, useEffect } from "react";
import { useParams } from 'react-router-dom';
import { connect } from 'react-redux';
import * as _ from 'lodash';
import yaml from 'js-yaml';
import ReactMarkdown from 'react-markdown';
import { getMetricLastValue } from '../../services/functions';
import NumericMetric from "../../modules/metric/numericMetric";
import VectorMetric from "../../modules/metric/vectorMetric";

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
    const [f2sfunction, setF2SFunction] = useState(null)
    const [metricFunctionCapacity, setMetricFunctionCapacity] = useState("unknown")
    const [metricFunctionReplicas, setMetricFunctionReplicas] = useState("unknown")
    const [metricFunctionTotalCompleted, setMetricFunctionTotalCompleted] = useState([])
    const [tab, setTab] = useState("metadata")

    // set current subscription as state
    useEffect(() => {
        // find specific subscription
        const functionId = routeParams.id;
        const f2sFunction = _.find(props.functions, f => { return f.uid === functionId })
        setF2SFunction(f2sFunction)
    }, [props.functions, routeParams.id])

    useEffect(() => {
        if(!f2sfunction) return
        getFunctionMetrics()
        const inter = setInterval(()=>{ getFunctionMetrics() }, 5000)
        return () => {
            clearInterval(inter);
        };
    }, [f2sfunction])

    async function getFunctionMetrics() 
    {
        if(!f2sfunction || !f2sfunction.hasOwnProperty("name")) return

        getMetricLastValue(`job:function_capacity_average:reqpersec{functionname=\"${f2sfunction.name}\"}`).then(data => {
            if(data.status && data.status == "success")
            setMetricFunctionCapacity(data.data.result[0].values[data.data.result[0].values.length -1][1])
        })
        getMetricLastValue(`kube_deployment_status_replicas_available{functionname=\"${f2sfunction.name}\"}`).then(data => {
            if(data.status && data.status == "success")
            setMetricFunctionReplicas(data.data.result[0].values[data.data.result[0].values.length -1][1])
        })

        getMetricLastValue(`sum by(functionuid)(increase(f2s_requests_completed_total{functionname=\"${f2sfunction.name}\"}[1m])%2B1)`).then(data => {
            if(!data.status || data.status !== "success") return
            setMetricFunctionTotalCompleted(data.data.result[0].values)
        })
    }

    if(!f2sfunction) return ""
    return (
        <React.Fragment>
            <h1 className='title'>F2S Function Details</h1>

            <div class="tabs">
                <ul>
                    <li className={tab==="metadata" ? "is-active" : ""}><a onClick={()=>setTab("metadata")}>Metadata</a></li>
                    <li className={tab==="specification" ? "is-active" : ""}><a onClick={()=>setTab("specification")}>Specification</a></li>
                    <li className={tab==="target" ? "is-active" : ""}><a onClick={()=>setTab("target")}>Target</a></li>
                    <li className={tab==="metrics" ? "is-active" : ""}><a onClick={()=>setTab("metrics")}>Metrics</a></li>
                    <li className={tab==="yaml" ? "is-active" : ""}><a onClick={()=>setTab("yaml")}>YAML Definition</a></li>
                </ul>
            </div>

            {/* Metadata */}
            {tab === "metadata" &&
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
            {tab === "specification" &&
            <React.Fragment>
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
            </React.Fragment>
            }

            {/* Target */}
            {tab === "target" &&
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
            {tab === "yaml" &&
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

            {/* Metrics */}
            {tab === "metrics" &&
                <React.Fragment>
                    <div className="columns">
                        <div className="column">
                            <NumericMetric title="Function Capacity" value={metricFunctionCapacity} units="req/s" />
                        </div>
                        <div className="column">
                            <NumericMetric title="Current Function Replicas" value={metricFunctionReplicas} units="replicas ready" />
                        </div>
                    </div>
                    <VectorMetric 
                        values={metricFunctionTotalCompleted}/>
                </React.Fragment>
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
