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
              </tr>
           </thead>
           <tbody>
              {
                  props.data &&
                  props.data.map(function (element) {
                     return <tr key={element.uuid}>
                       <td>{element.time}</td>
                       <td>{element.track}</td>
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