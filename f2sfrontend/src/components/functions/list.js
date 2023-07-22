import React, {useEffect} from 'react';
import { connect, useDispatch } from 'react-redux';
import { getAllFunctions } from '../../store/functionsSlice';
import DataTable from 'react-data-table-component';
import { NavLink, useParams } from 'react-router-dom';


function F2SFunctions(props) {
    const dispatch = useDispatch();

    const columns = [
        {
            name: 'uid',
            selector: row => row.uid
        },
        {
            name: 'name',
            selector: row => row.name
        },
        {
            name: 'endpoint',
            selector: row => row.spec.endpoint
        },
        {
            name: '',
            selector: row => <NavLink className="button is-primary" to={`/functions/${row.uid}/invoke`}>Invoke</NavLink>
        },
    ];


    useEffect(() => {
        dispatch( getAllFunctions())
    }, []);


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
    };
  }
  
export default connect(mapStateToProps)(F2SFunctions)
