import axios from 'axios';
import jwt from 'jwt-decode'
import { toast } from 'react-toastify';

function getHeaders() {
  const user = JSON.parse(localStorage.getItem('user'));

  if (user && user.jwt) {
    var unixtime = Math.round((new Date()).getTime() / 1000);
    if (user.exp < unixtime) {
      // Trigger refresh
      localStorage.removeItem('user');
      window.location.reload();
      // toast.warning("Session expired. Please log in again")
      // window.location.reload();
      return {};
    }

    return { Authorization: 'Bearer ' + user.jwt };
  } else {
    return {};
  }
}

function getUserUuid() {
  // TODO: move this to use Context values instead.
  const user = JSON.parse(localStorage.getItem('user'));

  if (user && user.uuid) {
    return user.uuid
  } else {
    return '';
  }
}

function handleErrorResp(error) {
  if (error.response) {
    if (error.response.status === 401)  {
      toast.error('Unauthorized')
    } else if (error.response.status === 429) {
      toast.error('Rate limited. Please try again shortly')
    } else {
      toast.error('An unknown error has occurred');
    }
  } else {
    toast.error('Failed to connect to API');
  }
  return {};
}

export const PostLogin = (formValues) => {
  return axios.post(process.env.REACT_APP_API_URL + "/api/v1/login", formValues)
    .then((response) => {
      if (response.data.token) {
        let expandedUser = jwt(response.data.token)
        let user = {
          jwt: response.data.token,
          uuid: expandedUser.sub,
          exp: expandedUser.exp,
          username: expandedUser.username,
          admin: expandedUser.admin,
          refresh_token: expandedUser.refresh_token,
          refresh_exp: expandedUser.refresh_exp,
        }

        toast.success('Successfully logged in.');
        return user;
      } else {
        toast.error(response.data.error ? response.data.error: 'An Unknown Error has occurred');
        return null
      }
    }).catch((error) => {
      if (error.response === 401)  {
        toast.error('Unauthorized')
      } else if (error.response === 429) {
        toast.error('Rate limited. Please try again shortly')
      } else {
        toast.error('Failed to connect');
      }
      return Promise.resolve();
    });
};

export const PostRefreshToken = (refreshToken) => {
  return axios.post(process.env.REACT_APP_API_URL + "/api/v1/refresh", {token: refreshToken})
    .then((response) => {
      if (response.data.token) {
        let expandedUser = jwt(response.data.token)
        let user = {
          jwt: response.data.token,
          uuid: expandedUser.sub,
          exp: expandedUser.exp,
          username: expandedUser.username,
          admin: expandedUser.admin,
          refresh_token: expandedUser.refresh_token,
          refresh_exp: expandedUser.refresh_exp,
        }

        return user;
      } else {
        toast.error(response.data.error ? response.data.error: 'An Unknown Error has occurred');
        return null
      }
    }).catch((error) => {
      if (error.response === 401)  {
        toast.error('Unauthorized')
      } else if (error.response === 429) {
        toast.error('Rate limited. Please try again shortly')
      } else {
        toast.error('Failed to connect');
      }
      return Promise.resolve();
    });
};


export const PostRegister = (formValues) => {
  return axios.post(process.env.REACT_APP_API_URL + "/api/v1/register", formValues)
    .then((response) => {
      if (response.data.message) {
        toast.success(response.data.message);

        return Promise.resolve();
      } else {
        toast.error(response.data.error ? response.data.error: 'An Unknown Error has occurred');

        return Promise.reject();
      }
    }).catch((error) => {
      handleErrorResp(error)
      return Promise.resolve();
  });
};

export const PostResetPassword = (formValues) => {
  return axios.post(process.env.REACT_APP_API_URL + "/api/v1/resetpassword", formValues)
    .then((response) => {
      if (response.data.message) {
        toast.success(response.data.message);

        return Promise.resolve();
      } else {
        toast.error(response.data.error ? response.data.error: 'An Unknown Error has occurred');

        return Promise.reject();
      }
    }).catch((error) => {
      handleErrorResp(error)
      return Promise.resolve();
    });
};

export const sendPasswordReset = (values) => {
  return axios.post(process.env.REACT_APP_API_URL + "/api/v1/sendreset", values).then(
    (data) => {
      return data.data;
  }).catch((error) => {
    return handleErrorResp(error)
  });
};

export const getStats = () => {
  return axios.get(process.env.REACT_APP_API_URL + "/api/v1/stats").then(
    (data) => {
      return data.data;
  }).catch((error) => {
    return handleErrorResp(error)
  });
};

export const getRecentScrobbles = (id) => {
  return axios.get(process.env.REACT_APP_API_URL + "/api/v1/user/" + id + "/scrobbles", { headers: getHeaders() })
    .then((data) => {
      return data.data;
    }).catch((error) => {
      return handleErrorResp(error)
    });
};

export const getConfigs = () => {
  return axios.get(process.env.REACT_APP_API_URL + "/api/v1/config", { headers: getHeaders() })
    .then((data) => {
      return data.data;
    }).catch((error) => {
      return handleErrorResp(error)
    });
};

