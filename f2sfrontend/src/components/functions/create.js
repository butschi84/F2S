import React, { useEffect } from 'react';
import useForm from '../../modules/forms/useForm';
import Select from '../../modules/forms/select'
import {createNewF2SFunction} from '../../store/functionsSlice';
import { connect, useDispatch } from 'react-redux';
import { v4 as uuidv4 } from 'uuid';
import { NavLink, useParams, useNavigate } from 'react-router-dom';

function F2SFunctionCreate(props) {
    const dispatch = useDispatch();
    const navigate = useNavigate()

    const { form, handleChange, resetForm, setForm, setInForm } = useForm({
		name: '',
        spec_endpoint: '',
        spec_method: 'GET',
        spec_description: '',
        target_containerimage: 'romanhuesler/f2snodejstest:latest',
        target_endpoint: '/',
        target_port: '9092',
        target_minreplicas: '0',
        target_maxreplicas: '1',
	});

    const methods = [
        {id: "GET", name: "GET"},
        {id: "POST", name: "POST"},
        {id: "PUT", name: "PUT"},
        {id: "DELETE", name: "DELETE"}
    ]

    useEffect(() => {
        const generatedFunctionName =  `my-function-${uuidv4()}`
        setInForm("name", generatedFunctionName)
        setInForm("spec_endpoint", `/${generatedFunctionName}`)
    }, [])

    async function save() {
        const f2sfunction = {
            name: form.name,
            spec: {
                endpoint: form.spec_endpoint,
                method: form.spec_method,
                description: form.spec_description,
            },
            target: {
                containerImage: form.target_containerimage,
                endpoint: form.target_endpoint,
                port: parseInt(form.target_port),
                minReplicas: parseInt(form.target_minreplicas),
                maxReplicas: parseInt(form.target_maxreplicas)
            }
        }
        await dispatch(createNewF2SFunction(f2sfunction));
        navigate('/f2sfunctions');
    }

    return (
        <React.Fragment>
            {/* METADATA */}
            <div className="card">
                <div className="card-content">

                    {/* name */}
                    <div className="content">
                        Function Name
                        <input 
                        className="input"
                        id="name"
                        name='name'
                        onChange={handleChange}
                        value={form.name} />

                    </div>
                </div>
            </div>

            {/* SPECIFICATION */}
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
                        id="spec_endpoint"
                        name="spec_endpoint"
                        onChange={handleChange}
                        value={form.spec_endpoint} />

                        {/* Method */}
                        <Select
                            className="input"
                            id="spec_method"
                            name="spec_method"
                            onChange={handleChange}
                            value={form.spec_method}
                            options={methods} />

                        {/* Description */}
                        Description
                        <input 
                        className="input"
                        id="spec_description"
                        name="spec_description"
                        onChange={handleChange}
                        value={form.spec_description} />
                    </div>
                </div>
            </div>

            {/* TARGET */}
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
                        id="target_containerimage"
                        name="target_containerimage"
                        onChange={handleChange}
                        value={form.target_containerimage} />

                        {/* Endpoint */}
                        Endpoint
                        <input 
                        className="input"
                        id="target_endpoint"
                        name="target_endpoint"
                        onChange={handleChange}
                        value={form.target_endpoint} />

                        {/* Port */}
                        Port
                        <input 
                        className="input"
                        id="target_port"
                        name="target_port"
                        onChange={handleChange}
                        value={form.target_port} />

                        {/* Maximum Replicas */}
                        Maximum Replicas
                        <input 
                        className="input"
                        id="target_maxreplicas"
                        name="target_maxreplicas"
                        onChange={handleChange}
                        value={form.target_maxreplicas} />

                        {/* Minimum Replicas */}
                        Minimum Replicas
                        <input 
                        className="input"
                        id="target_minreplicas"
                        name="target_minreplicas"
                        onChange={handleChange}
                        value={form.target_minreplicas} />
                    </div>
                </div>
            </div>

            <button 
                className='button is-primary'
                onClick={save}>Save</button>
        </React.Fragment>
    )
}

export default F2SFunctionCreate;