import "../App.css";
import ImportPicture from "../components/ImportPicture/ImportPicture";
import { stages } from "../stages";
import { useState } from "react";

function FormPage() {
  const [stage, changeStage] = useState(stages.upload);
  const [selectedImage, setSelectedImage] = useState(null);

  if (selectedImage && stage === stages.upload) changeStage(stages.analyze);

  console.log(stage);
  console.log(selectedImage);

  if (stage === stages.upload)
    return (
      <div className="form">
        <h1>Загрузка фото</h1>
        <ImportPicture
          selectedImage={selectedImage}
          onUpload={setSelectedImage}
        />
      </div>
    );
  else if (stage === stages.analyze)
    return (
      <div className="form">
        <h1>анализ счетчиков</h1>
        <img src={URL.createObjectURL(selectedImage)} alt="pic lost" />
        <button onClick={() => {changeStage(stages.upload); setSelectedImage(null)}}>поменять фото</button>
        <button onClick={() => {changeStage(stages.send)}}>отправить на анализ</button>
      </div>
    );
  else if (stage === stages.send)
    return (
      <div className="form">
        <h1>передача показаний</h1>
        <img src={URL.createObjectURL(selectedImage)} alt="pic lost" />
        <p>Текущие показания счетчика</p>
        <input type="number" />
        <button>передать показания</button>
      </div>
    );
  else return <p>Что-то не так...</p>;
}

export default FormPage;
