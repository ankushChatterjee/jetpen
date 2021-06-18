import React, { useState, useEffect } from "react";
import { List, Divider, Skeleton, Button, Modal } from "antd";
import config from "../../config.json";
import { getToken } from "../../utils/utils";
import {
  HistoryOutlined,
  PlusOutlined,
  HourglassOutlined,
} from "@ant-design/icons";
import LetterModal from "./LetterModal";

const LOADING = 0;
const ERROR = 1;
const DONE = 2;
const LIMIT = 10;

export default function LettersList({ nid, letters, setLetters }) {
  const [status, setStatus] = useState(LOADING);
  const [contentLoading, setContentLoading] = useState(false);
  const [cursor, setCursor] = useState("");
  const [allLoaded, setAllLoaded] = useState(false);
  const [modalState, setModalState] = useState({
    editable: false,
    letterId:"",
    visible:false,
    sub:""
  });

  const openModal = (id, isEditable, subject) => {
    setModalState({
      editable: isEditable,
      letterId: id,
      visible: true,
      sub: subject
    });
  }

  const onModalClose = (letterId, isPublished, subject) => {
    let lettersCopy = [...letters];
    let lid = 0, i=0;
    letters.forEach(l => {
      if(l.Id === letterId) {
        lid = i;
      }
      i++;
    });
    let letterCopy ={...lettersCopy[lid]};
    letterCopy.Subject = subject;
    letterCopy.IsPublished = isPublished;
    lettersCopy[lid] = letterCopy;
    setLetters(lettersCopy);
    setModalState({
      editable: false,
      letterId: "",
      visible: false,
      sub:""
    });
  }

  const closeModal = () => {
    setModalState({
      editable: false,
      letterId: "",
      visible: false,
      sub:""
    });
  }

  useEffect(() => {
    const token = getToken();
    fetch(config.BASE_URL + `/manage/letters/${nid}?limit=${LIMIT}`, {
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
        if (data.error) {
          setStatus(ERROR);
          return;
        }
        if (data.letters.length < LIMIT) {
          setAllLoaded(true);
        }
        setStatus(DONE);
        console.log(data);
        setLetters(data.letters);
        setCursor(data.nextCursor);
      })
      .catch(() => setStatus(ERROR));
  }, []);

  const ListHeader = (
    <div
      style={{
        display: "flex",
        justifyContent: "center",
        alignItems: "center",
      }}
    >
      <h3>
        <strong>
          <HistoryOutlined style={{ padding: "7px" }} />
          Past Newsletters
        </strong>
      </h3>
    </div>
  );

  const onLoadMore = (event) => {
    const token = getToken();
    setContentLoading(true);
    fetch(
      config.BASE_URL +
        `/manage/letters/${nid}?limit=${LIMIT}&cursor=${cursor}`,
      {
        headers: new Headers({
          Authorization: "Bearer " + token,
        }),
      }
    )
      .then((res) => {
        if (res.status != 200) {
          setStatus(ERROR);
          return;
        }
        return res.json();
      })
      .then((data) => {
        if (data.error) {
          setStatus(ERROR);
          return;
        }
        setStatus(DONE);
        const newData = data.letters.map((item) => {
          return { ...item, loaded: true };
        });
        if (data.letters.length < LIMIT) {
          setAllLoaded(true);
        }
        const newList = letters.concat(newData);
        console.log("new", newList);
        setLetters(newList);
        setContentLoading(false);
        setCursor(data.nextCursor);
      })
      .catch(() => {
        setStatus(ERROR);
      });
  };

  const LoadMoreButton = allLoaded ? null : (
    <div
      style={{
        textAlign: "center",
        marginTop: 12,
        height: 32,
        lineHeight: "32px",
        padding: "12px",
      }}
    >
      <Button
        loading={contentLoading}
        disabled={contentLoading}
        onClick={onLoadMore}
      >
        Load More
      </Button>
    </div>
  );

  return (
    <div>
      <Divider />
      <LetterModal {...modalState} onClose={onModalClose} onCancel={closeModal} nid={nid} />
      <List
        style={{
          paddingLeft: "3em",
          paddingRight: "3em",
        }}
        loading={status == LOADING}
        loadMore={LoadMoreButton}
        dataSource={letters}
        header={ListHeader}
        renderItem={(item) => (
          <List.Item key={item.Id}>
            <Skeleton avatar={false} title loading={item.loading} active>
              <List.Item.Meta title={item.Subject} />
              <div>
                {!item.IsPublished ? (
                  <HourglassOutlined
                    style={{ padding: "10px", color: "#999" }}
                  />
                ) : null}
                <Button icon={<PlusOutlined />} onClick={() => openModal(item.Id, !item.IsPublished, item.Subject)}>More</Button>
              </div>
            </Skeleton>
          </List.Item>
        )}
      />
    </div>
  );
}
