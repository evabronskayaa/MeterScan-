import React from "react";

const ImportPicture = ({selectedImage, onUpload}) => {
  return (
    <div>
      {!selectedImage && (
        <img src="img/default_picture.png" alt="no pic"/>
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
      <br />
      <br />
      <input
        type="file"
        name="myImage"
        onChange={(event) => {
          onUpload(event.target.files[0]);
        }}
      />
    </div>
  );
}

export default ImportPicture