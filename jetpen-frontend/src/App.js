import React, { useState, useEffect } from "react";
import "./App.less";
import {
  BrowserRouter as Router,
  Switch,
  Route,
  Redirect,
  useHistory,
} from "react-router-dom";
import Login from "./pages/login/Login";
import VerifyEmail from "./pages/verifyEmail/VerifyEmail";
import Dashboard from "./pages/dashboard/Dashboard";
import NewsletterDetail from "./pages/newsletter/NewsletterDetail";
import Subscribe from "./pages/sub/Subscribe";
import Unsubscribe from "./pages/sub/Unsubscribe";
import { getToken } from "./utils/utils";

const ONLY_LOGIN = 0;
const LOADING = 1;
const HOME = 2;

const App = () => {
  let status = HOME;
  if (getToken() == null) {
    status = ONLY_LOGIN;
  }

  const OnlyLoginRoutes = (
    <Switch>
      <Route exact path="/login" component={Login} />
      <Route exact path="/verify-email" component={VerifyEmail} />
      <Route exact path="/sub/:id" component={Subscribe} />
      <Route exact path="/unsub/" component={Unsubscribe} />
      <Redirect to="/login" />
    </Switch>
  );

  const AllRoutes = (
    <Switch>
      <Route path="/verify-email" exact component={VerifyEmail} />
      <Route path="/" exact component={Dashboard} />
      <Route path="/newsletter/:id" exact component={NewsletterDetail} />
      <Route exact path="/sub/:id" component={Subscribe} />
      <Route exact path="/unsub/" component={Unsubscribe} />
    </Switch>
  );

  let returnView = OnlyLoginRoutes;

  switch (status) {
    case ONLY_LOGIN:
      returnView = OnlyLoginRoutes;
      break;
    case HOME:
      returnView = AllRoutes;
      break;
  }

  return <Router>{returnView}</Router>;
};

export default App;
