import React, { useState } from "react";
import { useLocation } from "react-router-dom";
import { Button, Typography, Card } from "antd";
import config from "../../config.json";
import { openErrorNotification, openSuccessNotification } from "../../utils/utils";
import { Helmet } from "react-helmet";

export default function Unsubscribe() {
  const { Title } = Typography;
  const [loading, setLoading] = useState(false);
  const location = useLocation();

  const unsubscribe = () => {
    fetch(config.BASE_URL + "/sub/unsub/" + location.search)
      .then((res) => {
        if (res.status !== 200) {
          return res.json();
        }
      })
      .then((data) => {
        if (typeof data === "undefined") {
          setLoading(false);
          openSuccessNotification("Unsubscribed", "Hope to have you back someday");
        } else if(data.error) {
          openErrorNotification("Error unsubscribing", data.error);
        }
      });
  };

  return (
    <div className="center gradient-background" style={{ height: "100vh" }}>
      <Helmet>
        <meta charSet="utf-8" />
        <title>Are you sure to unsubscribe?</title>
      </Helmet>

      <Card
        title={
          <div>
            <Title level={2}>Are you sure to Unsubscribe?</Title>
          </div>
        }
      >
        <Button onClick={unsubscribe} loading={loading} size="large" type="primary">
          Unsubscribe
        </Button>
      </Card>
    </div>
  );
}
