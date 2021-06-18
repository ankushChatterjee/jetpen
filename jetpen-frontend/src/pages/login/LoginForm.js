import React, { useState, useEffect } from "react";
import { notification, Form, Input, Button } from "antd";
import config from "../../config.json";
import { setToken } from "../../utils/utils";
import {useMediaQuery} from "react-responsive";

export default function LoginForm() {
  const [form] = Form.useForm();
  const [loading, setLoading] = useState(false);
  const isSmallScreen = useMediaQuery({ maxWidth: 820 });

  const loginWithUsernameAndPassword = async (username, password) => {
    const formData = new FormData();
    formData.append("username", username);
    formData.append("password", password);
    const data = await fetch(config.BASE_URL + "/auth/login", {
      body: formData,
      method: "post",
    }).then((result) => {
      return result.json();
    });
    return data;
  };

  const openErrorNotification = (msg, desc) => {
    notification.error({
      message: msg,
      description: desc,
      placement: "bottomRight",
    });
  };
  const onFinish = () => {
    const username = form.getFieldValue("username");
    const password = form.getFieldValue("password");

    if (username.length == 0 || password.length == 0) {
      openErrorNotification(
        "Username or password is Empty",
        "Please enter proper data"
      );
      return;
    }
    setLoading(true);
    loginWithUsernameAndPassword(username, password)
      .then((data) => {
        setLoading(false);
        if (data.error) {
          openErrorNotification("Error Logging in", data.error);
        } else {
          setToken(data.token);
          window.location.assign("/");
        }
      })
      .catch((err) => {
        setLoading(false);
        openErrorNotification("Error Logging in", "");
      });
  };
  const onFinishFailed = (evt) => {
    console.log("fail", evt);
  };
  let loginButton = (
    <Button type="primary" htmlType="submit">
      Login
    </Button>
  );
  if (loading) {
    loginButton = (
      <Button type="primary" htmlType="submit" loading>
        Login
      </Button>
    );
  }
  return (
    <Form
      onFinish={onFinish}
      onFinishFailed={onFinishFailed}
      layout={isSmallScreen?"horizontal":"inline"}
      form={form}
    >
      <Form.Item  name="username">
        <Input placeholder="Username" />
      </Form.Item>
      <Form.Item name="password">
        <Input.Password placeholder="Password" />
      </Form.Item>
      <Form.Item>{loginButton}</Form.Item>
    </Form>
  );
}
