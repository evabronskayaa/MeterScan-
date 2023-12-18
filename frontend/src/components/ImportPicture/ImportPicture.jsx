import React, {useRef} from "react";
import './importPicture.scss';

const ImportPicture = ({selectedImage, onUpload}) => {
  const fileInputRef=useRef();

  return (
    <div>
      {!selectedImage && (
        <div className="container">
        <img className="image" src="img/default_picture.png" alt="no pic" />
        <input
          type="file"
          name="myImage"
          onChange={(event) => {
            onUpload(event.target.files[0]);
          }}
          ref={fileInputRef}
          hidden
        />
        <button className="button" onClick={()=>fileInputRef.current.click()}>
          <img src="img/img-placeholder.svg" alt="" />
          <span>загрузить фото с устройства</span>
        </button>
      </div>
      )}
    </div>
  );
}

export default ImportPicture