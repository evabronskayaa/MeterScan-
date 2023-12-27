import React from "react";
import "./HistoryPage.scss";

const HistoryPage = () => {
  return <div className="profile-container">
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
  </div>
};

export default HistoryPage;
