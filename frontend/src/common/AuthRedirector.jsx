import authService from "../services/auth.service";
import {Navigate, useLocation} from "react-router-dom";
import React from "react";

const pagesWithNoAuth = ['/login', '/register']

export const AuthRedirector = ({children}) => {
    const user = authService.getCurrentUser()
    const location = useLocation()
    const path = location.pathname

    if (!user) {
        if (!pagesWithNoAuth.includes(path)) {
            return <Navigate to='/login' state={{from: location}} replace/>
        }
    } else if (pagesWithNoAuth.includes(path) || (!user.verified && path !== '/')) {
        if (location.state?.from)
            return <Navigate to={location.state.from} replace/>
        return <Navigate to="/" replace/>
    }

    return children
}