import axios from "axios";
import authHeader from "./auth-header";

const API_URL = "http://localhost:8080/api/v1/";

class MLService {
  static value = 25;
  static async predict(image) {
    let form = new FormData();
    form.append("files", image);
    const response = (await axios.post(API_URL + "predict", form, {
      headers: authHeader(),
    })).data;

    const len = response[0]["results"].length;
    if (len === 1)
      this.value = response[0]["results"]["meter_readings"]
    else{
      const meterReadingsArray = response[0]["results"].map(item => item.meter_readings);
      this.value = meterReadingsArray;
    }
  }
}

export default MLService;
