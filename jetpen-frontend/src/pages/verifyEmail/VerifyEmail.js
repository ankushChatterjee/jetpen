import React, { useEffect, useState } from "react";
import { Result, Button, Spin } from "antd";
import { Link, useLocation } from "react-router-dom";
import config from "../../config.json";
import { Helmet } from "react-helmet";

function useQuery() {
  return new URLSearchParams(useLocation().search);
}

export default function VerifyEmail() {
  const [isLoading, setLoading] = useState(true);
  const [status, setStatus] = useState("error");
  const query = useQuery();

  useEffect(() => {
    fetch(
      config.BASE_URL +
        "/auth/email-verification?token=" +
        query.get("token") +
        "&username=" +
        query.get("username")
    ).then((res) => {
      if (res.status != 200) {
        setLoading(false);
      } else {
        setStatus("success");
        setLoading(true);
      }
    });
  });

  let resultView = (
    <div>
      <Result
        status={status}
        title={
          status == "error"
            ? "Sorry, The verification did not work"
            : "Welcome to JetPen!!"
        }
        subTitle={
          status == "error"
            ? "Please check the URL"
            : "Your EMail Is verified, lets start creating together"
        }
        extra={[
          <Link to="/">
            <Button type="primary" key="console">
              Login
            </Button>
          </Link>,
        ]}
      />
    </div>
  );
  if (isLoading) {
    resultView = <Spin />;
  }

  return (
    <div>
      <Helmet>
        <meta charSet="utf-8" />
        <title>Email Verification</title>
      </Helmet>
      {resultView}
    </div>
  );
}
