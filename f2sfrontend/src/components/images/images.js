import React from 'react';
import DataTable from 'react-data-table-component';

const samples = [
    {
        name: "NodeJS API Example",
        file: "sample_nodejs.zip"
    }
]

function F2SImages(props) {
    const columns = [
        {
            name: 'name',
            selector: row => row.name
        },
        {
            name: 'description',
            selector: row => "example nodejs web api with random response delay"
        },
        {
            name: 'Download Link',
            selector: row => <a href={`/${row.file}`}>{row.file}</a>
        }
    ]

    return (
        <React.Fragment>
            <h1 className='title'>
                Images
            </h1>

            <box className="box">
                These are some samples that help to get you started building your custom container images.
            </box>

            <box className="box"
>               <DataTable
                    columns={columns}
                    data={samples}
                    pagination={false}
                    persistTableHead/>
            </box>
        </React.Fragment>
    )
}

export default F2SImages