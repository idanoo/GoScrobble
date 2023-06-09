import React from "react";
import './TopTable.css'
import TopTableBox from './TopTableBox';
import ClipLoader from 'react-spinners/ClipLoader'

const TopTable = (props) => {
   if (!props.items || Object.keys(props.items).length < 1) {
      return (
         <span>Not enough data to show top {props.type}s.<br/></span>
      )
   }

   let tracks = props.items;

   if (props.loading) {
      return <div style={{textAlign: `center`}}>
      <span>Top {props.type}s</span>
      <div className="biggestWrapper">
         <div className="biggestBox"></div>
         <div className="biggestBox">
            <div style={{padding: `80px`}}>
               <ClipLoader color="#6AD7E5" size={150} />
            </div>
         </div>
         <div className="biggestBox"></div>
      </div>
   </div>
   }
   return (
      <div style={{textAlign: `center`}}>
         <span>Top {props.type}s {props.extraText && props.extraText}</span>
         <div className="biggestWrapper">
            <div className="biggestBox">
               <TopTableBox
                  size={300}
                  number="1"
                  title={tracks[1].name}
                  link={"/" + props.type + "/" + tracks[1].uuid}
                  uuid={tracks[1].img}
               />
            </div>
            { Object.keys(props.items).length > 5 &&
               <div className="biggestBox">
                     <TopTableBox
                        size={150}
                        number="2"
                        title={tracks[2].name}
                        link={"/" + props.type + "/" + tracks[2].uuid}
                        uuid={tracks[2].img}
                     />
                     <TopTableBox
                        size={150}
                        number="3"
                        title={tracks[3].name}
                        link={"/" + props.type + "/" + tracks[3].uuid}
                        uuid={tracks[3].img}
                     />
                     <TopTableBox
                        size={150}
                        number="4"
                        title={tracks[4].name}
                        link={"/" + props.type + "/" + tracks[4].uuid}
                        uuid={tracks[4].img}
                     />
                     <TopTableBox
                        size={150}
                        number="5"
                        title={tracks[5].name}
                        link={"/" + props.type + "/" + tracks[5].uuid}
                        uuid={tracks[5].img}
                     />
               </div>
            }
            { Object.keys(props.items).length >= 14 &&
               <div className="biggestBox">
                  <TopTableBox
                     size={100}
                     number="6"
                     title={tracks[6].name}
                     link={"/" + props.type + "/" + tracks[6].uuid}
                     uuid={tracks[6].img}
                  />
                  <TopTableBox
                     size={100}
                     number="7"
                     title={tracks[7].name}
                     link={"/" + props.type + "/" + tracks[7].uuid}
                     uuid={tracks[7].img}
                  />
                  <TopTableBox
                     size={100}
                     number="8"
                     title={tracks[8].name}
                     link={"/" + props.type + "/" + tracks[8].uuid}
                     uuid={tracks[8].img}
                  />
                  <TopTableBox
                     size={100}
                     number="9"
                     title={tracks[9].name}
                     link={"/" + props.type + "/" + tracks[9].uuid}
                     uuid={tracks[9].img}
                  />
                  <TopTableBox
                     size={100}
                     number="10"
                     title={tracks[10].name}
                     link={"/" + props.type + "/" + tracks[10].uuid}
                     uuid={tracks[10].img}
                  />
                  <TopTableBox
                     size={100}
                     number="11"
                     title={tracks[11].name}
                     link={"/" + props.type + "/" + tracks[11].uuid}
                     uuid={tracks[11].img}
                  />
                  <TopTableBox
                     size={100}
                     number="12"
                     title={tracks[12].name}
                     link={"/" + props.type + "/" + tracks[12].uuid}
                     uuid={tracks[12].img}
                  />
                  <TopTableBox
                     size={100}
                     number="13"
                     title={tracks[13].name}
                     link={"/" + props.type + "/" + tracks[13].uuid}
                     uuid={tracks[13].img}
                  />
                  <TopTableBox
                     size={100}
                     number="14"
                     title={tracks[14].name}
                     link={"/" + props.type + "/" + tracks[14].uuid}
                     uuid={tracks[14].img}
                  />
               </div>
            }
         </div>
      </div>
   );
}

export default TopTable;