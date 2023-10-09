/*
 * @Author: SpenserCai
 * @Date: 2023-10-02 21:33:34
 * @version: 
 * @LastEditors: SpenserCai
 * @LastEditTime: 2023-10-09 21:26:30
 * @Description: file content
 */
import request from "@/utils/request";

export function userinfo() {
  return request({
    url: "/user_info",
    method: "get"
  });
}

export function userhistory(body) {
  return request({
    url: "/user_history",
    method: "post",
    data: body
  });
}

export function communityhistory(body) {
  return request({
    url: "/community_history",
    method: "post",
    data: body
  });
}