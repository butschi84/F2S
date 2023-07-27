import React from 'react';
import useForm from '../../modules/forms/useForm';

function F2SFunctionCreate(props) {

    const { form, handleChange, resetForm, setForm } = useForm({
		name: '',
        spec_endpoint: '',
        spec_method: 'GET',
        spec_description: '',
        target_containerimage: '',
        targt_endpoint: '',
        target_port: '80',
        target_minreplicas: '0',
        target_maxreplicas: '1',
	});



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
                        value={form.spec_endpoint} />

                        {/* Method */}
                        Method
                        <input 
                        className="input"
                        id="spec_method"
                        name="spec_method"
                        value={form.spec_method} />

                        {/* Description */}
                        Description
                        <input 
                        className="input"
                        id="spec_description"
                        name="spec_description"
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
                        value={form.target_containerimage} />

                        {/* Endpoint */}
                        Endpoint
                        <input 
                        className="input"
                        id="target_endpoint"
                        name="target_endpoint"
                        value={form.target_endpoint} />

                        {/* Port */}
                        Port
                        <input 
                        className="input"
                        id="target_port"
                        name="target_port"
                        value={form.target_port} />

                        {/* Maximum Replicas */}
                        Maximum Replicas
                        <input 
                        className="input"
                        id="target_maxreplicas"
                        name="target_maxreplicas"
                        value={form.target_maxreplicas} />

                        {/* Minimum Replicas */}
                        Minimum Replicas
                        <input 
                        className="input"
                        id="target_minreplicas"
                        name="target_minreplicas"
                        value={form.target_minreplicas} />
                    </div>
                </div>
            </div>
        </React.Fragment>
    )
}

export default F2SFunctionCreate;