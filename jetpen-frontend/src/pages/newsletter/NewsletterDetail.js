import React, { useEffect, useState, useRef } from "react";
import { useParams } from "react-router-dom";
import config from "../../config.json";
import { getToken } from "../../utils/utils";
import { Layout } from "antd";
import CenteredSpin from "../../components/CenteredSpin";
import { Helmet } from "react-helmet";
import NewsletterEditor from "./NewsletterEditor";
import LettersList from "./LettersList";

const LOADING = 0;
const DONE = 1;
const ERROR = 2;

export default function NewsletterDetail() {
  const { id } = useParams();
  const [status, setStatus] = useState(LOADING);
  const [newsletterName, setName] = useState("");
  const [newsletterDescription, setDescription] = useState("");
  const [letters, setLetters] = useState([]);
  const { Footer } = Layout;

  useEffect(() => {
    const token = getToken();
    fetch(config.BASE_URL + `/manage/newsletter/${id}`, {
      headers: new Headers({
        Authorization: "Bearer " + token,
      }),
    })
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
    <div>
      <Helmet>
        <meta charSet="utf-8" />
        <title>{newsletterName}</title>
      </Helmet>
      <NewsletterEditor newsletterName={newsletterName} newsletterDescription={newsletterDescription} nid={id} letters={letters}  setLetters={setLetters} />
      <LettersList nid={id} letters={letters} setLetters={setLetters} />
      <Footer style={{ marginTop: "3em" }}>
        Best of luck for you next letter!
      </Footer>
    </div>
  );
}
