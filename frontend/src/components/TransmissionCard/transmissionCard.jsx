import "./transmissionCard.scss";
import { useState } from "react";

const TransmissionCard = (props) => {
  const [value, setValue] = useState(5);

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
      <button className="basic-button black-button" type="submit">
        передать показания
      </button>
    </div>
  );
};

export default TransmissionCard;
