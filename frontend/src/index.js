import React from "react";
import './index.css';
import ReactDOM from 'react-dom/client'
import RegisterPage from "./pages/RegisterPage/RegisterPage";
import LoginPage from "./pages/LoginPage/LoginPage";
import {
    createBrowserRouter,
    createRoutesFromElements,
    Navigate, Outlet,
    Route,
    RouterProvider, ScrollRestoration,
    useLocation
} from "react-router-dom";
import MainPage from "./pages/MainPage/MainPage";
import authService from "./services/auth.service";
import "./App.scss";
import "bootstrap/dist/css/bootstrap.min.css";

export const RemoveTrailingSlash = ({ ...rest }) => {
    const location = useLocation()

    // If the last character of the url is '/'
    if (location.pathname.match('/.*/$')) {
        return <Navigate replace { ...rest } to={ {
            pathname: location.pathname.replace(/\/+$/, ''),
            search: location.search
        } }/>
    } else return null
}

const pagesWithNoAuth = ['/login', '/register']

const AuthRedirector = ({children}) => {
    const user = authService.getCurrentUser()
    const location = useLocation()
    const path = location.pathname

    if (!user) {
        if (!pagesWithNoAuth.includes(path)) {
            return <Navigate to='/login' state={{from: location}} replace/>
        }
    } else if (pagesWithNoAuth.includes(path)) {
        if (location.state?.from)
            return <Navigate to={ location.state.from } replace/>
        return <Navigate to="/" replace/>
    }

    return children
}

const Root = () => {
    return <>
        <RemoveTrailingSlash/>
        <AuthRedirector>
            <Outlet/>
        </AuthRedirector>
        <ScrollRestoration/>
    </>
}

const router = createBrowserRouter(createRoutesFromElements(<Route path="/" Component={Root}>
    <Route path="/register" Component={RegisterPage}/>
    <Route path="/login" Component={LoginPage}/>
    <Route path="/" Component={MainPage}/>
</Route>))

ReactDOM.createRoot(document.getElementById('root')).render(<React.StrictMode>
    <RouterProvider router={router}/>
</React.StrictMode>)