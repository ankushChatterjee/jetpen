import React, { useEffect, useState } from "react";
import config from "../../config.json";
import { Helmet } from "react-helmet";
import { clearToken, getToken } from "../../utils/utils";
import {
  Button,
  PageHeader,
  Modal,
  Form,
  Input,
  notification,
  Space,
} from "antd";
import NewsletterList from "./NewsletterList";
import EditUserDrawer from "./EditUserDrawer";
import { useMediaQuery } from "react-responsive";
import { PlusOutlined, EditOutlined } from "@ant-design/icons";
import { useHistory } from "react-router-dom";

const REDIRECT = 0;
const DONE = 1;
const ERROR = 2;
const LOADING = 3;

export default function Dashboard() {
  const [form] = Form.useForm();
  const [status, setStatus] = useState(LOADING);
  const [userDetail, setUserDetail] = useState({
    name: "",
    email: "",
  });

  const [isModalVisible, setModalVisible] = useState(false);
  const [confirmLoading, setConfirmLoading] = useState(false);
  const [editDrawerVisible, setEditDrawerVisible] = useState(false);
  const isSmallScreen = useMediaQuery({ maxWidth: 820 });
  const token = getToken();
  const history = useHistory();
  const openErrorNotification = (msg, desc) => {
    notification.error({
      message: msg,
      description: desc,
      placement: "bottomRight",
    });
  };

  const openEditDrawer = () => {
    setEditDrawerVisible(true);
  };

  const closeEditDrawer = () => {
    setEditDrawerVisible(false);
  };

  useEffect(() => {
    if (token == null) {
      window.location.assign("/login");
    }
    fetch(config.BASE_URL + "/auth/get-user/", {
      headers: new Headers({
        Authorization: "Bearer " + token,
      }),
    })
      .then((res) => {
        if (res.status == 400) {
          clearToken();
          window.location.assign("/login");
          return;
        } else if (res.status != 200) {
          setStatus(ERROR);
          return;
        }
        return res.json();
      })
      .then((data) => {
        if (typeof data === "undefined") {
          return;
        }
        setUserDetail({
          name: data.name,
          email: data.email,
        });
        setStatus(DONE);
      });
  }, []);

  const createNewsLetter = (_) => {
    const name = form.getFieldValue("name");
    const description = form.getFieldValue("description");

    if (name.length == 0 || description.length == 0) {
      openErrorNotification("Fields are empty", "Please enter proper data");
      return;
    }
    setConfirmLoading(true);
    const formData = new FormData();
    formData.append("name", name);
    formData.append("description", description);

    fetch(config.BASE_URL + "/manage/newsletter/create", {
      method: "POST",
      body: formData,
      headers: new Headers({
        Authorization: "Bearer " + token,
      }),
    })
      .then((res) => {
        setConfirmLoading(false);
        return res.json();
      })
      .then((data) => {
        if (typeof data === "undefined") {
          return;
        }
        console.log(data);
        if (data.error) {
          console.log(data);
          form.resetFields();
          openErrorNotification("Error Registering", data.error);
        } else {
          history.push("/newsletter/" + data.id);
        }
      })
      .catch((err) => {
        console.log(err);
        openErrorNotification("Error adding newsletter", "");
        setConfirmLoading(false);
      });
  };

  if (status === ERROR) {
    return <h1>Cannot display the page, please try refreshing</h1>;
  }

  const PageHeaderExtra = (
    <div>
      {!isSmallScreen ? (
        <Space>
          <Button
            icon={<PlusOutlined />}
            onClick={() => setModalVisible(true)}
            size="large"
            type="text"
          >
            New Newsletter
          </Button>
          <strong>{userDetail.name}</strong>
          <Button icon={<EditOutlined />} onClick={openEditDrawer} />
        </Space>
      ) : (
        <Space>
          <Button
            icon={<PlusOutlined />}
            onClick={() => setModalVisible(true)}
            size="large"
            type="text"
          />
          <Button icon={<EditOutlined />} onClick={openEditDrawer} />
        </Space>
      )}
    </div>
  );

  return (
    <div>
      <Helmet>
        <meta charSet="utf-8" />
        <title>Jetpen Dashboard</title>
      </Helmet>
      <EditUserDrawer
        visible={editDrawerVisible}
        {...userDetail}
        onClose={closeEditDrawer}
        setParentUserDetails={setUserDetail}
      />
      <PageHeader
        title="Dashboard"
        subTitle="Your newsletters"
        extra={[PageHeaderExtra]}
      />
      <Modal
        title="Create new newsletter"
        visible={isModalVisible}
        onCancel={() => setModalVisible(false)}
        onOk={createNewsLetter}
        confirmLoading={confirmLoading}
      >
        <Form layout="vertical" form={form}>
          <Form.Item
            label="Newsletter name"
            name="name"
            rules={[
              { required: true, message: "Please input your newsletter name!" },
            ]}
          >
            <Input />
          </Form.Item>
          <Form.Item label="Description" name="description">
            <Input />
          </Form.Item>
        </Form>
      </Modal>
      <NewsletterList />
    </div>
  );
}
