import axios from "axios";
import authHeader from '../Services/auth-header';

class ApiService {
  async getStats() {
    return axios.get(process.env.REACT_APP_API_URL + "stats");
  }

  async getRecentScrobbles(id) {
    return axios.get(process.env.REACT_APP_API_URL + "user/" + id + "/scrobbles", { headers: authHeader() });
  }
}

export default new ApiService();