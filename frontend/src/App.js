import "./App.css";
import ImportPicture from "./components/ImportPicture/ImportPicture";
import { stages } from "./stages";
import { useState } from "react";

function App() {
  const [stage, changeStage] = useState(stages.upload);
  const [selectedImage, setSelectedImage] = useState(null);

  if (selectedImage && stage === stages.upload) changeStage(stages.analyze);

  console.log(stage);
  console.log(selectedImage);

  if (stage === stages.upload)
    return (
      <>
        <h1>Загрузка фото</h1>
        <ImportPicture
          selectedImage={selectedImage}
          onUpload={setSelectedImage}
        />
      </>
    );
  else if (stage === stages.analyze)
    return (
      <>
        <h1>анализ счетчиков</h1>
        <img src={URL.createObjectURL(selectedImage)} alt="pic lost" />
        <button onClick={() => {changeStage(stages.upload); setSelectedImage(null)}}>поменять фото</button>
        <button onClick={() => {changeStage(stages.send)}}>отправить на анализ</button>
      </>
    );
  else if (stage === stages.send)
    return (
      <>
        <h1>передача показаний</h1>
        <img src={URL.createObjectURL(selectedImage)} alt="pic lost" />
        <p>Текущие показания счетчика</p>
        <input type="number" />
        <button>передать показания</button>
      </>
    );
  else return <p>Что-то не так...</p>;
}

export default App;
