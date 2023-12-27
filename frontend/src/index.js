import React from "react";
import './index.css';
import ReactDOM from 'react-dom/client'
import RegisterPage from "./pages/RegisterPage/RegisterPage";
import LoginPage from "./pages/LoginPage/LoginPage";
import {createBrowserRouter, createRoutesFromElements, Route, RouterProvider} from "react-router-dom";
import HistoryPage from "./pages/HistoryPage/HistoryPage";
import MenuPage from "./pages/MenuPage/MenuPage";
import ProfilePage from "./pages/ProfilePage/ProfilePage";
import "./App.scss";
import "bootstrap/dist/css/bootstrap.min.css";
import MainLayout from "./components/MainLayout/MainLayout";
import RecognizePage from "./pages/RecognizePage/RecognizePage";
import Root from "./common/Root";

const router = createBrowserRouter(createRoutesFromElements(<Route path="/" Component={Root}>
    <Route path="/register" Component={RegisterPage}/>
    <Route path="/login" Component={LoginPage}/>

    <Route path="/" Component={MainLayout}>
        <Route path="/" Component={MenuPage}/>
        <Route path="/profile" Component={ProfilePage}/>
        <Route path="/recognize" Component={RecognizePage}/>
        <Route path="/history" Component={HistoryPage}/>
    </Route>
</Route>))

ReactDOM.createRoot(document.getElementById('root')).render(<React.StrictMode>
    <RouterProvider router={router}/>
</React.StrictMode>)