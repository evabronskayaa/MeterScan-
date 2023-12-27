import {RemoveTrailingSlash} from "./RemoveTrailingSlash";
import {AuthRedirector} from "./AuthRedirector";
import {Outlet, ScrollRestoration} from "react-router-dom";
import React, {useEffect, useState} from "react";
import authService from "../services/auth.service";
import {UserContext} from "./useUser";

const Root = () => {
    const [user, setUser] = useState(null)
    useEffect(() => {
        authService.me()?.then(setUser).catch(() => {
            setUser(null)
            authService.logout()
        })
    }, [])

    return <UserContext.Provider value={{user, setUser}}>
        <RemoveTrailingSlash/>
        <AuthRedirector>
            <Outlet/>
        </AuthRedirector>
        <ScrollRestoration/>
    </UserContext.Provider>
}

export default Root