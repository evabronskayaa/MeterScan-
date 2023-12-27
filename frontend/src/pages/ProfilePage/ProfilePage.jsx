import React from "react";
import "./ProfilePage.scss";
import Toggle from "react-toggle";
import "react-toggle/style.css";
import TimePicker from "react-time-picker";
import "react-time-picker/dist/TimePicker.css";
import "react-clock/dist/Clock.css";
import { useState } from "react";

const ProfilePage = () => {
  const [time, setTime] = useState("10:00");
  const [day, setDay] = useState([1]);
  const [sendNotifications, setSend] = useState(false);

  const CheckRange = (value) =>{
    console.log(typeof(value))
    if (value === "") return
    if (value<1 || value>28) window.alert("День отправки сообщений должен быть в пределах от 1 до 28, потому что месяцы разной длины");
    else setDay(value);
  }

  return (
    <div class="profile-container">
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

      <div className="profile-container border wide">
        <p>Напоминания о передаче показаний</p>
        <div className="profile-value-row">
          <Toggle
            backgroundColor="black"
            onChange={() => setSend(!sendNotifications)}
          />
          <span>Отправлять напоминания на почту</span>
        </div>
        <div className="profile-value-row">
          <TimePicker onChange={setTime} value={time} />
          <span>время отправки напоминания</span>
        </div>

        <div className="profile-value-row">
          <input type="number" max="28" min="1" onChange={(event) => CheckRange(event.target.value)} value={day} />
          <span>Дата отправки напоминания (от 1 до 28 числа)</span>
        </div>

        <button className="basic-button black-button">обновить настройки напоминаний</button>
      </div>
    </div>
  );
};

export default ProfilePage;