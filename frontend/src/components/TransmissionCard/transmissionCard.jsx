import "./transmissionCard.css"

const TransmissionCard = (props) => {
  return (
    <div className="container">
      <div className="numbers">
        <p>Текущие показания счетчика</p>
        <input className="input" type="number" value={props.number} />
      </div>
      <button className="basic-button black-button" type="submit">
        передать показания
      </button>
    </div>
  );
};

export default TransmissionCard;
