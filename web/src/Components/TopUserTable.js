import { Link } from 'react-router-dom';
import './TopUserTable.css'
import React, { useState, useEffect } from 'react';
import ScaleLoader from 'react-spinners/ScaleLoader';
import { getTopUsersForTrack } from '../Api/index'

const TopUserTable = (props) => {
    const [loading, setLoading] = useState(true);
    const [data, setData] = useState({});


    useEffect(() => {
        if (!props.uuid) {
          return false;
        }

        getTopUsersForTrack(props.uuid)
          .then(data => {
            setData(data);
            setLoading(false);
          })
      }, [props.uuid])

    if (loading) {
        return (
          <div className="pageWrapper">
            <ScaleLoader color="#6AD7E5" />
          </div>
        )
      }

    return (
        <div style={{
         width: `100%`,
         display: `flex`,
         flexWrap: `wrap`,
         marginLeft: `20px`,
         textAlign: `left`,
        }}>
         {
            data.items &&
            data.items.map(function (element) {
               return <div style={{width: `100%`, padding: `2px`}} key={"box" + props.uuid}>
                    <Link
                        key={"user" + element.user_uuid}
                        to={"/u/"+element.user_name}
                     >{element.user_name}</Link> ({element.count})
               </div>;

            })
         }
      </div>
    );
}

export default TopUserTable;