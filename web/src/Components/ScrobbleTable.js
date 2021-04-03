import React from "react";

const ScrobbleTable = (props) => {
    return (
        <div>
        <table border={1} cellPadding={5}>
           <thead>
              <tr>
                <td>Timestamp</td>
                <td>Track</td>
                <td>Artist</td>
                <td>Album</td>
                <td>Source</td>
              </tr>
           </thead>
           <tbody>
              {
                  props.data &&
                  props.data.map(function (element) {
                     let localTime = new Date(element.time);
                     return <tr key={element.uuid}>
                       <td>{localTime.toLocaleString()}</td>
                       <td>{element.track}</td>
                       <td>{element.artist}</td>
                       <td>{element.album}</td>
                       <td>{element.source}</td>
                     </tr>;
                  })
              }
           </tbody>
        </table>
      </div>
    );
}

export default ScrobbleTable;