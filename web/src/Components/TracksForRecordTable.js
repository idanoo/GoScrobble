import React, { useContext, useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import { getAlbumTracks, getArtistTracks } from '../Api/index'
import ScaleLoader from 'react-spinners/ScaleLoader';

const TracksForRecordTable = (props) => {
   const [tracks, setTracks] = useState({});
   const [loading, setLoading] = useState(true);

   useEffect(() => {
      if (props.albumuuid !== undefined) {
         getAlbumTracks(props.albumuuid)
         .then(data => {
            setTracks(data);
            setLoading(false);
         })
      } else if (props.artistuuid !== undefined) {
         getArtistTracks(props.artistuuid)
         .then(data => {
            setTracks(data);
            setLoading(false);
         })
      } 
   }, [props.albumuuid, props.artistuuid]);

   if (loading) {
      return (
        <div className="pageWrapper">
          <ScaleLoader color="#6AD7E5" />
        </div>
      )
    }

    console.log(tracks);

    return (
       <div style={{
         textAlign: `center`,
         }}>Tracks<br/>
        <div style={{
         border: `1px solid #FFFFFF`,
         width: `100%`,
         display: `flex`,
         flexWrap: `wrap`,
         minWidth: `300px`,
         maxWidth: `1200px`,
        }}>
         {
            tracks && tracks.tracks &&
            Object.keys(tracks.tracks).map(key => {
               return <div style={{borderBottom: `1px solid #CCC`, width: `100%`, padding: `2px`}} key={"box" + tracks.tracks[key].uuid}>
                     <Link
                        key={"track" + tracks.tracks[key].uuid}
                        to={"/track/"+tracks.tracks[key].uuid}
                  > {tracks.tracks[key].name}</Link>
               </div>;

            })
         }
      </div>
      </div>
    );
}

export default TracksForRecordTable;