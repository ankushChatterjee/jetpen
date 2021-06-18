import React, { useState, useEffect } from "react";
import { useParams } from "react-router-dom";
import { Button, Input, Typography, Card, Space } from "antd";
import config from "../../config.json";
import { getToken } from "../../utils/utils";
import { Helmet } from "react-helmet";
import CenteredSpin from "../../components/CenteredSpin";
import {
  openErrorNotification,
  openSuccessNotification,
} from "../../utils/utils";

const LOADING = 0;
const DONE = 1;
const ERROR = 2;

export default function Subscribe() {
  const { id } = useParams();
  const [status, setStatus] = useState(LOADING);
  const [newsletterName, setName] = useState("");
  const [description, setDescription] = useState("");
  const [email, setEmail] = useState("");
  const [subLoading, setSubLoading] = useState(false);
  const { Title, Paragraph } = Typography;

  useEffect(() => {
    const token = getToken();
    fetch(config.BASE_URL + `/manage/newsletter/${id}`)
      .then((res) => {
        if (res.status != 200) {
          setStatus(ERROR);
        }
        return res.json();
      })
      .then((data) => {
        if (typeof data === "undefined") {
          return;
        }
        if (data.error) {
          return;
        }
        console.log(data);
        setStatus(DONE);
        setName(data.Name);
        setDescription(data.Description);
      });
  }, []);

  const subscribe = () => {
    setSubLoading(true);
    const data = {
      nid: id,
      email: email,
    };
    fetch(config.BASE_URL + "/sub/", {
      method: "POST",
      body: JSON.stringify(data),
      headers: { "Content-Type": "application/json" },
    })
      .then((res) => {
        if (res.status !== 200) {
          setStatus(ERROR);
          return res.json();
        }
      })
      .then((data) => {
        if (typeof data === "undefined") {
          setSubLoading(false);
          setEmail("");
          openSuccessNotification("Subscribed!");
          return;
        }
        if (data.error) {
          openErrorNotification("Error Subscribing", data.error);
        } else {
          openSuccessNotification("Subscribed!");
        }
        setSubLoading(false);
        setEmail("");
      })
      .catch((err) => {
        console.log(err);
        openErrorNotification("Error Subscribing");
        setSubLoading(false);
      });
  };

  if (status === ERROR) {
    return <h1>Cannot display the page, please try refreshing</h1>;
  }

  if (status === LOADING) {
    return (
      <div>
        <Helmet>
          <meta charSet="utf-8" />
          <title>Jetpen</title>
        </Helmet>
        <CenteredSpin />
      </div>
    );
  }

  return (
    <div className="center gradient-background" style={{ height: "100vh" }}>
      <Helmet>
        <meta charSet="utf-8" />
        <title>Subscribe to {newsletterName}</title>
      </Helmet>

      <Card
        title={
          <div>
            <Title level={2}>Subscribe to {newsletterName}</Title>
            <Title level={5}>{description}</Title>
          </div>
        }
      >
        <Space>
          <Input
            size="large"
            value={email}
            onChange={(e) => {
              setEmail(e.target.value);
            }}
            placeholder="Enter your email address"
          />
          <Button
            size="large"
            loading={subLoading}
            type="primary"
            onClick={subscribe}
          >
            Subscribe
          </Button>
        </Space>
      </Card>
    </div>
  );
}
