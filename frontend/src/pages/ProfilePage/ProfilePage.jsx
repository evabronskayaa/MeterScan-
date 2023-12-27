import React, {useEffect, useState} from "react";
import "./ProfilePage.scss";
import Toggle from "react-toggle";
import "react-toggle/style.css";
import axios from "axios";
import authHeader from "../../services/auth-header";

const API_URL = "http://localhost/api/v1/";

const ProfilePage = () => {
  const [time, setTime] = useState(10);
  const [day, setDay] = useState(1);
  const [sendNotifications, setSend] = useState(false);

  useEffect(() => {
    axios.get(API_URL + "settings", {
      headers: authHeader()
    })
        .then(r => r.data)
        .then(data => {
          setSend(data.notification_enabled)
          if (data.notification_day_of_month !== 0) {
            setDay(data.notification_day_of_month)
          }
          if (data.notification_hour !== 0) {
            setTime(data.notification_hour)
          }
        })
  }, [])

  const checkRangeDayOfMonth = (value) =>{
    if (value === "") return
    if (value<1 || value>28) window.alert("День отправки сообщений должен быть в пределах от 1 до 28, потому что месяцы разной длины");
    else setDay(value);
  }

  const checkRangeHour = (value) =>{
    if (value === "") return
    if (value<0 || value>23) window.alert("Время отправки сообщений должно быть в пределах от 0 до 23");
    else setTime(value);
  }

  const updateSettings = () => {
    const form = new FormData()
    form.append("enabled", sendNotifications)
    form.append("day", day)
    form.append("hour", time)

    axios.put(API_URL + "settings/notification", form, {
      headers: authHeader()
    })
  }

  return <div class="profile-container">
    <div className="title">Личный кабинет</div>
    <div className="profile-container border wide">
      <p>Напоминания о передаче показаний</p>
      <div className="profile-value-row">
        <Toggle
            checked={sendNotifications}
            backgroundColor="black"
            onChange={() => setSend(!sendNotifications)}
        />
        <span>Отправлять напоминания на почту</span>
      </div>
      <div className="profile-value-row">
        <input type="number" max="23" min="0" onChange={e => checkRangeHour(e.target.value)} value={time}/>
        <span>время отправки напоминания (от 0 до 23)</span>
      </div>

      <div className="profile-value-row">
        <input type="number" max="28" min="1" onChange={e => checkRangeDayOfMonth(e.target.value)} value={day}/>
        <span>Дата отправки напоминания (от 1 до 28 числа)</span>
      </div>

      <button className="basic-button black-button" onClick={updateSettings}>обновить настройки напоминаний</button>
    </div>
  </div>;
};

export default ProfilePage;