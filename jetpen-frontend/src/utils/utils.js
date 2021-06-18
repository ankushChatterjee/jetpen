import jwt_decode from "jwt-decode";
import Cookies from "js-cookie";
import {notification} from "antd";

const JETPEN_TOKEN = "jetpenToken";

export function getToken() {
  const token = Cookies.get(JETPEN_TOKEN);
  if (typeof token !== "undefined") return token;
  else return null;
}

export function setToken(token) {
    Cookies.set(JETPEN_TOKEN, token, { expires: 14 });
}

export function clearToken() {
    Cookies.remove(JETPEN_TOKEN);
}

export function decodeToken(jwt) {
  const decoded = jwt_decode(jwt);
  console.log(decoded);
}


export const openErrorNotification = (msg, desc) => {
  notification.error({
    message: msg,
    description: desc,
    placement: "bottomRight",
  });
};

export const openSuccessNotification = (msg, desc) => {
  notification.success({
    message: msg,
    description: desc,
    placement: "bottomRight",
  });
};