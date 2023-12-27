import axios from "axios";
import authHeader from "./auth-header";

const API_URL = "http://localhost/api/v1/";

class AuthService {
  login(email, password, recaptcha) {
    return axios
      .post(API_URL + "sessions", {
        email,
        password,
        recaptcha
      })
      .then(response => {
        console.log(JSON.stringify(response.data))
        if (response.data.token) {
          localStorage.setItem("token", JSON.stringify(response.data.token));
          console.log(JSON.stringify(response.data))
        }

        return response.data;
      });
  }

  logout() {
    localStorage.removeItem("token");
  }

  register(email, password, recaptcha) {
    return axios.post(API_URL + "users", {
      email,
      password,
      recaptcha
    });
  }

  me() {
    if (localStorage.getItem("token")) {
      return axios.get(API_URL + "me", {
        headers: authHeader(),
      }).then(response => {
        return response.data;
      })
    } else {
      return null
    }
  }
}

const authService = new AuthService()

export default authService;