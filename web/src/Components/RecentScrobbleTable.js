import React from "react";
import { Link } from 'react-router-dom';

const RecentScrobbleTable = (props) => {
    return (
        <div style={{
         border: `1px solid #FFFFFF`,
         width: `100%`,
         display: `flex`,
         flexWrap: `wrap`,
         minWidth: `300px`,
         maxWidth: `1200px`,
        }}>
         {
            props.data &&
            props.data.map(function (element) {
               let localTime = new Date(element.time);
               return <div style={{borderBottom: `1px solid #CCC`, width: `100%`, padding: `2px`}} key={"box" + element.time}>
                    {localTime.toLocaleString()}
                     <Link
                        key={"track" + element.time}
                        to={"/track/"+element.track.uuid}
                  > {element.track.name}</Link> -&nbsp;
                     <Link
                        key={"artist" + element.time}
                        to={"/artist/"+element.artist.uuid}
                     >{element.artist.name}</Link>
                     &nbsp;by <Link
                        key={"user" + element.time}
                        to={"/u/"+element.user.name}
                     > {element.user.name}</Link>
               </div>;

            })
         }
      </div>
    );
}

export default RecentScrobbleTable;