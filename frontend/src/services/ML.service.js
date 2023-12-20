import axios from "axios";
import authHeader from "./auth-header";

const API_URL = "http://localhost/api/v1/";

class MLService {
  static async predict(image) {
    let form = new FormData();
    form.append("files", image);
    const response = (await axios.post(API_URL + "predictions", form, {
      headers: authHeader(),
    })).data;

    const len = response[0]["results"].length;
    if (len === 1)
      return response[0]["results"]["recognition"]
    else{
      return response[0]["results"].map(item => item.recognition);
    }
  }
}

export default MLService;
