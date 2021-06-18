import React, { useState, useEffect } from "react";
import { Drawer, Typography, Input, Modal, Space, Button } from "antd";
import {
  getToken,
  openErrorNotification,
  openSuccessNotification,
  clearToken,
} from "../../utils/utils";
import { ExclamationCircleOutlined } from "@ant-design/icons";
import config from "../../config.json";
import { useHistory } from "react-router-dom";

export default function EditUserDrawer({
  onClose,
  visible,
  email,
  name,
  setParentUserDetails,
}) {
  const { Paragraph } = Typography;
  const [editName, setEditName] = useState("");
  const [editPassword, setEditPassword] = useState("");
  const history = useHistory();

  function logout() {
    clearToken();
    history.go(0);
  }

  function confirmPasswordChange() {
    Modal.confirm({
      title: "Confirm password change",
      icon: <ExclamationCircleOutlined />,
      content:
        "Are you sure you want to change the password? You would be logged out",
      okText: "Yes",
      cancelText: "No",
      onOk: () => {
        editProfile(editPassword, "password");
      },
    });
  }

  const editProfile = (content, type) => {
    const url = config.BASE_URL + `/auth/edit-${type}/`;

    const data = {
      content: content,
    };
    const token = getToken();

    fetch(url, {
      method: "POST",
      body: JSON.stringify(data),
      headers: new Headers({
        Authorization: "Bearer " + token,
        "content-type": "application/json",
      }),
    })
      .then((res) => {
        if (res.status !== 200) {
          return res.json();
        }
      })
      .then((data) => {
        if (typeof data === "undefined") {
          if (type == "password") {
            logout();
          }
          openSuccessNotification(type + " Edited");
          if (type === "name") {
            setParentUserDetails({
              name: content,
              email: email,
            });
          }
        } else {
          openErrorNotification("Edit failed", data.error);
        }
      })
      .catch((err) => {
        console.error(err);
        openErrorNotification("Unable to edit");
      });
  };

  useEffect(() => {
    if (name) {
      setEditName(name);
    }
  }, [name]);
  return (
    <Drawer
      title="Profile"
      placement="bottom"
      closable={true}
      onClose={onClose}
      visible={visible}
    >
      <Paragraph
        editable={{
          onChange: (val) => {
            setEditName(val);
            editProfile(val, "name");
          },
        }}
      >
        {editName}
      </Paragraph>
      <Paragraph>{email}</Paragraph>
      <Space>
        <Input.Password
          onChange={(e) => setEditPassword(e.target.value)}
          placeholder="New password"
        />
        <Button onClick={confirmPasswordChange}>Change Password</Button>
        <Button type="primary" onClick={logout}>
          Logout
        </Button>
      </Space>
    </Drawer>
  );
}
