import axios from 'axios';
import jwt from 'jwt-decode'
import { toast } from 'react-toastify';

function getHeaders() {
  // Todo: move this to use Context values instead.
  const user = JSON.parse(localStorage.getItem('user'));

  if (user && user.jwt) {
    return { Authorization: 'Bearer ' + user.jwt };
  } else {
    return {};
  }
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
    })
    .catch(() => {
      toast.error('Failed to connect');
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
    })
    .catch(() => {
      toast.error('Failed to connect');
      return Promise.resolve();
    });
};

export const getStats = () => {
  return axios.get(process.env.REACT_APP_API_URL + "stats").then(
    (data) => {
      return data.data;
    }
  ).catch(() => {
    toast.error('Failed to connect');
    return {};
  });
};

export const getRecentScrobbles = (id) => {
  return axios.get(process.env.REACT_APP_API_URL + "user/" + id + "/scrobbles", { headers: getHeaders() })
    .then((data) => {
      return data.data;
    }).catch(() => {
      toast.error('Failed to connect');
      return {};
    });
};

export const getConfigs = () => {
  return axios.get(process.env.REACT_APP_API_URL + "config", { headers: getHeaders() })
    .then((data) => {
      return data.data;
    }).catch(() => {
      toast.error('Failed to connect');
      return {};
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
    })
    .catch(() => {
      toast.error('Error updating values');
    });
};

export const getProfile = (userName) => {
  return axios.get(process.env.REACT_APP_API_URL + "profile/" + userName, { headers: getHeaders() })
    .then((data) => {
      return data.data;
    }).catch(() => {
      toast.error('Failed to connect');
      return {};
    });
};

export const getUser = () => {
  return axios.get(process.env.REACT_APP_API_URL + "user", { headers: getHeaders() })
    .then((data) => {
      return data.data;
    }).catch(() => {
      toast.error('Failed to connect');
      return {};
    });
};
