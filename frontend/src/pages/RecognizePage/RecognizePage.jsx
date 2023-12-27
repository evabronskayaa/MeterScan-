import ImportPicture from "../../components/ImportPicture/ImportPicture";
import {stages} from "../../stages";
import {useCallback, useEffect, useState} from "react";
import {NavLink} from "react-router-dom";
import "slick-carousel/slick/slick.css";
import "slick-carousel/slick/slick-theme.css";
import Slider from "react-slick";
import TransmissionCard from "../../components/TransmissionCard/transmissionCard";
import "./RecognizePage.scss";
import MLService from "../../services/ML.service";

const RecognizePage = () => {
  const [stage, changeStage] = useState(stages.upload);
  const [selectedImage, setSelectedImage] = useState(null);
    const [imageWithScope, setImageWithScope] = useState(null)
  const [showModal, setShowModal] = useState(false);
    const [value, setValue] = useState([]);

    const [currentSlide, setCurrentSlide] = useState(0);

    const handleBeforeChange = (oldIndex, newIndex) => {
        setImageWithScope(null);
        setCurrentSlide(newIndex);
    };

    const applyRectangles = useCallback((img, rectangle) => {
      const canvas = document.createElement("canvas");
      const context = canvas.getContext("2d");
      const image = new Image();
      image.src = URL.createObjectURL(img)

      image.onload = () => {
        canvas.width = image.width;
        canvas.height = image.height;
        context.drawImage(image, 0, 0);

        context.strokeStyle = "lime";
        context.lineWidth = 3;
          context.strokeRect(rectangle.x1, rectangle.y1, rectangle.x2 - rectangle.x1, rectangle.y2 - rectangle.y1);
          context.fillStyle = "lime";

        canvas.toBlob((blob) => {
            setImageWithScope(blob);
        });
      };
    }, []
  );

  useEffect(() => {
    if (selectedImage && stage === stages.send && value[currentSlide]?.scope) {
      applyRectangles(selectedImage, value[currentSlide].scope
      );
    }
  }, [selectedImage, stage, value, applyRectangles, currentSlide]);

  if (selectedImage && stage === stages.upload) changeStage(stages.analyze);

  if (stage === stages.upload)
    return (
      <>
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
        <div className="form">
          <p className="title">передача показаний</p>
          <img
            className="image image-carousel"
            src={URL.createObjectURL(imageWithScope !== null ? imageWithScope : selectedImage)}
            alt="pic lost"
          />

            <Slider className="carousel border" dots accessibility={false} beforeChange={handleBeforeChange}>
            {value.map((item, index) => (
              <div key={index} className="carousel-item">
                  <TransmissionCard value={item} onRemove={() => {
                      setImageWithScope(null)
                      const newValue = value.filter(value => value.id !== item.id)
                      setValue(newValue)
                      if (newValue.length === 0) {
                          changeStage(stages.upload)
                          setSelectedImage(null)
                          setValue(null)
                      }
                  }}/>
              </div>
            ))}
            <div className="carousel-item">
              <p>Если вы подтвердили все показания, то самое время</p>
              <button
                className="basic-button black-button"
                onClick={() => {
                  changeStage(stages.sent);
                }}
              >
                Передать показания
              </button>
            </div>
          </Slider>
        </div>
      </>
    );
  else if (stage === stages.sent)
    return (
      <>
        <p>Показания переданы!</p>
        <NavLink to="/">Вернуться на главную</NavLink>
      </>
    );
  else return <p>Что-то не так...</p>;
};

export default RecognizePage;
