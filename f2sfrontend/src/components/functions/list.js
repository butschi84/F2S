import React, {useEffect} from 'react';
import { connect, useDispatch } from 'react-redux';
import { getAllFunctions } from '../../store/functionsSlice';
import DataTable from 'react-data-table-component';
import { NavLink } from 'react-router-dom';
import rocket from '../../images/rocket.svg';

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
            name: '',
            // selector: row => <NavLink className="button is-primary" to={`/f2sfunctions/${row.uid}/invoke`}>Invoke</NavLink>
            selector: row => <NavLink to={`/f2sfunctions/${row.uid}/invoke`}><img src="data:image/svg+xml;base64,PD94bWwgdmVyc2lvbj0iMS4wIiBlbmNvZGluZz0idXRmLTgiPz4KPCEtLSBHZW5lcmF0b3I6IEFkb2JlIElsbHVzdHJhdG9yIDI3LjguMSwgU1ZHIEV4cG9ydCBQbHVnLUluIC4gU1ZHIFZlcnNpb246IDYuMDAgQnVpbGQgMCkgIC0tPgo8c3ZnIHZlcnNpb249IjEuMSIgaWQ9IkViZW5lXzEiIHhtbG5zPSJodHRwOi8vd3d3LnczLm9yZy8yMDAwL3N2ZyIgeG1sbnM6eGxpbms9Imh0dHA6Ly93d3cudzMub3JnLzE5OTkveGxpbmsiIHg9IjBweCIgeT0iMHB4IgoJIHZpZXdCb3g9IjAgMCA1NTEuNzkgNTU1LjMyIiBzdHlsZT0iZW5hYmxlLWJhY2tncm91bmQ6bmV3IDAgMCA1NTEuNzkgNTU1LjMyOyIgeG1sOnNwYWNlPSJwcmVzZXJ2ZSI+CjxzdHlsZSB0eXBlPSJ0ZXh0L2NzcyI+Cgkuc3Qwe2ZpbGw6bm9uZTtzdHJva2U6I0U3MjI2MjtzdHJva2Utd2lkdGg6MTI7c3Ryb2tlLW1pdGVybGltaXQ6MTA7fQoJLnN0MDpob3ZlciB7IHN0cm9rZTogd2hpdGU7IH0KPC9zdHlsZT4KPHBhdGggY2xhc3M9InN0MCIgZD0iTTUzNi41Nyw1MC45OWMtMC43NS0xMC4xNy0yLjMxLTIwLTQuNzItMjkuNGMtOS4xOS0yLjExLTE4Ljc3LTMuNDMtMjguNjYtMy45NQoJYy03My4xMy0zLjg3LTE2My4wNSwzNS4yLTIzNS42NywxMDkuNjRjLTQyLjU2LDQzLjYzLTcyLjgyLDkyLjg5LTg5LjQ0LDE0MC44OWM3LjcyLDI2LjUzLDIxLjcsNTAuODMsNDIuMSw3MC43MwoJYzIxLjYsMjEuMDcsNDcuOTUsMzQuODIsNzYuNDgsNDEuNTFjNDYuNzYtMTcuOTcsOTQuMjktNDkuMDksMTM2LjA5LTkxLjk0QzUwNS4xLDIxNC4zMSw1NDEuOTIsMTIzLjg2LDUzNi41Nyw1MC45OXoiLz4KPGc+Cgk8cGF0aCBjbGFzcz0ic3QwIiBkPSJNMjMwLjUyLDE3MC42NGMtNDkuMTMtMjAuNDQtMTIxLjM2LDEuNTQtMTc2LjI2LDU3LjgzYy0xMy4zMSwxMy42NC0yNC41NywyOC4yNS0zMy43LDQzLjI1CgkJYzU5LjgyLTM0LjAyLDEyMi4xNS0zNy40NSwxNTcuMjUtMy4yMWwwLjI5LTAuM2MwLTAuMDEtMC4wMS0wLjAyLTAuMDEtMC4wM0MxODkuMzksMjM1LjUyLDIwNy4wMiwyMDIuMjksMjMwLjUyLDE3MC42NHoiLz4KCTxwYXRoIGNsYXNzPSJzdDAiIGQ9Ik0zOTIuNTYsMzI0Ljc5Yy0zMC45LDI0LjUxLTYzLjU4LDQzLjIxLTk1Ljg5LDU1LjYzYy0wLjAyLDAtMC4wNC0wLjAxLTAuMDYtMC4wMWwtMC4yNSwwLjI2CgkJYzM1LjEsMzQuMjQsMzMuMjIsOTYuNjMsMC42OSwxNTcuMjhjMTQuNzctOS41LDI5LjEtMjEuMTIsNDIuNC0zNC43NkMzOTQuNjMsNDQ2LjYsNDE0Ljc0LDM3My4zOSwzOTIuNTYsMzI0Ljc5eiIvPgoJPHBhdGggY2xhc3M9InN0MCIgZD0iTTIzMC41MiwxNzAuNjRjLTIzLjUxLDMxLjY1LTQxLjEzLDY0Ljg5LTUyLjQ0LDk3LjU0YzAsMC4wMSwwLjAxLDAuMDIsMC4wMSwwLjAzIi8+CjwvZz4KPGNpcmNsZSBjbGFzcz0ic3QwIiBjeD0iNDA4LjEyIiBjeT0iMTQ4LjQ0IiByPSI1Ni44MSIvPgo8cGF0aCBjbGFzcz0ic3QwIiBkPSJNMTc3LjE3LDQ1NC4wOGMtMzEuODYsMzIuNjYtMTI0Ljg4LDU5LjE0LTEyNC44OCw1OS4xNHMyNC4xNi05My42NSw1Ni4wMi0xMjYuMzFzNzMuMS00NC4xLDkyLjEyLTI1LjU1CglTMjA5LjAzLDQyMS40MiwxNzcuMTcsNDU0LjA4eiIvPgo8L3N2Zz4KCg==" alt="Embedded SVG" className="rocket" /></NavLink>
        }
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
