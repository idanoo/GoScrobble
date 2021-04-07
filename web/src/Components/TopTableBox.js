import React from "react";
import { Link } from 'react-router-dom';
import './TopTableBox.css'

const TopTableBox = (props) => {
   let img = 'https://www.foot.com/wp-content/uploads/2017/06/placeholder-square-300x300.jpg';
   if (props.img && props.img !== '') {
      img = props.img
   }

   return (
         <Link to={"/artist/"+props.uuid} float="left" >
            <div
               className="topTableBox"
               style={{
                  backgroundImage: `url(${img})`,
                  backgroundSize: `cover`,
                  backgroundPosition: `top center`,
                  width: `${props.size}px`,
                  height: `${props.size}px`,
                  float: `left`,
               }} >
               <div className="topOverlay">
                  <span className="topText">#{props.number} {props.title}</span>
               </div>
            </div>
         </Link>

   );
}

export default TopTableBox;