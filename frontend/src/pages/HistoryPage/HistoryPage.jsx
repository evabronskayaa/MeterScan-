import React, {useEffect, useState} from "react";
import "./HistoryPage.scss";
import MLService from "../../services/ML.service";

const HistoryPage = () => {
  const [predictions, setPredictions] = useState([])

  useEffect(() => {
    MLService.getPredictions().then(setPredictions)
  }, [])

  return <div>
    <div className="title">История показаний</div>

    {predictions.map((value, index) => {
      return <div key={index} className="value-row">
        <div className="value-row-container">{new Date(value.created_at * 1000).toLocaleString()}</div>
        {value.results && <div className="value-row-container">Показания:</div>}
        {value.results?.map((result) => {
          return <div className="value-row-container">{result.recognition}</div>
        })}
      </div>
    })}
  </div>
};

export default HistoryPage;
