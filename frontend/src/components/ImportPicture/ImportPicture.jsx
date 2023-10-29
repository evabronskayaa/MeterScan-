import React from "react";
import "./importPicture.css";

const ImportPicture = ({ selectedImage, onUpload }) => {
  return (
    <div>
      {!selectedImage && (
        <div className="container">
          <img src="img/default_picture.png" alt="no pic" />
          <input
            type="file"
            name="myImage"
            onChange={(event) => {
              onUpload(event.target.files[0]);
            }}
          />
          <button className="button">
            <img src="img/img-placeholder.svg" alt="" />
            <span>загрузить фото с устройства</span>
          </button>
        </div>
      )}
      {selectedImage && (
        <div>
          <img
            alt="not found"
            width={"250px"}
            src={URL.createObjectURL(selectedImage)}
          />
        </div>
      )}
    </div>
  );
};

export default ImportPicture;
