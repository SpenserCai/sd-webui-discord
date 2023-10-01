/*
 * @Author: SpenserCai
 * @Date: 2023-10-01 17:40:44
 * @version: 
 * @LastEditors: SpenserCai
 * @LastEditTime: 2023-10-01 17:59:20
 * @Description: file content
 */
import axios from "axios";
import js_cookie from "js-cookie";

// 创建 axios 实例
const service = axios.create({
  baseURL: "/api", // api base_url
  timeout: 30000 // 请求超时时间
});

// request 拦截器
service.interceptors.request.use(
 async (config) => {
    const token = js_cookie.get("token");
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
    if (res.code === 401) {
      window.location.href = service.defaults.baseURL + "/login";
    } else if (res.code === -100) {
        alert(res.msg);
    } else {
        if (res.code < 0) {
            alert(res.msg);
        }
        return res;
    }
  },
  error => {
    return Promise.reject(error);
  }
);

