import {createContext, useContext} from "react";

export const UserContext = createContext({
    user: null,
    setUser: (user) => {
    },
})

export default () => useContext(UserContext)