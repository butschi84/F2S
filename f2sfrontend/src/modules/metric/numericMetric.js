import React from 'react';

function NumericMetric(props) {
    return (
        <div className="card">
                <div className="card-content">
                    <div className="content">
                        {props.title}
                        <div className="columns">
                            <div className="column is-one-fifth">
                                <div className="metric_text">{Math.round(props.value * 100) / 100}</div>
                            </div>
                            <div className="column metric_units">
                                {props.units}
                            </div>
                        </div>
                    </div>
                </div>
            </div>
    )
}

export default NumericMetric