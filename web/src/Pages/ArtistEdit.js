import React, { useContext, useState, useEffect } from 'react';
import '../App.css';
import './ArtistEdit.css';
import { useHistory } from 'react-router-dom';
import ScaleLoader from 'react-spinners/ScaleLoader';
import { getArtist, uploadImage } from '../Api/index'
import { Link } from 'react-router-dom';
import AuthContext from '../Contexts/AuthContext';
import FileUploader from '../Components/FileUploader';

const ArtistEdit = (route) => {
  const history = useHistory();
  const { user } = useContext(AuthContext);

  const [loading, setLoading] = useState(true);
  const [artist, setArtist] = useState({});

  const [selectedFile, setSelectedFile] = useState(null);

  const submitForm = (e) => {
    e.preventDefault();

    const formData = new FormData();
    formData.append("name", "file");
    formData.append("file", selectedFile);
  
    uploadImage(formData, "artist", artist.uuid, history)
  };

  let artistUUID = false;
  if (route && route.match && route.match.params && route.match.params.uuid) {
    artistUUID = route.match.params.uuid;
  } else {
    artistUUID = false;
  }

  useEffect(() => {
    if (!artistUUID) {
      return false;
    }

    getArtist(artistUUID)
      .then(data => {
        setArtist(data);
        setLoading(false);
      })
  }, [artistUUID])

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

  if (!artistUUID || !artist) {
    return (
      <div className="pageWrapper">
        Unable to fetch user
      </div>
    )
  }

  return (
    <div className="pageWrapper">
      <h1 style={{margin: 0}}>
      {artist.name} {<Link
            key="editbuttonomg"
            to={"/artist/" + artistUUID}
          >unedit</Link>}
      </h1>
      <div className="pageBody" style={{width: `900px`, textAlign: `center`}}>
        <img src={process.env.REACT_APP_API_URL + "/img/" + artist.uuid + "_full.jpg"} alt={artist.name} style={{maxWidth: `300px`, maxHeight: `300px`}}/>

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

export default ArtistEdit;