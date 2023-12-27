import {Navigate, useLocation} from "react-router-dom";
import React from "react";
import useUser from "./useUser";

const pagesWithNoAuth = ['/login', '/register']

export const AuthRedirector = ({children}) => {
    const {user} = useUser()
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