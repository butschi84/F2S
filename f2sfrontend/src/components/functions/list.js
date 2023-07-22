import React, {useEffect} from 'react';
import { connect, useDispatch } from 'react-redux';
import { getAllFunctions } from '../../store/functionsSlice';

function F2SFunctions() {
    const dispatch = useDispatch();

    useEffect(() => {
        dispatch( getAllFunctions())
    }, []);


    return (
        <h1 className='title'>F2S Functions</h1>
    )
}
function mapStateToProps(state) {
    return { 
        functions: state.functionsSlice.functions,
    };
  }
  
export default connect(mapStateToProps)(F2SFunctions)
