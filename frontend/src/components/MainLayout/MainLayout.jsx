import React from "react";
import {NavLink, Outlet} from "react-router-dom";
import authService from "../../services/auth.service";
import "./MainLayout.scss";

const MainLayout = () => {
    const user = authService.getCurrentUser();
    const handleLogout = (e) => {
        authService.logout();
        window.location.reload();
    };

    return <div>
        <div className="header">
            <NavLink to="/">
                <div className="logo">MeterScan+</div>
            </NavLink>
            <div className="header-right-container">
                <div className="address">{user.email}</div>
                <NavLink to="/profile">
                    <button className="profile ml-2">Профиль</button>
                </NavLink>
                <button className="logout ml-2" onClick={handleLogout}>Выйти</button>
            </div>
        </div>
        <div>
            <Outlet/>
        </div>
    </div>
}

export default MainLayout