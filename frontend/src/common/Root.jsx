import {RemoveTrailingSlash} from "./RemoveTrailingSlash";
import {AuthRedirector} from "./AuthRedirector";
import {Outlet, ScrollRestoration} from "react-router-dom";
import React, {useEffect} from "react";
import authService from "../services/auth.service";

const Root = () => {
    useEffect(() => {
        authService.me()
    }, [])

    return <>
        <RemoveTrailingSlash/>
        <AuthRedirector>
            <Outlet/>
        </AuthRedirector>
        <ScrollRestoration/>
    </>
}

export default Root