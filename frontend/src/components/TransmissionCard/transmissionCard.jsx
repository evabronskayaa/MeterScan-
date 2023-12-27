import "./transmissionCard.scss";
import {useState} from "react";
import MLService from "../../services/ML.service";

const TransmissionCard = (props) => {
  const prediction = props.value
  const [value, setValue] = useState(prediction.recognition)
  const [clicked, setClick] = useState(true)

  const handleChange = (e) => {
    if (clicked) {
      const input = e.target.textContent;
      if (!isNaN(input)) {
        setValue(input);
      } else {
        e.target.textContent = value;
      }
    }
  }

  const handleConfirmation = (value) => {
    MLService.updatePredict(prediction.id, value).then(() => {
      setValue(value)
    })
  }

  const remove = () => {
    MLService.removePredict(prediction.id).then(props.onRemove)
  }

  return <div className="container">
    <div className="row-container justify">
      <div className="numbers">
        <p>Текущие показания счетчика</p>
        <div
            className="custom-input"
            contentEditable="true"
            onInput={handleChange}
        >
          {value}
        </div>
      </div>

      <div className="row-container delete border" onClick={remove}>
        <img src="./img/delete.svg" alt="удалить"/>
        <span>Удалить показание</span>
      </div>
    </div>
    <button
        className={
          clicked ? "basic-button black-button" : "basic-button white-button"
        }
        type="submit"
        onClick={() => {
          setClick(!clicked)
          handleConfirmation(clicked ? value : prediction.recognition)
        }}
    >
      {clicked ? "подтвердить показания" : "отменить подтверждение показаний"}
    </button>
  </div>
};

export default TransmissionCard;
