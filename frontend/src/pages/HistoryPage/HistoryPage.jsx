import React, {useEffect, useState} from "react";
import "./HistoryPage.scss";
import MLService from "../../services/ML.service";
import Plot from 'react-plotly.js';

const HistoryPage = () => {
  const [predictions, setPredictions] = useState([])

  useEffect(() => {
    MLService.getPredictions().then(data => setPredictions(data.sort((a, b) => a.created_at - b.created_at)))
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
        {value.results && value.results.length === 1 && index > 0 && predictions[index - 1].results && predictions[index - 1].results.length == 1 &&
            <div className="value-row-container">{"Разница по сравнению с предыдущим месяцем: " + (Number(value.results[0].recognition) - Number(predictions[index - 1].results[0].recognition))}</div>
        }
      </div>
    })}

    <Plot
        data={[
          {
            x: Array.from(predictions, prediction => new Date(prediction.created_at * 1000)),
            y: Array.from(predictions, prediction => prediction.results ? Number(prediction.results[0]?.recognition) : 0),
            type: 'scatter'
          },
        ]}
        layout={{
          xaxis: {
            title: {
              text: 'Дата',
            },
          },
          yaxis: {
            title: {
              text: 'Кубы',
            },
          }
        }}
    />
  </div>
};

export default HistoryPage;
