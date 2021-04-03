import axios from 'axios';
import jwt from 'jwt-decode'
import { toast } from 'react-toastify';

function getHeaders() {
  // TODO: move this to use Context values instead.
  const user = JSON.parse(localStorage.getItem('user'));

  if (user && user.jwt) {
    var unixtime = Math.round((new Date()).getTime() / 1000);
    if (user.exp < unixtime) {
      // TODO: Handle expiry nicer
      toast.warning("Session expired. Please log in again")
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
  return axios.post(process.env.REACT_APP_API_URL + "login", formValues)
    .then((response) => {
      if (response.data.token) {
        let expandedUser = jwt(response.data.token)
        let user = {
          jwt: response.data.token,
          uuid: expandedUser.sub,
          exp: expandedUser.exp,
          username: expandedUser.username,
          admin: expandedUser.admin
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

export const PostRegister = (formValues) => {
  return axios.post(process.env.REACT_APP_API_URL + "register", formValues)
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
  return axios.post(process.env.REACT_APP_API_URL + "resetpassword", formValues)
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
  return axios.post(process.env.REACT_APP_API_URL + "sendreset", values).then(
    (data) => {
      return data.data;
  }).catch((error) => {
    return handleErrorResp(error)
  });
};

export const getStats = () => {
  return axios.get(process.env.REACT_APP_API_URL + "stats").then(
    (data) => {
      return data.data;
  }).catch((error) => {
    return handleErrorResp(error)
  });
};

export const getRecentScrobbles = (id) => {
  return axios.get(process.env.REACT_APP_API_URL + "user/" + id + "/scrobbles", { headers: getHeaders() })
    .then((data) => {
      return data.data;
    }).catch((error) => {
      return handleErrorResp(error)
    });
};

export const getConfigs = () => {
  return axios.get(process.env.REACT_APP_API_URL + "config", { headers: getHeaders() })
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

  return axios.post(process.env.REACT_APP_API_URL + "config", values, { headers: getHeaders() })
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
  return axios.get(process.env.REACT_APP_API_URL + "profile/" + userName, { headers: getHeaders() })
    .then((data) => {
      return data.data;
    }).catch((error) => {
      return handleErrorResp(error)
    });
};

export const getUser = () => {
  return axios.get(process.env.REACT_APP_API_URL + "user", { headers: getHeaders() })
    .then((data) => {
      return data.data;
    }).catch((error) => {
      return handleErrorResp(error)
    });
};

export const patchUser = (values) => {
  return axios.patch(process.env.REACT_APP_API_URL + "user", values, { headers: getHeaders() })
    .then((data) => {
      return data.data;
    }).catch((error) => {
      return handleErrorResp(error)
    });
};

export const validateResetPassword = (tokenStr) => {
  return axios.get(process.env.REACT_APP_API_URL + "user/", { headers: getHeaders() })
    .then((data) => {
      return data.data;
    }).catch((error) => {
      return handleErrorResp(error)
    });
};

export const getSpotifyClientId = () => {
  return axios.get(process.env.REACT_APP_API_URL + "user/spotify", { headers: getHeaders() })
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
  return axios.delete(process.env.REACT_APP_API_URL + "user/spotify", { headers: getHeaders() })
  .then((data) => {
    toast.success(data.data.message);
    return true
  }).catch((error) => {
    return handleErrorResp(error)
  });
}


export const getServerInfo = () => {
  return axios.get(process.env.REACT_APP_API_URL + "serverinfo")
  .then((data) => {
    return data.data
  }).catch((error) => {
    return handleErrorResp(error)
  });
}
