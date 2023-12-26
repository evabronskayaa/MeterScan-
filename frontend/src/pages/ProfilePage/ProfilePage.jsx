import React from "react";
import "./ProfilePage.scss";

const ProfilePage = ({ selectedImage, onUpload }) => {
  return (
    <>
      <header style={{ display: "flex" }}>
        <div className="logo">MeterScan+</div>
        <div className="header-right-container">
          <div className="address">Курчатова 30, 78</div>
          <button className="profile">Профиль</button>
        </div>
      </header>
      <main class="profile-container">
        <div className="title">Личный кабинет</div>
        <div className="profile-value-row-with-button">
          <div className="profile-value-row">
            <div className="profile-value-row-container">л/с 4709040404</div>
            <div className="profile-value-row-container">23-12-2023</div>
            <div className="profile-value-row-container">горячая вода</div>
            <div className="profile-value-row-container">26 куб</div>
          </div>
          <div className="row-container delete border round">
            <img src="./img/delete.svg" alt="удалить" />
            <span>Удалить показание</span>
          </div>
        </div>

        <div className="profile-value-row-with-button">
          <div className="profile-value-row">
            <div className="profile-value-row-container">л/с 4709040404</div>
            <div className="profile-value-row-container">23-12-2023</div>
            <div className="profile-value-row-container">горячая вода</div>
            <div className="profile-value-row-container">26 куб</div>
          </div>
          <div className="row-container delete border round">
            <img src="./img/delete.svg" alt="удалить" />
            <span>Удалить показание</span>
          </div>
        </div>

        <div className="row-container delete border round">
          <img src="./img/icon-plus.png" alt="добавить" width={20} />
          <span>Добавить новый лицевой счет</span>
        </div>
      </main>
    </>
  );
};

export default ProfilePage;
