import React, { useContext, useState, useEffect } from 'react';
import '../App.css';
import './AlbumEdit.css';
import { useHistory } from 'react-router-dom';
import ScaleLoader from 'react-spinners/ScaleLoader';
import { getAlbum, uploadImage } from '../Api/index'
import { Link } from 'react-router-dom';
import AuthContext from '../Contexts/AuthContext';
import FileUploader from '../Components/FileUploader';

const AlbumEdit = (route) => {
  const history = useHistory();
  const { user } = useContext(AuthContext);

  const [loading, setLoading] = useState(true);
  const [album, setAlbum] = useState({});

  const [selectedFile, setSelectedFile] = useState(null);

  const submitForm = (e) => {
    e.preventDefault();

    const formData = new FormData();
    formData.append("name", "file");
    formData.append("file", selectedFile);
  
    uploadImage(formData, "album", album.uuid, history)
  };

  let albumUUID = false;
  if (route && route.match && route.match.params && route.match.params.uuid) {
    albumUUID = route.match.params.uuid;
  } else {
    albumUUID = false;
  }

  useEffect(() => {
    if (!albumUUID) {
      return false;
    }

    getAlbum(albumUUID)
      .then(data => {
        setAlbum(data);
        setLoading(false);
      })
  }, [albumUUID])

  if (!user) {
    history.push("/login")
  }

  if (user && !user.mod) {
    history.push("/Dashboard")
  }

  if (loading) {
    return (
      <div className="pageWrapper">
        <ScaleLoader color="#6AD7E5" />
      </div>
    )
  }

  if (!albumUUID || !album) {
    return (
      <div className="pageWrapper">
        Unable to fetch user
      </div>
    )
  }

  return (
    <div className="pageWrapper">
      <h1 style={{margin: 0}}>
      {album.name} {<Link
            key="editbuttonomg"
            to={"/album/" + albumUUID}
          >unedit</Link>}
      </h1>
      <div className="pageBody" style={{width: `900px`, textAlign: `center`}}>
        <img src={process.env.REACT_APP_API_URL + "/img/" + album.uuid + "_full.jpg"} alt={album.name} style={{maxWidth: `300px`, maxHeight: `300px`}}/>

        <form>
          <FileUploader
            onFileSelect={(file) => setSelectedFile(file)}
          />
          <button onClick={submitForm}>Submit</button>
        </form>

      </div>
    </div>
  )
}

export default AlbumEdit;