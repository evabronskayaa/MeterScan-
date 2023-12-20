import axios from "axios";
import authHeader from "./auth-header";

const API_URL = "http://localhost:8080/api/v1/";

class MLService {
  static async predict(image) {
    let form = new FormData();
    form.append("files", image);
    return [
      {
        image_name: "2091c4d9-9269-4963-8af0-9a3d416313a2",
        results: [
          {
            id: 20,
            scope: {
              x1: 312,
              y1: 792,
              x2: 524,
              y2: 841,
            },
            meter_readings: "00027",
            valid_meter_readings: "00027",
            metric: 0.9798002,
          },
          {
            id: 21,
            scope: {
              x1: 424,
              y1: 244,
              x2: 623,
              y2: 292,
            },
            meter_readings: "00025004",
            valid_meter_readings: "00025004",
            metric: 0.78123254,
          },
        ],
      },
    ];

    const response = (
      await axios.post(API_URL + "predictions", form, {
        headers: authHeader(),
      })
    ).data;

    const len = response[0]["results"].length;
    if (len === 1) return response[0]["results"]["meter_readings"];
    else {
      return response[0]["results"].map((item) => item.meter_readings);
    }
  }

  //TODO
  static async transmitData(data){
    
  }
}

export default MLService;
