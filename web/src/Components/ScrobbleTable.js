import React from "react";
import { Link } from 'react-router-dom';

const ScrobbleTable = (props) => {
    return (
        <div style={{width: `100%`}}>
        <table style={{width: `100%`}} border={1} cellPadding={5}>
           <thead>
              <tr>
                <td>Timestamp</td>
                <td>Track</td>
                <td>Artist</td>
                <td>Album</td>
              </tr>
           </thead>
           <tbody>
              {
                  props.data &&
                  props.data.map(function (element) {
                     let localTime = new Date(element.time);
                     return <tr key={element.uuid}>
                        <td>{localTime.toLocaleString()}</td>
                        <td>
                           <Link
                               key={element.track.uuid}
                              to={"/track/"+element.track.uuid}
                               >{element.track.name}
                           </Link>
                        </td>
                        <td>{element.artist}</td>
                        <td>{element.album}</td>
                     </tr>;
                  })
              }
           </tbody>
        </table>
      </div>
    );
}

export default ScrobbleTable;