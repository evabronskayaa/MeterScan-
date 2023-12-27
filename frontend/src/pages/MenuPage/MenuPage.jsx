import React from "react";
import "./MenuPage.scss";

const MenuPage = ({ selectedImage, onUpload }) => {
  return (
    <>
      <header style={{ display: "flex" }}>
        <div className="logo">MeterScan+</div>
        <div className="header-right-container">
          <div className="address">Курчатова 30, 78</div>
          <button className="profile">Профиль</button>
        </div>
      </header>
      <main>
        <div className="card">
            <img className="card-image" src="img/object-detection.svg"/>
          Распознать показания по фото
        </div>
        <div className="card">
            <img className="card-image" src="img/history.svg"/>
          История показаний
        </div>
        <div className="card">
            <img className="card-image" src="img/invoice.svg"/>
          Добавить лицевой счет
        </div>
        <div className="card">
            <img className="card-image" src="img/payment-security.png"/>
          Автоплатежи
        </div>
        <div className="card">
            <img className="card-image" src="img/plastic_cards.png"/>
          Карты
        </div>
        <div className="card">
            <img className="card-image" src="img/group.svg"/>
          Лицевые счета и группы
        </div>
        <div className="card">
            <img className="card-image" src="img/help.svg"/>
          Помощь
        </div>
      </main>
    </>
  );
};

export default MenuPage;
