import axios from "axios";
import jwt from 'jwt-decode' // import dependency

class AuthService {
  login(username, password) {
    return axios
      .post(process.env.REACT_APP_API_URL + "login", { username, password })
      .then((response) => {
        if (response.data.token) {
          let user = jwt(response.data.token)
          localStorage.setItem("jwt", response.data.token);
          localStorage.setItem("uuid", user.sub);
          localStorage.setItem("exp", user.exp);
        }

        return response.data;
      });
  }

  logout() {
    localStorage.removeItem("jwt");
    localStorage.removeItem("uuid");
    localStorage.removeItem("exp");
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
