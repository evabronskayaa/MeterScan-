import "./App.scss";
import ImportPicture from "./components/ImportPicture/ImportPicture";
import { stages } from "./stages";
import { useState } from "react";
import { uploadImage, number } from "./requests";
import "slick-carousel/slick/slick.css";
import "slick-carousel/slick/slick-theme.css";
import Slider from "react-slick";
import TransmissionCard from "./components/TransmissionCard/transmissionCard";

function App() {
  const [stage, changeStage] = useState(stages.upload);
  const [selectedImage, setSelectedImage] = useState(null);
  const [showModal, setShowModal] = useState(false);

  if (selectedImage && stage === stages.upload) changeStage(stages.analyze);

  console.log(stage);
  console.log(selectedImage);

  if (stage === stages.upload)
    return (
      <>
        <p className="title main-title">MeterScan+</p>
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
        {/* <p className="title">MeterScan+</p> */}
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
              changeStage(stages.send);
              uploadImage(selectedImage);
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
        <div className="form">
          <p className="title">передача показаний</p>
          <img
            className="image"
            src={URL.createObjectURL(selectedImage)}
            alt="pic lost"
          />
          {typeof number === 'object' ? (
            <>
            <Slider className="carousel" dots arrows>
              {number.map((item, index) => (
                <div key={index} className="carousel-item">
                <TransmissionCard number={item}/>
              </div>
              ))}
              
            </Slider>
            </>
          ) : (
            <>
              <TransmissionCard number={number}/>
            </>
          )}
        </div>
      </>
    );
  else return <p>Что-то не так...</p>;
}

export default App;
