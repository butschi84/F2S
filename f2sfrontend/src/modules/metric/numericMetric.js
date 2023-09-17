import React from 'react';

function NumericMetric(props) {
    return (
        <div className="card">
                <div className="card-content">
                    <div className="content">
                        {props.title}
                        <div style={{display: "flex"}}>
                            <div className="metric_text">{Math.round(props.value * 100) / 100}</div>
                            
                            <div className="metric_units">
                                {props.units}
                            </div>
                        </div>
                    </div>
                </div>
            </div>
    )
}

export default NumericMetric