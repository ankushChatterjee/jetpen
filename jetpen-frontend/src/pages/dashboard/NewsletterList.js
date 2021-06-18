import React, { useEffect, useState } from "react";
import config from "../../config.json";
import { getToken } from "../../utils/utils";
import { List, Button, Skeleton } from "antd";
import { ArrowRightOutlined } from "@ant-design/icons";
import {useMediaQuery} from "react-responsive";
import { Link } from "react-router-dom";

const LOADING = 0;
const ERROR = 1;
const DONE = 2;
const LIMIT = 10;

export default function NewsletterList() {
  const [status, setStatus] = useState(LOADING);
  const [newsletters, setNewsletters] = useState([]);
  const [contentLoading, setContentLoading] = useState(false);
  const [cursor, setCursor] = useState("");
  const [allLoaded, setAllLoaded] = useState(false);
  const isSmallScreen = useMediaQuery({ maxWidth: 820 });

  useEffect(() => {
    const token = getToken();
    fetch(config.BASE_URL + `/manage/newsletters?limit=${LIMIT}`, {
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
        if (data.newsletters.length < LIMIT) {
          setAllLoaded(true);
        }
        setStatus(DONE);
        console.log(data);
        setNewsletters(data.newsletters);
        setCursor(data.nextCursor);
      }).catch(() => setStatus(ERROR));
  }, []);

  if (status == ERROR) {
    return <h1>Error</h1>;
  }

  const onLoadMore = (event) => {
    const token = getToken();
    setContentLoading(true);
    fetch(
      config.BASE_URL + `/manage/newsletters?limit=${LIMIT}&cursor=${cursor}`,
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
        const newData = data.newsletters.map((item) => {
          return { ...item, loaded: true };
        });
        console.log("new1", newData);
        if (data.newsletters.length < LIMIT) {
          setAllLoaded(true);
        }
        const newList = newsletters.concat(newData);
        console.log("new", newList);
        setNewsletters(newList);
        setContentLoading(false);
        setCursor(data.nextCursor);
      }).catch(() => {
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
    <div style={{
      padding:isSmallScreen?"1em":"3em"
    }}>
      <List
        size="large"
        loading={status == LOADING}
        loadMore={LoadMoreButton}
        dataSource={newsletters}
        renderItem={(item) => (
          <List.Item key={item.Id}>
            <Skeleton avatar={false} title loading={item.loading} active>
              <List.Item.Meta
                title={<h2 style={{color:"#b51b62"}}>{item.Name}</h2>}
                description={item.Description}
              />
              <div>
                <Link to={`/newsletter/${item.Id}`}>
                  <Button
                    shape="circle"
                    icon={<ArrowRightOutlined />}
                    size="large"
                  />
                </Link>
              </div>
            </Skeleton>
          </List.Item>
        )}
      />
    </div>
  );
}
