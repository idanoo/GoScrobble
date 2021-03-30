import axios from "axios";

class ApiService {
  async getStats() {
    return axios.get(process.env.REACT_APP_API_URL + "stats");
  }
}

export default new ApiService();