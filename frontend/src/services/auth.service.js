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
          localStorage.setItem("user", JSON.stringify(response.data));
          console.log(JSON.stringify(response.data))
        }

        return response.data;
      });
  }

  logout() {
    localStorage.removeItem("user");
  }

  register(email, password, recaptcha) {
    return axios.post(API_URL + "users", {
      email,
      password,
      recaptcha
    });
  }

  me() {
    if (authService.getCurrentUser()) {
      axios.get(API_URL + "me", {
        headers: authHeader(),
      }).then(response => {
        return response.data;
      }).then(data => {
        const localStorageItem = JSON.parse(localStorage.getItem("user"))
        localStorageItem.user = data
        localStorage.setItem("user", JSON.stringify(localStorageItem));
      }).catch(() => {
        authService.logout()
      })
    }
  }

  getCurrentUser() {
    return JSON.parse(localStorage.getItem('user'))?.user;
  }
}

const authService = new AuthService()

export default authService;