import "./App.scss";
import MainPage from "./pages/MainPage/MainPage";
import LoginPage from "./pages/LoginPage/LoginPage";
import RegisterPage from "./pages/RegisterPage/RegisterPage";
import authService from "./services/auth.service";
import "bootstrap/dist/css/bootstrap.min.css";
import { useState } from "react";

function App() {
  const user = authService.getCurrentUser();
  var [page, setPage] = useState("register")

  if (user)
  return (<MainPage email={user}/>);
  else if (page === "register")
  return (<RegisterPage redirect={setPage}/>)
  else return (<LoginPage redirect={setPage}/>)
;
}

export default App;
