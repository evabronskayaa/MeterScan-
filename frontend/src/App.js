import "./App.scss";
import MainPage from "./pages/MainPage/MainPage";
import LoginPage from "./pages/LoginPage/LoginPage";
import authService from "./services/auth.service";

function App() {
  const user = authService.getCurrentUser();

  if (user)
  return (<MainPage email={user}/>);
  else 
  return (<LoginPage/>)
}

export default App;
