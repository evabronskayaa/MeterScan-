import React, { useState } from "react";
import { Link } from "react-router-dom";
import AuthService from "../../services/auth.service";
import "../../styles/form.scss";

const LoginPage = (props) => {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [loading, setLoading] = useState(false);
  const [message, setMessage] = useState("");

  const required = (value) => {
    if (!value) {
      return (
        <div className="alert alert-danger" role="alert">
          This field is required!
        </div>
      );
    }
  };

  const handleLogin = (e) => {
    e.preventDefault();
    setMessage("");
    setLoading(true);

    if (!email || !password) {
      setLoading(false);
      return setMessage("Both username and password are required.");
    }

    AuthService.login(email, password, "string").then(
      () => {
        window.location.reload();
        console.log("auth");
      },
      (error) => {
        const resMessage =
          (error.response &&
            error.response.data &&
            error.response.data.message) ||
          error.message ||
          error.toString();
        setLoading(false);
        setMessage(resMessage);
      }
    );
  };

  return (
    <div className="col-md-12 container">
      <p className="center">MeterScan+</p>
      <div className="card card-container">
        <p className="center">Логин</p>
        <form onSubmit={handleLogin}>
          <div className="form-group">
            <label htmlFor="email">Почта</label>
            <input
              type="text"
              className="form-control"
              name="username"
              value={email}
              onChange={(e) => setEmail(e.target.value)}
              validations={[required]}
            />
          </div>

          <div className="form-group">
            <label htmlFor="password">Пароль</label>
            <input
              type="password"
              className="form-control"
              name="password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              validations={[required]}
            />
          </div>

          <div className="form-group">
            <button className="btn btn-dark btn-block" disabled={loading}>
              {loading && (
                <span className="spinner-border spinner-border-sm"></span>
              )}
              <span>Войти</span>
            </button>
          </div>

          {message && (
            <div className="form-group">
              <div className="alert alert-danger" role="alert">
                {message}
              </div>
            </div>
          )}

          <input type="checkbox" style={{ display: "none" }} />
          
          <Link to="/register">
          <p className="no-acc center"
          onClick={() =>
            props.redirect("register")
          }
        >
          еще нет аккаунта?
        </p></Link>
        </form>
      </div>
        
    </div>
  );
};

export default LoginPage;
