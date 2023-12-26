import React from "react";
import "./HistoryPage.scss";

const HistoryPage = ({ selectedImage, onUpload }) => {
  return (
    <>
      <header style={{ display: "flex" }}>
        <div className="logo">MeterScan+</div>
        <div className="header-right-container">
          <div className="address">Курчатова 30, 78</div>
          <button className="profile">Профиль</button>
        </div>
      </header>
      <main className="profile-container">
        <div className="title">История показаний</div>
        <div className="value-row">
          <div className="value-row-container">л/с 4709040404</div>
          <div className="value-row-container">23-12-2023</div>
          <div className="value-row-container">горячая вода</div>
          <div className="value-row-container">26 куб</div>
        </div>

        <div className="value-row">
          <div className="value-row-container">л/с 4709040404</div>
          <div className="value-row-container">23-12-2023</div>
          <div className="value-row-container">горячая вода</div>
          <div className="value-row-container">26 куб</div>
        </div>
      </main>
    </>
  );
};

export default HistoryPage;
