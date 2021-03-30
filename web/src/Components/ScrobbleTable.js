import React from "react";

class ScrobbleTable extends React.Component {
  constructor(props) {
    super(props);

    this.state = {
      data: this.props.data,
    };
  }

  render() {
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
                  this.state.data && this.state.data.items &&
                  this.state.data.items.map(function (element) {
                     return <tr>
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
}

export default ScrobbleTable;