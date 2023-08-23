import React, {useEffect} from 'react';
import { connect, useDispatch } from 'react-redux';
import { getAllFunctions } from '../../store/functionsSlice';
import DataTable from 'react-data-table-component';
import { NavLink } from 'react-router-dom';


function F2SFunctions(props) {
    const dispatch = useDispatch();

    const columns = [
        {
            name: 'name',
            selector: row => <NavLink to={`/f2sfunctions/${row.uid}`}>{row.name}</NavLink>
        },
        {
            name: 'endpoint',
            selector: row => `${props.apiURL}/invoke${row.spec.endpoint}`
        },
        {
            name: 'method',
            selector: row => row.spec.method
        },
        {
            name: '',
            selector: row => <NavLink className="button is-primary" to={`/f2sfunctions/${row.uid}/invoke`}>Invoke</NavLink>
        },
        {
            name: '',
            selector: row => <NavLink className="button is-danger" to={`/f2sfunctions/${row.uid}/delete`}>Delete</NavLink>
        },
    ];


    useEffect(() => {
        dispatch(getAllFunctions());
    }, [dispatch]);


    return (
        <React.Fragment>
            <h1 className='title'>F2S Functions</h1>
            <DataTable
                columns={columns}
                data={props.functions}
                pagination
                persistTableHead
        />
        </React.Fragment>
    )
}
function mapStateToProps(state) {
    return { 
        functions: state.functionsSlice.functions,
        apiURL: state.connectivitySlice.apiURL
    };
  }
  
export default connect(mapStateToProps)(F2SFunctions)
