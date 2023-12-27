import axios from "axios";
import authHeader from "./auth-header";

const API_URL = "http://localhost/api/v1/";

class MLService {
  static async predict(image) {
    let form = new FormData();
    form.append("files", image);
    
    const response = (
      await axios.post(API_URL + "predictions", form, {
        headers: authHeader(),
      })
    ).data;

    return response[0]["results"];
  }

  static async updatePredict(id, meter_readings) {
    const form = new FormData();
    form.append("id", id)
    form.append("meter_readings", meter_readings)

    await axios.put(API_URL + "predictions", form, {
      headers: authHeader(),
    })
  }

  //TODO
  static async transmitData(data){
    
  }
}

export default MLService;
