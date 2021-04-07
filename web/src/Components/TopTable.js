import React from "react";
import './TopTable.css'
import TopTableBox from './TopTableBox';

const TopTable = (props) => {
    return (
      <div>
         <span>Top {props.type}</span>
         <div className="biggestWrapper">
            <div className="biggestBox">
               <TopTableBox
                  size={300}
                  number="1"
                  title="hot milk"
                  uuid="a2bcc230-f7be-4087-b49a-8c43d19ed316"
                  img="https://i.scdn.co/image/a397625e38fb671f1baa81997b4c1fd2670fcb10"
               />
            </div>
            <div className="biggestBox">
                  <TopTableBox
                     size={150}
                     number="2"
                     title="Pendulum"
                     uuid="a2bcc230-f7be-4087-b49a-8c43d19ed316"
                     img="https://i.scdn.co/image/0f476171f283207656e95e1005cea7040be475d7"
                  />
                  <TopTableBox size={150} number="3" title="Illenium" uuid="a2bcc230-f7be-4087-b49a-8c43d19ed316" />
                  <TopTableBox size={150} number="4" title="As It is" uuid="a2bcc230-f7be-4087-b49a-8c43d19ed316" />
                  <TopTableBox size={150} number="5" title="CHVRCHES" uuid="a2bcc230-f7be-4087-b49a-8c43d19ed316" />
            </div>
            <div className="biggestBox">
               <TopTableBox size={100} number="6" title="tester" uuid="a2bcc230-f7be-4087-b49a-8c43d19ed316" />
               <TopTableBox size={100} number="7" title="tester" uuid="a2bcc230-f7be-4087-b49a-8c43d19ed316" />
               <TopTableBox size={100} number="8" title="tester" uuid="a2bcc230-f7be-4087-b49a-8c43d19ed316" />
               <TopTableBox size={100} number="9" title="tester" uuid="a2bcc230-f7be-4087-b49a-8c43d19ed316" />
               <TopTableBox size={100} number="10" title="tester" uuid="a2bcc230-f7be-4087-b49a-8c43d19ed316" />
               <TopTableBox size={100} number="11" title="tester" uuid="a2bcc230-f7be-4087-b49a-8c43d19ed316" />
               <TopTableBox size={100} number="12" title="tester" uuid="a2bcc230-f7be-4087-b49a-8c43d19ed316" />
               <TopTableBox size={100} number="13" title="tester" uuid="a2bcc230-f7be-4087-b49a-8c43d19ed316" />
               <TopTableBox size={100} number="14" title="tester" uuid="a2bcc230-f7be-4087-b49a-8c43d19ed316" />
            </div>
         </div>
      </div>
    );
}

export default TopTable;