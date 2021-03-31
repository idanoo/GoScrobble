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
  // const { setLoading, setUser } = useContext(AuthContext);
  // setLoading(true)
  return axios.post(process.env.REACT_APP_API_URL + "login", formValues)
    .then((response) => {
      if (response.data.token) {
        let expandedUser = jwt(response.data.token)
        let user = {
          jwt: response.data.token,
          uuid: expandedUser.sub,
          exp: expandedUser.exp,
          username: expandedUser.username,
        }

        toast.success('Successfully logged in.');
        return user;
      } else {
        toast.error(response.data.error ? response.data.error: 'An Unknown Error has occurred');
        return null
      }
    })
    .catch(() => {
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
      return Promise.resolve();
    });
};

export const getStats = () => {
  return axios.get(process.env.REACT_APP_API_URL + "stats").then(
    (data) => {
      data.isLoading = false;
      return data.data;
    }
  );
};

export const getRecentScrobbles = (id) => {
  return axios.get(process.env.REACT_APP_API_URL + "user/" + id + "/scrobbles", { headers: getHeaders() })
    .then((data) => {
      return data.data;
    });
};