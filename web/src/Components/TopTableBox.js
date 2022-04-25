import React from "react";
import { Link } from 'react-router-dom';
import './TopTableBox.css'

const TopTableBox = (props) => {
   return (
         <Link to={props.link} float="left" >
            <div
               className="topTableBox"
               style={{
                  backgroundImage: `url(${process.env.REACT_APP_API_URL + "/img/" + props.uuid + "_300px.jpg"})`,
                  backgroundSize: `cover`,
                  backgroundPosition: `top center`,
                  width: `${props.size}px`,
                  height: `${props.size}px`,
                  float: `left`,
               }} >
               <div className="topOverlay" style={{ maxWidth: `${props.size-'5'}px` }}>
                  <span className="topText" style={{
                     fontSize: `${props.size === 300 ? '11pt' : (props.size === 150 ? '8pt': '8pt' )}`
                  }}>#{props.number} {props.title}</span>
               </div>
            </div>
         </Link>

   );
}

export default TopTableBox;