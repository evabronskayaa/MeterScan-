import axios from "axios";

const API_URL = "http://localhost:8080/api/v1/";

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

  getCurrentUser() {
    return JSON.parse(localStorage.getItem('user'))?.user?.email;
  }
}

const authService = new AuthService()

export default authService;