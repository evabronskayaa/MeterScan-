import React from "react";
import "./MenuPage.scss";
import { NavLink } from "react-router-dom";

const MenuPage = () => {
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
      <main>
        <NavLink to="/recognize">
          <div className="card">
            <img
              className="card-image"
              src="img/object-detection.svg"
              alt="no pic("
            />
            Распознать показания по фото
          </div>
        </NavLink>

        <NavLink to="/history">
          <div className="card">
            <img className="card-image" src="img/history.svg" alt="no pic(" />
            История показаний
          </div>
        </NavLink>
        <div className="card">
          <img className="card-image" src="img/invoice.svg" alt="no pic(" />
          Добавить лицевой счет
        </div>
        <div className="card">
          <img
            className="card-image"
            src="img/payment-security.png"
            alt="no pic("
          />
          Автоплатежи
        </div>
        <div className="card">
          <img
            className="card-image"
            src="img/plastic_cards.png"
            alt="no pic("
          />
          Карты
        </div>
        <div className="card">
          <img className="card-image" src="img/group.svg" alt="no pic(" />
          Лицевые счета и группы
        </div>
        <div className="card">
          <img className="card-image" src="img/help.svg" alt="no pic(" />
          Помощь
        </div>
      </main>
    </>
  );
};

export default MenuPage;
