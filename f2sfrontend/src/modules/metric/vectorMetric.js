import React, { useState, useEffect } from "react";
import * as Highcharts from 'highcharts';
import HighchartsReact from 'highcharts-react-official';
import * as _ from 'lodash'

function VectorMetric(props) {
    const [opts, setOptions] = useState(null)
    
    useEffect(() => {
        if(props.values.length === 0) return
        setOptions({
            title: {
                text: 'Function Invocations Total',
                align: 'left',
                style: {
                    color: 'white'
                }
            },
            colors: ["#E91E61"],
            chart: {
                backgroundColor: '#181C1F'
            },
            subtitle: {
                text: 'Total Number of Function Invocations',
                align: 'left'
            },
            time: {
                useUTC: false
            },
            yAxis: {
                gridLineColor: '#32363B',
                title: {
                    text: 'Total Completed Requests'
                },
                labels: {
                    style: {
                        color: 'white'
                    }
                }
            },
        
            xAxis: {
                type: 'datetime',
                tickInterval: 60 * 1000,
                gridLineColor: '#32363B',
                tickmarkPlacement: 'on',
                gridLineWidth: 1,
                labels: {
                    style: {
                        color: 'white'
                    }
                }
            },
        
            legend: {
                enabled: false
            },
        
            series: [{
                name: 'Manufacturing',
                data: _.map(props.values, v => {
                    return [
                        v[0] * 1000,
                        parseFloat(v[1])
                    ]
                })
            }],
        
            responsive: {
                rules: [{
                    condition: {
                        maxWidth: 500
                    }
                }]
            }
        })
    }, [props.values])

    return (
        <React.Fragment>
        {opts != null && opts.title && opts.series[0].data.length &&
            <div className="card">
                <div className="card-content">
                    <div className="content">
                    <HighchartsReact
                        highcharts={Highcharts}
                        options={opts}
                    />
                    </div>
                </div>
            </div>
        }
        </React.Fragment>
    )
}

export default VectorMetric

