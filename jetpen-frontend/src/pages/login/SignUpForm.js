import React, { useState } from "react";
import { notification, Form, Input, Button, Row, Col } from "antd";
import config from "../../config.json";
import {useMediaQuery} from "react-responsive";

function validateEmail(email) {
  if (
    /^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9-]+(?:\.[a-zA-Z0-9-]+)*$/.test(
      email
    )
  )
    return true;
  return false;
}

export default function SignUpForm() {
  const [form] = Form.useForm();
  const [isRegistered, setRegistered] = useState(false);
  const isSmallScreen = useMediaQuery({ maxWidth: 820 });

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
    const email = form.getFieldValue("email");
    const name = form.getFieldValue("name");

    if (!validateEmail(email)) {
      openErrorNotification("Error Registering", "Please enter proper email");
      return;
    }
    const formData = new FormData();
    formData.append("username", username);
    formData.append("password", password);
    formData.append("email", email);
    formData.append("name", name);
    fetch(config.BASE_URL + "/auth/register", {
      method: "post",
      body: formData,
    }).then((res) => {
      if(res.status != 200) {
        return res.json();  
      } else {
        setRegistered(true);
      }
    }).then((data) => {
      if(typeof data === "undefined") {
        return;
      }
      if(data.error) {
        form.resetFields();
        openErrorNotification("Error Registering", data.error);
      }
    });
  };
  const onFinishFailed = () => {
    openErrorNotification("Error Registering", "Please fill up the form");
  };

  let formView = (
    <div>
      <h1>Sign Up</h1>
      <Form
        onFinish={onFinish}
        onFinishFailed={onFinishFailed}
        wrapperCol={{ span: 24 }}
        layout="vertical"
        form={form}
      >
        <Form.Item
          rules={[{ required: true, message: "Please enter your username" }]}
          name="username"
          label="Username"
        >
          <Input placeholder="Username" />
        </Form.Item>
        <Form.Item
          rules={[{ required: true, message: "Please enter your name" }]}
          name="name"
          label="Name"
        >
          <Input placeholder="Name" />
        </Form.Item>
        <Form.Item
          rules={[{ required: true, message: "Please enter your email" }]}
          name="email"
          label="Email"
        >
          <Input placeholder="Email" />
        </Form.Item>
        <Form.Item
          rules={[{ required: true, message: "Please enter your password" }]}
          name="password"
          label="Password"
        >
          <Input.Password placeholder="Password" />
        </Form.Item>
        <Form.Item>
          <Button type="primary" htmlType="submit">
            Let's Go!
          </Button>
        </Form.Item>
      </Form>
    </div>
  );

  if (isRegistered) {
    formView = (
      <div>
        <h1>Welcome to Jetpen!</h1>
        <p>
          We are glad to have you! <br />
          Please check your email for verification. <br />
          <br />
          <strong>Thanks,</strong>
          <br />
          <strong>The Jetpen Team</strong>
          <br />
        </p>
      </div>
    );
  }

  return (
    <div className={!isSmallScreen?"pattern-background":""}>
      <Row>
        <Col span={isSmallScreen?24:8}>
          <div
            style={{
              height: "100%",
              background: "#f8f9fa",
              padding: "46px 44px",
            }}
          >
            {formView}
          </div>
        </Col>
      </Row>
    </div>
  );
}
