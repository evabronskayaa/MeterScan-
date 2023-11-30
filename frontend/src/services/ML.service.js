import axios from "axios";
import authHeader from "./auth-header";

const API_URL = "http://localhost:8080/api/v1/";

class MLService {
  static value = 25;
  static async predict(image) {
    let form = new FormData();
    form.append("files", image);
    this.value = await axios.post(API_URL + "predict", form, {
      headers: {
        Authorization: authHeader(), 
      },
    });
  }
}

export default MLService;
