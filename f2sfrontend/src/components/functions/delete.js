import React, { useState, useEffect } from "react";
import { useParams, useNavigate } from 'react-router-dom';
import * as _ from 'lodash';
import { connect, useDispatch } from 'react-redux';
import { deleteF2SFunction } from '../../store/functionsSlice';


function DeleteF2SFunction(props) {
    const routeParams = useParams();
    const dispatch = useDispatch();
    const [f2sfunction, setF2SFunction] = useState()
    const navigate = useNavigate()

    // set current subscription as state
    useEffect(() => {
        // find specific subscription
        const functionId = routeParams.id;
        setF2SFunction(_.find(props.functions, f => { return f.uid === functionId }))
    }, [props.functions, routeParams.id])

    function deletefunction() {
        dispatch(deleteF2SFunction(f2sfunction));
        navigate('/f2sfunctions');
    }

    return (
        <React.Fragment>
            <h2 className='title'>Delete Function</h2>
            Do you want to delete the function:<br />
            {f2sfunction &&
                <React.Fragment>
                    UID
                    <input className="input" readonly value={f2sfunction.uid} />
                    Name
                    <input className="input" readonly value={f2sfunction.name} />

                    <button type="submit" className="button is-danger" onClick={deletefunction}>Delete</button>
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
  
export default connect(mapStateToProps)(DeleteF2SFunction)
