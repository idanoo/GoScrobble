import axios from "axios";
import jwt from 'jwt-decode' // import dependency

class AuthService {
  login(username, password) {
    return axios
      .post(process.env.REACT_APP_API_URL + "login", { username, password })
      .then((response) => {
        if (response.data.token) {
          let expandedUser = jwt(response.data.token)
          let user = {
            jwt: response.data.token,
            uuid: expandedUser.sub,
            exp: expandedUser.exp,
          }
          localStorage.setItem('user', JSON.stringify(user))
        }

        return response.data;
      });
  }

  logout() {
    localStorage.removeItem("user");
  }

  register(username, email, password) {
    return axios.post(process.env.REACT_APP_API_URL + "register", {
      username,
      email,
      password,
    });
  }
}

export default new AuthService();
