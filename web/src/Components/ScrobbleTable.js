import React from "react";
import { Link } from 'react-router-dom';

const ScrobbleTable = (props) => {
    return (
        <div style={{
         border: `1px solid #FFFFFF`,
         width: `100%`,
         display: `flex`,
         flexWrap: `wrap`,
         minWidth: `300px`,
         maxWidth: `900px`,
        }}>
         {
            props.data &&
            props.data.map(function (element) {
               let localTime = new Date(element.time);
               return <div style={{borderBottom: `1px solid #CCC`, width: `100%`, padding: `2px`}} key={"box" + element.time}>
                    {localTime.toLocaleString()}<br/>
                    <Link
                        key={"artist" + element.time}
                        to={"/artist/"+element.artist.uuid}
                     >{element.artist.name}</Link> -
                     <Link
                        key={"track" + element.time}
                        to={"/track/"+element.track.uuid}
                     > {element.track.name}</Link>
               </div>;
            })
         }
      </div>
    );
}

export default ScrobbleTable;