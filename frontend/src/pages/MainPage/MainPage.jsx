import ImportPicture from "../../components/ImportPicture/ImportPicture";
import { stages } from "../../stages";
import { useState } from "react";
import "slick-carousel/slick/slick.css";
import "slick-carousel/slick/slick-theme.css";
import Slider from "react-slick";
import TransmissionCard from "../../components/TransmissionCard/transmissionCard";
import "./MainPage.scss";
import authService from "../../services/auth.service";
import MLService from "../../services/ML.service";

const MainPage = (props) => {
  const [stage, changeStage] = useState(stages.upload);
  const [selectedImage, setSelectedImage] = useState(null);
  const [showModal, setShowModal] = useState(false);
  const [value, setValue] = useState();

  const handleLogout = (e) => {
    authService.logout();
    window.location.reload();
  };

  const handleConfirmation = (index, newValue) => {
    const updatedValue = [...value];
    updatedValue[0].results[index].valid_meter_readings = newValue;
    setValue(updatedValue);
    console.log(value);
  };

  if (selectedImage && stage === stages.upload) changeStage(stages.analyze);

  if (stage === stages.upload)
    return (
      <>
        <p className="title main-title">MeterScan+</p>
        <div className="user">
          <span>{props.email}</span>
          <button className="logout" onClick={handleLogout}>
            Выйти
          </button>
        </div>
        <div className="form">
          <div className="row">
            <p className="title">Загрузка фото</p>
            <img
              src="img/tooltip.svg"
              alt="info"
              className="tooltip"
              onClick={() => setShowModal(true)}
            />
          </div>
          <ImportPicture
            selectedImage={selectedImage}
            onUpload={setSelectedImage}
          />
        </div>
        {showModal && (
          <div className="modal" onClick={() => setShowModal(false)}>
            <ul onClick={(e) => e.stopPropagation()}>
              <img
                src="img/close.svg"
                alt="close"
                onClick={() => setShowModal(false)}
              />
              <h3>Важно</h3>
              <li>Загружайте фото только с одним счетчиком</li>
              <li>когда фотографируете в темноте, включайте вспышку</li>
              <li>Фото должно быть не смазанным</li>
            </ul>
          </div>
        )}
      </>
    );
  else if (stage === stages.analyze)
    return (
      <>
        <p className="title main-title">MeterScan+</p>
        <div className="user">
          <span>{props.email}</span>
          <button className="logout" onClick={handleLogout}>
            Выйти
          </button>
        </div>
        <div className="form">
          <p className="title">анализ счетчиков</p>
          <img
            className="image"
            src={URL.createObjectURL(selectedImage)}
            alt="pic lost"
          />
          <button
            className="basic-button black-button"
            onClick={() => {
              MLService.predict(selectedImage).then((r) => {
                changeStage(stages.send);
                setValue(r);
              });
            }}
          >
            отправить на анализ
          </button>
          <button
            className="basic-button white-button"
            onClick={() => {
              changeStage(stages.upload);
              setSelectedImage(null);
            }}
          >
            поменять фото
          </button>
        </div>
      </>
    );
  else if (stage === stages.send)
    return (
      <>
        <p className="title main-title">MeterScan+</p>
        <div className="user">
          <span>{props.email}</span>
          <button className="logout" onClick={handleLogout}>
            Выйти
          </button>
        </div>
        <div className="form">
          <p className="title">передача показаний</p>
          <img
            className="image"
            src={URL.createObjectURL(selectedImage)}
            alt="pic lost"
          />

          <Slider className="carousel border" dots accessibility={false}>
            {value[0].results.map((item, index) => (
              <div key={index} className="carousel-item">
                <TransmissionCard
                  onConfirmation={(newValue) =>
                    handleConfirmation(index, newValue)
                  }
                  value={item.valid_meter_readings}
                />
              </div>
            ))}
            <div className="carousel-item">
              <p>Если вы подтвердили все показания, то самое время</p>
              <button className="basic-button black-button">
                Передать показания
              </button>
            </div>
          </Slider>
        </div>
      </>
    );
  else return <p>Что-то не так...</p>;
};

export default MainPage;
