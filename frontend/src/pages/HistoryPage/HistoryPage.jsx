import React from "react";
import "./HistoryPage.scss";
import { NavLink } from "react-router-dom";

const HistoryPage = ({ selectedImage, onUpload }) => {
  return (
    <>
      <header style={{ display: "flex" }}>
        <NavLink to="/">
          <div className="logo">MeterScan+</div>
        </NavLink>
        <div className="header-right-container">
          <div className="address">Курчатова 30, 78</div>
          <NavLink to="/profile">
            <button className="profile">Профиль</button>
          </NavLink>
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