export const postConfigs = (values, toggle) => {
  if (toggle) {
    values.REGISTRATION_ENABLED = "1"
  } else {
    values.REGISTRATION_ENABLED = "0"
  }

  return axios.post(process.env.REACT_APP_API_URL + "/api/v1/config", values, { headers: getHeaders() })
    .then((data) => {
      if (data.data && data.data.message) {
        toast.success(data.data.message);
      } else if (data.data && data.data.error) {
        toast.error(data.data.error);
      }
    }).catch((error) => {
      return handleErrorResp(error)
    });
};

export const getProfile = (userName) => {
  return axios.get(process.env.REACT_APP_API_URL + "/api/v1/profile/" + userName, { headers: getHeaders() })
    .then((data) => {
      return data.data;
    }).catch((error) => {
      return handleErrorResp(error)
    });
};

export const getUser = () => {
  return axios.get(process.env.REACT_APP_API_URL + "/api/v1/user", { headers: getHeaders() })
    .then((data) => {
      return data.data;
    }).catch((error) => {
      return handleErrorResp(error)
    });
};

export const patchUser = (values) => {
  return axios.patch(process.env.REACT_APP_API_URL + "/api/v1/user", values, { headers: getHeaders() })
    .then((data) => {
      return data.data;
    }).catch((error) => {
      return handleErrorResp(error)
    });
};

export const validateResetPassword = (tokenStr) => {
  return axios.get(process.env.REACT_APP_API_URL + "/api/v1/user/", { headers: getHeaders() })
    .then((data) => {
      return data.data;
    }).catch((error) => {
      return handleErrorResp(error)
    });
};

export const getSpotifyClientId = () => {
  return axios.get(process.env.REACT_APP_API_URL + "/api/v1/user/spotify", { headers: getHeaders() })
    .then((data) => {
      return data.data
    }).catch((error) => {
      return handleErrorResp(error)
    });
}

export const spotifyConnectionRequest = () => {
  return getSpotifyClientId().then((resp) => {
    var scopes = 'user-read-recently-played user-read-currently-playing';

    // Local, lets forward it to API
    let redirectUri = window.location.origin.toString()+ "/api/v1/link/spotify";

    // Stupid dev
    if (window.location.origin.toString() === "http://localhost:3000") {
      redirectUri = "http://localhost:42069/api/v1/link/spotify"
    }

    window.location = 'https://accounts.spotify.com/authorize' +
      '?response_type=code' +
      '&client_id=' + resp.token +
      '&scope=' + encodeURIComponent(scopes) +
      '&redirect_uri=' + encodeURIComponent(redirectUri) +
      '&state=' + getUserUuid();
  })
};

export const spotifyDisonnectionRequest = () => {
  return axios.delete(process.env.REACT_APP_API_URL + "/api/v1/user/spotify", { headers: getHeaders() })
  .then((data) => {
    toast.success(data.data.message);
    return true
  }).catch((error) => {
    return handleErrorResp(error)
  });
}

export const navidromeConnectionRequest = (values) => {
  return axios.post(process.env.REACT_APP_API_URL + "/api/v1/user/navidrome", values, { headers: getHeaders() })
  .then((data) => {
    toast.success(data.data.message);
    return true
  }).catch((error) => {
    return handleErrorResp(error)
  });
};

export const navidromeDisonnectionRequest = () => {
  return axios.delete(process.env.REACT_APP_API_URL + "/api/v1/user/navidrome", { headers: getHeaders() })
  .then((data) => {
    toast.success(data.data.message);
    return true
  }).catch((error) => {
    return handleErrorResp(error)
  });
}


export const getServerInfo = () => {
  return axios.get(process.env.REACT_APP_API_URL + "/api/v1/serverinfo")
  .then((data) => {
    return data.data
  }).catch((error) => {
    return handleErrorResp(error)
  });
}

export const getArtist = (uuid) => {
  return axios.get(process.env.REACT_APP_API_URL + "/api/v1/artists/" + uuid).then(
    (data) => {
      return data.data;
  }).catch((error) => {
    return handleErrorResp(error)
  });
};

export const getAlbum = (uuid) => {
  return axios.get(process.env.REACT_APP_API_URL + "/api/v1/albums/" + uuid).then(
    (data) => {
      return data.data;
  }).catch((error) => {
    return handleErrorResp(error)
  });
};

export const getTrack = (uuid) => {
  return axios.get(process.env.REACT_APP_API_URL + "/api/v1/tracks/" + uuid).then(
    (data) => {
      return data.data;
  }).catch((error) => {
    return handleErrorResp(error)
  });
};

export const getTopTracks = (uuid) => {
  return axios.get(process.env.REACT_APP_API_URL + "/api/v1/tracks/top/" + uuid).then(
    (data) => {
      return data.data;
  }).catch((error) => {
    return handleErrorResp(error)
  });
}

export const getTopArtists = (uuid) => {
  return axios.get(process.env.REACT_APP_API_URL + "/api/v1/artists/top/" + uuid).then(
    (data) => {
      return data.data;
  }).catch((error) => {
    return handleErrorResp(error)
  });
}