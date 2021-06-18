import React, { useState, useEffect } from "react";
import { Drawer, Typography, Input, Modal, Space, Button } from "antd";
import {
  getToken,
  openErrorNotification,
  openSuccessNotification,
} from "../../utils/utils";
import { DeleteOutlined } from "@ant-design/icons";
import { useHistory } from "react-router-dom";
import config from "../../config.json";

export default function EditNewsletterDrawer({
  visible,
  onClose,
  name,
  description,
  nid,
  updateNewsletterDetails,
}) {
  const [editName, setEditName] = useState("");
  const [editDescription, setEditDescription] = useState("");
  const { Paragraph, Title } = Typography;
  const history = useHistory();

  function confirmNewsletterDelete() {
    Modal.confirm({
      title: "Delete newsletter",
      icon: <DeleteOutlined />,
      content:
        "Are you sure you want to delete this newsletter. It cannot be overturned",
      okText: "Yes",
      cancelText: "No",
      onOk: () => {
        deleteNewsletter();
      },
    });
  }

  const deleteNewsletter = () => {
    const data = {
      nid: nid,
    };
    const token = getToken();

    fetch(config.BASE_URL + `/manage/newsletter/delete/`, {
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
          history.push("/");
        } else if (data.error) {
          openErrorNotification("Delete failed", data.error);
        }
      })
      .catch((err) => {
        console.error(err);
        openErrorNotification("Unable to edit");
      });
  };

  const editNewsletter = (content, type) => {
    const url = config.BASE_URL + `/manage/newsletter/edit-${type}/`;
    const data = {
      content: content,
      nid: nid,
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
          openSuccessNotification(type + " Edited");
          if (type === "name") {
            updateNewsletterDetails(content, editDescription);
          } else if (type === "description") {
            updateNewsletterDetails(editName, content);
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
    if (description) {
      setEditDescription(description);
    }
  }, [name, description]);

  return (
    <div>
      <Drawer
        title="Edit Newsletter"
        placement="bottom"
        closable={true}
        onClose={onClose}
        visible={visible}
      >
        <Title level={5}>Newsletter Name</Title>
        <Paragraph
          editable={{
            onChange: (val) => {
              setEditName(val);
              editNewsletter(val, "name");
            },
          }}
        >
          {editName}
        </Paragraph>
        <Title level={5}>Newsletter Description</Title>
        <Paragraph
          editable={{
            onChange: (val) => {
              setEditDescription(val);
              editNewsletter(val, "description");
            },
          }}
        >
          {editDescription}
        </Paragraph>
        <Button
          icon={<DeleteOutlined />}
          onClick={confirmNewsletterDelete}
          danger
          type="primary"
        >
          Delete Newsletter
        </Button>
      </Drawer>
    </div>
  );
}
