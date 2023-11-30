import "./transmissionCard.scss";
import { useState } from "react";

const TransmissionCard = (props) => {
  const [value, setValue] = useState(5);
  const [clicked, setClick] = useState(true);

  const handleChange = (e) => {
    const input = e.target.textContent;
    if (!isNaN(input)) {
      setValue(input);
    } else {
      e.target.textContent = value;
    }
  };

  return (
    <div className="container">
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

        <div className="row-container delete border">
          <img src="./img/delete.svg" alt="удалить" />
          <span>Удалить показание</span>
        </div>
      </div>
      <button
        className={
          clicked ? "basic-button black-button" : "basic-button white-button"
        }
        type="submit"
        onClick={() => setClick(!clicked)}
      >
        {clicked ? "подтвердить показания" : "отменить подтверждение показаний"}
      </button>
    </div>
  );
};

export default TransmissionCard;
