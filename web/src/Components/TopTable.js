import React from "react";
import './TopTable.css'
import TopTableBox from './TopTableBox';

const TopTable = (props) => {
   if (!props.items || !props.items.tracks) {
      return (
         <span>No data.</span>
      )
   }

   let tracks = props.items.tracks;

    return (
      <div>
         <span>Top {props.type}s</span>
         <div className="biggestWrapper">
            <div className="biggestBox">
               <TopTableBox
                  size={300}
                  number="1"
                  title={tracks[1].name}
                  link={"/" + props.type + "/" + tracks[1].uuid}
                  img={tracks[1].img}
               />
            </div>
            <div className="biggestBox">
                  <TopTableBox
                     size={150}
                     number="2"
                     title={tracks[2].name}
                     link={"/" + props.type + "/" + tracks[2].uuid}
                     img={tracks[2].img}
                  />
                  <TopTableBox
                     size={150}
                     number="3"
                     title={tracks[3].name}
                     link={"/" + props.type + "/" + tracks[3].uuid}
                     img={tracks[3].img}
                  />
                  <TopTableBox
                     size={150}
                     number="4"
                     title={tracks[4].name}
                     link={"/" + props.type + "/" + tracks[4].uuid}
                     img={tracks[4].img}
                  />
                  <TopTableBox
                     size={150}
                     number="5"
                     title={tracks[5].name}
                     link={"/" + props.type + "/" + tracks[5].uuid}
                     img={tracks[5].img}
                  />
            </div>
            <div className="biggestBox">
               <TopTableBox
                  size={100}
                  number="6"
                  title={tracks[6].name}
                  link={"/" + props.type + "/" + tracks[6].uuid}
                  img={tracks[6].img}
               />
               <TopTableBox
                  size={100}
                  number="7"
                  title={tracks[7].name}
                  link={"/" + props.type + "/" + tracks[7].uuid}
                  img={tracks[7].img}
               />
               <TopTableBox
                  size={100}
                  number="8"
                  title={tracks[8].name}
                  link={"/" + props.type + "/" + tracks[8].uuid}
                  img={tracks[8].img}
               />
               <TopTableBox
                  size={100}
                  number="9"
                  title={tracks[9].name}
                  link={"/" + props.type + "/" + tracks[9].uuid}
                  img={tracks[9].img}
               />
               <TopTableBox
                  size={100}
                  number="10"
                  title={tracks[10].name}
                  link={"/" + props.type + "/" + tracks[10].uuid}
                  img={tracks[10].img}
               />
               <TopTableBox
                  size={100}
                  number="11"
                  title={tracks[11].name}
                  link={"/" + props.type + "/" + tracks[11].uuid}
                  img={tracks[11].img}
               />
               <TopTableBox
                  size={100}
                  number="12"
                  title={tracks[12].name}
                  link={"/" + props.type + "/" + tracks[12].uuid}
                  img={tracks[12].img}
               />
               <TopTableBox
                  size={100}
                  number="13"
                  title={tracks[13].name}
                  link={"/" + props.type + "/" + tracks[13].uuid}
                  img={tracks[13].img}
               />
               <TopTableBox
                  size={100}
                  number="14"
                  title={tracks[14].name}
                  link={"/" + props.type + "/" + tracks[14].uuid}
                  img={tracks[14].img}
               />
            </div>
         </div>
      </div>
    );
}

export default TopTable;