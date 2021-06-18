import React, { useState, useEffect } from "react";
import { Modal, Input, Button, Spin, notification } from "antd";
import ReactQuill from "react-quill";
import { SendOutlined, SaveOutlined } from "@ant-design/icons";
import config from "../../config.json";
import { getToken } from "../../utils/utils";
import {openSuccessNotification, openErrorNotification} from "../../utils/utils";

const LOADING = 0;
const ERROR = 1;
const DONE = 2;

export default function LetterModal({
  visible,
  editable,
  letterId,
  sub,
  nid,
  onClose,
  onCancel
}) {
  const [status, setStatus] = useState(LOADING);
  const [value, setContent] = useState("");
  const [subject, setSubject] = useState(sub);
  let footerContent = [];

  const saveDraft = () => {
    const data = {
      id: letterId,
      content: value,
      subject: subject,
    };
    const token = getToken();
    fetch(config.BASE_URL + `/manage/letter/save/`, {
      method: "POST",
      body: JSON.stringify(data),
      headers: new Headers({
        Authorization: "Bearer " + token,
        "content-type": "application/json",
      }),
    })
      .then((res) => {
        if (res.status != 200) {
          return res.json();
        }
        return;
      })
      .then((data) => {
        if (typeof data === "undefined") {
          openSuccessNotification("Letter saved!");
          onClose(letterId, false, subject);
          return;
        }
        if (data.error) {
          openErrorNotification("Cannot Save letter", data.error);
        }
      });
  }
 
  const send = () => {
    const data = {
      id: letterId,
      content: value,
      subject: subject,
      isPublished: true,
      nid: nid
    };
    const token = getToken();
    fetch(config.BASE_URL + `/manage/letter/publish/`, {
      method: "POST",
      body: JSON.stringify(data),
      headers: new Headers({
        Authorization: "Bearer " + token,
        "content-type": "application/json",
      }),
    })
      .then((res) => {
        if (res.status != 200) {
          return res.json();
        }
        return;
      })
      .then((data) => {
        if (typeof data === "undefined") {
          openSuccessNotification("Sent Letter!");
          onClose(letterId, true, subject);
          return;
        }
        if (data.error) {
          openErrorNotification("Cannot Send letter", data.error);
        }
      });
  };

  useEffect(() => {
    if (sub) {
      setSubject(sub);
    }
  }, [sub]);

  useEffect(() => {
    if(!letterId)
      return;
    console.log("hello");
    const token = getToken();
    fetch(config.BASE_URL + `/manage/letter/${letterId}`, {
      headers: new Headers({
        Authorization: "Bearer " + token,
      }),
    })
      .then((res) => {
        if (res.status != 200) {
          setStatus(ERROR);
          return;
        }
        return res.json();
      })
      .then((data) => {
        if (typeof data === "undefined") {
          return;
        }
        setStatus(DONE);
        setContent(data.content);
      })
      .catch(() => setStatus(ERROR));
  },[letterId]);

  if (editable) {
    footerContent = [
      <Button icon={<SaveOutlined />} onClick={saveDraft}>Save draft</Button>,
      <Button icon={<SendOutlined />} type="primary" onClick={send}>
        Send
      </Button>,
    ];
  } else {
    footerContent = [<Button type="primary" onClick={() => onClose(letterId, true, subject)}>Close</Button>];
  }

  let modalContent = null;
  switch (status) {
    case ERROR:
      modalContent = <h1>There has been an error, try later</h1>;
      break;
    case DONE:
      modalContent = (
        <ReactQuill
          style={{
            height: "30vh",
          }}
          readOnly={!editable}
          theme="snow"
          value={value}
          onChange={setContent}
        />
      );
      break;
    case LOADING:
      modalContent = <Spin />;
  }

  return (
    <Modal
      title="Letter"
      visible={visible}
      onCancel={onCancel}
      onOk={false}
      confirmLoading={false}
      footer={footerContent}
    >
      <Input
        size="large"
        value={subject}
        onChange={(e) => {
          setSubject(e.target.value);
        }}
        bordered={editable}
        disabled={!editable}
        placeholder="Newsletter Subject"
      />

      {modalContent}
      <div style={{ height: "2em" }}></div>
    </Modal>
  );
}
