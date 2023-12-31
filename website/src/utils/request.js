/*
 * @Author: SpenserCai
 * @Date: 2023-10-01 17:40:44
 * @version: 
 * @LastEditors: SpenserCai
 * @LastEditTime: 2023-10-11 11:55:22
 * @Description: file content
 */
import axios from "axios";
import { notify } from "notiwind"
import Cookies from "js-cookie";
// 创建 axios 实例
const service = axios.create({
  baseURL: "/api", // api base_url
  timeout: 30000 // 请求超时时间
});

// request 拦截器
service.interceptors.request.use(
 async (config) => {
    const token = Cookies.get("token");
    if (token) {
      config.headers["Authorization"] = "Bearer " + token;
    }
    return config;
  },
  error => {
    // 请求错误处理
    Promise.reject(error);
  }
);

// response 拦截器, 如果时401则跳转到登录页面
service.interceptors.response.use(
  async (response) => {
    const res = response.data;
    if (res.code === -100) {
        notify({
            title: "Error",
            text: res.message,
            type: "error",
            group: "authenticated",
        }, 5000)
    } else {
        if (res.code < 0) {
            notify({
                title: "Error",
                text: res.message,
                type: "error",
                group: "authenticated",
            }, 5000)
        }
        return res;
    }
  },
  error => {
    if (error.response.status === 401) {
      window.location.href = service.defaults.baseURL + "/login";
    }
    return Promise.reject(error);
  }
);

export default service;

