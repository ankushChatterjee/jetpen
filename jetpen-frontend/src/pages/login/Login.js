import React, { useState } from "react";
import { Layout, Form, Input, PageHeader, Button, Row, Col } from "antd";
import SignUpForm from "./SignUpForm";
import LoginForm from "./LoginForm";
import { Helmet } from "react-helmet";

export default function Login() {
  const { Footer } = Layout;
  return (
    <div>
      <Helmet>
        <meta charSet="utf-8" />
        <title>Jetpen</title>
      </Helmet>
      <PageHeader
        title="JetPen"
        subTitle="The world's best junk email generator"
        extra={[<LoginForm />]}
      />
      <SignUpForm />
      <Footer
        style={{
          position: "fixed",
          bottom: "0",
          width: "100%",
        }}
      >
        Created By CodeCap
      </Footer>
    </div>
  );
}
