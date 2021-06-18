import React, { useState, useEffect } from "react";
import { Button, PageHeader, Input, Tooltip } from "antd";
import ReactQuill from "react-quill";
import "react-quill/dist/quill.snow.css";
import { getToken } from "../../utils/utils";
import {
  SendOutlined,
  SaveOutlined,
  ShareAltOutlined,
  EditOutlined
} from "@ant-design/icons";
import { useHistory } from "react-router-dom";
import config from "../../config.json";
import {
  openErrorNotification,
  openSuccessNotification,
} from "../../utils/utils";
import EditNewsletterDrawer from "./EditNewsletterDrawer";

function NewsletterEditor({ newsletterName, newsletterDescription, nid, letters, setLetters }) {
  const [value, setValue] = useState("");
  const [subject, setSubject] = useState("");
  const [editDrawerVisible, setEditDrawerVisible] = useState(false);
  const [name, setName] = useState("");
  const [description, setDescription] = useState("");
  const history = useHistory();

  const openEditDrawer = () => {
    setEditDrawerVisible(true);
  };

  const closeEditDrawer = () => {
    setEditDrawerVisible(false);
  };

  const updateNewsletterDetails = (name, description) => {
    console.log(name, description);
    setName(name);
    setDescription(description);
  }

  useEffect(() => {
    if(newsletterName) {
      setName(newsletterName);
    }
    if(newsletterDescription) {
      setDescription(newsletterDescription);
    }
  }, [newsletterName, newsletterDescription]);

  const send = (shouldSend) => {
    const data = {
      nid: nid,
      content: value,
      subject: subject,
      isPublished: shouldSend,
    };
    const token = getToken();
    fetch(config.BASE_URL + `/manage/letter/create/`, {
      method: "POST",
      body: JSON.stringify(data),
      headers: new Headers({
        Authorization: "Bearer " + token,
        "content-type": "application/json",
      }),
    })
      .then((res) => {
        return res.json();
      })
      .then((data) => {
        console.log(data);
        if (data.error) {
          openErrorNotification("Cannot Send letter", data.error);
        } else {
          if (shouldSend) openSuccessNotification("Sent Letter!");
          else openSuccessNotification("Letter saved");
          const newLetterObj = {
            Id: data.id,
            Nid: nid,
            IsPublished: shouldSend,
            Subject: subject,
          }
          const lettersCopy = [...letters]
          lettersCopy.unshift(newLetterObj);
          setLetters(lettersCopy);
        }
      }).catch( err => {
        openErrorNotification("Cannot complete operation");
      });
  };

  return (
    <PageHeader
      onBack={() => history.goBack()}
      title={name}
      subTitle={description}
      extra={[
        <Tooltip placement="topLeft" title="Copy Link to subscribe">
          <Button
            icon={<ShareAltOutlined />}
            onClick={() => {
              navigator.clipboard.writeText(`${config.HOSTNAME}/sub/${nid}`);
              openSuccessNotification("Copied link for your newsletter");
            }}
            type="link"
          />
        </Tooltip>,
        <Button icon={<EditOutlined />} onClick={openEditDrawer} />,
        <Button icon={<SaveOutlined />} onClick={() => send(false)}>
          Save as draft
        </Button>,
        <Button
          icon={<SendOutlined />}
          onClick={() => send(true)}
          type="primary"
        >
          Send
        </Button>,
      ]}
    >
      <EditNewsletterDrawer updateNewsletterDetails={updateNewsletterDetails} visible={editDrawerVisible}  onClose={closeEditDrawer} name={newsletterName} description={newsletterDescription} nid={nid} />
      <div>
        <Input
          size="large"
          value={subject}
          onChange={(e) => {
            setSubject(e.target.value);
          }}
          placeholder="Newsletter Subject"
        />
        <ReactQuill
          style={{
            height: "50vh",
          }}
          theme="snow"
          value={value}
          onChange={setValue}
        />
      </div>
    </PageHeader>
  );
}

export default NewsletterEditor;
